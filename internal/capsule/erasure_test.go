package capsule

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func erasureFixture(t *testing.T, size int) (string, Manifest, Asset, string) {
	t.Helper()
	root := t.TempDir()
	assetPath := writeLargeFixture(t, root, "model.gguf", size)
	output := filepath.Join(root, "demo.ark")
	manifest, err := PackWithOptions("demo", output, []AssetSource{{Role: "model", Path: assetPath}}, time.Now(), PackOptions{ChunkSize: MinChunkSize})
	if err != nil {
		t.Fatalf("PackWithOptions() error = %v", err)
	}
	return root, manifest, manifest.Assets[0], output
}

func TestErasureRecoversObjectDeletedEntirely(t *testing.T) {
	t.Parallel()
	root, manifest, asset, capsuleRoot := erasureFixture(t, 3*int(ErasureBlockSize)+321)
	shardDirectory := filepath.Join(root, "shards")
	objectPath := filepath.Join(capsuleRoot, "objects", asset.SHA256)

	plan, err := ProtectAsset(manifest.CapsuleID, asset, objectPath, shardDirectory, 4, 2)
	if err != nil {
		t.Fatalf("ProtectAsset() error = %v", err)
	}
	if plan.ShardSize%plan.BlockSize != 0 || plan.ShardSize*4 < asset.Size {
		t.Fatalf("plan = %+v, want shards covering the asset", plan)
	}

	// The donor is gone: the object itself and two shards are removed.
	if err := os.Remove(objectPath); err != nil {
		t.Fatalf("remove object: %v", err)
	}
	if err := os.Remove(filepath.Join(shardDirectory, ShardFileName(0))); err != nil {
		t.Fatalf("remove shard 0: %v", err)
	}
	if err := os.Remove(filepath.Join(shardDirectory, ShardFileName(3))); err != nil {
		t.Fatalf("remove shard 3: %v", err)
	}

	report, err := RecoverAsset(manifest, plan, shardDirectory, capsuleRoot)
	if err != nil {
		t.Fatalf("RecoverAsset() error = %v", err)
	}
	if len(report.MissingShards) != 2 || report.Bytes != asset.Size {
		t.Fatalf("report = %+v, want two missing shards", report)
	}
	if _, err := Verify(capsuleRoot); err != nil {
		t.Fatalf("Verify() after recovery error = %v", err)
	}
}

func TestErasureRefusesBelowThreshold(t *testing.T) {
	t.Parallel()
	root, manifest, asset, capsuleRoot := erasureFixture(t, 2*int(ErasureBlockSize))
	shardDirectory := filepath.Join(root, "shards")
	objectPath := filepath.Join(capsuleRoot, "objects", asset.SHA256)
	plan, err := ProtectAsset(manifest.CapsuleID, asset, objectPath, shardDirectory, 3, 1)
	if err != nil {
		t.Fatalf("ProtectAsset() error = %v", err)
	}
	if err := os.Remove(objectPath); err != nil {
		t.Fatalf("remove object: %v", err)
	}
	for _, index := range []int{0, 2} {
		if err := os.Remove(filepath.Join(shardDirectory, ShardFileName(index))); err != nil {
			t.Fatalf("remove shard %d: %v", index, err)
		}
	}
	if _, err := RecoverAsset(manifest, plan, shardDirectory, capsuleRoot); err == nil || !strings.Contains(err.Error(), "required") {
		t.Fatalf("RecoverAsset() error = %v, want threshold refusal", err)
	}
	if _, err := os.Stat(objectPath); !os.IsNotExist(err) {
		t.Fatalf("object exists after refused recovery: %v", err)
	}
}

// TestErasureRejectsConsistentlyForgedShards proves the plan is not a trust
// anchor: shards whose digests match a rewritten plan still fail, because the
// reconstructed bytes are checked against the signed manifest.
func TestErasureRejectsConsistentlyForgedShards(t *testing.T) {
	t.Parallel()
	root, manifest, asset, capsuleRoot := erasureFixture(t, 2*int(ErasureBlockSize))
	shardDirectory := filepath.Join(root, "shards")
	objectPath := filepath.Join(capsuleRoot, "objects", asset.SHA256)
	plan, err := ProtectAsset(manifest.CapsuleID, asset, objectPath, shardDirectory, 2, 2)
	if err != nil {
		t.Fatalf("ProtectAsset() error = %v", err)
	}
	if err := os.Remove(objectPath); err != nil {
		t.Fatalf("remove object: %v", err)
	}
	forged := make([]byte, plan.ShardSize)
	for index := range forged {
		forged[index] = byte(index % 7)
	}
	shardPath := filepath.Join(shardDirectory, ShardFileName(1))
	if err := os.WriteFile(shardPath, forged, 0o644); err != nil {
		t.Fatalf("forge shard: %v", err)
	}
	digest := sha256.Sum256(forged)
	plan.ShardDigests[1] = hex.EncodeToString(digest[:])

	if _, err := RecoverAsset(manifest, plan, shardDirectory, capsuleRoot); err == nil || !strings.Contains(err.Error(), "recovered bytes rejected") {
		t.Fatalf("RecoverAsset() error = %v, want signed manifest rejection", err)
	}
	if _, err := os.Stat(objectPath); !os.IsNotExist(err) {
		t.Fatalf("object exists after rejected recovery: %v", err)
	}
}

func TestErasurePlanValidationRejectsMismatchedPlans(t *testing.T) {
	t.Parallel()
	root, manifest, asset, capsuleRoot := erasureFixture(t, int(ErasureBlockSize))
	shardDirectory := filepath.Join(root, "shards")
	plan, err := ProtectAsset(manifest.CapsuleID, asset, filepath.Join(capsuleRoot, "objects", asset.SHA256), shardDirectory, 2, 1)
	if err != nil {
		t.Fatalf("ProtectAsset() error = %v", err)
	}
	cases := map[string]func(ErasurePlan) ErasurePlan{
		"capsule ID does not match": func(altered ErasurePlan) ErasurePlan {
			altered.CapsuleID = "sha256:" + strings.Repeat("0", 64)
			return altered
		},
		"size does not match": func(altered ErasurePlan) ErasurePlan {
			altered.AssetSize = asset.Size + 1
			return altered
		},
		"invalid block size": func(altered ErasurePlan) ErasurePlan {
			altered.BlockSize = 3
			return altered
		},
		"digest count does not match": func(altered ErasurePlan) ErasurePlan {
			altered.ShardDigests = altered.ShardDigests[:1]
			return altered
		},
	}
	for name, alter := range cases {
		if _, err := validateErasurePlan(manifest, alter(plan)); err == nil {
			t.Fatalf("validateErasurePlan(%s) error = nil, want rejection", name)
		}
	}
}
