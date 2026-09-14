package capsule

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func auditFixture(t *testing.T, chunks int) (Manifest, Asset, ChunkTreeFile, string) {
	t.Helper()
	root := t.TempDir()
	assetPath := writeLargeFixture(t, root, "model.gguf", chunks*int(MinChunkSize))
	output := filepath.Join(root, "demo.ark")
	manifest, err := PackWithOptions("demo", output, []AssetSource{{Role: "model", Path: assetPath}}, time.Now(), PackOptions{ChunkSize: MinChunkSize})
	if err != nil {
		t.Fatalf("PackWithOptions() error = %v", err)
	}
	asset := manifest.Assets[0]
	objectPath := filepath.Join(output, "objects", asset.SHA256)
	tree, err := BuildChunkTreeFile(manifest.CapsuleID, asset, objectPath)
	if err != nil {
		t.Fatalf("BuildChunkTreeFile() error = %v", err)
	}
	return manifest, asset, tree, objectPath
}

func TestAuditPassesAndBoundsHiddenDamage(t *testing.T) {
	t.Parallel()
	_, asset, tree, objectPath := auditFixture(t, 40)
	result, err := AuditAsset(objectPath, asset, tree, AuditPolicy{Tolerance: 0.25, Confidence: 0.95}, nil)
	if err != nil {
		t.Fatalf("AuditAsset() error = %v", err)
	}
	if len(result.Failed) != 0 {
		t.Fatalf("failed = %v, want none", result.Failed)
	}
	if result.Sampled <= 0 || result.Sampled > result.ChunkCount {
		t.Fatalf("sampled = %d, want within 1..%d", result.Sampled, result.ChunkCount)
	}
	// A partial audit must never claim the object is fully intact.
	if result.MinIntactFraction >= 1 && !result.Exhaustive {
		t.Fatalf("partial audit claimed full integrity: %+v", result)
	}
	if result.MaxBadChunks < 0 || result.MaxBadChunks >= result.ChunkCount {
		t.Fatalf("maxBadChunks = %d, want a useful bound", result.MaxBadChunks)
	}
}

func TestExhaustiveAuditProvesNoDamage(t *testing.T) {
	t.Parallel()
	_, asset, tree, objectPath := auditFixture(t, 8)
	result, err := AuditAsset(objectPath, asset, tree, AuditPolicy{Samples: 8, Confidence: 0.99}, nil)
	if err != nil {
		t.Fatalf("AuditAsset() error = %v", err)
	}
	if !result.Exhaustive || result.MaxBadChunks != 0 || result.MinIntactFraction != 1 {
		t.Fatalf("result = %+v, want exhaustive proof", result)
	}
}

func TestAuditReportsObservedDamageWithoutClaimingABound(t *testing.T) {
	t.Parallel()
	_, asset, tree, objectPath := auditFixture(t, 6)
	corruptChunk(t, objectPath, 2*MinChunkSize)
	result, err := AuditAsset(objectPath, asset, tree, AuditPolicy{Samples: 6, Confidence: 0.99}, nil)
	if err != nil {
		t.Fatalf("AuditAsset() error = %v", err)
	}
	if len(result.Failed) != 1 || result.Failed[0] != 2 {
		t.Fatalf("failed = %v, want chunk 2", result.Failed)
	}
	if result.MinIntactFraction != 0 || result.MaxBadChunks != result.ChunkCount {
		t.Fatalf("result = %+v, want no integrity claim after observed damage", result)
	}
}

func TestAuditDetectsTruncatedObject(t *testing.T) {
	t.Parallel()
	_, asset, tree, objectPath := auditFixture(t, 5)
	if err := os.Truncate(objectPath, MinChunkSize); err != nil {
		t.Fatalf("truncate object: %v", err)
	}
	result, err := AuditAsset(objectPath, asset, tree, AuditPolicy{Samples: 5, Confidence: 0.99}, nil)
	if err != nil {
		t.Fatalf("AuditAsset() error = %v", err)
	}
	if len(result.Failed) != 4 {
		t.Fatalf("failed = %v, want four missing chunks", result.Failed)
	}
}

func TestRetentionLogRoundTripAndSummary(t *testing.T) {
	t.Parallel()
	_, asset, tree, objectPath := auditFixture(t, 4)
	logPath := filepath.Join(t.TempDir(), "retention.log")
	clean, err := AuditAsset(objectPath, asset, tree, AuditPolicy{Samples: 4, Confidence: 0.99}, nil)
	if err != nil {
		t.Fatalf("AuditAsset() error = %v", err)
	}
	first := NewAuditRecord(clean, time.Date(2026, 9, 14, 1, 0, 0, 0, time.UTC))
	if err := AppendAuditRecord(logPath, first); err != nil {
		t.Fatalf("AppendAuditRecord() error = %v", err)
	}

	corruptChunk(t, objectPath, 0)
	damaged, err := AuditAsset(objectPath, asset, tree, AuditPolicy{Samples: 4, Confidence: 0.99}, nil)
	if err != nil {
		t.Fatalf("AuditAsset() error = %v", err)
	}
	second := NewAuditRecord(damaged, time.Date(2026, 9, 15, 1, 0, 0, 0, time.UTC))
	if err := AppendAuditRecord(logPath, second); err != nil {
		t.Fatalf("AppendAuditRecord() error = %v", err)
	}

	records, err := LoadAuditRecords(logPath)
	if err != nil {
		t.Fatalf("LoadAuditRecords() error = %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("len(records) = %d, want 2", len(records))
	}
	summary := SummarizeAuditRecords(records)
	if summary.Clean != 1 || summary.Damaged != 1 || summary.Records != 2 {
		t.Fatalf("summary = %+v, want one clean and one damaged", summary)
	}
	if summary.FirstChecked != "2026-09-14T01:00:00Z" || summary.LastChecked != "2026-09-15T01:00:00Z" {
		t.Fatalf("summary window = %+v", summary)
	}
}

func TestLoadAuditRecordsRejectsMalformedLines(t *testing.T) {
	t.Parallel()
	directory := t.TempDir()
	cases := map[string]string{
		"duplicate field": `{"schema_version":"arkmesh.audit/v0alpha1","recorded_at":"2026-09-14T01:00:00Z","recorded_at":"2026-09-15T01:00:00Z"}`,
		"unknown field":   `{"schema_version":"arkmesh.audit/v0alpha1","recorded_at":"2026-09-14T01:00:00Z","surprise":1}`,
		"bad schema":      `{"schema_version":"other","recorded_at":"2026-09-14T01:00:00Z"}`,
		"bad timestamp":   `{"schema_version":"arkmesh.audit/v0alpha1","recorded_at":"yesterday"}`,
	}
	for name, line := range cases {
		path := filepath.Join(directory, strings.ReplaceAll(name, " ", "-")+".log")
		if err := os.WriteFile(path, []byte(line+"\n"), 0o644); err != nil {
			t.Fatalf("write log: %v", err)
		}
		if _, err := LoadAuditRecords(path); err == nil {
			t.Fatalf("LoadAuditRecords(%s) error = nil, want rejection", name)
		}
	}
}

func TestAppendAuditRecordRejectsInvalidRecords(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "retention.log")
	if err := AppendAuditRecord(path, AuditRecord{SchemaVersion: "other", RecordedAt: "2026-09-14T01:00:00Z"}); err == nil {
		t.Fatal("AppendAuditRecord() error = nil, want schema rejection")
	}
	if err := AppendAuditRecord(path, AuditRecord{SchemaVersion: AuditSchemaVersion, RecordedAt: "tomorrow"}); err == nil {
		t.Fatal("AppendAuditRecord() error = nil, want timestamp rejection")
	}
}
