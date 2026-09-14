package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"arkmesh/internal/capsule"
)

func runChunks(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return errors.New("usage: arkmesh chunks <inspect|prove|check>")
	}
	switch args[0] {
	case "inspect":
		if len(args) != 2 {
			return errors.New("usage: arkmesh chunks inspect CAPSULE")
		}
		manifest, err := capsule.Verify(args[1])
		if err != nil {
			return err
		}
		fmt.Fprintf(stdout, "capsule: %s\n", manifest.CapsuleID)
		for _, asset := range manifest.Assets {
			if asset.ChunkRoot == "" {
				fmt.Fprintf(stdout, "%s %s: no chunk commitment\n", asset.Role, asset.Name)
				continue
			}
			fmt.Fprintf(stdout, "%s %s: size=%d chunk_size=%d chunks=%d root=%s\n",
				asset.Role, asset.Name, asset.Size, asset.ChunkSize,
				capsule.ChunkCount(asset.Size, asset.ChunkSize), asset.ChunkRoot)
		}
		return nil
	case "prove":
		flags := flag.NewFlagSet("chunks prove", flag.ContinueOnError)
		flags.SetOutput(stderr)
		assetDigest := flags.String("asset", "", "asset SHA-256 digest to prove")
		index := flags.Int("index", -1, "chunk index to prove")
		output := flags.String("out", "", "new chunk proof file")
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		if flags.NArg() != 1 || strings.TrimSpace(*assetDigest) == "" || *index < 0 || strings.TrimSpace(*output) == "" {
			return errors.New("usage: arkmesh chunks prove --asset DIGEST --index N --out FILE CAPSULE")
		}
		manifest, err := capsule.Verify(flags.Arg(0))
		if err != nil {
			return err
		}
		asset, err := findAsset(manifest, *assetDigest)
		if err != nil {
			return err
		}
		objectPath := filepath.Join(flags.Arg(0), "objects", asset.SHA256)
		proof, err := capsule.BuildChunkProof(manifest.CapsuleID, asset, objectPath, *index)
		if err != nil {
			return err
		}
		if err := capsule.WriteChunkProof(*output, proof); err != nil {
			return err
		}
		fmt.Fprintf(stdout, "created chunk proof\npath: %s\nasset: %s\nchunk: %d of %d\nroot: %s\n",
			*output, proof.AssetSHA256, proof.ChunkIndex, proof.ChunkCount, proof.ChunkRoot)
		return nil
	case "check":
		flags := flag.NewFlagSet("chunks check", flag.ContinueOnError)
		flags.SetOutput(stderr)
		proofPath := flags.String("proof", "", "chunk proof file to verify")
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		if flags.NArg() != 1 || strings.TrimSpace(*proofPath) == "" {
			return errors.New("usage: arkmesh chunks check --proof FILE CAPSULE")
		}
		manifest, err := capsule.Load(flags.Arg(0))
		if err != nil {
			return err
		}
		if err := capsule.ValidateManifest(manifest); err != nil {
			return err
		}
		proof, err := capsule.LoadChunkProof(*proofPath)
		if err != nil {
			return err
		}
		if err := capsule.VerifyChunkProof(manifest, proof); err != nil {
			return err
		}
		fmt.Fprintf(stdout, "chunk proof verified\ncapsule: %s\nasset: %s\nchunk: %d of %d\n",
			manifest.CapsuleID, proof.AssetSHA256, proof.ChunkIndex, proof.ChunkCount)
		return nil
	case "tree":
		flags := flag.NewFlagSet("chunks tree", flag.ContinueOnError)
		flags.SetOutput(stderr)
		assetDigest := flags.String("asset", "", "asset SHA-256 digest to export")
		output := flags.String("out", "", "new chunk tree file")
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		if flags.NArg() != 1 || strings.TrimSpace(*assetDigest) == "" || strings.TrimSpace(*output) == "" {
			return errors.New("usage: arkmesh chunks tree --asset DIGEST --out FILE CAPSULE")
		}
		manifest, err := capsule.Verify(flags.Arg(0))
		if err != nil {
			return err
		}
		asset, err := findAsset(manifest, *assetDigest)
		if err != nil {
			return err
		}
		tree, err := capsule.BuildChunkTreeFile(manifest.CapsuleID, asset, filepath.Join(flags.Arg(0), "objects", asset.SHA256))
		if err != nil {
			return err
		}
		if err := capsule.WriteChunkTreeFile(*output, tree); err != nil {
			return err
		}
		fmt.Fprintf(stdout, "exported chunk tree\npath: %s\nasset: %s\nchunks: %d\nroot: %s\n", *output, tree.AssetSHA256, tree.ChunkCount, tree.ChunkRoot)
		return nil
	case "scan":
		flags := flag.NewFlagSet("chunks scan", flag.ContinueOnError)
		flags.SetOutput(stderr)
		treePath := flags.String("tree", "", "authenticated chunk tree file")
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		if flags.NArg() != 1 || strings.TrimSpace(*treePath) == "" {
			return errors.New("usage: arkmesh chunks scan --tree FILE CAPSULE")
		}
		asset, tree, objectPath, err := loadRepairInputs(flags.Arg(0), *treePath)
		if err != nil {
			return err
		}
		damage, err := capsule.ScanChunkDamage(objectPath, asset, tree)
		if err != nil {
			return err
		}
		if len(damage) == 0 {
			fmt.Fprintf(stdout, "no chunk damage\nasset: %s\nchunks: %d\n", asset.SHA256, tree.ChunkCount)
			return nil
		}
		for _, item := range damage {
			fmt.Fprintf(stdout, "damaged chunk %d offset=%d length=%d reason=%s\n", item.Index, item.Offset, item.Length, item.Reason)
		}
		return fmt.Errorf("%d of %d chunks are damaged", len(damage), tree.ChunkCount)
	case "repair":
		flags := flag.NewFlagSet("chunks repair", flag.ContinueOnError)
		flags.SetOutput(stderr)
		treePath := flags.String("tree", "", "authenticated chunk tree file")
		sourcePath := flags.String("source", "", "donor copy of the object")
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		if flags.NArg() != 1 || strings.TrimSpace(*treePath) == "" || strings.TrimSpace(*sourcePath) == "" {
			return errors.New("usage: arkmesh chunks repair --tree FILE --source FILE CAPSULE")
		}
		asset, tree, objectPath, err := loadRepairInputs(flags.Arg(0), *treePath)
		if err != nil {
			return err
		}
		damage, err := capsule.ScanChunkDamage(objectPath, asset, tree)
		if err != nil {
			return err
		}
		repaired, err := capsule.RepairChunks(objectPath, asset, tree, *sourcePath, damage)
		if err != nil {
			return err
		}
		if _, err := capsule.Verify(flags.Arg(0)); err != nil {
			return fmt.Errorf("capsule still fails verification after repair: %w", err)
		}
		fmt.Fprintf(stdout, "repaired capsule\nasset: %s\nchunks repaired: %d of %d\n", asset.SHA256, repaired, tree.ChunkCount)
		return nil
	default:
		return fmt.Errorf("unknown chunks command %q", args[0])
	}
}

// loadRepairInputs reads a manifest that may describe a damaged object, so it
// validates the manifest without reading object bytes, then authenticates the
// supplied tree against that manifest.
func loadRepairInputs(capsuleRoot, treePath string) (capsule.Asset, capsule.ChunkTreeFile, string, error) {
	manifest, err := capsule.Load(capsuleRoot)
	if err != nil {
		return capsule.Asset{}, capsule.ChunkTreeFile{}, "", err
	}
	if err := capsule.ValidateManifest(manifest); err != nil {
		return capsule.Asset{}, capsule.ChunkTreeFile{}, "", err
	}
	tree, err := capsule.LoadChunkTreeFile(treePath)
	if err != nil {
		return capsule.Asset{}, capsule.ChunkTreeFile{}, "", err
	}
	asset, err := capsule.VerifyChunkTreeFile(manifest, tree)
	if err != nil {
		return capsule.Asset{}, capsule.ChunkTreeFile{}, "", err
	}
	return asset, tree, filepath.Join(capsuleRoot, "objects", asset.SHA256), nil
}

func findAsset(manifest capsule.Manifest, digest string) (capsule.Asset, error) {
	for _, asset := range manifest.Assets {
		if asset.SHA256 == digest {
			return asset, nil
		}
	}
	return capsule.Asset{}, fmt.Errorf("capsule has no asset %q", digest)
}
