package identity

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
)

const (
	PublicSchemaVersion  = "arkmesh.identity/v0alpha1"
	PrivateSchemaVersion = "arkmesh.private-identity/v0alpha1"
	Algorithm            = "ed25519"
	PublicFileName       = "identity.json"
	PrivateFileName      = "identity.key"
)

type PublicIdentity struct {
	SchemaVersion string `json:"schema_version"`
	Algorithm     string `json:"algorithm"`
	KeyID         string `json:"key_id"`
	PublicKey     string `json:"public_key"`
}

type PrivateIdentity struct {
	SchemaVersion string `json:"schema_version"`
	Algorithm     string `json:"algorithm"`
	KeyID         string `json:"key_id"`
	PublicKey     string `json:"public_key"`
	PrivateKey    string `json:"private_key"`
}

type CreatedIdentity struct {
	Public      PublicIdentity
	PublicPath  string
	PrivatePath string
}

func Create(directory string, random io.Reader) (CreatedIdentity, error) {
	if directory == "" {
		return CreatedIdentity{}, errors.New("identity directory is required")
	}
	if random == nil {
		return CreatedIdentity{}, errors.New("random source is required")
	}
	if _, err := os.Stat(directory); err == nil {
		return CreatedIdentity{}, fmt.Errorf("identity directory already exists: %s", directory)
	} else if !errors.Is(err, os.ErrNotExist) {
		return CreatedIdentity{}, fmt.Errorf("inspect identity directory: %w", err)
	}
	if err := os.MkdirAll(directory, 0o700); err != nil {
		return CreatedIdentity{}, fmt.Errorf("create identity directory: %w", err)
	}
	if err := os.Chmod(directory, 0o700); err != nil {
		_ = os.RemoveAll(directory)
		return CreatedIdentity{}, fmt.Errorf("secure identity directory: %w", err)
	}
	complete := false
	defer func() {
		if !complete {
			_ = os.RemoveAll(directory)
		}
	}()

	publicKey, privateKey, err := ed25519.GenerateKey(random)
	if err != nil {
		return CreatedIdentity{}, fmt.Errorf("generate Ed25519 identity: %w", err)
	}
	public := NewPublic(publicKey)
	private := PrivateIdentity{
		SchemaVersion: PrivateSchemaVersion,
		Algorithm:     Algorithm,
		KeyID:         public.KeyID,
		PublicKey:     public.PublicKey,
		PrivateKey:    base64.StdEncoding.EncodeToString(privateKey),
	}

	publicPath := filepath.Join(directory, PublicFileName)
	privatePath := filepath.Join(directory, PrivateFileName)
	if err := writeJSON(publicPath, public, 0o644); err != nil {
		return CreatedIdentity{}, err
	}
	if err := writeJSON(privatePath, private, 0o600); err != nil {
		return CreatedIdentity{}, err
	}

	complete = true
	return CreatedIdentity{
		Public:      public,
		PublicPath:  publicPath,
		PrivatePath: privatePath,
	}, nil
}

func NewPublic(key ed25519.PublicKey) PublicIdentity {
	copyOfKey := append(ed25519.PublicKey(nil), key...)
	return PublicIdentity{
		SchemaVersion: PublicSchemaVersion,
		Algorithm:     Algorithm,
		KeyID:         KeyID(copyOfKey),
		PublicKey:     base64.StdEncoding.EncodeToString(copyOfKey),
	}
}

func KeyID(key ed25519.PublicKey) string {
	digest := sha256.Sum256(key)
	return "ed25519:" + hex.EncodeToString(digest[:])
}

func LoadPublic(path string) (PublicIdentity, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return PublicIdentity{}, fmt.Errorf("read public identity: %w", err)
	}
	var public PublicIdentity
	if err := decodeStrict(data, &public); err != nil {
		return PublicIdentity{}, fmt.Errorf("decode public identity: %w", err)
	}
	if _, err := public.Key(); err != nil {
		return PublicIdentity{}, err
	}
	return public, nil
}

func LoadPrivate(path string) (PrivateIdentity, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return PrivateIdentity{}, fmt.Errorf("read private identity: %w", err)
	}
	var private PrivateIdentity
	if err := decodeStrict(data, &private); err != nil {
		return PrivateIdentity{}, fmt.Errorf("decode private identity: %w", err)
	}
	if _, _, err := private.Keys(); err != nil {
		return PrivateIdentity{}, err
	}
	return private, nil
}

func (identity PublicIdentity) Key() (ed25519.PublicKey, error) {
	if identity.SchemaVersion != PublicSchemaVersion {
		return nil, fmt.Errorf("unsupported public identity schema %q", identity.SchemaVersion)
	}
	if identity.Algorithm != Algorithm {
		return nil, fmt.Errorf("unsupported identity algorithm %q", identity.Algorithm)
	}
	decoded, err := base64.StdEncoding.DecodeString(identity.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("decode Ed25519 public key: %w", err)
	}
	if len(decoded) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("invalid Ed25519 public key length %d", len(decoded))
	}
	publicKey := ed25519.PublicKey(decoded)
	if identity.KeyID != KeyID(publicKey) {
		return nil, errors.New("public identity key ID does not match public key")
	}
	return publicKey, nil
}

func (identity PrivateIdentity) Keys() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	if identity.SchemaVersion != PrivateSchemaVersion {
		return nil, nil, fmt.Errorf("unsupported private identity schema %q", identity.SchemaVersion)
	}
	if identity.Algorithm != Algorithm {
		return nil, nil, fmt.Errorf("unsupported identity algorithm %q", identity.Algorithm)
	}
	publicBytes, err := base64.StdEncoding.DecodeString(identity.PublicKey)
	if err != nil {
		return nil, nil, fmt.Errorf("decode Ed25519 public key: %w", err)
	}
	if len(publicBytes) != ed25519.PublicKeySize {
		return nil, nil, fmt.Errorf("invalid Ed25519 public key length %d", len(publicBytes))
	}
	privateBytes, err := base64.StdEncoding.DecodeString(identity.PrivateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("decode Ed25519 private key: %w", err)
	}
	if len(privateBytes) != ed25519.PrivateKeySize {
		return nil, nil, fmt.Errorf("invalid Ed25519 private key length %d", len(privateBytes))
	}
	publicKey := ed25519.PublicKey(publicBytes)
	privateKey := ed25519.PrivateKey(privateBytes)
	derivedPublic := privateKey.Public().(ed25519.PublicKey)
	if !bytes.Equal(publicKey, derivedPublic) {
		return nil, nil, errors.New("private identity public key does not match private key")
	}
	if identity.KeyID != KeyID(publicKey) {
		return nil, nil, errors.New("private identity key ID does not match public key")
	}
	return publicKey, privateKey, nil
}

func (identity PrivateIdentity) Public() (PublicIdentity, error) {
	publicKey, _, err := identity.Keys()
	if err != nil {
		return PublicIdentity{}, err
	}
	return NewPublic(publicKey), nil
}

func writeJSON(path string, value any, permissions os.FileMode) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("encode %s: %w", filepath.Base(path), err)
	}
	data = append(data, '\n')
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, permissions)
	if err != nil {
		return fmt.Errorf("create %s: %w", filepath.Base(path), err)
	}
	if _, err := file.Write(data); err != nil {
		_ = file.Close()
		return fmt.Errorf("write %s: %w", filepath.Base(path), err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close %s: %w", filepath.Base(path), err)
	}
	return nil
}

func decodeStrict(data []byte, value any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(value); err != nil {
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
