package capsule

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestChunkTreeLocalizesAndRepairsExactChunks(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	assetPath := writeLargeFixture(t, root, "model.gguf", 5*int(MinChunkSize)+7)
	donorPath := filepath.Join(root, "donor.bin")
	original, err := os.ReadFile(assetPath)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	if err := os.WriteFile(donorPath, original, 0o644); err != nil {
		t.Fatalf("write donor: %v", err)
	}
	output := filepath.Join(root, "demo.ark")
	manifest, err := PackWithOptions("demo", output, []AssetSource{{Role: "model", Path: assetPath}}, time.Now(), PackOptions{ChunkSize: MinChunkSize})
	if err != nil {
		t.Fatalf("PackWithOptions() error = %v", err)
	}
	asset := manifest.Assets[0]
	objectPath := filepath.Join(output, "objects", asset.SHA256)

	tree, err := BuildChunkTreeFile(manifest.CapsuleID, asset, objectPath)
	if err != nil {
		t.Fatalf("BuildChunkTreeFile() error = %v", err)
	}
	if _, err := VerifyChunkTreeFile(manifest, tree); err != nil {
		t.Fatalf("VerifyChunkTreeFile() error = %v", err)
	}

	corruptChunk(t, objectPath, 1*MinChunkSize)
	corruptChunk(t, objectPath, 4*MinChunkSize)
	damage, err := ScanChunkDamage(objectPath, asset, tree)
	if err != nil {
		t.Fatalf("ScanChunkDamage() error = %v", err)
	}
	if len(damage) != 2 || damage[0].Index != 1 || damage[1].Index != 4 {
		t.Fatalf("damage = %+v, want chunks 1 and 4", damage)
	}

	repaired, err := RepairChunks(objectPath, asset, tree, donorPath, damage)
	if err != nil {
		t.Fatalf("RepairChunks() error = %v", err)
	}
	if repaired != 2 {
		t.Fatalf("repaired = %d, want 2", repaired)
	}
	if _, err := Verify(output); err != nil {
		t.Fatalf("Verify() after repair error = %v", err)
	}
}

func TestRepairRejectsHostileDonor(t *testing.T) {
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
	tree, err := BuildChunkTreeFile(manifest.CapsuleID, asset, objectPath)
	if err != nil {
		t.Fatalf("BuildChunkTreeFile() error = %v", err)
	}
	corruptChunk(t, objectPath, 0)
	damage, err := ScanChunkDamage(objectPath, asset, tree)
	if err != nil {
		t.Fatalf("ScanChunkDamage() error = %v", err)
	}
	hostile := filepath.Join(root, "hostile.bin")
	if err := os.WriteFile(hostile, make([]byte, asset.Size), 0o644); err != nil {
		t.Fatalf("write hostile donor: %v", err)
	}
	if _, err := RepairChunks(objectPath, asset, tree, hostile, damage); err == nil || !strings.Contains(err.Error(), "does not match the committed tree") {
		t.Fatalf("RepairChunks() error = %v, want donor rejection", err)
	}
}

func TestVerifyChunkTreeFileRejectsForgedLeaves(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	assetPath := writeLargeFixture(t, root, "model.gguf", 3*int(MinChunkSize))
	output := filepath.Join(root, "demo.ark")
	manifest, err := PackWithOptions("demo", output, []AssetSource{{Role: "model", Path: assetPath}}, time.Now(), PackOptions{ChunkSize: MinChunkSize})
	if err != nil {
		t.Fatalf("PackWithOptions() error = %v", err)
	}
	asset := manifest.Assets[0]
	tree, err := BuildChunkTreeFile(manifest.CapsuleID, asset, filepath.Join(output, "objects", asset.SHA256))
	if err != nil {
		t.Fatalf("BuildChunkTreeFile() error = %v", err)
	}
	forged := tree
	forged.Leaves = append([]string(nil), tree.Leaves...)
	forged.Leaves[1] = strings.Repeat("0", 64)
	if _, err := VerifyChunkTreeFile(manifest, forged); err == nil || !strings.Contains(err.Error(), "do not reconstruct") {
		t.Fatalf("VerifyChunkTreeFile() error = %v, want forged leaf rejection", err)
	}

	short := tree
	short.Leaves = tree.Leaves[:2]
	if _, err := VerifyChunkTreeFile(manifest, short); err == nil || !strings.Contains(err.Error(), "want 3") {
		t.Fatalf("VerifyChunkTreeFile() error = %v, want leaf count rejection", err)
	}
}

func TestParallelVerifyReportsFirstAssetInManifestOrder(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	knowledge := writeLargeFixture(t, root, "manual.txt", int(MinChunkSize))
	model := writeLargeFixture(t, root, "model.gguf", int(MinChunkSize)+5)
	output := filepath.Join(root, "demo.ark")
	manifest, err := PackWithOptions("demo", output, []AssetSource{
		{Role: "model", Path: model},
		{Role: "knowledge", Path: knowledge},
	}, time.Now(), PackOptions{ChunkSize: MinChunkSize})
	if err != nil {
		t.Fatalf("PackWithOptions() error = %v", err)
	}
	if manifest.Assets[0].Role != "knowledge" {
		t.Fatalf("assets = %+v, want knowledge first", manifest.Assets)
	}
	for _, asset := range manifest.Assets {
		corruptChunk(t, filepath.Join(output, "objects", asset.SHA256), 0)
	}
	for attempt := 0; attempt < 5; attempt++ {
		_, err := Verify(output)
		if err == nil || !strings.Contains(err.Error(), "manual.txt") {
			t.Fatalf("Verify() error = %v, want first manifest asset reported", err)
		}
	}
}

func corruptChunk(t *testing.T, path string, offset int64) {
	t.Helper()
	file, err := os.OpenFile(path, os.O_RDWR, 0o644)
	if err != nil {
		t.Fatalf("open object: %v", err)
	}
	defer file.Close()
	if _, err := file.WriteAt([]byte{0xff, 0xfe, 0xfd}, offset); err != nil {
		t.Fatalf("corrupt object: %v", err)
	}
}
