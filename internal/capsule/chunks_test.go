package capsule

import (
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestPackCommitsChunkRootAndVerifies(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	asset := writeLargeFixture(t, root, "model.gguf", 3*int(MinChunkSize)+17)
	output := filepath.Join(root, "demo.ark")
	manifest, err := PackWithOptions("demo", output, []AssetSource{{Role: "model", Path: asset}}, time.Now(), PackOptions{ChunkSize: MinChunkSize})
	if err != nil {
		t.Fatalf("PackWithOptions() error = %v", err)
	}
	if manifest.Assets[0].ChunkSize != MinChunkSize || len(manifest.Assets[0].ChunkRoot) != 64 {
		t.Fatalf("asset = %+v, want chunk commitment", manifest.Assets[0])
	}
	if got, want := ChunkCount(manifest.Assets[0].Size, MinChunkSize), 4; got != want {
		t.Fatalf("chunk count = %d, want %d", got, want)
	}
	if _, err := Verify(output); err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
}

func TestVerifyDetectsChunkRootMismatch(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	asset := writeLargeFixture(t, root, "model.gguf", 2*int(MinChunkSize))
	output := filepath.Join(root, "demo.ark")
	manifest, err := PackWithOptions("demo", output, []AssetSource{{Role: "model", Path: asset}}, time.Now(), PackOptions{ChunkSize: MinChunkSize})
	if err != nil {
		t.Fatalf("PackWithOptions() error = %v", err)
	}
	manifest.Assets[0].ChunkRoot = strings.Repeat("0", 64)
	manifest.CapsuleID = capsuleID(manifest)
	if err := writeManifest(output, manifest); err != nil {
		t.Fatalf("writeManifest() error = %v", err)
	}
	if _, err := Verify(output); err == nil || !strings.Contains(err.Error(), "chunk root mismatch") {
		t.Fatalf("Verify() error = %v, want chunk root mismatch", err)
	}
}

func TestChunkRootPromotesUnpairedNodeInsteadOfDuplicating(t *testing.T) {
	t.Parallel()
	leaves := [][32]byte{chunkLeafDigest([]byte("a")), chunkLeafDigest([]byte("b")), chunkLeafDigest([]byte("c"))}
	duplicated := append(append([][32]byte(nil), leaves...), leaves[2])
	if chunkRoot(leaves) == chunkRoot(duplicated) {
		t.Fatal("three leaves and a duplicated third leaf share one root")
	}
	if chunkRoot(leaves[:1]) != hex.EncodeToString(func() []byte { digest := leaves[0]; return digest[:] }()) {
		t.Fatal("single leaf root is not the leaf digest")
	}
}

func TestChunkProofRoundTripAndRejections(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	assetPath := writeLargeFixture(t, root, "model.gguf", 5*int(MinChunkSize)+3)
	output := filepath.Join(root, "demo.ark")
	manifest, err := PackWithOptions("demo", output, []AssetSource{{Role: "model", Path: assetPath}}, time.Now(), PackOptions{ChunkSize: MinChunkSize})
	if err != nil {
		t.Fatalf("PackWithOptions() error = %v", err)
	}
	asset := manifest.Assets[0]
	objectPath := filepath.Join(output, "objects", asset.SHA256)
	count := ChunkCount(asset.Size, asset.ChunkSize)
	for index := 0; index < count; index++ {
		proof, err := BuildChunkProof(manifest.CapsuleID, asset, objectPath, index)
		if err != nil {
			t.Fatalf("BuildChunkProof(%d) error = %v", index, err)
		}
		if err := VerifyChunkProof(manifest, proof); err != nil {
			t.Fatalf("VerifyChunkProof(%d) error = %v", index, err)
		}
	}

	valid, err := BuildChunkProof(manifest.CapsuleID, asset, objectPath, 2)
	if err != nil {
		t.Fatalf("BuildChunkProof() error = %v", err)
	}

	tampered := valid
	tampered.Chunk = hex.EncodeToString(make([]byte, MinChunkSize))
	if err := VerifyChunkProof(manifest, tampered); err == nil || !strings.Contains(err.Error(), "does not reconstruct") {
		t.Fatalf("VerifyChunkProof(tampered chunk) error = %v, want reconstruction failure", err)
	}

	moved := valid
	moved.ChunkIndex = 3
	if err := VerifyChunkProof(manifest, moved); err == nil || !strings.Contains(err.Error(), "does not reconstruct") {
		t.Fatalf("VerifyChunkProof(moved index) error = %v, want reconstruction failure", err)
	}

	restated := valid
	restated.ChunkRoot = strings.Repeat("0", 64)
	if err := VerifyChunkProof(manifest, restated); err == nil || !strings.Contains(err.Error(), "does not match the signed commitment") {
		t.Fatalf("VerifyChunkProof(restated root) error = %v, want signed commitment rejection", err)
	}

	padded := valid
	padded.Path = append(append([]string(nil), valid.Path...), strings.Repeat("0", 64))
	if err := VerifyChunkProof(manifest, padded); err == nil || !strings.Contains(err.Error(), "too long") {
		t.Fatalf("VerifyChunkProof(padded path) error = %v, want path length rejection", err)
	}

	short := valid
	short.Chunk = hex.EncodeToString([]byte("truncated"))
	if err := VerifyChunkProof(manifest, short); err == nil || !strings.Contains(err.Error(), "want") {
		t.Fatalf("VerifyChunkProof(short chunk) error = %v, want length rejection", err)
	}
}

func TestBuildChunkProofRejectsCorruptedObject(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	assetPath := writeLargeFixture(t, root, "model.gguf", 2*int(MinChunkSize))
	output := filepath.Join(root, "demo.ark")
	manifest, err := PackWithOptions("demo", output, []AssetSource{{Role: "model", Path: assetPath}}, time.Now(), PackOptions{ChunkSize: MinChunkSize})
	if err != nil {
		t.Fatalf("PackWithOptions() error = %v", err)
	}
	asset := manifest.Assets[0]
	objectPath := filepath.Join(output, "objects", asset.SHA256)
	corrupted := make([]byte, asset.Size)
	if err := os.WriteFile(objectPath, corrupted, 0o644); err != nil {
		t.Fatalf("corrupt object: %v", err)
	}
	if _, err := BuildChunkProof(manifest.CapsuleID, asset, objectPath, 0); err == nil || !strings.Contains(err.Error(), "does not match its signed commitment") {
		t.Fatalf("BuildChunkProof() error = %v, want commitment rejection", err)
	}
}

func writeLargeFixture(t *testing.T, root, name string, size int) string {
	t.Helper()
	path := filepath.Join(root, name)
	data := make([]byte, size)
	for index := range data {
		data[index] = byte(index % 251)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write fixture %q: %v", name, err)
	}
	return path
}

func TestVerifyPreservesCapsulesWithoutChunkCommitments(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	assetPath := writeLargeFixture(t, root, "model.gguf", 2*int(MinChunkSize))
	output := filepath.Join(root, "legacy.ark")
	manifest, err := PackWithOptions("legacy", output, []AssetSource{{Role: "model", Path: assetPath}}, time.Now(), PackOptions{ChunkSize: MinChunkSize})
	if err != nil {
		t.Fatalf("PackWithOptions() error = %v", err)
	}
	for index := range manifest.Assets {
		manifest.Assets[index].ChunkSize = 0
		manifest.Assets[index].ChunkRoot = ""
	}
	manifest.CapsuleID = capsuleID(manifest)
	if err := writeManifest(output, manifest); err != nil {
		t.Fatalf("writeManifest() error = %v", err)
	}
	verified, err := Verify(output)
	if err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
	if verified.Assets[0].ChunkRoot != "" || verified.CapsuleID != manifest.CapsuleID {
		t.Fatalf("verified manifest = %+v, want legacy shape preserved", verified)
	}
}
