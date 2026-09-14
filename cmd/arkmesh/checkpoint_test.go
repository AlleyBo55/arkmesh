package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckpointRollbackWorkflow(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	oldIdentity := filepath.Join(root, "old-author")
	newIdentity := filepath.Join(root, "new-author")
	parentDirectory := filepath.Join(root, "parent.ark")
	childDirectory := filepath.Join(root, "child.ark")
	checkpointPath := filepath.Join(root, "checkpoint.json")
	revocationsPath := filepath.Join(root, "revocations.json")
	parentAsset := filepath.Join(root, "parent.gguf")
	childAsset := filepath.Join(root, "child.gguf")
	if err := os.WriteFile(parentAsset, []byte("checkpoint parent"), 0o644); err != nil {
		t.Fatalf("write parent fixture: %v", err)
	}
	if err := os.WriteFile(childAsset, []byte("checkpoint child"), 0o644); err != nil {
		t.Fatalf("write child fixture: %v", err)
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

	runOK("identity", "create", "--out", oldIdentity)
	runOK("identity", "create", "--out", newIdentity)
	oldPrivate := filepath.Join(oldIdentity, "identity.key")
	oldPublic := filepath.Join(oldIdentity, "identity.json")
	newPrivate := filepath.Join(newIdentity, "identity.key")
	newPublic := filepath.Join(newIdentity, "identity.json")

	runOK(
		"pack", "--name", "checkpoint-parent", "--out", parentDirectory,
		"--asset", "model="+parentAsset,
		"--signing-key", oldPrivate,
	)
	runOK(
		"pack", "--name", "checkpoint-child", "--out", childDirectory,
		"--asset", "model="+childAsset,
		"--signing-key", newPrivate,
		"--parent", parentDirectory,
		"--rotation-key", oldPrivate,
	)

	created := runOK(
		"checkpoint", "create", "--out", checkpointPath,
		"--trust", oldPublic,
		parentDirectory,
	)
	if !strings.Contains(created, "created local checkpoint") {
		t.Fatalf("checkpoint create stdout = %q", created)
	}

	matchedParent := runOK(
		"verify", "--checkpoint", checkpointPath, "--require-trusted",
		parentDirectory,
	)
	if !strings.Contains(matchedParent, "checkpoint: matched") {
		t.Fatalf("parent verify stdout = %q, want checkpoint match", matchedParent)
	}

	advanced := runOK(
		"checkpoint", "advance", "--checkpoint", checkpointPath,
		"--parent", parentDirectory,
		childDirectory,
	)
	if !strings.Contains(advanced, "lineage: verified_key_rotation") {
		t.Fatalf("checkpoint advance stdout = %q, want rotation lineage", advanced)
	}

	matchedChild := runOK(
		"verify", "--checkpoint", checkpointPath, "--require-trusted",
		childDirectory,
	)
	if !strings.Contains(matchedChild, "checkpoint: matched") {
		t.Fatalf("child verify stdout = %q, want checkpoint match", matchedChild)
	}

	stdout.Reset()
	stderr.Reset()
	if code := run([]string{
		"verify", "--checkpoint", checkpointPath,
		parentDirectory,
	}, &stdout, &stderr); code == 0 || !strings.Contains(stderr.String(), "does not match checkpoint head") {
		t.Fatalf("rollback verify code = %d, stderr = %q, want checkpoint rejection", code, stderr.String())
	}

	runOK(
		"identity", "revoke", "--identity", newPublic,
		"--out", revocationsPath,
	)
	stdout.Reset()
	stderr.Reset()
	if code := run([]string{
		"verify", "--checkpoint", checkpointPath,
		"--revocations", revocationsPath,
		childDirectory,
	}, &stdout, &stderr); code == 0 || !strings.Contains(stderr.String(), "revoked by local policy") {
		t.Fatalf("revoked checkpoint verify code = %d, stderr = %q, want revocation rejection", code, stderr.String())
	}
}
