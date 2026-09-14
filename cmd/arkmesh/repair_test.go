package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestChunkRepairWorkflow(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	capsuleDirectory := filepath.Join(root, "demo.ark")
	treePath := filepath.Join(root, "tree.json")
	donorPath := filepath.Join(root, "donor.bin")
	assetPath := filepath.Join(root, "model.bin")
	data := make([]byte, 5*4096+9)
	for index := range data {
		data[index] = byte(index % 89)
	}
	if err := os.WriteFile(assetPath, data, 0o644); err != nil {
		t.Fatalf("write asset: %v", err)
	}
	if err := os.WriteFile(donorPath, data, 0o644); err != nil {
		t.Fatalf("write donor: %v", err)
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

	runOK("pack", "--name", "repairable", "--out", capsuleDirectory,
		"--asset", "model="+assetPath, "--chunk-size", "4096")
	digest := chunkedAssetDigest(t, capsuleDirectory)
	exported := runOK("chunks", "tree", "--asset", digest, "--out", treePath, capsuleDirectory)
	if !strings.Contains(exported, "chunks: 6") {
		t.Fatalf("chunks tree stdout = %q, want chunk count", exported)
	}

	objectPath := filepath.Join(capsuleDirectory, "objects", digest)
	file, err := os.OpenFile(objectPath, os.O_RDWR, 0o644)
	if err != nil {
		t.Fatalf("open object: %v", err)
	}
	if _, err := file.WriteAt([]byte{0, 1, 2, 3}, 3*4096); err != nil {
		t.Fatalf("corrupt object: %v", err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("close object: %v", err)
	}

	stdout.Reset()
	stderr.Reset()
	if code := run([]string{"chunks", "scan", "--tree", treePath, capsuleDirectory}, &stdout, &stderr); code == 0 {
		t.Fatalf("chunks scan code = 0, want damage detection")
	}
	if !strings.Contains(stdout.String(), "damaged chunk 3") {
		t.Fatalf("chunks scan stdout = %q, want damaged chunk 3", stdout.String())
	}

	repaired := runOK("chunks", "repair", "--tree", treePath, "--source", donorPath, capsuleDirectory)
	if !strings.Contains(repaired, "chunks repaired: 1 of 6") {
		t.Fatalf("chunks repair stdout = %q, want single chunk repair", repaired)
	}
	runOK("verify", capsuleDirectory)

	clean := runOK("chunks", "scan", "--tree", treePath, capsuleDirectory)
	if !strings.Contains(clean, "no chunk damage") {
		t.Fatalf("chunks scan stdout = %q, want clean scan", clean)
	}

	forgedTree := filepath.Join(root, "forged.json")
	original, err := os.ReadFile(treePath)
	if err != nil {
		t.Fatalf("read tree: %v", err)
	}
	forged := bytes.Replace(original, []byte(`"leaves": [`), []byte(`"leaves": [
    "`+strings.Repeat("0", 64)+`",`), 1)
	if err := os.WriteFile(forgedTree, forged, 0o644); err != nil {
		t.Fatalf("write forged tree: %v", err)
	}
	stdout.Reset()
	stderr.Reset()
	if code := run([]string{"chunks", "scan", "--tree", forgedTree, capsuleDirectory}, &stdout, &stderr); code == 0 || !strings.Contains(stderr.String(), "want 6") {
		t.Fatalf("forged tree scan code = %d, stderr = %q, want leaf count rejection", code, stderr.String())
	}
}
