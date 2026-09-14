package capsule

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"arkmesh/internal/identity"
)

const (
	AuthoritySchemaVersion = "arkmesh.authority/v0alpha1"
	AuthorityFileName      = "authority.json"
	AuthorityActionRotate  = "rotate"
)

type AuthorityTransition struct {
	SchemaVersion    string                  `json:"schema_version"`
	Action           string                  `json:"action"`
	ParentCapsuleID  string                  `json:"parent_capsule_id"`
	ChildCapsuleID   string                  `json:"child_capsule_id"`
	PreviousSignerID string                  `json:"previous_signer_id"`
	NextSigner       identity.PublicIdentity `json:"next_signer"`
	Signature        string                  `json:"signature"`
}

func AuthorizeKeyRotation(childRoot string, parent Manifest, parentAuthenticity Authenticity, child Manifest, childAuthenticity Authenticity, previousSigner identity.PrivateIdentity) (AuthorityTransition, error) {
	if err := validateLineageEndpoints(child, childAuthenticity, parent, parentAuthenticity); err != nil {
		return AuthorityTransition{}, err
	}
	if childAuthenticity.SignerID == parentAuthenticity.SignerID {
		return AuthorityTransition{}, errors.New("key rotation requires a different child signer")
	}
	previousPublic, err := previousSigner.Public()
	if err != nil {
		return AuthorityTransition{}, fmt.Errorf("validate previous signer: %w", err)
	}
	if previousPublic.KeyID != parentAuthenticity.SignerID {
		return AuthorityTransition{}, fmt.Errorf("rotation signer %s does not match parent signer %s", previousPublic.KeyID, parentAuthenticity.SignerID)
	}
	_, privateKey, err := previousSigner.Keys()
	if err != nil {
		return AuthorityTransition{}, fmt.Errorf("load previous signer keys: %w", err)
	}
	transition := AuthorityTransition{
		SchemaVersion:    AuthoritySchemaVersion,
		Action:           AuthorityActionRotate,
		ParentCapsuleID:  parent.CapsuleID,
		ChildCapsuleID:   child.CapsuleID,
		PreviousSignerID: parentAuthenticity.SignerID,
		NextSigner:       childAuthenticity.Signer,
	}
	transition.Signature = base64.StdEncoding.EncodeToString(ed25519.Sign(privateKey, authorityPayload(transition)))
	if err := writeAuthority(childRoot, transition); err != nil {
		return AuthorityTransition{}, err
	}
	return transition, nil
}

func verifyKeyRotation(childRoot string, parent Manifest, parentAuthenticity Authenticity, child Manifest, childAuthenticity Authenticity) (Lineage, error) {
	path := filepath.Join(childRoot, AuthorityFileName)
	info, err := os.Lstat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Lineage{}, errors.New("different child signer requires a key rotation authority")
		}
		return Lineage{}, fmt.Errorf("inspect key rotation authority: %w", err)
	}
	if !info.Mode().IsRegular() {
		return Lineage{}, errors.New("key rotation authority is not a regular file")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Lineage{}, fmt.Errorf("read key rotation authority: %w", err)
	}
	var transition AuthorityTransition
	if err := decodeAuthority(data, &transition); err != nil {
		return Lineage{}, fmt.Errorf("decode key rotation authority: %w", err)
	}
	if transition.SchemaVersion != AuthoritySchemaVersion {
		return Lineage{}, fmt.Errorf("unsupported authority schema %q", transition.SchemaVersion)
	}
	if transition.Action != AuthorityActionRotate {
		return Lineage{}, fmt.Errorf("unsupported authority action %q", transition.Action)
	}
	if transition.ParentCapsuleID != parent.CapsuleID || transition.ParentCapsuleID != child.ParentCapsuleID {
		return Lineage{}, errors.New("authority parent capsule ID does not match lineage")
	}
	if transition.ChildCapsuleID != child.CapsuleID {
		return Lineage{}, errors.New("authority child capsule ID does not match child")
	}
	if transition.PreviousSignerID != parentAuthenticity.SignerID {
		return Lineage{}, errors.New("authority previous signer does not match parent signer")
	}
	nextKey, err := transition.NextSigner.Key()
	if err != nil {
		return Lineage{}, fmt.Errorf("validate authority next signer: %w", err)
	}
	childKey, err := childAuthenticity.Signer.Key()
	if err != nil {
		return Lineage{}, fmt.Errorf("validate child signer: %w", err)
	}
	if transition.NextSigner.KeyID != childAuthenticity.SignerID || !bytes.Equal(nextKey, childKey) {
		return Lineage{}, errors.New("authority next signer does not match child signer")
	}
	previousKey, err := parentAuthenticity.Signer.Key()
	if err != nil {
		return Lineage{}, fmt.Errorf("validate parent signer: %w", err)
	}
	signature, err := base64.StdEncoding.DecodeString(transition.Signature)
	if err != nil {
		return Lineage{}, fmt.Errorf("decode authority signature: %w", err)
	}
	if len(signature) != ed25519.SignatureSize {
		return Lineage{}, fmt.Errorf("invalid authority signature length %d", len(signature))
	}
	if !ed25519.Verify(previousKey, authorityPayload(transition), signature) {
		return Lineage{}, errors.New("invalid key rotation authority signature")
	}
	return Lineage{
		Status:          LineageKeyRotation,
		ParentCapsuleID: parent.CapsuleID,
		SignerID:        childAuthenticity.SignerID,
	}, nil
}

func authorityPayload(transition AuthorityTransition) []byte {
	return []byte(
		"arkmesh.capsule.authority/v0alpha1\n" +
			transition.Action + "\n" +
			transition.ParentCapsuleID + "\n" +
			transition.ChildCapsuleID + "\n" +
			transition.PreviousSignerID + "\n" +
			transition.NextSigner.KeyID + "\n",
	)
}

func writeAuthority(root string, transition AuthorityTransition) error {
	data, err := json.MarshalIndent(transition, "", "  ")
	if err != nil {
		return fmt.Errorf("encode key rotation authority: %w", err)
	}
	data = append(data, '\n')
	path := filepath.Join(root, AuthorityFileName)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return fmt.Errorf("create key rotation authority: %w", err)
	}
	if _, err := file.Write(data); err != nil {
		_ = file.Close()
		return fmt.Errorf("write key rotation authority: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close key rotation authority: %w", err)
	}
	return nil
}

func decodeAuthority(data []byte, transition *AuthorityTransition) error {
	if err := rejectDuplicateJSONKeys(data); err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(transition); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		if err == nil {
			return errors.New("unexpected trailing JSON value")
		}
		return err
	}
	return nil
}
