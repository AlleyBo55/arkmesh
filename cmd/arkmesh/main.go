package main

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"arkmesh/internal/capsule"
	"arkmesh/internal/identity"
)

type repeatedFlags []string

func (values *repeatedFlags) String() string {
	return strings.Join(*values, ",")
}

func (values *repeatedFlags) Set(value string) error {
	*values = append(*values, value)
	return nil
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		printUsage(stderr)
		return 2
	}

	var err error
	switch args[0] {
	case "identity":
		err = runIdentity(args[1:], stdout, stderr)
	case "checkpoint":
		err = runCheckpoint(args[1:], stdout, stderr)
	case "recovery":
		err = runRecovery(args[1:], stdout, stderr)
	case "chunks":
		err = runChunks(args[1:], stdout, stderr)
	case "erasure":
		err = runErasure(args[1:], stdout, stderr)
	case "audit":
		err = runAudit(args[1:], stdout, stderr)
	case "pack":
		err = runPack(args[1:], stdout, stderr)
	case "inspect":
		err = runInspect(args[1:], stdout)
	case "verify":
		err = runVerify(args[1:], stdout, stderr)
	case "help", "-h", "--help":
		printUsage(stdout)
		return 0
	default:
		err = fmt.Errorf("unknown command %q", args[0])
	}
	if err != nil {
		fmt.Fprintf(stderr, "arkmesh: %v\n", err)
		return 1
	}
	return 0
}

func runIdentity(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return errors.New("usage: arkmesh identity <create|show|revoke>")
	}
	switch args[0] {
	case "create":
		flags := flag.NewFlagSet("identity create", flag.ContinueOnError)
		flags.SetOutput(stderr)
		output := flags.String("out", "", "new identity directory")
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		if flags.NArg() != 0 || strings.TrimSpace(*output) == "" {
			return errors.New("usage: arkmesh identity create --out DIR")
		}
		created, err := identity.Create(*output, rand.Reader)
		if err != nil {
			return err
		}
		fmt.Fprintf(stdout, "created %s\npublic: %s\nprivate: %s\n", created.Public.KeyID, created.PublicPath, created.PrivatePath)
		fmt.Fprintln(stdout, "Keep the private identity secret and backed up offline.")
		return nil
	case "show":
		if len(args) != 2 {
			return errors.New("usage: arkmesh identity show <public-identity.json>")
		}
		public, err := identity.LoadPublic(args[1])
		if err != nil {
			return err
		}
		encoder := json.NewEncoder(stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(public)
	case "revoke":
		flags := flag.NewFlagSet("identity revoke", flag.ContinueOnError)
		flags.SetOutput(stderr)
		output := flags.String("out", "", "new local revocation policy file")
		var identityPaths repeatedFlags
		flags.Var(&identityPaths, "identity", "public identity to revoke (repeatable)")
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		if flags.NArg() != 0 || strings.TrimSpace(*output) == "" || len(identityPaths) == 0 {
			return errors.New("usage: arkmesh identity revoke --identity PUBLIC_IDENTITY [--identity ...] --out FILE")
		}
		publicIdentities := make([]identity.PublicIdentity, 0, len(identityPaths))
		for _, path := range identityPaths {
			public, err := identity.LoadPublic(path)
			if err != nil {
				return fmt.Errorf("load revoked identity %q: %w", path, err)
			}
			publicIdentities = append(publicIdentities, public)
		}
		set, err := identity.WriteRevocationSet(*output, publicIdentities)
		if err != nil {
			return err
		}
		fmt.Fprintf(stdout, "created local revocation policy\npath: %s\nrevoked identities: %d\n", *output, len(set.RevokedKeyIDs))
		return nil
	default:
		return fmt.Errorf("unknown identity command %q", args[0])
	}
}

func runPack(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("pack", flag.ContinueOnError)
	flags.SetOutput(stderr)
	name := flags.String("name", "", "human-readable capsule name")
	output := flags.String("out", "", "output capsule directory")
	signingKey := flags.String("signing-key", "", "private identity used to sign the capsule")
	rotationKey := flags.String("rotation-key", "", "parent private identity authorizing a new child signer")
	recoveryPolicyPath := flags.String("recovery-policy", "", "recovery policy that will authorize a new child signer")
	chunkSize := flags.Int64("chunk-size", 0, "chunk commitment size in bytes (default 1 MiB)")
	parentPath := flags.String("parent", "", "verified parent capsule directory")
	var assets repeatedFlags
	flags.Var(&assets, "asset", "asset as role=/path/to/file (repeatable)")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return errors.New("pack accepts only flags; use repeated --asset role=path")
	}
	if strings.TrimSpace(*output) == "" {
		return errors.New("--out is required")
	}

	var signer *identity.PrivateIdentity
	if strings.TrimSpace(*signingKey) != "" {
		loaded, err := identity.LoadPrivate(*signingKey)
		if err != nil {
			return err
		}
		signer = &loaded
	}
	var rotationSigner *identity.PrivateIdentity
	if strings.TrimSpace(*rotationKey) != "" {
		loaded, err := identity.LoadPrivate(*rotationKey)
		if err != nil {
			return fmt.Errorf("load rotation key: %w", err)
		}
		rotationSigner = &loaded
	}

	parentCapsuleID := ""
	var parentManifest capsule.Manifest
	var parentAuthenticity capsule.Authenticity
	var recoveryPolicy capsule.RecoveryPolicy
	isRotation := false
	isRecovery := false
	recoveryRequested := strings.TrimSpace(*recoveryPolicyPath) != ""
	if rotationSigner != nil && recoveryRequested {
		return errors.New("--rotation-key and --recovery-policy are mutually exclusive")
	}
	if strings.TrimSpace(*parentPath) != "" {
		if signer == nil {
			return errors.New("--parent requires --signing-key")
		}
		var err error
		parentManifest, parentAuthenticity, err = capsule.VerifyAuthenticated(*parentPath, capsule.VerifyOptions{RequireSignature: true})
		if err != nil {
			return fmt.Errorf("verify parent capsule: %w", err)
		}
		publicSigner, err := signer.Public()
		if err != nil {
			return fmt.Errorf("validate child signer: %w", err)
		}
		if parentAuthenticity.SignerID == publicSigner.KeyID {
			if rotationSigner != nil || recoveryRequested {
				return errors.New("same-author child must not include rotation or recovery authority")
			}
		} else if recoveryRequested {
			recoveryPolicy, err = capsule.LoadRecoveryPolicy(*recoveryPolicyPath)
			if err != nil {
				return err
			}
			isRecovery = true
		} else {
			if rotationSigner == nil {
				return fmt.Errorf("different child signer %s requires --rotation-key or --recovery-policy", publicSigner.KeyID)
			}
			rotationPublic, err := rotationSigner.Public()
			if err != nil {
				return fmt.Errorf("validate rotation signer: %w", err)
			}
			if rotationPublic.KeyID != parentAuthenticity.SignerID {
				return fmt.Errorf("rotation signer %s does not match parent signer %s", rotationPublic.KeyID, parentAuthenticity.SignerID)
			}
			isRotation = true
		}
		parentCapsuleID = parentManifest.CapsuleID
	} else if rotationSigner != nil || recoveryRequested {
		return errors.New("rotation and recovery authority require --parent")
	}

	sources := make([]capsule.AssetSource, 0, len(assets))
	for _, value := range assets {
		role, path, ok := strings.Cut(value, "=")
		if !ok || strings.TrimSpace(role) == "" || strings.TrimSpace(path) == "" {
			return fmt.Errorf("invalid --asset %q; expected role=/path/to/file", value)
		}
		sources = append(sources, capsule.AssetSource{Role: role, Path: path})
	}

	manifest, err := capsule.PackWithOptions(*name, *output, sources, time.Now(), capsule.PackOptions{
		ParentCapsuleID: parentCapsuleID,
		ChunkSize:       *chunkSize,
	})
	if err != nil {
		return err
	}
	fmt.Fprintf(stdout, "packed %s\nassets: %d\npath: %s\n", manifest.CapsuleID, len(manifest.Assets), *output)
	if manifest.ParentCapsuleID != "" {
		fmt.Fprintf(stdout, "parent: %s\n", manifest.ParentCapsuleID)
	}
	if signer != nil {
		envelope, err := capsule.Sign(*output, *signer)
		if err != nil {
			return err
		}
		fmt.Fprintf(stdout, "signed by: %s\n", envelope.Signer.KeyID)
	}
	if isRotation {
		childManifest, childAuthenticity, err := capsule.VerifyAuthenticated(*output, capsule.VerifyOptions{RequireSignature: true})
		if err != nil {
			return fmt.Errorf("verify rotated child: %w", err)
		}
		transition, err := capsule.AuthorizeKeyRotation(*output, parentManifest, parentAuthenticity, childManifest, childAuthenticity, *rotationSigner)
		if err != nil {
			return err
		}
		fmt.Fprintf(stdout, "key rotation: %s -> %s\n", transition.PreviousSignerID, transition.NextSigner.KeyID)
	}
	if isRecovery {
		fmt.Fprintf(stdout, "recovery pending: %s\n", recoveryPolicy.PolicyID)
	}
	return nil
}

func runInspect(args []string, stdout io.Writer) error {
	if len(args) != 1 {
		return errors.New("usage: arkmesh inspect <capsule-directory>")
	}
	manifest, err := capsule.Load(args[0])
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(manifest)
}

func runVerify(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("verify", flag.ContinueOnError)
	flags.SetOutput(stderr)
	var trustPaths repeatedFlags
	flags.Var(&trustPaths, "trust", "trusted public identity file (repeatable)")
	var revocationPaths repeatedFlags
	flags.Var(&revocationPaths, "revocations", "local revocation policy file (repeatable)")
	parentPath := flags.String("parent", "", "parent capsule directory used to verify ancestry")
	checkpointPath := flags.String("checkpoint", "", "local checkpoint that must match this capsule")
	requireSignature := flags.Bool("require-signature", false, "reject unsigned capsules")
	requireTrusted := flags.Bool("require-trusted", false, "reject capsules without a trusted root or parent")
	requireLineage := flags.Bool("require-lineage", false, "reject descendants whose parent was not checked")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return errors.New("usage: arkmesh verify [flags] <capsule-directory>")
	}

	trusted, err := loadTrustedIdentities(trustPaths)
	if err != nil {
		return err
	}
	revocations, err := loadRevocations(revocationPaths)
	if err != nil {
		return err
	}

	checkpointSupplied := strings.TrimSpace(*checkpointPath) != ""
	childOptions := capsule.VerifyOptions{
		Trusted:          trusted,
		Revocations:      revocations,
		RequireSignature: *requireSignature || *requireTrusted || checkpointSupplied,
	}
	manifest, authenticity, err := capsule.VerifyAuthenticated(flags.Arg(0), childOptions)
	if err != nil {
		return err
	}

	checkpointMatched := false
	if checkpointSupplied {
		checkpoint, err := capsule.LoadCheckpoint(*checkpointPath)
		if err != nil {
			return err
		}
		if err := capsule.VerifyCheckpoint(checkpoint, manifest, authenticity); err != nil {
			return err
		}
		checkpointMatched = true
	}

	lineage := capsule.DescribeLineage(manifest)
	parentSupplied := strings.TrimSpace(*parentPath) != ""
	if manifest.ParentCapsuleID == "" {
		if parentSupplied {
			return errors.New("root capsule does not declare a parent")
		}
		if *requireTrusted && authenticity.Status != capsule.AuthenticityTrusted && !checkpointMatched {
			return fmt.Errorf("capsule signer %s is not trusted", authenticity.SignerID)
		}
	} else if !parentSupplied {
		if *requireLineage {
			return fmt.Errorf("capsule declares parent %s but no --parent was supplied", manifest.ParentCapsuleID)
		}
		if *requireTrusted && authenticity.Status != capsule.AuthenticityTrusted && !checkpointMatched {
			return fmt.Errorf("capsule signer %s is not trusted", authenticity.SignerID)
		}
	} else {
		parentOptions := capsule.VerifyOptions{
			Trusted:          trusted,
			Revocations:      revocations,
			RequireSignature: true,
			RequireTrusted:   *requireTrusted && !checkpointMatched,
		}
		parentManifest, parentAuthenticity, err := capsule.VerifyAuthenticated(*parentPath, parentOptions)
		if err != nil {
			return fmt.Errorf("verify parent capsule: %w", err)
		}
		lineage, err = capsule.CheckLineageWithAuthority(flags.Arg(0), manifest, authenticity, parentManifest, parentAuthenticity)
		if err != nil {
			return err
		}
	}

	fmt.Fprintf(stdout, "verified %s\nassets: %d\nsignature: %s\nlineage: %s\n", manifest.CapsuleID, len(manifest.Assets), authenticity.Status, lineage.Status)
	if authenticity.SignerID != "" {
		fmt.Fprintf(stdout, "signer: %s\n", authenticity.SignerID)
	}
	if checkpointMatched {
		fmt.Fprintln(stdout, "checkpoint: matched")
	}
	if lineage.ParentCapsuleID != "" {
		fmt.Fprintf(stdout, "parent: %s\n", lineage.ParentCapsuleID)
	}
	return nil
}

func printUsage(writer io.Writer) {
	fmt.Fprintln(writer, `ArkMesh v0alpha1

Usage:
  arkmesh identity create --out DIR
  arkmesh identity show PUBLIC_IDENTITY
  arkmesh identity revoke --identity PUBLIC_IDENTITY [--identity ...] --out FILE
  arkmesh checkpoint create --out FILE --trust PUBLIC_IDENTITY [--revocations FILE] CAPSULE
  arkmesh checkpoint show FILE
  arkmesh checkpoint advance --checkpoint FILE --parent PARENT_CAPSULE [--recovery-policy POLICY] [--revocations FILE] CHILD_CAPSULE
  arkmesh recovery policy create --threshold N --identity PUBLIC_IDENTITY [--identity ...] --out FILE
  arkmesh recovery approve --policy POLICY --signing-key PRIVATE_IDENTITY --parent PARENT --out FILE CHILD
  arkmesh recovery assemble --policy POLICY --parent PARENT --approval FILE [--approval ...] [--revocations FILE] CHILD
  arkmesh chunks inspect CAPSULE
  arkmesh chunks prove --asset DIGEST --index N --out FILE CAPSULE
  arkmesh chunks check --proof FILE CAPSULE
  arkmesh chunks tree --asset DIGEST --out FILE CAPSULE
  arkmesh chunks scan --tree FILE CAPSULE
  arkmesh chunks repair --tree FILE --source FILE CAPSULE
  arkmesh erasure protect --asset DIGEST [--data K] [--parity M] --shards DIR --plan FILE CAPSULE
  arkmesh erasure recover --plan FILE --shards DIR CAPSULE
  arkmesh audit sample --tree FILE [--tolerance F] [--confidence F] [--samples N] [--log FILE] CAPSULE
  arkmesh audit history --log FILE
  arkmesh pack --name NAME --out DIR --asset role=/path/to/file [--asset ...] [--signing-key PRIVATE_IDENTITY] [--parent PARENT_CAPSULE] [--rotation-key PARENT_PRIVATE_IDENTITY] [--recovery-policy POLICY] [--chunk-size BYTES]
  arkmesh inspect DIR
  arkmesh verify [--trust PUBLIC_IDENTITY] [--revocations FILE] [--checkpoint FILE] [--require-signature] [--require-trusted] [--parent PARENT_CAPSULE] [--require-lineage] DIR

Integrity verification remains available for unsigned capsules. Signed chunk roots allow single chunk possession proofs, damage localization, and repair that verifies donor chunks before writing them. Reed Solomon shards rebuild a missing object without any donor holding it, and reconstructed bytes must match the signed manifest. Sampled audits report a retrievability bound with explicit confidence rather than claiming intactness. Parent-signed transitions authorize exact planned rotations. Threshold recovery requires distinct approvals under explicit local policy. Local checkpoints reject rollback to another capsule head. Local revocation policy overrides trust.`)
}
