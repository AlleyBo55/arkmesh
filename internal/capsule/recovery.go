package capsule

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"arkmesh/internal/identity"
)

const (
	RecoveryPolicySchemaVersion   = "arkmesh.recovery-policy/v0alpha1"
	RecoveryApprovalSchemaVersion = "arkmesh.recovery-approval/v0alpha1"
	RecoverySchemaVersion         = "arkmesh.recovery/v0alpha1"
	RecoveryAction                = "recover"
	RecoveryFileName              = "recovery.json"
)

type RecoveryPolicy struct {
	SchemaVersion string                    `json:"schema_version"`
	PolicyID      string                    `json:"policy_id"`
	Threshold     int                       `json:"threshold"`
	Members       []identity.PublicIdentity `json:"members"`
}

type RecoveryApproval struct {
	SchemaVersion    string                  `json:"schema_version"`
	PolicyID         string                  `json:"policy_id"`
	Action           string                  `json:"action"`
	ParentCapsuleID  string                  `json:"parent_capsule_id"`
	ChildCapsuleID   string                  `json:"child_capsule_id"`
	PreviousSignerID string                  `json:"previous_signer_id"`
	NextSigner       identity.PublicIdentity `json:"next_signer"`
	ApproverID       string                  `json:"approver_id"`
	Signature        string                  `json:"signature"`
}

type RecoverySignature struct {
	ApproverID string `json:"approver_id"`
	Signature  string `json:"signature"`
}

type RecoveryTransition struct {
	SchemaVersion    string                  `json:"schema_version"`
	PolicyID         string                  `json:"policy_id"`
	Action           string                  `json:"action"`
	ParentCapsuleID  string                  `json:"parent_capsule_id"`
	ChildCapsuleID   string                  `json:"child_capsule_id"`
	PreviousSignerID string                  `json:"previous_signer_id"`
	NextSigner       identity.PublicIdentity `json:"next_signer"`
	Approvals        []RecoverySignature     `json:"approvals"`
}

func NewRecoveryPolicy(threshold int, members []identity.PublicIdentity) (RecoveryPolicy, error) {
	if len(members) < 2 {
		return RecoveryPolicy{}, errors.New("recovery policy requires at least two members")
	}
	if threshold < 2 || threshold > len(members) {
		return RecoveryPolicy{}, fmt.Errorf("recovery threshold must be between 2 and %d", len(members))
	}
	sorted := append([]identity.PublicIdentity(nil), members...)
	for _, member := range sorted {
		if _, err := member.Key(); err != nil {
			return RecoveryPolicy{}, fmt.Errorf("validate recovery member: %w", err)
		}
	}
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].KeyID < sorted[j].KeyID })
	for index := 1; index < len(sorted); index++ {
		if sorted[index-1].KeyID == sorted[index].KeyID {
			return RecoveryPolicy{}, fmt.Errorf("duplicate recovery member %s", sorted[index].KeyID)
		}
	}
	policy := RecoveryPolicy{SchemaVersion: RecoveryPolicySchemaVersion, Threshold: threshold, Members: sorted}
	policy.PolicyID = recoveryPolicyID(policy)
	return policy, nil
}

func WriteRecoveryPolicy(path string, policy RecoveryPolicy) error {
	if err := policy.Validate(); err != nil {
		return err
	}
	return writeExclusiveJSON(path, policy, "recovery policy")
}

func LoadRecoveryPolicy(path string) (RecoveryPolicy, error) {
	var policy RecoveryPolicy
	if err := loadStrictRegularJSON(path, &policy, "recovery policy"); err != nil {
		return RecoveryPolicy{}, err
	}
	if err := policy.Validate(); err != nil {
		return RecoveryPolicy{}, err
	}
	return policy, nil
}

func (policy RecoveryPolicy) Validate() error {
	if policy.SchemaVersion != RecoveryPolicySchemaVersion {
		return fmt.Errorf("unsupported recovery policy schema %q", policy.SchemaVersion)
	}
	if len(policy.Members) < 2 || policy.Threshold < 2 || policy.Threshold > len(policy.Members) {
		return errors.New("recovery policy requires at least two members and a valid threshold of at least two")
	}
	previous := ""
	for index, member := range policy.Members {
		if _, err := member.Key(); err != nil {
			return fmt.Errorf("validate recovery member: %w", err)
		}
		if index > 0 && member.KeyID <= previous {
			return errors.New("recovery members must be unique and sorted by key ID")
		}
		previous = member.KeyID
	}
	if policy.PolicyID != recoveryPolicyID(policy) {
		return errors.New("recovery policy ID does not match policy contents")
	}
	return nil
}

func ApproveRecovery(path string, policy RecoveryPolicy, parent Manifest, parentAuth Authenticity, child Manifest, childAuth Authenticity, approver identity.PrivateIdentity) (RecoveryApproval, error) {
	if err := policy.Validate(); err != nil {
		return RecoveryApproval{}, err
	}
	if err := validateLineageEndpoints(child, childAuth, parent, parentAuth); err != nil {
		return RecoveryApproval{}, err
	}
	if childAuth.SignerID == parentAuth.SignerID {
		return RecoveryApproval{}, errors.New("recovery requires a different child signer")
	}
	public, err := approver.Public()
	if err != nil {
		return RecoveryApproval{}, fmt.Errorf("validate recovery approver: %w", err)
	}
	if _, ok := recoveryMember(policy, public.KeyID); !ok {
		return RecoveryApproval{}, fmt.Errorf("approver %s is not a recovery policy member", public.KeyID)
	}
	_, privateKey, err := approver.Keys()
	if err != nil {
		return RecoveryApproval{}, fmt.Errorf("load recovery approver keys: %w", err)
	}
	approval := RecoveryApproval{
		SchemaVersion: RecoveryApprovalSchemaVersion,
		PolicyID:      policy.PolicyID, Action: RecoveryAction,
		ParentCapsuleID: parent.CapsuleID, ChildCapsuleID: child.CapsuleID,
		PreviousSignerID: parentAuth.SignerID, NextSigner: childAuth.Signer,
		ApproverID: public.KeyID,
	}
	approval.Signature = base64.StdEncoding.EncodeToString(ed25519.Sign(privateKey, recoveryPayload(approval.PolicyID, approval.Action, approval.ParentCapsuleID, approval.ChildCapsuleID, approval.PreviousSignerID, approval.NextSigner.KeyID)))
	if err := writeExclusiveJSON(path, approval, "recovery approval"); err != nil {
		return RecoveryApproval{}, err
	}
	return approval, nil
}

func LoadRecoveryApproval(path string) (RecoveryApproval, error) {
	var approval RecoveryApproval
	if err := loadStrictRegularJSON(path, &approval, "recovery approval"); err != nil {
		return RecoveryApproval{}, err
	}
	return approval, nil
}

func AssembleRecovery(childRoot string, policy RecoveryPolicy, parent Manifest, parentAuth Authenticity, child Manifest, childAuth Authenticity, approvals []RecoveryApproval, revocations identity.RevocationSet) (RecoveryTransition, error) {
	if err := policy.Validate(); err != nil {
		return RecoveryTransition{}, err
	}
	if err := revocations.Validate(); err != nil {
		return RecoveryTransition{}, fmt.Errorf("invalid revocation policy: %w", err)
	}
	if err := validateLineageEndpoints(child, childAuth, parent, parentAuth); err != nil {
		return RecoveryTransition{}, err
	}
	if childAuth.SignerID == parentAuth.SignerID {
		return RecoveryTransition{}, errors.New("recovery requires a different child signer")
	}
	if _, err := os.Lstat(filepath.Join(childRoot, AuthorityFileName)); err == nil {
		return RecoveryTransition{}, errors.New("child must not include both rotation authority and recovery transition")
	} else if !errors.Is(err, os.ErrNotExist) {
		return RecoveryTransition{}, fmt.Errorf("inspect key rotation authority: %w", err)
	}
	transition := RecoveryTransition{
		SchemaVersion: RecoverySchemaVersion, PolicyID: policy.PolicyID, Action: RecoveryAction,
		ParentCapsuleID: parent.CapsuleID, ChildCapsuleID: child.CapsuleID,
		PreviousSignerID: parentAuth.SignerID, NextSigner: childAuth.Signer,
	}
	seen := map[string]struct{}{}
	for _, approval := range approvals {
		if err := verifyRecoveryApproval(policy, transition, approval); err != nil {
			return RecoveryTransition{}, err
		}
		if _, exists := seen[approval.ApproverID]; exists {
			return RecoveryTransition{}, fmt.Errorf("duplicate recovery approval from %s", approval.ApproverID)
		}
		seen[approval.ApproverID] = struct{}{}
		if !revocations.Contains(approval.ApproverID) {
			transition.Approvals = append(transition.Approvals, RecoverySignature{ApproverID: approval.ApproverID, Signature: approval.Signature})
		}
	}
	sort.Slice(transition.Approvals, func(i, j int) bool { return transition.Approvals[i].ApproverID < transition.Approvals[j].ApproverID })
	if len(transition.Approvals) < policy.Threshold {
		return RecoveryTransition{}, fmt.Errorf("recovery has %d valid unrevoked approvals, requires %d", len(transition.Approvals), policy.Threshold)
	}
	if err := writeExclusiveJSON(filepath.Join(childRoot, RecoveryFileName), transition, "recovery transition"); err != nil {
		return RecoveryTransition{}, err
	}
	return transition, nil
}

func verifyThresholdRecovery(childRoot string, policy RecoveryPolicy, revocations identity.RevocationSet, parent Manifest, parentAuth Authenticity, child Manifest, childAuth Authenticity) (Lineage, error) {
	if err := policy.Validate(); err != nil {
		return Lineage{}, err
	}
	if err := revocations.Validate(); err != nil {
		return Lineage{}, fmt.Errorf("invalid revocation policy: %w", err)
	}
	var transition RecoveryTransition
	if err := loadStrictRegularJSON(filepath.Join(childRoot, RecoveryFileName), &transition, "recovery transition"); err != nil {
		return Lineage{}, err
	}
	if transition.SchemaVersion != RecoverySchemaVersion || transition.Action != RecoveryAction {
		return Lineage{}, errors.New("unsupported recovery transition")
	}
	if transition.PolicyID != policy.PolicyID || transition.ParentCapsuleID != parent.CapsuleID || transition.ParentCapsuleID != child.ParentCapsuleID || transition.ChildCapsuleID != child.CapsuleID || transition.PreviousSignerID != parentAuth.SignerID {
		return Lineage{}, errors.New("recovery transition does not match exact lineage edge")
	}
	nextKey, err := transition.NextSigner.Key()
	if err != nil {
		return Lineage{}, fmt.Errorf("validate recovery next signer: %w", err)
	}
	childKey, err := childAuth.Signer.Key()
	if err != nil {
		return Lineage{}, fmt.Errorf("validate child signer: %w", err)
	}
	if transition.NextSigner.KeyID != childAuth.SignerID || !bytes.Equal(nextKey, childKey) {
		return Lineage{}, errors.New("recovery next signer does not match child signer")
	}
	valid := 0
	previous := ""
	for _, approval := range transition.Approvals {
		if approval.ApproverID <= previous {
			return Lineage{}, errors.New("recovery approvals must be unique and sorted")
		}
		previous = approval.ApproverID
		member, ok := recoveryMember(policy, approval.ApproverID)
		if !ok {
			return Lineage{}, fmt.Errorf("recovery approver %s is not a policy member", approval.ApproverID)
		}
		signature, err := decodeRecoverySignature(approval.Signature)
		if err != nil {
			return Lineage{}, err
		}
		key, _ := member.Key()
		if !ed25519.Verify(key, recoveryPayload(transition.PolicyID, transition.Action, transition.ParentCapsuleID, transition.ChildCapsuleID, transition.PreviousSignerID, transition.NextSigner.KeyID), signature) {
			return Lineage{}, fmt.Errorf("invalid recovery approval from %s", approval.ApproverID)
		}
		if !revocations.Contains(approval.ApproverID) {
			valid++
		}
	}
	if valid < policy.Threshold {
		return Lineage{}, fmt.Errorf("recovery has %d valid unrevoked approvals, requires %d", valid, policy.Threshold)
	}
	return Lineage{Status: LineageThresholdRecovery, ParentCapsuleID: parent.CapsuleID, SignerID: childAuth.SignerID}, nil
}

func verifyRecoveryApproval(policy RecoveryPolicy, transition RecoveryTransition, approval RecoveryApproval) error {
	if approval.SchemaVersion != RecoveryApprovalSchemaVersion || approval.PolicyID != transition.PolicyID || approval.Action != transition.Action || approval.ParentCapsuleID != transition.ParentCapsuleID || approval.ChildCapsuleID != transition.ChildCapsuleID || approval.PreviousSignerID != transition.PreviousSignerID || approval.NextSigner.KeyID != transition.NextSigner.KeyID {
		return fmt.Errorf("recovery approval from %s does not match exact transition", approval.ApproverID)
	}
	member, ok := recoveryMember(policy, approval.ApproverID)
	if !ok {
		return fmt.Errorf("approver %s is not a recovery policy member", approval.ApproverID)
	}
	signature, err := decodeRecoverySignature(approval.Signature)
	if err != nil {
		return err
	}
	key, _ := member.Key()
	if !ed25519.Verify(key, recoveryPayload(approval.PolicyID, approval.Action, approval.ParentCapsuleID, approval.ChildCapsuleID, approval.PreviousSignerID, approval.NextSigner.KeyID), signature) {
		return fmt.Errorf("invalid recovery approval from %s", approval.ApproverID)
	}
	return nil
}

func recoveryMember(policy RecoveryPolicy, keyID string) (identity.PublicIdentity, bool) {
	index := sort.Search(len(policy.Members), func(i int) bool { return policy.Members[i].KeyID >= keyID })
	if index < len(policy.Members) && policy.Members[index].KeyID == keyID {
		return policy.Members[index], true
	}
	return identity.PublicIdentity{}, false
}

func recoveryPolicyID(policy RecoveryPolicy) string {
	payload := "arkmesh.recovery.policy/v0alpha1\n" + strconv.Itoa(policy.Threshold) + "\n"
	for _, member := range policy.Members {
		payload += member.KeyID + "\n"
	}
	digest := sha256.Sum256([]byte(payload))
	return "sha256:" + hex.EncodeToString(digest[:])
}

func recoveryPayload(policyID, action, parentID, childID, previousSignerID, nextSignerID string) []byte {
	return []byte("arkmesh.capsule.recovery/v0alpha1\n" + policyID + "\n" + action + "\n" + parentID + "\n" + childID + "\n" + previousSignerID + "\n" + nextSignerID + "\n")
}

func decodeRecoverySignature(value string) ([]byte, error) {
	signature, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, fmt.Errorf("decode recovery signature: %w", err)
	}
	if len(signature) != ed25519.SignatureSize {
		return nil, fmt.Errorf("invalid recovery signature length %d", len(signature))
	}
	return signature, nil
}

func writeExclusiveJSON(path string, value any, description string) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("encode %s: %w", description, err)
	}
	data = append(data, '\n')
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return fmt.Errorf("create %s: %w", description, err)
	}
	if _, err := file.Write(data); err != nil {
		_ = file.Close()
		_ = os.Remove(path)
		return fmt.Errorf("write %s: %w", description, err)
	}
	if err := file.Close(); err != nil {
		_ = os.Remove(path)
		return fmt.Errorf("close %s: %w", description, err)
	}
	return nil
}

func loadStrictRegularJSON(path string, value any, description string) error {
	info, err := os.Lstat(path)
	if err != nil {
		return fmt.Errorf("inspect %s: %w", description, err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", description)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", description, err)
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(value); err != nil {
		return fmt.Errorf("decode %s: %w", description, err)
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		if err == nil {
			return fmt.Errorf("decode %s: unexpected trailing JSON value", description)
		}
		return fmt.Errorf("decode %s: %w", description, err)
	}
	return nil
}
