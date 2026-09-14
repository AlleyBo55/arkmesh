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

func TestKeyRotationAuthorizesDifferentChildSigner(t *testing.T) {
	t.Parallel()
	fixture := rotationFixture(t)
	lineage, err := CheckLineageWithAuthority(fixture.childRoot, fixture.child, fixture.childAuth, fixture.parent, fixture.parentAuth)
	if err != nil {
		t.Fatalf("CheckLineageWithAuthority() error = %v", err)
	}
	if lineage.Status != LineageKeyRotation || lineage.SignerID != fixture.childAuth.SignerID {
		t.Fatalf("lineage = %+v, want verified key rotation", lineage)
	}
}

func TestKeyRotationRejectsMissingAuthority(t *testing.T) {
	t.Parallel()
	fixture := rotationFixture(t)
	if err := os.Remove(filepath.Join(fixture.childRoot, AuthorityFileName)); err != nil {
		t.Fatalf("remove authority: %v", err)
	}
	if _, err := CheckLineageWithAuthority(fixture.childRoot, fixture.child, fixture.childAuth, fixture.parent, fixture.parentAuth); err == nil || !strings.Contains(err.Error(), "requires a key rotation authority") {
		t.Fatalf("CheckLineageWithAuthority() error = %v, want missing authority", err)
	}
}

func TestKeyRotationRejectsSymlinkAuthority(t *testing.T) {
	t.Parallel()
	fixture := rotationFixtureWithoutAuthority(t)
	path := filepath.Join(fixture.childRoot, AuthorityFileName)
	if err := os.Symlink("missing-authority.json", path); err != nil {
		t.Fatalf("create authority symlink: %v", err)
	}
	if _, err := CheckLineageWithAuthority(fixture.childRoot, fixture.child, fixture.childAuth, fixture.parent, fixture.parentAuth); err == nil || !strings.Contains(err.Error(), "not a regular file") {
		t.Fatalf("CheckLineageWithAuthority() error = %v, want nonregular authority", err)
	}
}
func TestKeyRotationRejectsModifiedSignature(t *testing.T) {
	t.Parallel()
	fixture := rotationFixture(t)
	path := filepath.Join(fixture.childRoot, AuthorityFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read authority: %v", err)
	}
	var transition AuthorityTransition
	if err := json.Unmarshal(data, &transition); err != nil {
		t.Fatalf("decode authority: %v", err)
	}
	transition.Signature = base64.StdEncoding.EncodeToString(make([]byte, 64))
	data, err = json.Marshal(transition)
	if err != nil {
		t.Fatalf("encode authority: %v", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("rewrite authority: %v", err)
	}
	if _, err := CheckLineageWithAuthority(fixture.childRoot, fixture.child, fixture.childAuth, fixture.parent, fixture.parentAuth); err == nil || !strings.Contains(err.Error(), "invalid key rotation authority signature") {
		t.Fatalf("CheckLineageWithAuthority() error = %v, want invalid signature", err)
	}
}

func TestAuthorizeKeyRotationRejectsWrongPreviousKey(t *testing.T) {
	t.Parallel()
	fixture := rotationFixtureWithoutAuthority(t)
	other, err := identity.Create(filepath.Join(t.TempDir(), "other"), rand.Reader)
	if err != nil {
		t.Fatalf("create other identity: %v", err)
	}
	otherPrivate, err := identity.LoadPrivate(other.PrivatePath)
	if err != nil {
		t.Fatalf("load other identity: %v", err)
	}
	if _, err := AuthorizeKeyRotation(fixture.childRoot, fixture.parent, fixture.parentAuth, fixture.child, fixture.childAuth, otherPrivate); err == nil || !strings.Contains(err.Error(), "does not match parent signer") {
		t.Fatalf("AuthorizeKeyRotation() error = %v, want previous signer mismatch", err)
	}
}

func TestKeyRotationAuthorityCannotBeReusedForAnotherChild(t *testing.T) {
	t.Parallel()
	fixture := rotationFixture(t)
	otherAsset := writeFixture(t, t.TempDir(), "other.gguf", "other child")
	otherRoot := filepath.Join(t.TempDir(), "other.ark")
	otherChild, err := PackWithOptions("other", otherRoot, []AssetSource{{Role: "model", Path: otherAsset}}, time.Now(), PackOptions{
		ParentCapsuleID: fixture.parent.CapsuleID,
	})
	if err != nil {
		t.Fatalf("pack other child: %v", err)
	}
	if _, err := Sign(otherRoot, fixture.childPrivate); err != nil {
		t.Fatalf("sign other child: %v", err)
	}
	otherChild, otherAuth, err := VerifyAuthenticated(otherRoot, VerifyOptions{})
	if err != nil {
		t.Fatalf("verify other child: %v", err)
	}
	authorityData, err := os.ReadFile(filepath.Join(fixture.childRoot, AuthorityFileName))
	if err != nil {
		t.Fatalf("read original authority: %v", err)
	}
	if err := os.WriteFile(filepath.Join(otherRoot, AuthorityFileName), authorityData, 0o644); err != nil {
		t.Fatalf("copy authority: %v", err)
	}
	if _, err := CheckLineageWithAuthority(otherRoot, otherChild, otherAuth, fixture.parent, fixture.parentAuth); err == nil || !strings.Contains(err.Error(), "child capsule ID does not match") {
		t.Fatalf("CheckLineageWithAuthority() error = %v, want child ID mismatch", err)
	}
}

type rotationTestFixture struct {
	parent        Manifest
	parentAuth    Authenticity
	parentPrivate identity.PrivateIdentity
	childRoot     string
	child         Manifest
	childAuth     Authenticity
	childPrivate  identity.PrivateIdentity
}

func rotationFixture(t *testing.T) rotationTestFixture {
	t.Helper()
	fixture := rotationFixtureWithoutAuthority(t)
	if _, err := AuthorizeKeyRotation(fixture.childRoot, fixture.parent, fixture.parentAuth, fixture.child, fixture.childAuth, fixture.parentPrivate); err != nil {
		t.Fatalf("AuthorizeKeyRotation() error = %v", err)
	}
	return fixture
}

func rotationFixtureWithoutAuthority(t *testing.T) rotationTestFixture {
	t.Helper()
	root := t.TempDir()
	oldCreated, err := identity.Create(filepath.Join(root, "old"), rand.Reader)
	if err != nil {
		t.Fatalf("create old identity: %v", err)
	}
	newCreated, err := identity.Create(filepath.Join(root, "new"), rand.Reader)
	if err != nil {
		t.Fatalf("create new identity: %v", err)
	}
	oldPrivate, err := identity.LoadPrivate(oldCreated.PrivatePath)
	if err != nil {
		t.Fatalf("load old identity: %v", err)
	}
	newPrivate, err := identity.LoadPrivate(newCreated.PrivatePath)
	if err != nil {
		t.Fatalf("load new identity: %v", err)
	}
	parentAsset := writeFixture(t, root, "parent.gguf", "parent")
	parentRoot := filepath.Join(root, "parent.ark")
	parent, err := Pack("parent", parentRoot, []AssetSource{{Role: "model", Path: parentAsset}}, time.Now())
	if err != nil {
		t.Fatalf("pack parent: %v", err)
	}
	if _, err := Sign(parentRoot, oldPrivate); err != nil {
		t.Fatalf("sign parent: %v", err)
	}
	parent, parentAuth, err := VerifyAuthenticated(parentRoot, VerifyOptions{})
	if err != nil {
		t.Fatalf("verify parent: %v", err)
	}
	childAsset := writeFixture(t, root, "child.gguf", "child")
	childRoot := filepath.Join(root, "child.ark")
	child, err := PackWithOptions("child", childRoot, []AssetSource{{Role: "model", Path: childAsset}}, time.Now(), PackOptions{
		ParentCapsuleID: parent.CapsuleID,
	})
	if err != nil {
		t.Fatalf("pack child: %v", err)
	}
	if _, err := Sign(childRoot, newPrivate); err != nil {
		t.Fatalf("sign child: %v", err)
	}
	child, childAuth, err := VerifyAuthenticated(childRoot, VerifyOptions{})
	if err != nil {
		t.Fatalf("verify child: %v", err)
	}
	return rotationTestFixture{
		parent: parent, parentAuth: parentAuth, parentPrivate: oldPrivate,
		childRoot: childRoot, child: child, childAuth: childAuth, childPrivate: newPrivate,
	}
}
