package capsule

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
)

const ChunkTreeSchemaVersion = "arkmesh.chunk-tree/v0alpha1"

// ChunkTreeFile carries every leaf digest of one object so a damaged replica can
// find exactly which chunks are wrong.
//
// It needs no signature of its own. Every leaf is checked by recomputing the root
// and comparing it with the chunk root inside the signed manifest, so a forged
// tree cannot survive verification.
type ChunkTreeFile struct {
	SchemaVersion string   `json:"schema_version"`
	CapsuleID     string   `json:"capsule_id"`
	AssetSHA256   string   `json:"asset_sha256"`
	ChunkSize     int64    `json:"chunk_size"`
	ChunkCount    int      `json:"chunk_count"`
	ChunkRoot     string   `json:"chunk_root"`
	Leaves        []string `json:"leaves"`
}

// ChunkDamage names one chunk that does not match the committed tree.
type ChunkDamage struct {
	Index  int
	Offset int64
	Length int64
	Reason string
}

func BuildChunkTreeFile(capsuleID string, asset Asset, objectPath string) (ChunkTreeFile, error) {
	if asset.ChunkRoot == "" || !ValidChunkSize(asset.ChunkSize) {
		return ChunkTreeFile{}, errors.New("asset does not carry a chunk commitment")
	}
	digest, size, tree, err := ScanObject(objectPath, asset.ChunkSize)
	if err != nil {
		return ChunkTreeFile{}, err
	}
	if digest != asset.SHA256 || size != asset.Size || tree.Root != asset.ChunkRoot {
		return ChunkTreeFile{}, errors.New("stored object does not match its signed commitment")
	}
	leaves := make([]string, 0, len(tree.leaves))
	for _, leaf := range tree.leaves {
		leaves = append(leaves, hex.EncodeToString(leaf[:]))
	}
	return ChunkTreeFile{
		SchemaVersion: ChunkTreeSchemaVersion,
		CapsuleID:     capsuleID,
		AssetSHA256:   asset.SHA256,
		ChunkSize:     asset.ChunkSize,
		ChunkCount:    tree.Count,
		ChunkRoot:     tree.Root,
		Leaves:        leaves,
	}, nil
}

func WriteChunkTreeFile(path string, tree ChunkTreeFile) error {
	return writeExclusiveJSON(path, tree, "chunk tree")
}

func LoadChunkTreeFile(path string) (ChunkTreeFile, error) {
	var tree ChunkTreeFile
	if err := loadStrictRegularJSON(path, &tree, "chunk tree"); err != nil {
		return ChunkTreeFile{}, err
	}
	return tree, nil
}

// VerifyChunkTreeFile authenticates a tree against the signed manifest and
// returns the asset it describes.
func VerifyChunkTreeFile(manifest Manifest, tree ChunkTreeFile) (Asset, error) {
	if tree.SchemaVersion != ChunkTreeSchemaVersion {
		return Asset{}, fmt.Errorf("unsupported chunk tree schema %q", tree.SchemaVersion)
	}
	if tree.CapsuleID != manifest.CapsuleID {
		return Asset{}, errors.New("chunk tree capsule ID does not match manifest")
	}
	var asset Asset
	found := false
	for _, candidate := range manifest.Assets {
		if candidate.SHA256 == tree.AssetSHA256 {
			asset = candidate
			found = true
			break
		}
	}
	if !found {
		return Asset{}, fmt.Errorf("capsule has no asset %q", tree.AssetSHA256)
	}
	if asset.ChunkRoot == "" || !ValidChunkSize(asset.ChunkSize) {
		return Asset{}, errors.New("asset does not carry a chunk commitment")
	}
	if tree.ChunkSize != asset.ChunkSize || tree.ChunkRoot != asset.ChunkRoot {
		return Asset{}, errors.New("chunk tree does not match the signed commitment")
	}
	count := ChunkCount(asset.Size, asset.ChunkSize)
	if tree.ChunkCount != count || len(tree.Leaves) != count {
		return Asset{}, fmt.Errorf("chunk tree declares %d leaves, want %d", len(tree.Leaves), count)
	}
	leaves := make([][32]byte, 0, count)
	for index, leaf := range tree.Leaves {
		decoded, err := hex.DecodeString(leaf)
		if err != nil {
			return Asset{}, fmt.Errorf("decode leaf %d: %w", index, err)
		}
		if len(decoded) != 32 {
			return Asset{}, fmt.Errorf("leaf %d has invalid length %d", index, len(decoded))
		}
		var digest [32]byte
		copy(digest[:], decoded)
		leaves = append(leaves, digest)
	}
	if chunkRoot(leaves) != asset.ChunkRoot {
		return Asset{}, errors.New("chunk tree leaves do not reconstruct the signed chunk root")
	}
	return asset, nil
}

// ScanChunkDamage reports exactly which chunks of a local object disagree with
// an authenticated tree. An empty result means the object is intact.
func ScanChunkDamage(objectPath string, asset Asset, tree ChunkTreeFile) ([]ChunkDamage, error) {
	file, err := os.Open(objectPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buffer := make([]byte, asset.ChunkSize)
	var damage []ChunkDamage
	for index := 0; index < tree.ChunkCount; index++ {
		offset := int64(index) * asset.ChunkSize
		length := expectedChunkLength(asset.Size, asset.ChunkSize, index)
		read, err := file.ReadAt(buffer[:length], offset)
		if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
			return nil, err
		}
		if int64(read) != length {
			damage = append(damage, ChunkDamage{Index: index, Offset: offset, Length: length, Reason: "missing bytes"})
			continue
		}
		if hex.EncodeToString(chunkLeafDigestBytes(buffer[:length])) != tree.Leaves[index] {
			damage = append(damage, ChunkDamage{Index: index, Offset: offset, Length: length, Reason: "digest mismatch"})
		}
	}
	return damage, nil
}

// RepairChunks copies only the damaged chunks from a donor object. Every donor
// chunk is verified against the authenticated tree BEFORE it is written, so a
// hostile donor cannot use repair to inject bytes.
func RepairChunks(objectPath string, asset Asset, tree ChunkTreeFile, donorPath string, damage []ChunkDamage) (int, error) {
	if len(damage) == 0 {
		return 0, nil
	}
	donor, err := os.Open(donorPath)
	if err != nil {
		return 0, err
	}
	defer donor.Close()
	target, err := os.OpenFile(objectPath, os.O_RDWR, 0o644)
	if err != nil {
		return 0, err
	}
	defer target.Close()

	buffer := make([]byte, asset.ChunkSize)
	repaired := 0
	for _, item := range damage {
		if item.Index < 0 || item.Index >= tree.ChunkCount {
			return repaired, fmt.Errorf("chunk index %d is out of range", item.Index)
		}
		length := expectedChunkLength(asset.Size, asset.ChunkSize, item.Index)
		offset := int64(item.Index) * asset.ChunkSize
		read, err := donor.ReadAt(buffer[:length], offset)
		if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
			return repaired, err
		}
		if int64(read) != length {
			return repaired, fmt.Errorf("donor is missing chunk %d", item.Index)
		}
		if hex.EncodeToString(chunkLeafDigestBytes(buffer[:length])) != tree.Leaves[item.Index] {
			return repaired, fmt.Errorf("donor chunk %d does not match the committed tree", item.Index)
		}
		if _, err := target.WriteAt(buffer[:length], offset); err != nil {
			return repaired, err
		}
		repaired++
	}
	if err := target.Truncate(asset.Size); err != nil {
		return repaired, err
	}
	if err := target.Sync(); err != nil {
		return repaired, err
	}
	return repaired, nil
}

func chunkLeafDigestBytes(chunk []byte) []byte {
	digest := chunkLeafDigest(chunk)
	return digest[:]
}
