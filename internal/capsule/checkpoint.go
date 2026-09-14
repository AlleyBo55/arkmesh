package capsule

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"arkmesh/internal/identity"
)

const CheckpointSchemaVersion = "arkmesh.checkpoint/v0alpha1"

type Checkpoint struct {
	SchemaVersion string `json:"schema_version"`
	CapsuleID     string `json:"capsule_id"`
	SignerID      string `json:"signer_id"`
}

func NewCheckpoint(manifest Manifest, authenticity Authenticity) (Checkpoint, error) {
	if !validCapsuleID(manifest.CapsuleID) {
		return Checkpoint{}, fmt.Errorf("invalid checkpoint capsule ID %q", manifest.CapsuleID)
	}
	if !isSignedAuthenticity(authenticity) {
		return Checkpoint{}, errors.New("checkpoint capsule must have a valid signature")
	}
	if _, err := authenticity.Signer.Key(); err != nil {
		return Checkpoint{}, fmt.Errorf("validate checkpoint signer: %w", err)
	}
	if authenticity.SignerID != authenticity.Signer.KeyID {
		return Checkpoint{}, errors.New("checkpoint signer ID does not match verified signer")
	}
	return Checkpoint{
		SchemaVersion: CheckpointSchemaVersion,
		CapsuleID:     manifest.CapsuleID,
		SignerID:      authenticity.SignerID,
	}, nil
}

func WriteCheckpoint(path string, checkpoint Checkpoint) error {
	if err := checkpoint.Validate(); err != nil {
		return err
	}
	data, err := checkpointJSON(checkpoint)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return fmt.Errorf("create checkpoint: %w", err)
	}
	completed := false
	defer func() {
		if !completed {
			_ = os.Remove(path)
		}
	}()
	if err := writeAndSync(file, data); err != nil {
		_ = file.Close()
		return fmt.Errorf("write checkpoint: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close checkpoint: %w", err)
	}
	completed = true
	return nil
}

func LoadCheckpoint(path string) (Checkpoint, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return Checkpoint{}, fmt.Errorf("inspect checkpoint: %w", err)
	}
	if !info.Mode().IsRegular() {
		return Checkpoint{}, errors.New("checkpoint is not a regular file")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Checkpoint{}, fmt.Errorf("read checkpoint: %w", err)
	}
	var checkpoint Checkpoint
	if err := decodeCheckpoint(data, &checkpoint); err != nil {
		return Checkpoint{}, fmt.Errorf("decode checkpoint: %w", err)
	}
	if err := checkpoint.Validate(); err != nil {
		return Checkpoint{}, err
	}
	return checkpoint, nil
}

func VerifyCheckpoint(checkpoint Checkpoint, manifest Manifest, authenticity Authenticity) error {
	if err := checkpoint.Validate(); err != nil {
		return err
	}
	if !isSignedAuthenticity(authenticity) {
		return errors.New("checkpointed capsule must have a valid signature")
	}
	if manifest.CapsuleID != checkpoint.CapsuleID {
		return fmt.Errorf("capsule %s does not match checkpoint head %s", manifest.CapsuleID, checkpoint.CapsuleID)
	}
	if authenticity.SignerID != checkpoint.SignerID {
		return fmt.Errorf("capsule signer %s does not match checkpoint signer %s", authenticity.SignerID, checkpoint.SignerID)
	}
	return nil
}

func AdvanceCheckpointFile(path, childRoot string, parent Manifest, parentAuthenticity Authenticity, child Manifest, childAuthenticity Authenticity) (Checkpoint, Lineage, error) {
	current, err := LoadCheckpoint(path)
	if err != nil {
		return Checkpoint{}, Lineage{}, err
	}
	if err := VerifyCheckpoint(current, parent, parentAuthenticity); err != nil {
		return Checkpoint{}, Lineage{}, fmt.Errorf("verify checkpoint parent: %w", err)
	}
	lineage, err := CheckLineageWithAuthority(childRoot, child, childAuthenticity, parent, parentAuthenticity)
	if err != nil {
		return Checkpoint{}, Lineage{}, fmt.Errorf("verify checkpoint advancement: %w", err)
	}
	next, err := NewCheckpoint(child, childAuthenticity)
	if err != nil {
		return Checkpoint{}, Lineage{}, err
	}
	if err := replaceCheckpoint(path, current, next); err != nil {
		return Checkpoint{}, Lineage{}, err
	}
	return next, lineage, nil
}

func AdvanceCheckpointFileWithRecovery(path, childRoot string, parent Manifest, parentAuthenticity Authenticity, child Manifest, childAuthenticity Authenticity, policy RecoveryPolicy, revocations identity.RevocationSet) (Checkpoint, Lineage, error) {
	current, err := LoadCheckpoint(path)
	if err != nil {
		return Checkpoint{}, Lineage{}, err
	}
	if err := VerifyCheckpoint(current, parent, parentAuthenticity); err != nil {
		return Checkpoint{}, Lineage{}, fmt.Errorf("verify checkpoint parent: %w", err)
	}
	lineage, err := CheckLineageWithRecovery(childRoot, child, childAuthenticity, parent, parentAuthenticity, policy, revocations)
	if err != nil {
		return Checkpoint{}, Lineage{}, fmt.Errorf("verify checkpoint recovery: %w", err)
	}
	next, err := NewCheckpoint(child, childAuthenticity)
	if err != nil {
		return Checkpoint{}, Lineage{}, err
	}
	if err := replaceCheckpoint(path, current, next); err != nil {
		return Checkpoint{}, Lineage{}, err
	}
	return next, lineage, nil
}

func (checkpoint Checkpoint) Validate() error {
	if checkpoint.SchemaVersion != CheckpointSchemaVersion {
		return fmt.Errorf("unsupported checkpoint schema %q", checkpoint.SchemaVersion)
	}
	if !validCapsuleID(checkpoint.CapsuleID) {
		return fmt.Errorf("invalid checkpoint capsule ID %q", checkpoint.CapsuleID)
	}
	if !identity.ValidKeyID(checkpoint.SignerID) {
		return fmt.Errorf("invalid checkpoint signer ID %q", checkpoint.SignerID)
	}
	return nil
}

func replaceCheckpoint(path string, expected, next Checkpoint) error {
	lockPath := path + ".lock"
	lock, err := os.OpenFile(lockPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return errors.New("checkpoint is locked by another advancement")
		}
		return fmt.Errorf("create checkpoint lock: %w", err)
	}
	if err := lock.Close(); err != nil {
		_ = os.Remove(lockPath)
		return fmt.Errorf("close checkpoint lock: %w", err)
	}
	defer os.Remove(lockPath)

	data, err := checkpointJSON(next)
	if err != nil {
		return err
	}
	directory := filepath.Dir(path)
	temporary, err := os.CreateTemp(directory, ".arkmesh-checkpoint-*")
	if err != nil {
		return fmt.Errorf("create temporary checkpoint: %w", err)
	}
	temporaryPath := temporary.Name()
	removeTemporary := true
	defer func() {
		if removeTemporary {
			_ = os.Remove(temporaryPath)
		}
	}()
	if err := temporary.Chmod(0o644); err != nil {
		_ = temporary.Close()
		return fmt.Errorf("set checkpoint permissions: %w", err)
	}
	if err := writeAndSync(temporary, data); err != nil {
		_ = temporary.Close()
		return fmt.Errorf("write temporary checkpoint: %w", err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close temporary checkpoint: %w", err)
	}
	current, err := LoadCheckpoint(path)
	if err != nil {
		return err
	}
	if current != expected {
		return errors.New("checkpoint changed while advancement was being prepared")
	}
	if err := os.Rename(temporaryPath, path); err != nil {
		return fmt.Errorf("replace checkpoint: %w", err)
	}
	removeTemporary = false
	return nil
}

func checkpointJSON(checkpoint Checkpoint) ([]byte, error) {
	data, err := json.MarshalIndent(checkpoint, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("encode checkpoint: %w", err)
	}
	return append(data, '\n'), nil
}

func decodeCheckpoint(data []byte, checkpoint *Checkpoint) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(checkpoint); err != nil {
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

func writeAndSync(file *os.File, data []byte) error {
	if _, err := file.Write(data); err != nil {
		return err
	}
	return file.Sync()
}
