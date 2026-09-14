package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestChunkProofWorkflow(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	capsuleDirectory := filepath.Join(root, "demo.ark")
	verifierDirectory := filepath.Join(root, "verifier.ark")
	proofPath := filepath.Join(root, "proof.json")
	assetPath := filepath.Join(root, "model.bin")
	data := make([]byte, 3*4096+11)
	for index := range data {
		data[index] = byte(index % 97)
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

	runOK("pack", "--name", "chunked", "--out", capsuleDirectory,
		"--asset", "model="+assetPath, "--chunk-size", "4096")

	inspected := runOK("chunks", "inspect", capsuleDirectory)
	if !strings.Contains(inspected, "chunks=4") || !strings.Contains(inspected, "chunk_size=4096") {
		t.Fatalf("chunks inspect stdout = %q, want chunk metadata", inspected)
	}

	digest := chunkedAssetDigest(t, capsuleDirectory)
	proved := runOK("chunks", "prove", "--asset", digest, "--index", "2", "--out", proofPath, capsuleDirectory)
	if !strings.Contains(proved, "chunk: 2 of 4") {
		t.Fatalf("chunks prove stdout = %q, want chunk position", proved)
	}

	// A verifier that holds only the signed manifest, and none of the objects,
	// must still be able to check the proof.
	if err := os.MkdirAll(verifierDirectory, 0o755); err != nil {
		t.Fatalf("create verifier directory: %v", err)
	}
	manifest, err := os.ReadFile(filepath.Join(capsuleDirectory, "manifest.json"))
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	if err := os.WriteFile(filepath.Join(verifierDirectory, "manifest.json"), manifest, 0o644); err != nil {
		t.Fatalf("write verifier manifest: %v", err)
	}
	checked := runOK("chunks", "check", "--proof", proofPath, verifierDirectory)
	if !strings.Contains(checked, "chunk proof verified") {
		t.Fatalf("chunks check stdout = %q, want verification", checked)
	}

	tamperedProof := filepath.Join(root, "tampered.json")
	original, err := os.ReadFile(proofPath)
	if err != nil {
		t.Fatalf("read proof: %v", err)
	}
	if err := os.WriteFile(tamperedProof, bytes.Replace(original, []byte(`"chunk_index": 2`), []byte(`"chunk_index": 1`), 1), 0o644); err != nil {
		t.Fatalf("write tampered proof: %v", err)
	}
	stdout.Reset()
	stderr.Reset()
	if code := run([]string{"chunks", "check", "--proof", tamperedProof, verifierDirectory}, &stdout, &stderr); code == 0 || !strings.Contains(stderr.String(), "does not reconstruct") {
		t.Fatalf("tampered check code = %d, stderr = %q, want reconstruction failure", code, stderr.String())
	}
}

func chunkedAssetDigest(t *testing.T, capsuleDirectory string) string {
	t.Helper()
	entries, err := os.ReadDir(filepath.Join(capsuleDirectory, "objects"))
	if err != nil {
		t.Fatalf("read objects: %v", err)
	}
	for _, entry := range entries {
		if len(entry.Name()) == 64 {
			return entry.Name()
		}
	}
	t.Fatal("no stored object found")
	return ""
}
