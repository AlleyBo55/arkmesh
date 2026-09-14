package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAuditSamplingWorkflow(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	capsuleDirectory := filepath.Join(root, "demo.ark")
	treePath := filepath.Join(root, "tree.json")
	logPath := filepath.Join(root, "retention.log")
	assetPath := filepath.Join(root, "model.bin")
	data := make([]byte, 64*4096)
	for index := range data {
		data[index] = byte(index % 241)
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

	runOK("pack", "--name", "audited", "--out", capsuleDirectory,
		"--asset", "model="+assetPath, "--chunk-size", "4096")
	digest := chunkedAssetDigest(t, capsuleDirectory)
	runOK("chunks", "tree", "--asset", digest, "--out", treePath, capsuleDirectory)

	sampled := runOK("audit", "sample", "--tree", treePath,
		"--tolerance", "0.05", "--confidence", "0.99", "--log", logPath, capsuleDirectory)
	if !strings.Contains(sampled, "failed: 0") {
		t.Fatalf("audit stdout = %q, want a clean audit", sampled)
	}
	if !strings.Contains(sampled, "with 99% confidence") || !strings.Contains(sampled, "at most") {
		t.Fatalf("audit stdout = %q, want an explicit confidence bound", sampled)
	}
	if strings.Contains(sampled, "at least 100.00%") {
		t.Fatalf("audit stdout = %q, must not claim full integrity from a sample", sampled)
	}

	exhaustive := runOK("audit", "sample", "--tree", treePath, "--samples", "64", capsuleDirectory)
	if !strings.Contains(exhaustive, "every chunk was checked") {
		t.Fatalf("exhaustive audit stdout = %q, want full check statement", exhaustive)
	}

	object, err := os.OpenFile(filepath.Join(capsuleDirectory, "objects", digest), os.O_RDWR, 0o644)
	if err != nil {
		t.Fatalf("open object: %v", err)
	}
	if _, err := object.WriteAt([]byte{9, 9, 9}, 7*4096); err != nil {
		t.Fatalf("corrupt object: %v", err)
	}
	if err := object.Close(); err != nil {
		t.Fatalf("close object: %v", err)
	}

	stdout.Reset()
	stderr.Reset()
	if code := run([]string{"audit", "sample", "--tree", treePath, "--samples", "64", "--log", logPath, capsuleDirectory}, &stdout, &stderr); code == 0 {
		t.Fatalf("damaged audit code = 0, want failure")
	}
	if !strings.Contains(stdout.String(), "failed chunk 7") {
		t.Fatalf("damaged audit stdout = %q, want failed chunk 7", stdout.String())
	}
	if !strings.Contains(stdout.String(), "no retrievability bound is claimed") {
		t.Fatalf("damaged audit stdout = %q, want withheld bound", stdout.String())
	}

	history := runOK("audit", "history", "--log", logPath)
	if !strings.Contains(history, "audits: 2") || !strings.Contains(history, "clean: 1") || !strings.Contains(history, "damaged: 1") {
		t.Fatalf("history stdout = %q, want two recorded audits", history)
	}
	if !strings.Contains(history, "cannot prove that checks were not omitted") {
		t.Fatalf("history stdout = %q, want the honesty note", history)
	}
}
