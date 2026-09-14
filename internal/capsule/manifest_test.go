package capsule

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestPackAndVerify(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	model := writeFixture(t, root, "tiny.gguf", "model bytes")
	knowledge := writeFixture(t, root, "manual.txt", "clean water instructions")
	output := filepath.Join(root, "demo.ark")
	created := time.Date(2026, 9, 14, 2, 0, 0, 0, time.UTC)

	manifest, err := Pack("field-assistant", output, []AssetSource{
		{Role: "model", Path: model},
		{Role: "knowledge", Path: knowledge},
	}, created)
	if err != nil {
		t.Fatalf("Pack() error = %v", err)
	}
	if !strings.HasPrefix(manifest.CapsuleID, "sha256:") {
		t.Fatalf("CapsuleID = %q, want sha256 prefix", manifest.CapsuleID)
	}
	if len(manifest.Assets) != 2 {
		t.Fatalf("len(Assets) = %d, want 2", len(manifest.Assets))
	}

	verified, err := Verify(output)
	if err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
	if verified.CapsuleID != manifest.CapsuleID {
		t.Fatalf("verified CapsuleID = %q, want %q", verified.CapsuleID, manifest.CapsuleID)
	}
}

func TestCapsuleIDIsDeterministicForSameManifest(t *testing.T) {
	t.Parallel()
	manifest := Manifest{
		SchemaVersion: SchemaVersion,
		Name:          "test",
		CreatedAt:     "2026-09-14T02:00:00Z",
		Assets: []Asset{{
			Role:   "model",
			Name:   "model.gguf",
			Size:   3,
			SHA256: strings.Repeat("a", 64),
		}},
	}
	first := capsuleID(manifest)
	manifest.CapsuleID = "ignored-while-calculating"
	second := capsuleID(manifest)
	if first != second {
		t.Fatalf("capsuleID() = %q then %q", first, second)
	}
}

func TestVerifyDetectsTamperedObject(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	model := writeFixture(t, root, "tiny.gguf", "original")
	output := filepath.Join(root, "demo.ark")
	manifest, err := Pack("demo", output, []AssetSource{{Role: "model", Path: model}}, time.Now())
	if err != nil {
		t.Fatalf("Pack() error = %v", err)
	}

	objectPath := filepath.Join(output, "objects", manifest.Assets[0].SHA256)
	if err := os.WriteFile(objectPath, []byte("tampered"), 0o644); err != nil {
		t.Fatalf("tamper object: %v", err)
	}
	if _, err := Verify(output); err == nil || !strings.Contains(err.Error(), "digest mismatch") {
		t.Fatalf("Verify() error = %v, want digest mismatch", err)
	}
}

func TestVerifyDetectsMissingObject(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	model := writeFixture(t, root, "tiny.gguf", "model")
	output := filepath.Join(root, "demo.ark")
	manifest, err := Pack("demo", output, []AssetSource{{Role: "model", Path: model}}, time.Now())
	if err != nil {
		t.Fatalf("Pack() error = %v", err)
	}

	if err := os.Remove(filepath.Join(output, "objects", manifest.Assets[0].SHA256)); err != nil {
		t.Fatalf("remove object: %v", err)
	}
	if _, err := Verify(output); err == nil {
		t.Fatal("Verify() error = nil, want missing object error")
	}
}

func TestVerifyRejectsManifestMutation(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	model := writeFixture(t, root, "tiny.gguf", "model")
	output := filepath.Join(root, "demo.ark")
	manifest, err := Pack("demo", output, []AssetSource{{Role: "model", Path: model}}, time.Now())
	if err != nil {
		t.Fatalf("Pack() error = %v", err)
	}

	manifest.Name = "renamed-without-new-id"
	data, err := json.Marshal(manifest)
	if err != nil {
		t.Fatalf("marshal manifest: %v", err)
	}
	if err := os.WriteFile(filepath.Join(output, "manifest.json"), data, 0o644); err != nil {
		t.Fatalf("rewrite manifest: %v", err)
	}
	if _, err := Verify(output); err == nil || !strings.Contains(err.Error(), "capsule ID") {
		t.Fatalf("Verify() error = %v, want capsule ID mismatch", err)
	}
}

func TestPackRejectsInvalidRole(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	model := writeFixture(t, root, "tiny.gguf", "model")
	_, err := Pack("demo", filepath.Join(root, "demo.ark"), []AssetSource{
		{Role: "../../model", Path: model},
	}, time.Now())
	if err == nil || !strings.Contains(err.Error(), "invalid asset role") {
		t.Fatalf("Pack() error = %v, want invalid role", err)
	}
}

func writeFixture(t *testing.T, root, name, content string) string {
	t.Helper()
	path := filepath.Join(root, name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	return path
}

func TestVerifyRejectsUnknownManifestField(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	model := writeFixture(t, root, "tiny.gguf", "model")
	output := filepath.Join(root, "demo.ark")
	if _, err := Pack("demo", output, []AssetSource{{Role: "model", Path: model}}, time.Now()); err != nil {
		t.Fatalf("Pack() error = %v", err)
	}

	manifestPath := filepath.Join(output, "manifest.json")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	changed := strings.Replace(string(data), "{\n", "{\n  \"unexpected\": true,\n", 1)
	if err := os.WriteFile(manifestPath, []byte(changed), 0o644); err != nil {
		t.Fatalf("rewrite manifest: %v", err)
	}
	if _, err := Verify(output); err == nil || !strings.Contains(err.Error(), "unknown field") {
		t.Fatalf("Verify() error = %v, want unknown field", err)
	}
}

func TestVerifyRejectsSymlinkObject(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	model := writeFixture(t, root, "tiny.gguf", "model")
	output := filepath.Join(root, "demo.ark")
	manifest, err := Pack("demo", output, []AssetSource{{Role: "model", Path: model}}, time.Now())
	if err != nil {
		t.Fatalf("Pack() error = %v", err)
	}

	objectPath := filepath.Join(output, "objects", manifest.Assets[0].SHA256)
	if err := os.Remove(objectPath); err != nil {
		t.Fatalf("remove packed object: %v", err)
	}
	if err := os.Symlink(model, objectPath); err != nil {
		t.Fatalf("create object symlink: %v", err)
	}
	if _, err := Verify(output); err == nil || !strings.Contains(err.Error(), "not a regular file") {
		t.Fatalf("Verify() error = %v, want non-regular object", err)
	}
}
