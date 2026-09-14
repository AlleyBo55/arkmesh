package capsule

import (
	"errors"
	"fmt"
)

type LineageStatus string

const (
	LineageRoot             LineageStatus = "root"
	LineageParentNotChecked LineageStatus = "parent_not_checked"
	LineageVerified         LineageStatus = "verified_same_author"
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
	if child.ParentCapsuleID == "" {
		return Lineage{}, errors.New("child capsule does not declare a parent")
	}
	if child.ParentCapsuleID != parent.CapsuleID {
		return Lineage{}, fmt.Errorf("declared parent %s does not match supplied capsule %s", child.ParentCapsuleID, parent.CapsuleID)
	}
	if !isSignedAuthenticity(childAuthenticity) {
		return Lineage{}, errors.New("child capsule must have a valid signature")
	}
	if !isSignedAuthenticity(parentAuthenticity) {
		return Lineage{}, errors.New("parent capsule must have a valid signature")
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

func isSignedAuthenticity(authenticity Authenticity) bool {
	return authenticity.SignerID != "" && (authenticity.Status == AuthenticityUnknown || authenticity.Status == AuthenticityTrusted)
}
