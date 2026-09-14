package capsule

import (
	"strings"
	"testing"

	"arkmesh/internal/identity"
)

func TestProtocolPayloadVectors(t *testing.T) {
	t.Parallel()
	parentID := "sha256:" + strings.Repeat("a", 64)
	childID := "sha256:" + strings.Repeat("b", 64)
	oldSigner := "ed25519:" + strings.Repeat("c", 64)
	newSigner := "ed25519:" + strings.Repeat("d", 64)
	policyID := "sha256:" + strings.Repeat("e", 64)

	if got, want := string(signaturePayload(childID)), "arkmesh.capsule.signature/v0alpha1\n"+childID+"\n"; got != want {
		t.Fatalf("signature payload = %q, want %q", got, want)
	}
	authority := AuthorityTransition{Action: AuthorityActionRotate, ParentCapsuleID: parentID, ChildCapsuleID: childID, PreviousSignerID: oldSigner, NextSigner: identity.PublicIdentity{KeyID: newSigner}}
	if got, want := string(authorityPayload(authority)), "arkmesh.capsule.authority/v0alpha1\nrotate\n"+parentID+"\n"+childID+"\n"+oldSigner+"\n"+newSigner+"\n"; got != want {
		t.Fatalf("authority payload = %q, want %q", got, want)
	}
	if got, want := string(recoveryPayload(policyID, RecoveryAction, parentID, childID, oldSigner, newSigner)), "arkmesh.capsule.recovery/v0alpha1\n"+policyID+"\nrecover\n"+parentID+"\n"+childID+"\n"+oldSigner+"\n"+newSigner+"\n"; got != want {
		t.Fatalf("recovery payload = %q, want %q", got, want)
	}
}

func TestRecoveryPolicyIDVector(t *testing.T) {
	t.Parallel()
	policy := RecoveryPolicy{Threshold: 2, Members: []identity.PublicIdentity{
		{KeyID: "ed25519:" + strings.Repeat("0", 64)},
		{KeyID: "ed25519:" + strings.Repeat("1", 64)},
	}}
	if got, want := recoveryPolicyID(policy), "sha256:c4873a51b3e37fcc48436365711e3d75bb74cd0d86d9b9890e4bffce5d731c43"; got != want {
		t.Fatalf("recovery policy ID = %s, want %s", got, want)
	}
}
