package capsule

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"arkmesh/internal/identity"
)

func TestSignAndVerifyTrustedCapsule(t *testing.T) {
	t.Parallel()
	root, private, public := signedFixture(t)
	manifest, authenticity, err := VerifyAuthenticated(root, VerifyOptions{
		Trusted:        []identity.PublicIdentity{public},
		RequireTrusted: true,
	})
	if err != nil {
		t.Fatalf("VerifyAuthenticated() error = %v", err)
	}
	if manifest.Name != "signed-demo" {
		t.Fatalf("manifest.Name = %q, want signed-demo", manifest.Name)
	}
	if authenticity.Status != AuthenticityTrusted || authenticity.SignerID != public.KeyID {
		t.Fatalf("authenticity = %+v, want trusted signer %q", authenticity, public.KeyID)
	}
	if _, _, err := private.Keys(); err != nil {
		t.Fatalf("fixture private identity invalid: %v", err)
	}
}

func TestVerifyReportsValidUnknownAuthor(t *testing.T) {
	t.Parallel()
	root, _, public := signedFixture(t)
	_, authenticity, err := VerifyAuthenticated(root, VerifyOptions{})
	if err != nil {
		t.Fatalf("VerifyAuthenticated() error = %v", err)
	}
	if authenticity.Status != AuthenticityUnknown || authenticity.SignerID != public.KeyID {
		t.Fatalf("authenticity = %+v, want unknown author %q", authenticity, public.KeyID)
	}
}

func TestVerifyPreservesUnsignedCompatibility(t *testing.T) {
	t.Parallel()
	root := unsignedFixture(t)
	_, authenticity, err := VerifyAuthenticated(root, VerifyOptions{})
	if err != nil {
		t.Fatalf("VerifyAuthenticated() error = %v", err)
	}
	if authenticity.Status != AuthenticityUnsigned {
		t.Fatalf("authenticity.Status = %q, want unsigned", authenticity.Status)
	}
	if _, _, err := VerifyAuthenticated(root, VerifyOptions{RequireSignature: true}); err == nil || !strings.Contains(err.Error(), "unsigned") {
		t.Fatalf("required signature error = %v, want unsigned", err)
	}
}

func TestVerifyRejectsInvalidSignature(t *testing.T) {
	t.Parallel()
	root, _, _ := signedFixture(t)
	path := filepath.Join(root, SignatureFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read signature: %v", err)
	}
	var envelope SignatureEnvelope
	if err := json.Unmarshal(data, &envelope); err != nil {
		t.Fatalf("decode signature: %v", err)
	}
	envelope.Signature = base64.StdEncoding.EncodeToString(make([]byte, 64))
	data, err = json.Marshal(envelope)
	if err != nil {
		t.Fatalf("encode changed signature: %v", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("rewrite signature: %v", err)
	}
	if _, _, err := VerifyAuthenticated(root, VerifyOptions{}); err == nil || !strings.Contains(err.Error(), "invalid capsule signature") {
		t.Fatalf("VerifyAuthenticated() error = %v, want invalid signature", err)
	}
}

func TestVerifyRejectsUntrustedSignerWhenRequired(t *testing.T) {
	t.Parallel()
	root, _, _ := signedFixture(t)
	otherDirectory := filepath.Join(t.TempDir(), "other")
	other, err := identity.Create(otherDirectory, rand.Reader)
	if err != nil {
		t.Fatalf("create other identity: %v", err)
	}
	if _, _, err := VerifyAuthenticated(root, VerifyOptions{
		Trusted:        []identity.PublicIdentity{other.Public},
		RequireTrusted: true,
	}); err == nil || !strings.Contains(err.Error(), "is not trusted") {
		t.Fatalf("VerifyAuthenticated() error = %v, want untrusted signer", err)
	}
}

func TestSignRefusesSecondSignature(t *testing.T) {
	t.Parallel()
	root, private, _ := signedFixture(t)
	if _, err := Sign(root, private); err == nil || !strings.Contains(err.Error(), "create capsule signature") {
		t.Fatalf("second Sign() error = %v, want existing signature error", err)
	}
}

func signedFixture(t *testing.T) (string, identity.PrivateIdentity, identity.PublicIdentity) {
	t.Helper()
	root := unsignedFixture(t)
	identityDirectory := filepath.Join(t.TempDir(), "author")
	created, err := identity.Create(identityDirectory, rand.Reader)
	if err != nil {
		t.Fatalf("identity.Create() error = %v", err)
	}
	private, err := identity.LoadPrivate(created.PrivatePath)
	if err != nil {
		t.Fatalf("identity.LoadPrivate() error = %v", err)
	}
	if _, err := Sign(root, private); err != nil {
		t.Fatalf("Sign() error = %v", err)
	}
	return root, private, created.Public
}

func unsignedFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	model := writeFixture(t, root, "tiny.gguf", "model")
	output := filepath.Join(root, "demo.ark")
	if _, err := Pack("signed-demo", output, []AssetSource{{Role: "model", Path: model}}, time.Now()); err != nil {
		t.Fatalf("Pack() error = %v", err)
	}
	return output
}

func TestVerifyRejectsLocallyRevokedSigner(t *testing.T) {
	t.Parallel()
	root, _, public := signedFixture(t)
	revocations, err := identity.NewRevocationSet([]identity.PublicIdentity{public})
	if err != nil {
		t.Fatalf("NewRevocationSet() error = %v", err)
	}
	if _, _, err := VerifyAuthenticated(root, VerifyOptions{
		Trusted:     []identity.PublicIdentity{public},
		Revocations: revocations,
	}); err == nil || !strings.Contains(err.Error(), "revoked by local policy") {
		t.Fatalf("VerifyAuthenticated() error = %v, want revoked signer", err)
	}
}
