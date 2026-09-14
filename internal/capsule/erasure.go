package capsule

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"

	"arkmesh/internal/erasure"
)

const (
	ErasureSchemaVersion = "arkmesh.erasure/v0alpha1"
	ErasureBlockSize     = int64(1) << 16
)

// ErasurePlan describes how one asset was split into data and parity shards.
//
// The plan carries no signature and needs none. Recovery reconstructs bytes and
// then verifies them against the digest, size, and chunk root inside the SIGNED
// manifest before anything is written into the capsule. A hostile plan or a
// hostile shard set can waste work, but cannot introduce accepted bytes.
type ErasurePlan struct {
	SchemaVersion string   `json:"schema_version"`
	CapsuleID     string   `json:"capsule_id"`
	AssetSHA256   string   `json:"asset_sha256"`
	AssetSize     int64    `json:"asset_size"`
	DataShards    int      `json:"data_shards"`
	ParityShards  int      `json:"parity_shards"`
	ShardSize     int64    `json:"shard_size"`
	BlockSize     int64    `json:"block_size"`
	ShardDigests  []string `json:"shard_digests"`
}

type RecoveryReport struct {
	MissingShards    []int
	ReconstructedFor string
	Bytes            int64
}

func ShardFileName(index int) string {
	return fmt.Sprintf("shard-%03d", index)
}

func WriteErasurePlan(path string, plan ErasurePlan) error {
	return writeExclusiveJSON(path, plan, "erasure plan")
}

func LoadErasurePlan(path string) (ErasurePlan, error) {
	var plan ErasurePlan
	if err := loadStrictRegularJSON(path, &plan, "erasure plan"); err != nil {
		return ErasurePlan{}, err
	}
	return plan, nil
}

// ProtectAsset writes data and parity shards for one asset, streaming one stripe
// at a time so memory stays bounded regardless of asset size.
func ProtectAsset(capsuleID string, asset Asset, objectPath, shardDirectory string, dataShards, parityShards int) (ErasurePlan, error) {
	if asset.Size <= 0 {
		return ErasurePlan{}, errors.New("cannot protect an empty asset")
	}
	if err := verifyObjectAgainstAsset(objectPath, asset); err != nil {
		return ErasurePlan{}, err
	}
	codec, err := erasure.New(dataShards, parityShards)
	if err != nil {
		return ErasurePlan{}, err
	}
	blocks := (asset.Size + int64(dataShards)*ErasureBlockSize - 1) / (int64(dataShards) * ErasureBlockSize)
	shardSize := blocks * ErasureBlockSize

	object, err := os.Open(objectPath)
	if err != nil {
		return ErasurePlan{}, err
	}
	defer object.Close()
	if err := os.MkdirAll(shardDirectory, 0o755); err != nil {
		return ErasurePlan{}, fmt.Errorf("create shard directory: %w", err)
	}

	total := codec.TotalShards()
	files := make([]*os.File, total)
	hashers := make([]hash.Hash, total)
	for index := 0; index < total; index++ {
		path := filepath.Join(shardDirectory, ShardFileName(index))
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if err != nil {
			closeShards(files)
			return ErasurePlan{}, fmt.Errorf("create shard %d: %w", index, err)
		}
		files[index] = file
		hashers[index] = sha256.New()
	}

	shards := make([][]byte, total)
	for index := range shards {
		shards[index] = make([]byte, ErasureBlockSize)
	}
	for block := int64(0); block < blocks; block++ {
		for index := 0; index < dataShards; index++ {
			target := shards[index]
			for offset := range target {
				target[offset] = 0
			}
			offset := int64(index)*shardSize + block*ErasureBlockSize
			if offset < asset.Size {
				length := ErasureBlockSize
				if offset+length > asset.Size {
					length = asset.Size - offset
				}
				if _, err := object.ReadAt(target[:length], offset); err != nil && !errors.Is(err, io.EOF) {
					closeShards(files)
					return ErasurePlan{}, err
				}
			}
		}
		if err := codec.Encode(shards); err != nil {
			closeShards(files)
			return ErasurePlan{}, err
		}
		for index := 0; index < total; index++ {
			if _, err := files[index].Write(shards[index]); err != nil {
				closeShards(files)
				return ErasurePlan{}, fmt.Errorf("write shard %d: %w", index, err)
			}
			if _, err := hashers[index].Write(shards[index]); err != nil {
				closeShards(files)
				return ErasurePlan{}, err
			}
		}
	}
	for index, file := range files {
		if err := file.Close(); err != nil {
			return ErasurePlan{}, fmt.Errorf("close shard %d: %w", index, err)
		}
	}

	plan := ErasurePlan{
		SchemaVersion: ErasureSchemaVersion,
		CapsuleID:     capsuleID,
		AssetSHA256:   asset.SHA256,
		AssetSize:     asset.Size,
		DataShards:    dataShards,
		ParityShards:  parityShards,
		ShardSize:     shardSize,
		BlockSize:     ErasureBlockSize,
		ShardDigests:  make([]string, 0, total),
	}
	for index := 0; index < total; index++ {
		plan.ShardDigests = append(plan.ShardDigests, hex.EncodeToString(hashers[index].Sum(nil)))
	}
	return plan, nil
}

// RecoverAsset rebuilds a missing or damaged object from surviving shards. The
// result is verified against the signed manifest before it replaces anything.
func RecoverAsset(manifest Manifest, plan ErasurePlan, shardDirectory, capsuleRoot string) (RecoveryReport, error) {
	asset, err := validateErasurePlan(manifest, plan)
	if err != nil {
		return RecoveryReport{}, err
	}
	codec, err := erasure.New(plan.DataShards, plan.ParityShards)
	if err != nil {
		return RecoveryReport{}, err
	}
	total := codec.TotalShards()

	files := make([]*os.File, total)
	present := make([]bool, total)
	var missing []int
	defer func() {
		for _, file := range files {
			if file != nil {
				_ = file.Close()
			}
		}
	}()
	for index := 0; index < total; index++ {
		path := filepath.Join(shardDirectory, ShardFileName(index))
		digest, size, err := hashFile(path)
		if err != nil || size != plan.ShardSize || digest != plan.ShardDigests[index] {
			missing = append(missing, index)
			continue
		}
		file, err := os.Open(path)
		if err != nil {
			missing = append(missing, index)
			continue
		}
		files[index] = file
		present[index] = true
	}
	usable := 0
	for _, ok := range present {
		if ok {
			usable++
		}
	}
	if usable < plan.DataShards {
		return RecoveryReport{}, fmt.Errorf("%d usable shards, %d required", usable, plan.DataShards)
	}

	objectsDirectory := filepath.Join(capsuleRoot, "objects")
	if err := os.MkdirAll(objectsDirectory, 0o755); err != nil {
		return RecoveryReport{}, err
	}
	temporary, err := os.CreateTemp(objectsDirectory, ".recovering-*")
	if err != nil {
		return RecoveryReport{}, err
	}
	temporaryPath := temporary.Name()
	success := false
	defer func() {
		if !success {
			_ = os.Remove(temporaryPath)
		}
	}()

	blocks := plan.ShardSize / plan.BlockSize
	shards := make([][]byte, total)
	for index := range shards {
		shards[index] = make([]byte, plan.BlockSize)
	}
	for block := int64(0); block < blocks; block++ {
		blockPresent := make([]bool, total)
		for index := 0; index < total; index++ {
			for offset := range shards[index] {
				shards[index][offset] = 0
			}
			if !present[index] {
				continue
			}
			if _, err := files[index].ReadAt(shards[index], block*plan.BlockSize); err != nil && !errors.Is(err, io.EOF) {
				_ = temporary.Close()
				return RecoveryReport{}, err
			}
			blockPresent[index] = true
		}
		if err := codec.Reconstruct(shards, blockPresent); err != nil {
			_ = temporary.Close()
			return RecoveryReport{}, err
		}
		for index := 0; index < plan.DataShards; index++ {
			offset := int64(index)*plan.ShardSize + block*plan.BlockSize
			if offset >= plan.AssetSize {
				continue
			}
			length := plan.BlockSize
			if offset+length > plan.AssetSize {
				length = plan.AssetSize - offset
			}
			if _, err := temporary.WriteAt(shards[index][:length], offset); err != nil {
				_ = temporary.Close()
				return RecoveryReport{}, err
			}
		}
	}
	if err := temporary.Truncate(plan.AssetSize); err != nil {
		_ = temporary.Close()
		return RecoveryReport{}, err
	}
	if err := temporary.Sync(); err != nil {
		_ = temporary.Close()
		return RecoveryReport{}, err
	}
	if err := temporary.Close(); err != nil {
		return RecoveryReport{}, err
	}

	if err := verifyObjectAgainstAsset(temporaryPath, asset); err != nil {
		return RecoveryReport{}, fmt.Errorf("recovered bytes rejected: %w", err)
	}
	if err := os.Rename(temporaryPath, filepath.Join(objectsDirectory, asset.SHA256)); err != nil {
		return RecoveryReport{}, err
	}
	success = true
	return RecoveryReport{MissingShards: missing, ReconstructedFor: asset.SHA256, Bytes: asset.Size}, nil
}

func validateErasurePlan(manifest Manifest, plan ErasurePlan) (Asset, error) {
	if plan.SchemaVersion != ErasureSchemaVersion {
		return Asset{}, fmt.Errorf("unsupported erasure plan schema %q", plan.SchemaVersion)
	}
	if plan.CapsuleID != manifest.CapsuleID {
		return Asset{}, errors.New("erasure plan capsule ID does not match manifest")
	}
	var asset Asset
	found := false
	for _, candidate := range manifest.Assets {
		if candidate.SHA256 == plan.AssetSHA256 {
			asset = candidate
			found = true
			break
		}
	}
	if !found {
		return Asset{}, fmt.Errorf("capsule has no asset %q", plan.AssetSHA256)
	}
	if plan.AssetSize != asset.Size {
		return Asset{}, errors.New("erasure plan size does not match the signed asset")
	}
	if plan.DataShards < 1 || plan.ParityShards < 1 || plan.DataShards+plan.ParityShards > erasure.MaxShards {
		return Asset{}, errors.New("erasure plan has invalid shard counts")
	}
	if plan.BlockSize <= 0 || plan.BlockSize > ErasureBlockSize || plan.BlockSize&(plan.BlockSize-1) != 0 {
		return Asset{}, fmt.Errorf("erasure plan has invalid block size %d", plan.BlockSize)
	}
	if plan.ShardSize <= 0 || plan.ShardSize%plan.BlockSize != 0 {
		return Asset{}, fmt.Errorf("erasure plan has invalid shard size %d", plan.ShardSize)
	}
	if plan.ShardSize*int64(plan.DataShards) < asset.Size {
		return Asset{}, errors.New("erasure plan shards cannot hold the signed asset")
	}
	if len(plan.ShardDigests) != plan.DataShards+plan.ParityShards {
		return Asset{}, errors.New("erasure plan digest count does not match shard count")
	}
	for index, digest := range plan.ShardDigests {
		if !digestPattern.MatchString(digest) {
			return Asset{}, fmt.Errorf("erasure plan shard %d has an invalid digest", index)
		}
	}
	return asset, nil
}

// verifyObjectAgainstAsset checks stored bytes against every commitment the
// signed manifest carries for that asset.
func verifyObjectAgainstAsset(objectPath string, asset Asset) error {
	if asset.ChunkRoot != "" {
		digest, size, tree, err := ScanObject(objectPath, asset.ChunkSize)
		if err != nil {
			return err
		}
		if digest != asset.SHA256 || size != asset.Size {
			return errors.New("object does not match its signed digest")
		}
		if tree.Root != asset.ChunkRoot {
			return errors.New("object does not match its signed chunk root")
		}
		return nil
	}
	digest, size, err := hashFile(objectPath)
	if err != nil {
		return err
	}
	if digest != asset.SHA256 || size != asset.Size {
		return errors.New("object does not match its signed digest")
	}
	return nil
}

func closeShards(files []*os.File) {
	for _, file := range files {
		if file != nil {
			_ = file.Close()
		}
	}
}
