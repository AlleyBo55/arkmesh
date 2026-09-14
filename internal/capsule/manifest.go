package capsule

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

const SchemaVersion = "arkmesh.capsule/v0alpha1"

var (
	digestPattern = regexp.MustCompile(`^[0-9a-f]{64}$`)
	rolePattern   = regexp.MustCompile(`^[a-z][a-z0-9_-]*$`)
)

type AssetSource struct {
	Role string
	Path string
}

type Asset struct {
	Role      string `json:"role"`
	Name      string `json:"name"`
	Size      int64  `json:"size"`
	SHA256    string `json:"sha256"`
	ChunkSize int64  `json:"chunk_size,omitempty"`
	ChunkRoot string `json:"chunk_root,omitempty"`
}

type Manifest struct {
	SchemaVersion   string  `json:"schema_version"`
	CapsuleID       string  `json:"capsule_id"`
	Name            string  `json:"name"`
	CreatedAt       string  `json:"created_at"`
	ParentCapsuleID string  `json:"parent_capsule_id,omitempty"`
	Assets          []Asset `json:"assets"`
}

type PackOptions struct {
	ParentCapsuleID string
	ChunkSize       int64
}

func Pack(name, outputDir string, sources []AssetSource, now time.Time) (Manifest, error) {
	return PackWithOptions(name, outputDir, sources, now, PackOptions{})
}

func PackWithOptions(name, outputDir string, sources []AssetSource, now time.Time, options PackOptions) (Manifest, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Manifest{}, errors.New("capsule name is required")
	}
	if len(sources) == 0 {
		return Manifest{}, errors.New("at least one asset is required")
	}
	parentCapsuleID := strings.TrimSpace(options.ParentCapsuleID)
	if parentCapsuleID != "" && !validCapsuleID(parentCapsuleID) {
		return Manifest{}, fmt.Errorf("invalid parent capsule ID %q", options.ParentCapsuleID)
	}
	chunkSize := options.ChunkSize
	if chunkSize == 0 {
		chunkSize = DefaultChunkSize
	}
	if !ValidChunkSize(chunkSize) {
		return Manifest{}, fmt.Errorf("invalid chunk size %d", options.ChunkSize)
	}
	if _, err := os.Stat(filepath.Join(outputDir, "manifest.json")); err == nil {
		return Manifest{}, fmt.Errorf("capsule already exists at %s", outputDir)
	} else if !errors.Is(err, os.ErrNotExist) {
		return Manifest{}, fmt.Errorf("inspect output: %w", err)
	}

	objectsDir := filepath.Join(outputDir, "objects")
	if err := os.MkdirAll(objectsDir, 0o755); err != nil {
		return Manifest{}, fmt.Errorf("create objects directory: %w", err)
	}

	assets := make([]Asset, 0, len(sources))
	for _, source := range sources {
		asset, err := storeAsset(objectsDir, source, chunkSize)
		if err != nil {
			return Manifest{}, err
		}
		assets = append(assets, asset)
	}

	sort.Slice(assets, func(i, j int) bool {
		if assets[i].Role != assets[j].Role {
			return assets[i].Role < assets[j].Role
		}
		if assets[i].Name != assets[j].Name {
			return assets[i].Name < assets[j].Name
		}
		return assets[i].SHA256 < assets[j].SHA256
	})

	manifest := Manifest{
		SchemaVersion:   SchemaVersion,
		Name:            name,
		CreatedAt:       now.UTC().Format(time.RFC3339),
		ParentCapsuleID: parentCapsuleID,
		Assets:          assets,
	}
	manifest.CapsuleID = capsuleID(manifest)

	if err := writeManifest(outputDir, manifest); err != nil {
		return Manifest{}, err
	}
	return manifest, nil
}

func Load(root string) (Manifest, error) {
	file, err := os.Open(filepath.Join(root, "manifest.json"))
	if err != nil {
		return Manifest{}, fmt.Errorf("read manifest: %w", err)
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return Manifest{}, fmt.Errorf("read manifest: %w", err)
	}
	if err := rejectDuplicateJSONKeys(data); err != nil {
		return Manifest{}, fmt.Errorf("decode manifest: %w", err)
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var manifest Manifest
	if err := decoder.Decode(&manifest); err != nil {
		return Manifest{}, fmt.Errorf("decode manifest: %w", err)
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		if err == nil {
			return Manifest{}, errors.New("decode manifest: unexpected trailing JSON value")
		}
		return Manifest{}, fmt.Errorf("decode manifest: %w", err)
	}
	return manifest, nil
}

func Verify(root string) (Manifest, error) {
	manifest, err := Load(root)
	if err != nil {
		return Manifest{}, err
	}
	if err := validateManifest(manifest); err != nil {
		return Manifest{}, err
	}

	for _, asset := range manifest.Assets {
		objectPath := filepath.Join(root, "objects", asset.SHA256)
		info, err := os.Lstat(objectPath)
		if err != nil {
			return Manifest{}, fmt.Errorf("asset %q (%s): %w", asset.Name, asset.Role, err)
		}
		if !info.Mode().IsRegular() {
			return Manifest{}, fmt.Errorf("asset %q (%s): object is not a regular file", asset.Name, asset.Role)
		}
		if info.Size() != asset.Size {
			return Manifest{}, fmt.Errorf("asset %q (%s): size mismatch: manifest=%d actual=%d", asset.Name, asset.Role, asset.Size, info.Size())
		}
		if asset.ChunkRoot != "" {
			digest, _, tree, err := ScanObject(objectPath, asset.ChunkSize)
			if err != nil {
				return Manifest{}, fmt.Errorf("asset %q (%s): %w", asset.Name, asset.Role, err)
			}
			if digest != asset.SHA256 {
				return Manifest{}, fmt.Errorf("asset %q (%s): digest mismatch", asset.Name, asset.Role)
			}
			if tree.Count != ChunkCount(asset.Size, asset.ChunkSize) {
				return Manifest{}, fmt.Errorf("asset %q (%s): chunk count mismatch", asset.Name, asset.Role)
			}
			if tree.Root != asset.ChunkRoot {
				return Manifest{}, fmt.Errorf("asset %q (%s): chunk root mismatch", asset.Name, asset.Role)
			}
			continue
		}
		digest, _, err := hashFile(objectPath)
		if err != nil {
			return Manifest{}, fmt.Errorf("asset %q (%s): %w", asset.Name, asset.Role, err)
		}
		if digest != asset.SHA256 {
			return Manifest{}, fmt.Errorf("asset %q (%s): digest mismatch", asset.Name, asset.Role)
		}
	}

	return manifest, nil
}

func validCapsuleID(value string) bool {
	algorithm, digest, found := strings.Cut(value, ":")
	return found && algorithm == "sha256" && digestPattern.MatchString(digest)
}

func ValidateManifest(manifest Manifest) error {
	return validateManifest(manifest)
}

func validateManifest(manifest Manifest) error {
	if manifest.SchemaVersion != SchemaVersion {
		return fmt.Errorf("unsupported schema version %q", manifest.SchemaVersion)
	}
	if strings.TrimSpace(manifest.Name) == "" {
		return errors.New("manifest name is required")
	}
	if _, err := time.Parse(time.RFC3339, manifest.CreatedAt); err != nil {
		return fmt.Errorf("invalid created_at: %w", err)
	}
	if len(manifest.Assets) == 0 {
		return errors.New("manifest must contain at least one asset")
	}
	if manifest.ParentCapsuleID != "" && !validCapsuleID(manifest.ParentCapsuleID) {
		return fmt.Errorf("invalid parent capsule ID %q", manifest.ParentCapsuleID)
	}
	if manifest.CapsuleID != capsuleID(manifest) {
		return errors.New("capsule ID does not match manifest")
	}
	for i, asset := range manifest.Assets {
		if !rolePattern.MatchString(asset.Role) {
			return fmt.Errorf("asset %d has invalid role %q", i, asset.Role)
		}
		if strings.TrimSpace(asset.Name) == "" {
			return fmt.Errorf("asset %d has an empty name", i)
		}
		if asset.Size < 0 {
			return fmt.Errorf("asset %q has a negative size", asset.Name)
		}
		if !digestPattern.MatchString(asset.SHA256) {
			return fmt.Errorf("asset %q has an invalid SHA-256 digest", asset.Name)
		}
		if (asset.ChunkSize != 0) != (asset.ChunkRoot != "") {
			return fmt.Errorf("asset %q must declare both chunk size and chunk root", asset.Name)
		}
		if asset.ChunkRoot != "" {
			if !ValidChunkSize(asset.ChunkSize) {
				return fmt.Errorf("asset %q has an invalid chunk size %d", asset.Name, asset.ChunkSize)
			}
			if !digestPattern.MatchString(asset.ChunkRoot) {
				return fmt.Errorf("asset %q has an invalid chunk root", asset.Name)
			}
		}
	}
	return nil
}

func storeAsset(objectsDir string, source AssetSource, chunkSize int64) (Asset, error) {
	role := strings.TrimSpace(source.Role)
	if !rolePattern.MatchString(role) {
		return Asset{}, fmt.Errorf("invalid asset role %q", source.Role)
	}
	if strings.TrimSpace(source.Path) == "" {
		return Asset{}, errors.New("asset path is required")
	}
	info, err := os.Stat(source.Path)
	if err != nil {
		return Asset{}, fmt.Errorf("inspect asset %q: %w", source.Path, err)
	}
	if !info.Mode().IsRegular() {
		return Asset{}, fmt.Errorf("asset %q is not a regular file", source.Path)
	}

	temp, err := os.CreateTemp(objectsDir, ".packing-*")
	if err != nil {
		return Asset{}, fmt.Errorf("create temporary object: %w", err)
	}
	tempPath := temp.Name()
	keepTemp := false
	defer func() {
		if !keepTemp {
			_ = os.Remove(tempPath)
		}
	}()

	sourceFile, err := os.Open(source.Path)
	if err != nil {
		_ = temp.Close()
		return Asset{}, fmt.Errorf("open asset %q: %w", source.Path, err)
	}

	hasher := sha256.New()
	size, copyErr := io.Copy(io.MultiWriter(temp, hasher), sourceFile)
	closeSourceErr := sourceFile.Close()
	closeTempErr := temp.Close()
	if copyErr != nil {
		return Asset{}, fmt.Errorf("copy asset %q: %w", source.Path, copyErr)
	}
	if closeSourceErr != nil {
		return Asset{}, fmt.Errorf("close asset %q: %w", source.Path, closeSourceErr)
	}
	if closeTempErr != nil {
		return Asset{}, fmt.Errorf("close object for %q: %w", source.Path, closeTempErr)
	}

	digest := hex.EncodeToString(hasher.Sum(nil))
	finalPath := filepath.Join(objectsDir, digest)
	if _, err := os.Stat(finalPath); errors.Is(err, os.ErrNotExist) {
		if err := os.Rename(tempPath, finalPath); err != nil {
			return Asset{}, fmt.Errorf("store object for %q: %w", source.Path, err)
		}
		keepTemp = true
	} else if err != nil {
		return Asset{}, fmt.Errorf("inspect stored object for %q: %w", source.Path, err)
	}

	storedDigest, storedSize, tree, err := ScanObject(finalPath, chunkSize)
	if err != nil {
		return Asset{}, fmt.Errorf("commit chunks for %q: %w", source.Path, err)
	}
	if storedDigest != digest || storedSize != size {
		return Asset{}, fmt.Errorf("stored object for %q changed during packing", source.Path)
	}

	return Asset{
		Role:      role,
		Name:      filepath.Base(source.Path),
		Size:      size,
		SHA256:    digest,
		ChunkSize: chunkSize,
		ChunkRoot: tree.Root,
	}, nil
}

func writeManifest(root string, manifest Manifest) error {
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("encode manifest: %w", err)
	}
	data = append(data, '\n')
	temp, err := os.CreateTemp(root, ".manifest-*")
	if err != nil {
		return fmt.Errorf("create temporary manifest: %w", err)
	}
	tempPath := temp.Name()
	defer os.Remove(tempPath)
	if _, err := temp.Write(data); err != nil {
		_ = temp.Close()
		return fmt.Errorf("write manifest: %w", err)
	}
	if err := temp.Close(); err != nil {
		return fmt.Errorf("close manifest: %w", err)
	}
	if err := os.Rename(tempPath, filepath.Join(root, "manifest.json")); err != nil {
		return fmt.Errorf("store manifest: %w", err)
	}
	return nil
}

func capsuleID(manifest Manifest) string {
	manifest.CapsuleID = ""
	data, err := json.Marshal(manifest)
	if err != nil {
		panic(fmt.Sprintf("marshal manifest identity: %v", err))
	}
	digest := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(digest[:])
}

func hashFile(path string) (string, int64, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", 0, err
	}
	defer file.Close()
	hasher := sha256.New()
	size, err := io.Copy(hasher, file)
	if err != nil {
		return "", 0, err
	}
	return hex.EncodeToString(hasher.Sum(nil)), size, nil
}
