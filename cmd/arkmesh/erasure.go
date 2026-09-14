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

func runErasure(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return errors.New("usage: arkmesh erasure <protect|recover>")
	}
	switch args[0] {
	case "protect":
		flags := flag.NewFlagSet("erasure protect", flag.ContinueOnError)
		flags.SetOutput(stderr)
		assetDigest := flags.String("asset", "", "asset SHA-256 digest to protect")
		dataShards := flags.Int("data", 4, "number of data shards")
		parityShards := flags.Int("parity", 2, "number of parity shards")
		shardDirectory := flags.String("shards", "", "new directory for shard files")
		planPath := flags.String("plan", "", "new erasure plan file")
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		if flags.NArg() != 1 || strings.TrimSpace(*assetDigest) == "" || strings.TrimSpace(*shardDirectory) == "" || strings.TrimSpace(*planPath) == "" {
			return errors.New("usage: arkmesh erasure protect --asset DIGEST [--data K] [--parity M] --shards DIR --plan FILE CAPSULE")
		}
		manifest, err := capsule.Verify(flags.Arg(0))
		if err != nil {
			return err
		}
		asset, err := findAsset(manifest, *assetDigest)
		if err != nil {
			return err
		}
		plan, err := capsule.ProtectAsset(manifest.CapsuleID, asset, filepath.Join(flags.Arg(0), "objects", asset.SHA256), *shardDirectory, *dataShards, *parityShards)
		if err != nil {
			return err
		}
		if err := capsule.WriteErasurePlan(*planPath, plan); err != nil {
			return err
		}
		fmt.Fprintf(stdout, "protected asset\nplan: %s\nshards: %s\nlayout: %d data + %d parity\nshard size: %d\ntolerates losing: %d shards\n",
			*planPath, *shardDirectory, plan.DataShards, plan.ParityShards, plan.ShardSize, plan.ParityShards)
		return nil
	case "recover":
		flags := flag.NewFlagSet("erasure recover", flag.ContinueOnError)
		flags.SetOutput(stderr)
		planPath := flags.String("plan", "", "erasure plan file")
		shardDirectory := flags.String("shards", "", "directory holding surviving shards")
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		if flags.NArg() != 1 || strings.TrimSpace(*planPath) == "" || strings.TrimSpace(*shardDirectory) == "" {
			return errors.New("usage: arkmesh erasure recover --plan FILE --shards DIR CAPSULE")
		}
		manifest, err := capsule.Load(flags.Arg(0))
		if err != nil {
			return err
		}
		if err := capsule.ValidateManifest(manifest); err != nil {
			return err
		}
		plan, err := capsule.LoadErasurePlan(*planPath)
		if err != nil {
			return err
		}
		report, err := capsule.RecoverAsset(manifest, plan, *shardDirectory, flags.Arg(0))
		if err != nil {
			return err
		}
		fmt.Fprintf(stdout, "recovered asset\nasset: %s\nbytes: %d\nunusable shards: %d\n", report.ReconstructedFor, report.Bytes, len(report.MissingShards))
		return nil
	default:
		return fmt.Errorf("unknown erasure command %q", args[0])
	}
}
