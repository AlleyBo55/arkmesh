package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestErasureRecoveryWorkflow(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	capsuleDirectory := filepath.Join(root, "demo.ark")
	shardDirectory := filepath.Join(root, "shards")
	planPath := filepath.Join(root, "plan.json")
	assetPath := filepath.Join(root, "model.bin")
	data := make([]byte, 200000)
	for index := range data {
		data[index] = byte(index % 251)
	}
	if err := os.WriteFile(assetPath, data, 0o644); err != nil {
		t.Fatalf("write asset: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	runOK := func(args ...string) string {
		t.Helper()
		stdout.Reset()
		stderr.Reset()
		if code := run(args, &stdout, &stderr); code != 0 {
			t.Fatalf("run(%q) code = %d, stderr = %q", args, code, stderr.String())
		}
		return stdout.String()
	}

	runOK("identity", "create", "--out", filepath.Join(root, "author"))
	runOK("pack", "--name", "protected", "--out", capsuleDirectory,
		"--asset", "model="+assetPath,
		"--signing-key", filepath.Join(root, "author", "identity.key"))
	digest := chunkedAssetDigest(t, capsuleDirectory)

	protected := runOK("erasure", "protect", "--asset", digest,
		"--data", "4", "--parity", "2",
		"--shards", shardDirectory, "--plan", planPath, capsuleDirectory)
	if !strings.Contains(protected, "4 data + 2 parity") {
		t.Fatalf("erasure protect stdout = %q, want shard layout", protected)
	}

	// Lose the object entirely plus two shards: no donor holds the missing bytes.
	if err := os.Remove(filepath.Join(capsuleDirectory, "objects", digest)); err != nil {
		t.Fatalf("remove object: %v", err)
	}
	for _, index := range []string{"shard-001", "shard-004"} {
		if err := os.Remove(filepath.Join(shardDirectory, index)); err != nil {
			t.Fatalf("remove %s: %v", index, err)
		}
	}

	recovered := runOK("erasure", "recover", "--plan", planPath, "--shards", shardDirectory, capsuleDirectory)
	if !strings.Contains(recovered, "unusable shards: 2") {
		t.Fatalf("erasure recover stdout = %q, want two unusable shards", recovered)
	}
	verified := runOK("verify", "--trust", filepath.Join(root, "author", "identity.json"), "--require-trusted", capsuleDirectory)
	if !strings.Contains(verified, "signature: valid_trusted_author") {
		t.Fatalf("verify stdout = %q, want trusted signature after recovery", verified)
	}

	rebuilt, err := os.ReadFile(filepath.Join(capsuleDirectory, "objects", digest))
	if err != nil {
		t.Fatalf("read recovered object: %v", err)
	}
	if !bytes.Equal(rebuilt, data) {
		t.Fatal("recovered object does not match the original bytes")
	}

	// A third loss crosses the parity threshold and must be refused.
	if err := os.Remove(filepath.Join(shardDirectory, "shard-002")); err != nil {
		t.Fatalf("remove shard-002: %v", err)
	}
	if err := os.Remove(filepath.Join(capsuleDirectory, "objects", digest)); err != nil {
		t.Fatalf("remove object again: %v", err)
	}
	stdout.Reset()
	stderr.Reset()
	if code := run([]string{"erasure", "recover", "--plan", planPath, "--shards", shardDirectory, capsuleDirectory}, &stdout, &stderr); code == 0 || !strings.Contains(stderr.String(), "required") {
		t.Fatalf("recover code = %d, stderr = %q, want threshold refusal", code, stderr.String())
	}
}
