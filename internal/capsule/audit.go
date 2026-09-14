package capsule

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"arkmesh/internal/audit"
	"arkmesh/internal/strictjson"
)

const AuditSchemaVersion = "arkmesh.audit/v0alpha1"

type AuditPolicy struct {
	Tolerance  float64
	Confidence float64
	Samples    int
}

// AuditResult is the outcome of one sampled retrievability check.
//
// A clean result never means "intact". It means no more than MaxBadChunks can be
// missing at the stated confidence, which is what sampling can actually support.
type AuditResult struct {
	CapsuleID         string
	AssetSHA256       string
	ChunkCount        int
	Sampled           int
	Failed            []int
	Tolerance         float64
	Confidence        float64
	MaxBadChunks      int
	MinIntactFraction float64
	Exhaustive        bool
}

// AuditRecord is one line of the retention log.
type AuditRecord struct {
	SchemaVersion     string  `json:"schema_version"`
	RecordedAt        string  `json:"recorded_at"`
	CapsuleID         string  `json:"capsule_id"`
	AssetSHA256       string  `json:"asset_sha256"`
	ChunkCount        int     `json:"chunk_count"`
	Sampled           int     `json:"sampled"`
	Failed            int     `json:"failed"`
	Tolerance         float64 `json:"tolerance"`
	Confidence        float64 `json:"confidence"`
	MaxBadChunks      int     `json:"max_bad_chunks"`
	MinIntactFraction float64 `json:"min_intact_fraction"`
	Exhaustive        bool    `json:"exhaustive"`
}

type RetentionSummary struct {
	Records      int
	Clean        int
	Damaged      int
	FirstChecked string
	LastChecked  string
	TotalSampled int
}

// AuditAsset samples chunks at random and compares them with an authenticated
// chunk tree. The caller must have already verified the tree against the signed
// manifest, which is what makes a leaf comparison meaningful.
func AuditAsset(objectPath string, asset Asset, tree ChunkTreeFile, policy AuditPolicy, random io.Reader) (AuditResult, error) {
	if tree.ChunkCount <= 0 || len(tree.Leaves) != tree.ChunkCount {
		return AuditResult{}, errors.New("chunk tree does not describe any chunks")
	}
	confidence := policy.Confidence
	if confidence == 0 {
		confidence = 0.99
	}
	samples := policy.Samples
	if samples <= 0 {
		tolerance := policy.Tolerance
		if tolerance == 0 {
			tolerance = 0.01
		}
		derived, err := audit.SamplesForFraction(tree.ChunkCount, tolerance, confidence)
		if err != nil {
			return AuditResult{}, err
		}
		samples = derived
		policy.Tolerance = tolerance
	}
	if samples > tree.ChunkCount {
		samples = tree.ChunkCount
	}
	indexes, err := audit.SelectSamples(tree.ChunkCount, samples, random)
	if err != nil {
		return AuditResult{}, err
	}

	file, err := os.Open(objectPath)
	if err != nil {
		return AuditResult{}, err
	}
	defer file.Close()

	buffer := make([]byte, asset.ChunkSize)
	var failed []int
	for _, index := range indexes {
		offset := int64(index) * asset.ChunkSize
		length := expectedChunkLength(asset.Size, asset.ChunkSize, index)
		read, err := file.ReadAt(buffer[:length], offset)
		if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
			return AuditResult{}, err
		}
		if int64(read) != length || hex.EncodeToString(chunkLeafDigestBytes(buffer[:length])) != tree.Leaves[index] {
			failed = append(failed, index)
		}
	}

	result := AuditResult{
		CapsuleID:   tree.CapsuleID,
		AssetSHA256: asset.SHA256,
		ChunkCount:  tree.ChunkCount,
		Sampled:     len(indexes),
		Failed:      failed,
		Tolerance:   policy.Tolerance,
		Confidence:  confidence,
		Exhaustive:  len(indexes) == tree.ChunkCount,
	}
	if len(failed) > 0 {
		// Damage was observed, so no upper bound on hidden damage is claimed.
		result.MaxBadChunks = tree.ChunkCount
		result.MinIntactFraction = 0
		return result, nil
	}
	bound, err := audit.MaxUndetectedBadChunks(tree.ChunkCount, len(indexes), confidence)
	if err != nil {
		return AuditResult{}, err
	}
	result.MaxBadChunks = bound
	result.MinIntactFraction = audit.MinIntactFraction(tree.ChunkCount, bound)
	return result, nil
}

func NewAuditRecord(result AuditResult, now time.Time) AuditRecord {
	return AuditRecord{
		SchemaVersion:     AuditSchemaVersion,
		RecordedAt:        now.UTC().Format(time.RFC3339),
		CapsuleID:         result.CapsuleID,
		AssetSHA256:       result.AssetSHA256,
		ChunkCount:        result.ChunkCount,
		Sampled:           result.Sampled,
		Failed:            len(result.Failed),
		Tolerance:         result.Tolerance,
		Confidence:        result.Confidence,
		MaxBadChunks:      result.MaxBadChunks,
		MinIntactFraction: result.MinIntactFraction,
		Exhaustive:        result.Exhaustive,
	}
}

// AppendAuditRecord adds one line to the retention log. The log is append only
// by convention, not by enforcement: it records what an operator observed and
// cannot prove that observations were not omitted.
func AppendAuditRecord(path string, record AuditRecord) error {
	if record.SchemaVersion != AuditSchemaVersion {
		return fmt.Errorf("unsupported audit schema %q", record.SchemaVersion)
	}
	if _, err := time.Parse(time.RFC3339, record.RecordedAt); err != nil {
		return fmt.Errorf("invalid audit timestamp: %w", err)
	}
	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("encode audit record: %w", err)
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return fmt.Errorf("open retention log: %w", err)
	}
	if _, err := file.Write(append(data, '\n')); err != nil {
		_ = file.Close()
		return fmt.Errorf("write retention log: %w", err)
	}
	if err := file.Sync(); err != nil {
		_ = file.Close()
		return fmt.Errorf("sync retention log: %w", err)
	}
	return file.Close()
}

func LoadAuditRecords(path string) ([]AuditRecord, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, fmt.Errorf("inspect retention log: %w", err)
	}
	if !info.Mode().IsRegular() {
		return nil, errors.New("retention log is not a regular file")
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var records []AuditRecord
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	line := 0
	for scanner.Scan() {
		line++
		raw := bytes.TrimSpace(scanner.Bytes())
		if len(raw) == 0 {
			continue
		}
		if err := strictjson.RejectDuplicateKeys(bytes.NewReader(raw)); err != nil {
			return nil, fmt.Errorf("retention log line %d: %w", line, err)
		}
		decoder := json.NewDecoder(bytes.NewReader(raw))
		decoder.DisallowUnknownFields()
		var record AuditRecord
		if err := decoder.Decode(&record); err != nil {
			return nil, fmt.Errorf("retention log line %d: %w", line, err)
		}
		if record.SchemaVersion != AuditSchemaVersion {
			return nil, fmt.Errorf("retention log line %d: unsupported schema %q", line, record.SchemaVersion)
		}
		if _, err := time.Parse(time.RFC3339, record.RecordedAt); err != nil {
			return nil, fmt.Errorf("retention log line %d: invalid timestamp: %w", line, err)
		}
		records = append(records, record)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read retention log: %w", err)
	}
	return records, nil
}

func SummarizeAuditRecords(records []AuditRecord) RetentionSummary {
	summary := RetentionSummary{Records: len(records)}
	for index, record := range records {
		if record.Failed == 0 {
			summary.Clean++
		} else {
			summary.Damaged++
		}
		summary.TotalSampled += record.Sampled
		if index == 0 || record.RecordedAt < summary.FirstChecked {
			summary.FirstChecked = record.RecordedAt
		}
		if index == 0 || record.RecordedAt > summary.LastChecked {
			summary.LastChecked = record.RecordedAt
		}
	}
	return summary
}
