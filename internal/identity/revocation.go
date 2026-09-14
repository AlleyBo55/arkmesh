package identity

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
)

const RevocationSchemaVersion = "arkmesh.revocations/v0alpha1"

type RevocationSet struct {
	SchemaVersion string   `json:"schema_version"`
	RevokedKeyIDs []string `json:"revoked_key_ids"`
}

func NewRevocationSet(identities []PublicIdentity) (RevocationSet, error) {
	keyIDs := make([]string, 0, len(identities))
	for _, public := range identities {
		if _, err := public.Key(); err != nil {
			return RevocationSet{}, fmt.Errorf("validate revoked identity: %w", err)
		}
		keyIDs = append(keyIDs, public.KeyID)
	}
	return newRevocationSetFromKeyIDs(keyIDs)
}

func WriteRevocationSet(path string, identities []PublicIdentity) (RevocationSet, error) {
	set, err := NewRevocationSet(identities)
	if err != nil {
		return RevocationSet{}, err
	}
	if len(set.RevokedKeyIDs) == 0 {
		return RevocationSet{}, errors.New("at least one revoked identity is required")
	}
	if err := writeJSON(path, set, 0o644); err != nil {
		return RevocationSet{}, err
	}
	return set, nil
}

func LoadRevocationSet(path string) (RevocationSet, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return RevocationSet{}, fmt.Errorf("read revocation set: %w", err)
	}
	var set RevocationSet
	if err := decodeStrict(data, &set); err != nil {
		return RevocationSet{}, fmt.Errorf("decode revocation set: %w", err)
	}
	if err := set.Validate(); err != nil {
		return RevocationSet{}, err
	}
	return set, nil
}

func MergeRevocationSets(sets ...RevocationSet) (RevocationSet, error) {
	var keyIDs []string
	for _, set := range sets {
		if err := set.Validate(); err != nil {
			return RevocationSet{}, err
		}
		keyIDs = append(keyIDs, set.RevokedKeyIDs...)
	}
	return newRevocationSetFromKeyIDs(keyIDs)
}

func (set RevocationSet) Validate() error {
	if set.SchemaVersion == "" && len(set.RevokedKeyIDs) == 0 {
		return nil
	}
	if set.SchemaVersion != RevocationSchemaVersion {
		return fmt.Errorf("unsupported revocation schema %q", set.SchemaVersion)
	}
	previous := ""
	for index, keyID := range set.RevokedKeyIDs {
		if !ValidKeyID(keyID) {
			return fmt.Errorf("invalid revoked key ID %q", keyID)
		}
		if index > 0 && keyID <= previous {
			return errors.New("revoked key IDs must be unique and sorted")
		}
		previous = keyID
	}
	return nil
}

func (set RevocationSet) Contains(keyID string) bool {
	index := sort.SearchStrings(set.RevokedKeyIDs, keyID)
	return index < len(set.RevokedKeyIDs) && set.RevokedKeyIDs[index] == keyID
}

func ValidKeyID(value string) bool {
	algorithm, digest, found := strings.Cut(value, ":")
	if !found || algorithm != Algorithm || len(digest) != 64 {
		return false
	}
	decoded, err := hex.DecodeString(digest)
	return err == nil && len(decoded) == 32 && strings.ToLower(digest) == digest
}

func newRevocationSetFromKeyIDs(keyIDs []string) (RevocationSet, error) {
	unique := make(map[string]struct{}, len(keyIDs))
	for _, keyID := range keyIDs {
		if !ValidKeyID(keyID) {
			return RevocationSet{}, fmt.Errorf("invalid revoked key ID %q", keyID)
		}
		unique[keyID] = struct{}{}
	}
	sorted := make([]string, 0, len(unique))
	for keyID := range unique {
		sorted = append(sorted, keyID)
	}
	sort.Strings(sorted)
	return RevocationSet{
		SchemaVersion: RevocationSchemaVersion,
		RevokedKeyIDs: sorted,
	}, nil
}
