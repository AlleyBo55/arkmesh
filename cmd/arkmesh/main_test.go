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
