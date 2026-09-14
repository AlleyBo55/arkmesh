package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSignedCapsuleWorkflow(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	identityDirectory := filepath.Join(root, "author")
	capsuleDirectory := filepath.Join(root, "demo.ark")
	modelPath := filepath.Join(root, "model.gguf")
	if err := os.WriteFile(modelPath, []byte("small model fixture"), 0o644); err != nil {
		t.Fatalf("write model fixture: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if code := run([]string{"identity", "create", "--out", identityDirectory}, &stdout, &stderr); code != 0 {
		t.Fatalf("identity create code = %d, stderr = %q", code, stderr.String())
	}
	privatePath := filepath.Join(identityDirectory, "identity.key")
	publicPath := filepath.Join(identityDirectory, "identity.json")

	stdout.Reset()
	stderr.Reset()
	if code := run([]string{
		"pack",
		"--name", "signed-demo",
		"--out", capsuleDirectory,
		"--asset", "model=" + modelPath,
		"--signing-key", privatePath,
	}, &stdout, &stderr); code != 0 {
		t.Fatalf("pack code = %d, stderr = %q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "signed by: ed25519:") {
		t.Fatalf("pack stdout = %q, want signer", stdout.String())
	}

	stdout.Reset()
	stderr.Reset()
	if code := run([]string{
		"verify",
		"--trust", publicPath,
		"--require-trusted",
		capsuleDirectory,
	}, &stdout, &stderr); code != 0 {
		t.Fatalf("verify code = %d, stderr = %q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "signature: valid_trusted_author") {
		t.Fatalf("verify stdout = %q, want trusted signature", stdout.String())
	}
}

func TestLineageWorkflow(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	identityDirectory := filepath.Join(root, "author")
	parentDirectory := filepath.Join(root, "parent.ark")
	childDirectory := filepath.Join(root, "child.ark")
	parentAsset := filepath.Join(root, "parent.gguf")
	childAsset := filepath.Join(root, "child.gguf")
	if err := os.WriteFile(parentAsset, []byte("parent model"), 0o644); err != nil {
		t.Fatalf("write parent fixture: %v", err)
	}
	if err := os.WriteFile(childAsset, []byte("child model"), 0o644); err != nil {
		t.Fatalf("write child fixture: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if code := run([]string{"identity", "create", "--out", identityDirectory}, &stdout, &stderr); code != 0 {
		t.Fatalf("identity create code = %d, stderr = %q", code, stderr.String())
	}
	privatePath := filepath.Join(identityDirectory, "identity.key")
	publicPath := filepath.Join(identityDirectory, "identity.json")

	stdout.Reset()
	stderr.Reset()
	if code := run([]string{
		"pack",
		"--name", "parent",
		"--out", parentDirectory,
		"--asset", "model=" + parentAsset,
		"--signing-key", privatePath,
	}, &stdout, &stderr); code != 0 {
		t.Fatalf("parent pack code = %d, stderr = %q", code, stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	if code := run([]string{
		"pack",
		"--name", "child",
		"--out", childDirectory,
		"--asset", "model=" + childAsset,
		"--signing-key", privatePath,
		"--parent", parentDirectory,
	}, &stdout, &stderr); code != 0 {
		t.Fatalf("child pack code = %d, stderr = %q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "parent: sha256:") {
		t.Fatalf("child pack stdout = %q, want parent ID", stdout.String())
	}

	stdout.Reset()
	stderr.Reset()
	if code := run([]string{
		"verify",
		"--trust", publicPath,
		"--require-trusted",
		"--parent", parentDirectory,
		"--require-lineage",
		childDirectory,
	}, &stdout, &stderr); code != 0 {
		t.Fatalf("lineage verify code = %d, stderr = %q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "lineage: verified_same_author") {
		t.Fatalf("lineage verify stdout = %q, want verified lineage", stdout.String())
	}
}

func TestKeyRotationAndRevocationWorkflow(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	oldIdentity := filepath.Join(root, "old-author")
	newIdentity := filepath.Join(root, "new-author")
	parentDirectory := filepath.Join(root, "parent.ark")
	childDirectory := filepath.Join(root, "child.ark")
	parentAsset := filepath.Join(root, "parent.gguf")
	childAsset := filepath.Join(root, "child.gguf")
	revocationsPath := filepath.Join(root, "revocations.json")
	if err := os.WriteFile(parentAsset, []byte("parent model"), 0o644); err != nil {
		t.Fatalf("write parent fixture: %v", err)
	}
	if err := os.WriteFile(childAsset, []byte("rotated child model"), 0o644); err != nil {
		t.Fatalf("write child fixture: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	for _, directory := range []string{oldIdentity, newIdentity} {
		stdout.Reset()
		stderr.Reset()
		if code := run([]string{"identity", "create", "--out", directory}, &stdout, &stderr); code != 0 {
			t.Fatalf("identity create %q code = %d, stderr = %q", directory, code, stderr.String())
		}
	}
	oldPrivate := filepath.Join(oldIdentity, "identity.key")
	oldPublic := filepath.Join(oldIdentity, "identity.json")
	newPrivate := filepath.Join(newIdentity, "identity.key")
	newPublic := filepath.Join(newIdentity, "identity.json")

	stdout.Reset()
	stderr.Reset()
	if code := run([]string{
		"pack", "--name", "parent", "--out", parentDirectory,
		"--asset", "model=" + parentAsset,
		"--signing-key", oldPrivate,
	}, &stdout, &stderr); code != 0 {
		t.Fatalf("parent pack code = %d, stderr = %q", code, stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	if code := run([]string{
		"pack", "--name", "rotated-child", "--out", childDirectory,
		"--asset", "model=" + childAsset,
		"--signing-key", newPrivate,
		"--parent", parentDirectory,
		"--rotation-key", oldPrivate,
	}, &stdout, &stderr); code != 0 {
		t.Fatalf("rotated child pack code = %d, stderr = %q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "key rotation: ed25519:") {
		t.Fatalf("rotated child stdout = %q, want key rotation", stdout.String())
	}

	stdout.Reset()
	stderr.Reset()
	if code := run([]string{
		"verify", "--trust", oldPublic, "--require-trusted",
		"--parent", parentDirectory, "--require-lineage", childDirectory,
	}, &stdout, &stderr); code != 0 {
		t.Fatalf("rotated verify code = %d, stderr = %q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "lineage: verified_key_rotation") {
		t.Fatalf("rotated verify stdout = %q, want key rotation lineage", stdout.String())
	}

	stdout.Reset()
	stderr.Reset()
	if code := run([]string{
		"identity", "revoke", "--identity", newPublic, "--out", revocationsPath,
	}, &stdout, &stderr); code != 0 {
		t.Fatalf("identity revoke code = %d, stderr = %q", code, stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	if code := run([]string{
		"verify", "--trust", oldPublic, "--require-trusted",
		"--revocations", revocationsPath,
		"--parent", parentDirectory, "--require-lineage", childDirectory,
	}, &stdout, &stderr); code == 0 || !strings.Contains(stderr.String(), "revoked by local policy") {
		t.Fatalf("revoked verify code = %d, stderr = %q, want revocation rejection", code, stderr.String())
	}
}
