package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	"arkmesh/internal/capsule"
	"arkmesh/internal/identity"
)

func runCheckpoint(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return errors.New("usage: arkmesh checkpoint <create|show|advance>")
	}
	switch args[0] {
	case "create":
		flags := flag.NewFlagSet("checkpoint create", flag.ContinueOnError)
		flags.SetOutput(stderr)
		output := flags.String("out", "", "new local checkpoint file")
		var trustPaths repeatedFlags
		flags.Var(&trustPaths, "trust", "trusted public identity file (repeatable)")
		var revocationPaths repeatedFlags
		flags.Var(&revocationPaths, "revocations", "local revocation policy file (repeatable)")
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		if flags.NArg() != 1 || strings.TrimSpace(*output) == "" || len(trustPaths) == 0 {
			return errors.New("usage: arkmesh checkpoint create --out FILE --trust PUBLIC_IDENTITY [--revocations FILE] CAPSULE")
		}
		trusted, err := loadTrustedIdentities(trustPaths)
		if err != nil {
			return err
		}
		revocations, err := loadRevocations(revocationPaths)
		if err != nil {
			return err
		}
		manifest, authenticity, err := capsule.VerifyAuthenticated(flags.Arg(0), capsule.VerifyOptions{
			Trusted:          trusted,
			Revocations:      revocations,
			RequireSignature: true,
			RequireTrusted:   true,
		})
		if err != nil {
			return err
		}
		checkpoint, err := capsule.NewCheckpoint(manifest, authenticity)
		if err != nil {
			return err
		}
		if err := capsule.WriteCheckpoint(*output, checkpoint); err != nil {
			return err
		}
		fmt.Fprintf(stdout, "created local checkpoint\npath: %s\ncapsule: %s\nsigner: %s\n", *output, checkpoint.CapsuleID, checkpoint.SignerID)
		return nil
	case "show":
		if len(args) != 2 {
			return errors.New("usage: arkmesh checkpoint show FILE")
		}
		checkpoint, err := capsule.LoadCheckpoint(args[1])
		if err != nil {
			return err
		}
		encoder := json.NewEncoder(stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(checkpoint)
	case "advance":
		flags := flag.NewFlagSet("checkpoint advance", flag.ContinueOnError)
		flags.SetOutput(stderr)
		checkpointPath := flags.String("checkpoint", "", "local checkpoint file to advance")
		parentPath := flags.String("parent", "", "checkpointed parent capsule directory")
		var revocationPaths repeatedFlags
		flags.Var(&revocationPaths, "revocations", "local revocation policy file (repeatable)")
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		if flags.NArg() != 1 || strings.TrimSpace(*checkpointPath) == "" || strings.TrimSpace(*parentPath) == "" {
			return errors.New("usage: arkmesh checkpoint advance --checkpoint FILE --parent PARENT_CAPSULE [--revocations FILE] CHILD_CAPSULE")
		}
		revocations, err := loadRevocations(revocationPaths)
		if err != nil {
			return err
		}
		parentManifest, parentAuthenticity, err := capsule.VerifyAuthenticated(*parentPath, capsule.VerifyOptions{
			Revocations:      revocations,
			RequireSignature: true,
		})
		if err != nil {
			return fmt.Errorf("verify checkpoint parent: %w", err)
		}
		childManifest, childAuthenticity, err := capsule.VerifyAuthenticated(flags.Arg(0), capsule.VerifyOptions{
			Revocations:      revocations,
			RequireSignature: true,
		})
		if err != nil {
			return fmt.Errorf("verify checkpoint child: %w", err)
		}
		next, lineage, err := capsule.AdvanceCheckpointFile(*checkpointPath, flags.Arg(0), parentManifest, parentAuthenticity, childManifest, childAuthenticity)
		if err != nil {
			return err
		}
		fmt.Fprintf(stdout, "advanced local checkpoint\npath: %s\ncapsule: %s\nsigner: %s\nlineage: %s\n", *checkpointPath, next.CapsuleID, next.SignerID, lineage.Status)
		return nil
	default:
		return fmt.Errorf("unknown checkpoint command %q", args[0])
	}
}

func loadTrustedIdentities(paths repeatedFlags) ([]identity.PublicIdentity, error) {
	trusted := make([]identity.PublicIdentity, 0, len(paths))
	for _, path := range paths {
		public, err := identity.LoadPublic(path)
		if err != nil {
			return nil, fmt.Errorf("load trusted identity %q: %w", path, err)
		}
		trusted = append(trusted, public)
	}
	return trusted, nil
}

func loadRevocations(paths repeatedFlags) (identity.RevocationSet, error) {
	sets := make([]identity.RevocationSet, 0, len(paths))
	for _, path := range paths {
		set, err := identity.LoadRevocationSet(path)
		if err != nil {
			return identity.RevocationSet{}, fmt.Errorf("load revocation policy %q: %w", path, err)
		}
		sets = append(sets, set)
	}
	return identity.MergeRevocationSets(sets...)
}
