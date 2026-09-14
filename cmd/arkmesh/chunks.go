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
	default:
		return fmt.Errorf("unknown chunks command %q", args[0])
	}
}

func findAsset(manifest capsule.Manifest, digest string) (capsule.Asset, error) {
	for _, asset := range manifest.Assets {
		if asset.SHA256 == digest {
			return asset, nil
		}
	}
	return capsule.Asset{}, fmt.Errorf("capsule has no asset %q", digest)
}
