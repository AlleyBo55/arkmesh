package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestThresholdRecoveryWorkflow(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	oldIdentity := filepath.Join(root, "old")
	newIdentity := filepath.Join(root, "new")
	recoveryIdentities := []string{
		filepath.Join(root, "recovery-1"),
		filepath.Join(root, "recovery-2"),
		filepath.Join(root, "recovery-3"),
	}
	parent := filepath.Join(root, "parent.ark")
	child := filepath.Join(root, "child.ark")
	policy := filepath.Join(root, "recovery-policy.json")
	checkpoint := filepath.Join(root, "checkpoint.json")
	approvalOne := filepath.Join(root, "approval-1.json")
	approvalTwo := filepath.Join(root, "approval-2.json")
	parentAsset := filepath.Join(root, "parent.bin")
	childAsset := filepath.Join(root, "child.bin")
	if err := os.WriteFile(parentAsset, []byte("parent capability"), 0o644); err != nil {
		t.Fatalf("write parent asset: %v", err)
	}
	if err := os.WriteFile(childAsset, []byte("recovered capability"), 0o644); err != nil {
		t.Fatalf("write child asset: %v", err)
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

	for _, directory := range append([]string{oldIdentity, newIdentity}, recoveryIdentities...) {
		runOK("identity", "create", "--out", directory)
	}
	runOK(
		"recovery", "policy", "create", "--threshold", "2", "--out", policy,
		"--identity", filepath.Join(recoveryIdentities[0], "identity.json"),
		"--identity", filepath.Join(recoveryIdentities[1], "identity.json"),
		"--identity", filepath.Join(recoveryIdentities[2], "identity.json"),
	)
	runOK(
		"pack", "--name", "parent", "--out", parent,
		"--asset", "model="+parentAsset,
		"--signing-key", filepath.Join(oldIdentity, "identity.key"),
	)
	runOK(
		"checkpoint", "create", "--out", checkpoint,
		"--trust", filepath.Join(oldIdentity, "identity.json"),
		parent,
	)
	pending := runOK(
		"pack", "--name", "recovered-child", "--out", child,
		"--asset", "model="+childAsset,
		"--signing-key", filepath.Join(newIdentity, "identity.key"),
		"--parent", parent,
		"--recovery-policy", policy,
	)
	if !strings.Contains(pending, "recovery pending: sha256:") {
		t.Fatalf("recovery pack stdout = %q, want pending policy", pending)
	}

	for index, output := range []string{approvalOne, approvalTwo} {
		runOK(
			"recovery", "approve", "--policy", policy,
			"--signing-key", filepath.Join(recoveryIdentities[index], "identity.key"),
			"--parent", parent, "--out", output,
			child,
		)
	}
	assembled := runOK(
		"recovery", "assemble", "--policy", policy, "--parent", parent,
		"--approval", approvalOne, "--approval", approvalTwo,
		child,
	)
	if !strings.Contains(assembled, "approvals: 2") {
		t.Fatalf("recovery assemble stdout = %q, want two approvals", assembled)
	}
	advanced := runOK(
		"checkpoint", "advance", "--checkpoint", checkpoint,
		"--parent", parent, "--recovery-policy", policy,
		child,
	)
	if !strings.Contains(advanced, "lineage: verified_threshold_recovery") {
		t.Fatalf("checkpoint advance stdout = %q, want threshold recovery", advanced)
	}
	verified := runOK("verify", "--checkpoint", checkpoint, "--require-trusted", child)
	if !strings.Contains(verified, "checkpoint: matched") {
		t.Fatalf("recovered child verify stdout = %q, want checkpoint match", verified)
	}
}
