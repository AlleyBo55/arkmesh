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
