package capsule

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type LineageStatus string

const (
	LineageRoot             LineageStatus = "root"
	LineageParentNotChecked LineageStatus = "parent_not_checked"
	LineageVerified         LineageStatus = "verified_same_author"
	LineageKeyRotation      LineageStatus = "verified_key_rotation"
)

type Lineage struct {
	Status          LineageStatus
	ParentCapsuleID string
	SignerID        string
}

func DescribeLineage(manifest Manifest) Lineage {
	if manifest.ParentCapsuleID == "" {
		return Lineage{Status: LineageRoot}
	}
	return Lineage{
		Status:          LineageParentNotChecked,
		ParentCapsuleID: manifest.ParentCapsuleID,
	}
}

func CheckLineage(child Manifest, childAuthenticity Authenticity, parent Manifest, parentAuthenticity Authenticity) (Lineage, error) {
	if err := validateLineageEndpoints(child, childAuthenticity, parent, parentAuthenticity); err != nil {
		return Lineage{}, err
	}
	if childAuthenticity.SignerID != parentAuthenticity.SignerID {
		return Lineage{}, fmt.Errorf("child signer %s is not authorized by parent signer %s", childAuthenticity.SignerID, parentAuthenticity.SignerID)
	}
	return Lineage{
		Status:          LineageVerified,
		ParentCapsuleID: parent.CapsuleID,
		SignerID:        childAuthenticity.SignerID,
	}, nil
}

func CheckLineageWithAuthority(childRoot string, child Manifest, childAuthenticity Authenticity, parent Manifest, parentAuthenticity Authenticity) (Lineage, error) {
	if err := validateLineageEndpoints(child, childAuthenticity, parent, parentAuthenticity); err != nil {
		return Lineage{}, err
	}
	if childAuthenticity.SignerID == parentAuthenticity.SignerID {
		if _, err := os.Lstat(filepath.Join(childRoot, AuthorityFileName)); err == nil {
			return Lineage{}, errors.New("same-author child must not include a key rotation authority")
		} else if !errors.Is(err, os.ErrNotExist) {
			return Lineage{}, fmt.Errorf("inspect key rotation authority: %w", err)
		}
		return CheckLineage(child, childAuthenticity, parent, parentAuthenticity)
	}
	return verifyKeyRotation(childRoot, parent, parentAuthenticity, child, childAuthenticity)
}

func validateLineageEndpoints(child Manifest, childAuthenticity Authenticity, parent Manifest, parentAuthenticity Authenticity) error {
	if child.ParentCapsuleID == "" {
		return errors.New("child capsule does not declare a parent")
	}
	if child.ParentCapsuleID != parent.CapsuleID {
		return fmt.Errorf("declared parent %s does not match supplied capsule %s", child.ParentCapsuleID, parent.CapsuleID)
	}
	if !isSignedAuthenticity(childAuthenticity) {
		return errors.New("child capsule must have a valid signature")
	}
	if !isSignedAuthenticity(parentAuthenticity) {
		return errors.New("parent capsule must have a valid signature")
	}
	return nil
}

func isSignedAuthenticity(authenticity Authenticity) bool {
	return authenticity.SignerID != "" && (authenticity.Status == AuthenticityUnknown || authenticity.Status == AuthenticityTrusted)
}
