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
	SignatureSchemaVersion = "arkmesh.signature/v0alpha1"
	SignatureFileName      = "signature.json"
)

type AuthenticityStatus string

const (
	AuthenticityUnsigned AuthenticityStatus = "unsigned"
	AuthenticityUnknown  AuthenticityStatus = "valid_unknown_author"
	AuthenticityTrusted  AuthenticityStatus = "valid_trusted_author"
)

type SignatureEnvelope struct {
	SchemaVersion string                  `json:"schema_version"`
	CapsuleID     string                  `json:"capsule_id"`
	Signer        identity.PublicIdentity `json:"signer"`
	Signature     string                  `json:"signature"`
}

type Authenticity struct {
	Status   AuthenticityStatus
	SignerID string
}

type VerifyOptions struct {
	Trusted          []identity.PublicIdentity
	RequireSignature bool
	RequireTrusted   bool
}

func Sign(root string, signer identity.PrivateIdentity) (SignatureEnvelope, error) {
	manifest, err := Verify(root)
	if err != nil {
		return SignatureEnvelope{}, err
	}
	publicIdentity, err := signer.Public()
	if err != nil {
		return SignatureEnvelope{}, fmt.Errorf("validate signer: %w", err)
	}
	_, privateKey, err := signer.Keys()
	if err != nil {
		return SignatureEnvelope{}, fmt.Errorf("load signer keys: %w", err)
	}
	signature := ed25519.Sign(privateKey, signaturePayload(manifest.CapsuleID))
	envelope := SignatureEnvelope{
		SchemaVersion: SignatureSchemaVersion,
		CapsuleID:     manifest.CapsuleID,
		Signer:        publicIdentity,
		Signature:     base64.StdEncoding.EncodeToString(signature),
	}
	if err := writeSignature(root, envelope); err != nil {
		return SignatureEnvelope{}, err
	}
	return envelope, nil
}

func VerifyAuthenticated(root string, options VerifyOptions) (Manifest, Authenticity, error) {
	manifest, err := Verify(root)
	if err != nil {
		return Manifest{}, Authenticity{}, err
	}
	authenticity, err := verifySignature(root, manifest, options)
	if err != nil {
		return Manifest{}, Authenticity{}, err
	}
	return manifest, authenticity, nil
}

func verifySignature(root string, manifest Manifest, options VerifyOptions) (Authenticity, error) {
	data, err := os.ReadFile(filepath.Join(root, SignatureFileName))
	if errors.Is(err, os.ErrNotExist) {
		if options.RequireSignature || options.RequireTrusted {
			return Authenticity{}, errors.New("capsule is unsigned")
		}
		return Authenticity{Status: AuthenticityUnsigned}, nil
	}
	if err != nil {
		return Authenticity{}, fmt.Errorf("read capsule signature: %w", err)
	}

	var envelope SignatureEnvelope
	if err := decodeSignature(data, &envelope); err != nil {
		return Authenticity{}, fmt.Errorf("decode capsule signature: %w", err)
	}
	if envelope.SchemaVersion != SignatureSchemaVersion {
		return Authenticity{}, fmt.Errorf("unsupported signature schema %q", envelope.SchemaVersion)
	}
	if envelope.CapsuleID != manifest.CapsuleID {
		return Authenticity{}, errors.New("signature capsule ID does not match manifest")
	}
	publicKey, err := envelope.Signer.Key()
	if err != nil {
		return Authenticity{}, fmt.Errorf("validate signature signer: %w", err)
	}
	signature, err := base64.StdEncoding.DecodeString(envelope.Signature)
	if err != nil {
		return Authenticity{}, fmt.Errorf("decode Ed25519 signature: %w", err)
	}
	if len(signature) != ed25519.SignatureSize {
		return Authenticity{}, fmt.Errorf("invalid Ed25519 signature length %d", len(signature))
	}
	if !ed25519.Verify(publicKey, signaturePayload(manifest.CapsuleID), signature) {
		return Authenticity{}, errors.New("invalid capsule signature")
	}

	trusted, err := isTrusted(envelope.Signer, options.Trusted)
	if err != nil {
		return Authenticity{}, err
	}
	if trusted {
		return Authenticity{Status: AuthenticityTrusted, SignerID: envelope.Signer.KeyID}, nil
	}
	if options.RequireTrusted {
		return Authenticity{}, fmt.Errorf("capsule signer %s is not trusted", envelope.Signer.KeyID)
	}
	return Authenticity{Status: AuthenticityUnknown, SignerID: envelope.Signer.KeyID}, nil
}

func isTrusted(signer identity.PublicIdentity, trusted []identity.PublicIdentity) (bool, error) {
	signerKey, err := signer.Key()
	if err != nil {
		return false, err
	}
	for _, candidate := range trusted {
		candidateKey, err := candidate.Key()
		if err != nil {
			return false, fmt.Errorf("invalid trusted identity %q: %w", candidate.KeyID, err)
		}
		if candidate.KeyID == signer.KeyID && bytes.Equal(candidateKey, signerKey) {
			return true, nil
		}
	}
	return false, nil
}

func signaturePayload(capsuleID string) []byte {
	return []byte("arkmesh.capsule.signature/v0alpha1\n" + capsuleID + "\n")
}

func writeSignature(root string, envelope SignatureEnvelope) error {
	data, err := json.MarshalIndent(envelope, "", "  ")
	if err != nil {
		return fmt.Errorf("encode capsule signature: %w", err)
	}
	data = append(data, '\n')
	path := filepath.Join(root, SignatureFileName)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return fmt.Errorf("create capsule signature: %w", err)
	}
	if _, err := file.Write(data); err != nil {
		_ = file.Close()
		return fmt.Errorf("write capsule signature: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close capsule signature: %w", err)
	}
	return nil
}

func decodeSignature(data []byte, envelope *SignatureEnvelope) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(envelope); err != nil {
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
