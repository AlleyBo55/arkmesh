package identity

import (
	"crypto/rand"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateAndLoadIdentity(t *testing.T) {
	t.Parallel()
	directory := filepath.Join(t.TempDir(), "author")
	created, err := Create(directory, rand.Reader)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	public, err := LoadPublic(created.PublicPath)
	if err != nil {
		t.Fatalf("LoadPublic() error = %v", err)
	}
	private, err := LoadPrivate(created.PrivatePath)
	if err != nil {
		t.Fatalf("LoadPrivate() error = %v", err)
	}
	derivedPublic, err := private.Public()
	if err != nil {
		t.Fatalf("PrivateIdentity.Public() error = %v", err)
	}
	if public.KeyID != created.Public.KeyID || derivedPublic.KeyID != public.KeyID {
		t.Fatalf("identity IDs differ: created=%q public=%q private=%q", created.Public.KeyID, public.KeyID, derivedPublic.KeyID)
	}

	info, err := os.Stat(created.PrivatePath)
	if err != nil {
		t.Fatalf("stat private identity: %v", err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("private identity mode = %o, want 600", got)
	}
}

func TestCreateRefusesExistingDirectory(t *testing.T) {
	t.Parallel()
	directory := filepath.Join(t.TempDir(), "author")
	if _, err := Create(directory, rand.Reader); err != nil {
		t.Fatalf("first Create() error = %v", err)
	}
	if _, err := Create(directory, rand.Reader); err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("second Create() error = %v, want already exists", err)
	}
}

func TestLoadPublicRejectsChangedKeyID(t *testing.T) {
	t.Parallel()
	directory := filepath.Join(t.TempDir(), "author")
	created, err := Create(directory, rand.Reader)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	public := created.Public
	public.KeyID = "ed25519:" + strings.Repeat("0", 64)
	data, err := json.Marshal(public)
	if err != nil {
		t.Fatalf("marshal public identity: %v", err)
	}
	if err := os.WriteFile(created.PublicPath, data, 0o644); err != nil {
		t.Fatalf("rewrite public identity: %v", err)
	}
	if _, err := LoadPublic(created.PublicPath); err == nil || !strings.Contains(err.Error(), "does not match") {
		t.Fatalf("LoadPublic() error = %v, want key ID mismatch", err)
	}
}

func TestLoadPrivateRejectsMismatchedPublicKey(t *testing.T) {
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
	private, err := LoadPrivate(first.PrivatePath)
	if err != nil {
		t.Fatalf("load first private identity: %v", err)
	}
	private.PublicKey = second.Public.PublicKey
	private.KeyID = second.Public.KeyID
	data, err := json.Marshal(private)
	if err != nil {
		t.Fatalf("marshal private identity: %v", err)
	}
	if err := os.WriteFile(first.PrivatePath, data, 0o600); err != nil {
		t.Fatalf("rewrite private identity: %v", err)
	}
	if _, err := LoadPrivate(first.PrivatePath); err == nil || !strings.Contains(err.Error(), "does not match private key") {
		t.Fatalf("LoadPrivate() error = %v, want public/private mismatch", err)
	}
}

func TestLoadPublicRejectsDuplicateField(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "identity.json")
	data := `{"schema_version":"arkmesh.identity/v0alpha1","algorithm":"ed25519","algorithm":"other","key_id":"ed25519:` + strings.Repeat("0", 64) + `","public_key":""}`
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatalf("write identity: %v", err)
	}
	if _, err := LoadPublic(path); err == nil || !strings.Contains(err.Error(), "duplicate JSON field") {
		t.Fatalf("LoadPublic() error = %v, want duplicate field rejection", err)
	}
}
