package identity

import (
	"crypto/rand"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteAndLoadRevocationSet(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	first, err := Create(filepath.Join(root, "first"), rand.Reader)
	if err != nil {
		t.Fatalf("create first identity: %v", err)
	}
	second, err := Create(filepath.Join(root, "second"), rand.Reader)
	if err != nil {
		t.Fatalf("create second identity: %v", err)
	}
	path := filepath.Join(root, "revocations.json")
	written, err := WriteRevocationSet(path, []PublicIdentity{second.Public, first.Public, first.Public})
	if err != nil {
		t.Fatalf("WriteRevocationSet() error = %v", err)
	}
	if len(written.RevokedKeyIDs) != 2 {
		t.Fatalf("len(RevokedKeyIDs) = %d, want 2", len(written.RevokedKeyIDs))
	}
	loaded, err := LoadRevocationSet(path)
	if err != nil {
		t.Fatalf("LoadRevocationSet() error = %v", err)
	}
	if !loaded.Contains(first.Public.KeyID) || !loaded.Contains(second.Public.KeyID) {
		t.Fatalf("loaded set does not contain both identities: %+v", loaded)
	}
}

func TestLoadRevocationSetRejectsUnsortedKeys(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	first, err := Create(filepath.Join(root, "first"), rand.Reader)
	if err != nil {
		t.Fatalf("create first identity: %v", err)
	}
	second, err := Create(filepath.Join(root, "second"), rand.Reader)
	if err != nil {
		t.Fatalf("create second identity: %v", err)
	}
	ids := []string{first.Public.KeyID, second.Public.KeyID}
	if ids[0] < ids[1] {
		ids[0], ids[1] = ids[1], ids[0]
	}
	data, err := json.Marshal(RevocationSet{SchemaVersion: RevocationSchemaVersion, RevokedKeyIDs: ids})
	if err != nil {
		t.Fatalf("marshal revocations: %v", err)
	}
	path := filepath.Join(root, "revocations.json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write revocations: %v", err)
	}
	if _, err := LoadRevocationSet(path); err == nil || !strings.Contains(err.Error(), "unique and sorted") {
		t.Fatalf("LoadRevocationSet() error = %v, want sorted error", err)
	}
}

func TestLoadRevocationSetRejectsInvalidKeyID(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "revocations.json")
	data, err := json.Marshal(RevocationSet{
		SchemaVersion: RevocationSchemaVersion,
		RevokedKeyIDs: []string{"ed25519:not-a-key"},
	})
	if err != nil {
		t.Fatalf("marshal revocations: %v", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write revocations: %v", err)
	}
	if _, err := LoadRevocationSet(path); err == nil || !strings.Contains(err.Error(), "invalid revoked key ID") {
		t.Fatalf("LoadRevocationSet() error = %v, want invalid key ID", err)
	}
}

func TestWriteRevocationSetRefusesExistingFile(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	created, err := Create(filepath.Join(root, "author"), rand.Reader)
	if err != nil {
		t.Fatalf("create identity: %v", err)
	}
	path := filepath.Join(root, "revocations.json")
	if _, err := WriteRevocationSet(path, []PublicIdentity{created.Public}); err != nil {
		t.Fatalf("first WriteRevocationSet() error = %v", err)
	}
	if _, err := WriteRevocationSet(path, []PublicIdentity{created.Public}); err == nil || !strings.Contains(err.Error(), "create revocations.json") {
		t.Fatalf("second WriteRevocationSet() error = %v, want existing file", err)
	}
}
