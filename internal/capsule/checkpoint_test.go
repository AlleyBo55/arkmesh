package capsule

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckpointRoundTripAndExactMatch(t *testing.T) {
	t.Parallel()
	fixture := rotationFixtureWithoutAuthority(t)
	checkpoint, err := NewCheckpoint(fixture.parent, fixture.parentAuth)
	if err != nil {
		t.Fatalf("NewCheckpoint() error = %v", err)
	}
	path := filepath.Join(t.TempDir(), "checkpoint.json")
	if err := WriteCheckpoint(path, checkpoint); err != nil {
		t.Fatalf("WriteCheckpoint() error = %v", err)
	}
	loaded, err := LoadCheckpoint(path)
	if err != nil {
		t.Fatalf("LoadCheckpoint() error = %v", err)
	}
	if loaded != checkpoint {
		t.Fatalf("LoadCheckpoint() = %+v, want %+v", loaded, checkpoint)
	}
	if err := VerifyCheckpoint(loaded, fixture.parent, fixture.parentAuth); err != nil {
		t.Fatalf("VerifyCheckpoint() error = %v", err)
	}
	if err := VerifyCheckpoint(loaded, fixture.child, fixture.childAuth); err == nil || !strings.Contains(err.Error(), "does not match checkpoint head") {
		t.Fatalf("VerifyCheckpoint() error = %v, want capsule mismatch", err)
	}
}

func TestWriteCheckpointRefusesExistingFile(t *testing.T) {
	t.Parallel()
	fixture := rotationFixtureWithoutAuthority(t)
	checkpoint, err := NewCheckpoint(fixture.parent, fixture.parentAuth)
	if err != nil {
		t.Fatalf("NewCheckpoint() error = %v", err)
	}
	path := filepath.Join(t.TempDir(), "checkpoint.json")
	if err := WriteCheckpoint(path, checkpoint); err != nil {
		t.Fatalf("first WriteCheckpoint() error = %v", err)
	}
	if err := WriteCheckpoint(path, checkpoint); err == nil || !strings.Contains(err.Error(), "file exists") {
		t.Fatalf("second WriteCheckpoint() error = %v, want existing-file refusal", err)
	}
}

func TestVerifyCheckpointRejectsDifferentSigner(t *testing.T) {
	t.Parallel()
	fixture := rotationFixtureWithoutAuthority(t)
	checkpoint, err := NewCheckpoint(fixture.parent, fixture.parentAuth)
	if err != nil {
		t.Fatalf("NewCheckpoint() error = %v", err)
	}
	checkpoint.SignerID = fixture.childAuth.SignerID
	if err := VerifyCheckpoint(checkpoint, fixture.parent, fixture.parentAuth); err == nil || !strings.Contains(err.Error(), "does not match checkpoint signer") {
		t.Fatalf("VerifyCheckpoint() error = %v, want signer mismatch", err)
	}
}

func TestAdvanceCheckpointThroughKeyRotationRejectsRollback(t *testing.T) {
	t.Parallel()
	fixture := rotationFixture(t)
	path := filepath.Join(t.TempDir(), "checkpoint.json")
	checkpoint, err := NewCheckpoint(fixture.parent, fixture.parentAuth)
	if err != nil {
		t.Fatalf("NewCheckpoint() error = %v", err)
	}
	if err := WriteCheckpoint(path, checkpoint); err != nil {
		t.Fatalf("WriteCheckpoint() error = %v", err)
	}
	next, lineage, err := AdvanceCheckpointFile(path, fixture.childRoot, fixture.parent, fixture.parentAuth, fixture.child, fixture.childAuth)
	if err != nil {
		t.Fatalf("AdvanceCheckpointFile() error = %v", err)
	}
	if next.CapsuleID != fixture.child.CapsuleID || next.SignerID != fixture.childAuth.SignerID {
		t.Fatalf("next checkpoint = %+v, want child head", next)
	}
	if lineage.Status != LineageKeyRotation {
		t.Fatalf("lineage status = %s, want %s", lineage.Status, LineageKeyRotation)
	}
	loaded, err := LoadCheckpoint(path)
	if err != nil {
		t.Fatalf("LoadCheckpoint() error = %v", err)
	}
	if err := VerifyCheckpoint(loaded, fixture.parent, fixture.parentAuth); err == nil || !strings.Contains(err.Error(), "does not match checkpoint head") {
		t.Fatalf("VerifyCheckpoint(parent) error = %v, want rollback rejection", err)
	}
	if err := VerifyCheckpoint(loaded, fixture.child, fixture.childAuth); err != nil {
		t.Fatalf("VerifyCheckpoint(child) error = %v", err)
	}
}

func TestAdvanceCheckpointRejectsParentThatIsNotCurrentHead(t *testing.T) {
	t.Parallel()
	fixture := rotationFixture(t)
	path := filepath.Join(t.TempDir(), "checkpoint.json")
	checkpoint, err := NewCheckpoint(fixture.child, fixture.childAuth)
	if err != nil {
		t.Fatalf("NewCheckpoint() error = %v", err)
	}
	if err := WriteCheckpoint(path, checkpoint); err != nil {
		t.Fatalf("WriteCheckpoint() error = %v", err)
	}
	if _, _, err := AdvanceCheckpointFile(path, fixture.childRoot, fixture.parent, fixture.parentAuth, fixture.child, fixture.childAuth); err == nil || !strings.Contains(err.Error(), "does not match checkpoint head") {
		t.Fatalf("AdvanceCheckpointFile() error = %v, want stale parent rejection", err)
	}
}

func TestLoadCheckpointRejectsUnknownField(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "checkpoint.json")
	data := `{"schema_version":"arkmesh.checkpoint/v0alpha1","capsule_id":"sha256:` + strings.Repeat("0", 64) + `","signer_id":"ed25519:` + strings.Repeat("0", 64) + `","unknown":true}`
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatalf("write checkpoint: %v", err)
	}
	if _, err := LoadCheckpoint(path); err == nil || !strings.Contains(err.Error(), "unknown field") {
		t.Fatalf("LoadCheckpoint() error = %v, want unknown field rejection", err)
	}
}

func TestLoadCheckpointRejectsSymlink(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	target := filepath.Join(root, "target.json")
	if err := os.WriteFile(target, []byte("{}"), 0o644); err != nil {
		t.Fatalf("write target: %v", err)
	}
	path := filepath.Join(root, "checkpoint.json")
	if err := os.Symlink(target, path); err != nil {
		t.Fatalf("create checkpoint symlink: %v", err)
	}
	if _, err := LoadCheckpoint(path); err == nil || !strings.Contains(err.Error(), "not a regular file") {
		t.Fatalf("LoadCheckpoint() error = %v, want symlink rejection", err)
	}
}
