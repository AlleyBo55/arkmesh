package capsule

import (
	"crypto/rand"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"arkmesh/internal/identity"
)

func TestThresholdRecoveryAuthorizesExactChild(t *testing.T) {
	t.Parallel()
	fixture := rotationFixtureWithoutAuthority(t)
	policy, privateKeys := recoveryPolicyFixture(t, 2, 3)
	approvals := make([]RecoveryApproval, 0, len(privateKeys))
	for index, privateKey := range privateKeys {
		approval, err := ApproveRecovery(filepath.Join(t.TempDir(), "approval.json"), policy, fixture.parent, fixture.parentAuth, fixture.child, fixture.childAuth, privateKey)
		if err != nil {
			t.Fatalf("ApproveRecovery(%d) error = %v", index, err)
		}
		approvals = append(approvals, approval)
	}
	revocations, err := identity.NewRevocationSet([]identity.PublicIdentity{policy.Members[0]})
	if err != nil {
		t.Fatalf("NewRevocationSet() error = %v", err)
	}
	transition, err := AssembleRecovery(fixture.childRoot, policy, fixture.parent, fixture.parentAuth, fixture.child, fixture.childAuth, approvals, revocations)
	if err != nil {
		t.Fatalf("AssembleRecovery() error = %v", err)
	}
	if len(transition.Approvals) != 2 {
		t.Fatalf("transition approvals = %d, want 2 unrevoked approvals", len(transition.Approvals))
	}
	lineage, err := CheckLineageWithRecovery(fixture.childRoot, fixture.child, fixture.childAuth, fixture.parent, fixture.parentAuth, policy, revocations)
	if err != nil {
		t.Fatalf("CheckLineageWithRecovery() error = %v", err)
	}
	if lineage.Status != LineageThresholdRecovery || lineage.SignerID != fixture.childAuth.SignerID {
		t.Fatalf("lineage = %+v, want threshold recovery", lineage)
	}
}

func TestRecoveryRejectsInsufficientApprovals(t *testing.T) {
	t.Parallel()
	fixture := rotationFixtureWithoutAuthority(t)
	policy, privateKeys := recoveryPolicyFixture(t, 2, 2)
	approval, err := ApproveRecovery(filepath.Join(t.TempDir(), "approval.json"), policy, fixture.parent, fixture.parentAuth, fixture.child, fixture.childAuth, privateKeys[0])
	if err != nil {
		t.Fatalf("ApproveRecovery() error = %v", err)
	}
	if _, err := AssembleRecovery(fixture.childRoot, policy, fixture.parent, fixture.parentAuth, fixture.child, fixture.childAuth, []RecoveryApproval{approval}, identity.RevocationSet{}); err == nil || !strings.Contains(err.Error(), "requires 2") {
		t.Fatalf("AssembleRecovery() error = %v, want threshold rejection", err)
	}
}

func TestRecoveryPolicyRejectsSingleApprovalThreshold(t *testing.T) {
	t.Parallel()
	policy, _ := recoveryPolicyFixture(t, 2, 2)
	if _, err := NewRecoveryPolicy(1, policy.Members); err == nil || !strings.Contains(err.Error(), "between 2") {
		t.Fatalf("NewRecoveryPolicy(1) error = %v, want single-key rejection", err)
	}
}

func TestRecoveryRejectsModifiedApproval(t *testing.T) {
	t.Parallel()
	fixture := rotationFixtureWithoutAuthority(t)
	policy, privateKeys := recoveryPolicyFixture(t, 2, 2)
	approvals := make([]RecoveryApproval, 0, 2)
	for _, privateKey := range privateKeys {
		approval, err := ApproveRecovery(filepath.Join(t.TempDir(), "approval.json"), policy, fixture.parent, fixture.parentAuth, fixture.child, fixture.childAuth, privateKey)
		if err != nil {
			t.Fatalf("ApproveRecovery() error = %v", err)
		}
		approvals = append(approvals, approval)
	}
	approvals[0].ChildCapsuleID = "sha256:" + strings.Repeat("0", 64)
	if _, err := AssembleRecovery(fixture.childRoot, policy, fixture.parent, fixture.parentAuth, fixture.child, fixture.childAuth, approvals, identity.RevocationSet{}); err == nil || !strings.Contains(err.Error(), "does not match exact transition") {
		t.Fatalf("AssembleRecovery() error = %v, want exact transition rejection", err)
	}
}

func TestRecoveryRequiresExplicitPolicy(t *testing.T) {
	t.Parallel()
	fixture := rotationFixtureWithoutAuthority(t)
	policy, privateKeys := recoveryPolicyFixture(t, 2, 2)
	var approvals []RecoveryApproval
	for _, privateKey := range privateKeys {
		approval, err := ApproveRecovery(filepath.Join(t.TempDir(), "approval.json"), policy, fixture.parent, fixture.parentAuth, fixture.child, fixture.childAuth, privateKey)
		if err != nil {
			t.Fatalf("ApproveRecovery() error = %v", err)
		}
		approvals = append(approvals, approval)
	}
	if _, err := AssembleRecovery(fixture.childRoot, policy, fixture.parent, fixture.parentAuth, fixture.child, fixture.childAuth, approvals, identity.RevocationSet{}); err != nil {
		t.Fatalf("AssembleRecovery() error = %v", err)
	}
	if _, err := CheckLineageWithAuthority(fixture.childRoot, fixture.child, fixture.childAuth, fixture.parent, fixture.parentAuth); err == nil || !strings.Contains(err.Error(), "explicit recovery policy") {
		t.Fatalf("CheckLineageWithAuthority() error = %v, want explicit policy requirement", err)
	}
}

func TestLoadRecoveryPolicyRejectsUnknownField(t *testing.T) {
	t.Parallel()
	policy, _ := recoveryPolicyFixture(t, 2, 2)
	data, err := json.Marshal(policy)
	if err != nil {
		t.Fatalf("marshal policy: %v", err)
	}
	data = append(data[:len(data)-1], []byte(`,"unknown":true}`)...)
	path := filepath.Join(t.TempDir(), "policy.json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write policy: %v", err)
	}
	if _, err := LoadRecoveryPolicy(path); err == nil || !strings.Contains(err.Error(), "unknown field") {
		t.Fatalf("LoadRecoveryPolicy() error = %v, want unknown field rejection", err)
	}
}

func recoveryPolicyFixture(t *testing.T, threshold, count int) (RecoveryPolicy, []identity.PrivateIdentity) {
	t.Helper()
	members := make([]identity.PublicIdentity, 0, count)
	privateKeys := make([]identity.PrivateIdentity, 0, count)
	for index := 0; index < count; index++ {
		created, err := identity.Create(filepath.Join(t.TempDir(), "recovery"), rand.Reader)
		if err != nil {
			t.Fatalf("create recovery identity %d: %v", index, err)
		}
		privateKey, err := identity.LoadPrivate(created.PrivatePath)
		if err != nil {
			t.Fatalf("load recovery identity %d: %v", index, err)
		}
		members = append(members, created.Public)
		privateKeys = append(privateKeys, privateKey)
	}
	policy, err := NewRecoveryPolicy(threshold, members)
	if err != nil {
		t.Fatalf("NewRecoveryPolicy() error = %v", err)
	}
	return policy, privateKeys
}
