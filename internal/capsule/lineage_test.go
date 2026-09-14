package capsule

import (
	"crypto/rand"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"arkmesh/internal/identity"
)

func TestEmptyParentPreservesLegacyManifestShape(t *testing.T) {
	t.Parallel()
	manifest := Manifest{
		SchemaVersion: SchemaVersion,
		Name:          "legacy",
		CreatedAt:     "2026-09-14T02:00:00Z",
		Assets: []Asset{{
			Role:   "model",
			Name:   "model.gguf",
			Size:   1,
			SHA256: strings.Repeat("a", 64),
		}},
	}
	data, err := json.Marshal(manifest)
	if err != nil {
		t.Fatalf("marshal manifest: %v", err)
	}
	if strings.Contains(string(data), "parent_capsule_id") {
		t.Fatalf("legacy manifest unexpectedly contains parent: %s", data)
	}
}

func TestPackWithParentBindsParentIntoCapsuleID(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	model := writeFixture(t, root, "model.gguf", "same model")
	created := time.Date(2026, 9, 14, 2, 0, 0, 0, time.UTC)
	base, err := Pack("demo", filepath.Join(root, "base.ark"), []AssetSource{{Role: "model", Path: model}}, created)
	if err != nil {
		t.Fatalf("Pack() error = %v", err)
	}
	child, err := PackWithOptions("demo", filepath.Join(root, "child.ark"), []AssetSource{{Role: "model", Path: model}}, created, PackOptions{
		ParentCapsuleID: base.CapsuleID,
	})
	if err != nil {
		t.Fatalf("PackWithOptions() error = %v", err)
	}
	if child.ParentCapsuleID != base.CapsuleID {
		t.Fatalf("child parent = %q, want %q", child.ParentCapsuleID, base.CapsuleID)
	}
	if child.CapsuleID == base.CapsuleID {
		t.Fatal("child and base capsule IDs match; parent was not bound")
	}
}

func TestCheckLineageAcceptsSameAuthor(t *testing.T) {
	t.Parallel()
	parent, parentAuth, child, childAuth := lineageFixture(t, false)
	lineage, err := CheckLineage(child, childAuth, parent, parentAuth)
	if err != nil {
		t.Fatalf("CheckLineage() error = %v", err)
	}
	if lineage.Status != LineageVerified || lineage.ParentCapsuleID != parent.CapsuleID || lineage.SignerID != childAuth.SignerID {
		t.Fatalf("lineage = %+v, want verified same author", lineage)
	}
}

func TestCheckLineageRejectsDifferentAuthor(t *testing.T) {
	t.Parallel()
	parent, parentAuth, child, childAuth := lineageFixture(t, true)
	if _, err := CheckLineage(child, childAuth, parent, parentAuth); err == nil || !strings.Contains(err.Error(), "not authorized") {
		t.Fatalf("CheckLineage() error = %v, want unauthorized signer", err)
	}
}

func TestCheckLineageRejectsWrongParent(t *testing.T) {
	t.Parallel()
	parent, parentAuth, child, childAuth := lineageFixture(t, false)
	parent.CapsuleID = "sha256:" + strings.Repeat("0", 64)
	if _, err := CheckLineage(child, childAuth, parent, parentAuth); err == nil || !strings.Contains(err.Error(), "does not match") {
		t.Fatalf("CheckLineage() error = %v, want parent mismatch", err)
	}
}

func TestCheckLineageRejectsUnsignedParent(t *testing.T) {
	t.Parallel()
	parent, _, child, childAuth := lineageFixture(t, false)
	if _, err := CheckLineage(child, childAuth, parent, Authenticity{Status: AuthenticityUnsigned}); err == nil || !strings.Contains(err.Error(), "parent capsule must have") {
		t.Fatalf("CheckLineage() error = %v, want unsigned parent", err)
	}
}

func lineageFixture(t *testing.T, differentChildAuthor bool) (Manifest, Authenticity, Manifest, Authenticity) {
	t.Helper()
	root := t.TempDir()
	firstCreated, err := identity.Create(filepath.Join(root, "first-author"), rand.Reader)
	if err != nil {
		t.Fatalf("create first identity: %v", err)
	}
	firstPrivate, err := identity.LoadPrivate(firstCreated.PrivatePath)
	if err != nil {
		t.Fatalf("load first identity: %v", err)
	}

	parentAsset := writeFixture(t, root, "parent.gguf", "parent")
	parentRoot := filepath.Join(root, "parent.ark")
	parent, err := Pack("parent", parentRoot, []AssetSource{{Role: "model", Path: parentAsset}}, time.Now())
	if err != nil {
		t.Fatalf("pack parent: %v", err)
	}
	if _, err := Sign(parentRoot, firstPrivate); err != nil {
		t.Fatalf("sign parent: %v", err)
	}
	parent, parentAuth, err := VerifyAuthenticated(parentRoot, VerifyOptions{})
	if err != nil {
		t.Fatalf("verify parent: %v", err)
	}

	childPrivate := firstPrivate
	if differentChildAuthor {
		secondCreated, err := identity.Create(filepath.Join(root, "second-author"), rand.Reader)
		if err != nil {
			t.Fatalf("create second identity: %v", err)
		}
		childPrivate, err = identity.LoadPrivate(secondCreated.PrivatePath)
		if err != nil {
			t.Fatalf("load second identity: %v", err)
		}
	}
	childAsset := writeFixture(t, root, "child.gguf", "child")
	childRoot := filepath.Join(root, "child.ark")
	child, err := PackWithOptions("child", childRoot, []AssetSource{{Role: "model", Path: childAsset}}, time.Now(), PackOptions{
		ParentCapsuleID: parent.CapsuleID,
	})
	if err != nil {
		t.Fatalf("pack child: %v", err)
	}
	if _, err := Sign(childRoot, childPrivate); err != nil {
		t.Fatalf("sign child: %v", err)
	}
	child, childAuth, err := VerifyAuthenticated(childRoot, VerifyOptions{})
	if err != nil {
		t.Fatalf("verify child: %v", err)
	}
	return parent, parentAuth, child, childAuth
}
