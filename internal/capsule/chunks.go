package capsule

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
)

const (
	ChunkProofSchemaVersion = "arkmesh.chunk-proof/v0alpha1"
	DefaultChunkSize        = int64(1) << 20
	MinChunkSize            = int64(4096)
	MaxChunkSize            = int64(1) << 26
)

const (
	chunkLeafDomain  = "arkmesh.chunk.leaf/v0alpha1\n"
	chunkNodeDomain  = "arkmesh.chunk.node/v0alpha1\n"
	chunkEmptyDomain = "arkmesh.chunk.empty/v0alpha1\n"
)

// ChunkTree is a Merkle commitment over fixed size chunks of one stored object.
//
// Unpaired nodes are promoted unchanged instead of being duplicated. Duplicating
// an unpaired node lets two different chunk sequences produce one root, which is
// the malleability weakness reported against Bitcoin's Merkle construction.
type ChunkTree struct {
	ChunkSize int64
	Count     int
	Root      string
	leaves    [][32]byte
}

// ChunkProof proves that one exact chunk belongs to a signed chunk root.
// The chunk bytes are carried, so verification also proves possession.
type ChunkProof struct {
	SchemaVersion string   `json:"schema_version"`
	CapsuleID     string   `json:"capsule_id"`
	AssetSHA256   string   `json:"asset_sha256"`
	ChunkSize     int64    `json:"chunk_size"`
	ChunkCount    int      `json:"chunk_count"`
	ChunkIndex    int      `json:"chunk_index"`
	ChunkRoot     string   `json:"chunk_root"`
	Chunk         string   `json:"chunk"`
	Path          []string `json:"path"`
}

func ValidChunkSize(size int64) bool {
	return size >= MinChunkSize && size <= MaxChunkSize && size&(size-1) == 0
}

func WriteChunkProof(path string, proof ChunkProof) error {
	return writeExclusiveJSON(path, proof, "chunk proof")
}

func LoadChunkProof(path string) (ChunkProof, error) {
	var proof ChunkProof
	if err := loadStrictRegularJSON(path, &proof, "chunk proof"); err != nil {
		return ChunkProof{}, err
	}
	return proof, nil
}

func ChunkCount(size, chunkSize int64) int {
	if size <= 0 {
		return 0
	}
	return int((size + chunkSize - 1) / chunkSize)
}

// ScanObject streams one object once, returning both its whole content digest
// and its chunk tree, so verification never needs a second read.
func ScanObject(path string, chunkSize int64) (string, int64, ChunkTree, error) {
	if !ValidChunkSize(chunkSize) {
		return "", 0, ChunkTree{}, fmt.Errorf("invalid chunk size %d", chunkSize)
	}
	file, err := os.Open(path)
	if err != nil {
		return "", 0, ChunkTree{}, err
	}
	defer file.Close()

	whole := sha256.New()
	buffer := make([]byte, chunkSize)
	var size int64
	var leaves [][32]byte
	for {
		read, err := io.ReadFull(file, buffer)
		if read > 0 {
			chunk := buffer[:read]
			whole.Write(chunk)
			leaves = append(leaves, chunkLeafDigest(chunk))
			size += int64(read)
		}
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			break
		}
		if err != nil {
			return "", 0, ChunkTree{}, err
		}
	}
	tree := ChunkTree{ChunkSize: chunkSize, Count: len(leaves), leaves: leaves}
	tree.Root = chunkRoot(leaves)
	return hex.EncodeToString(whole.Sum(nil)), size, tree, nil
}

func BuildChunkProof(capsuleID string, asset Asset, objectPath string, index int) (ChunkProof, error) {
	if asset.ChunkRoot == "" || !ValidChunkSize(asset.ChunkSize) {
		return ChunkProof{}, errors.New("asset does not carry a chunk commitment")
	}
	digest, size, tree, err := ScanObject(objectPath, asset.ChunkSize)
	if err != nil {
		return ChunkProof{}, err
	}
	if digest != asset.SHA256 || size != asset.Size || tree.Root != asset.ChunkRoot {
		return ChunkProof{}, errors.New("stored object does not match its signed commitment")
	}
	if index < 0 || index >= tree.Count {
		return ChunkProof{}, fmt.Errorf("chunk index %d is out of range", index)
	}
	path, err := chunkProofPath(tree.leaves, index)
	if err != nil {
		return ChunkProof{}, err
	}
	chunk, err := readChunk(objectPath, asset.ChunkSize, index)
	if err != nil {
		return ChunkProof{}, err
	}
	return ChunkProof{
		SchemaVersion: ChunkProofSchemaVersion,
		CapsuleID:     capsuleID,
		AssetSHA256:   asset.SHA256,
		ChunkSize:     asset.ChunkSize,
		ChunkCount:    tree.Count,
		ChunkIndex:    index,
		ChunkRoot:     tree.Root,
		Chunk:         hex.EncodeToString(chunk),
		Path:          path,
	}, nil
}

// VerifyChunkProof checks a proof against the capsule's own signed manifest.
// Every field an attacker could restate is compared with the signed asset, so
// the proof itself never defines what counts as correct.
func VerifyChunkProof(manifest Manifest, proof ChunkProof) error {
	if proof.SchemaVersion != ChunkProofSchemaVersion {
		return fmt.Errorf("unsupported chunk proof schema %q", proof.SchemaVersion)
	}
	if proof.CapsuleID != manifest.CapsuleID {
		return errors.New("chunk proof capsule ID does not match manifest")
	}
	var asset *Asset
	for index := range manifest.Assets {
		if manifest.Assets[index].SHA256 == proof.AssetSHA256 {
			asset = &manifest.Assets[index]
			break
		}
	}
	if asset == nil {
		return fmt.Errorf("capsule has no asset %q", proof.AssetSHA256)
	}
	if asset.ChunkRoot == "" || !ValidChunkSize(asset.ChunkSize) {
		return errors.New("asset does not carry a chunk commitment")
	}
	if proof.ChunkSize != asset.ChunkSize || proof.ChunkRoot != asset.ChunkRoot {
		return errors.New("chunk proof does not match the signed commitment")
	}
	count := ChunkCount(asset.Size, asset.ChunkSize)
	if proof.ChunkCount != count {
		return fmt.Errorf("chunk proof count %d does not match signed count %d", proof.ChunkCount, count)
	}
	if proof.ChunkIndex < 0 || proof.ChunkIndex >= count {
		return fmt.Errorf("chunk index %d is out of range", proof.ChunkIndex)
	}
	chunk, err := hex.DecodeString(proof.Chunk)
	if err != nil {
		return fmt.Errorf("decode chunk bytes: %w", err)
	}
	expected := expectedChunkLength(asset.Size, asset.ChunkSize, proof.ChunkIndex)
	if int64(len(chunk)) != expected {
		return fmt.Errorf("chunk %d has length %d, want %d", proof.ChunkIndex, len(chunk), expected)
	}
	return verifyChunkPath(asset.ChunkRoot, count, proof.ChunkIndex, chunkLeafDigest(chunk), proof.Path)
}

func expectedChunkLength(size, chunkSize int64, index int) int64 {
	remaining := size - int64(index)*chunkSize
	if remaining < chunkSize {
		return remaining
	}
	return chunkSize
}

func readChunk(path string, chunkSize int64, index int) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	buffer := make([]byte, chunkSize)
	read, err := file.ReadAt(buffer, int64(index)*chunkSize)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}
	if read == 0 {
		return nil, fmt.Errorf("chunk %d is empty", index)
	}
	return buffer[:read], nil
}

func chunkLeafDigest(chunk []byte) [32]byte {
	hasher := sha256.New()
	hasher.Write([]byte(chunkLeafDomain))
	hasher.Write(chunk)
	var digest [32]byte
	copy(digest[:], hasher.Sum(nil))
	return digest
}

func chunkNodeDigest(left, right [32]byte) [32]byte {
	hasher := sha256.New()
	hasher.Write([]byte(chunkNodeDomain))
	hasher.Write(left[:])
	hasher.Write(right[:])
	var digest [32]byte
	copy(digest[:], hasher.Sum(nil))
	return digest
}

func chunkRoot(leaves [][32]byte) string {
	if len(leaves) == 0 {
		digest := sha256.Sum256([]byte(chunkEmptyDomain))
		return hex.EncodeToString(digest[:])
	}
	level := append([][32]byte(nil), leaves...)
	for len(level) > 1 {
		next := make([][32]byte, 0, (len(level)+1)/2)
		for index := 0; index < len(level); index += 2 {
			if index+1 == len(level) {
				next = append(next, level[index])
				continue
			}
			next = append(next, chunkNodeDigest(level[index], level[index+1]))
		}
		level = next
	}
	return hex.EncodeToString(level[0][:])
}

func chunkProofPath(leaves [][32]byte, index int) ([]string, error) {
	if index < 0 || index >= len(leaves) {
		return nil, fmt.Errorf("chunk index %d is out of range", index)
	}
	var path []string
	level := append([][32]byte(nil), leaves...)
	position := index
	for len(level) > 1 {
		sibling := position ^ 1
		if sibling < len(level) {
			path = append(path, hex.EncodeToString(level[sibling][:]))
		}
		next := make([][32]byte, 0, (len(level)+1)/2)
		for cursor := 0; cursor < len(level); cursor += 2 {
			if cursor+1 == len(level) {
				next = append(next, level[cursor])
				continue
			}
			next = append(next, chunkNodeDigest(level[cursor], level[cursor+1]))
		}
		level = next
		position /= 2
	}
	return path, nil
}

// verifyChunkPath derives sibling direction from the index and the level sizes
// implied by the signed chunk count, so no attacker supplied direction bits exist.
func verifyChunkPath(root string, count, index int, leaf [32]byte, path []string) error {
	current := leaf
	position := index
	levelSize := count
	consumed := 0
	for levelSize > 1 {
		sibling := position ^ 1
		if sibling < levelSize {
			if consumed >= len(path) {
				return errors.New("chunk proof path is too short")
			}
			decoded, err := hex.DecodeString(path[consumed])
			if err != nil {
				return fmt.Errorf("decode chunk proof path: %w", err)
			}
			if len(decoded) != sha256.Size {
				return fmt.Errorf("invalid chunk proof path digest length %d", len(decoded))
			}
			var node [32]byte
			copy(node[:], decoded)
			if position%2 == 0 {
				current = chunkNodeDigest(current, node)
			} else {
				current = chunkNodeDigest(node, current)
			}
			consumed++
		}
		position /= 2
		levelSize = (levelSize + 1) / 2
	}
	if consumed != len(path) {
		return errors.New("chunk proof path is too long")
	}
	if hex.EncodeToString(current[:]) != root {
		return errors.New("chunk proof does not reconstruct the signed chunk root")
	}
	return nil
}
