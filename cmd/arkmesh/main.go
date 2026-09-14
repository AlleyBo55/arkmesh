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
		return errors.New("usage: arkmesh identity <create|show>")
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

	parentCapsuleID := ""
	if strings.TrimSpace(*parentPath) != "" {
		if signer == nil {
			return errors.New("--parent requires --signing-key")
		}
		parentManifest, parentAuthenticity, err := capsule.VerifyAuthenticated(*parentPath, capsule.VerifyOptions{RequireSignature: true})
		if err != nil {
			return fmt.Errorf("verify parent capsule: %w", err)
		}
		publicSigner, err := signer.Public()
		if err != nil {
			return fmt.Errorf("validate child signer: %w", err)
		}
		if parentAuthenticity.SignerID != publicSigner.KeyID {
			return fmt.Errorf("child signer %s is not authorized by parent signer %s", publicSigner.KeyID, parentAuthenticity.SignerID)
		}
		parentCapsuleID = parentManifest.CapsuleID
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
	parentPath := flags.String("parent", "", "parent capsule directory used to verify ancestry")
	requireSignature := flags.Bool("require-signature", false, "reject unsigned capsules")
	requireTrusted := flags.Bool("require-trusted", false, "reject capsules not signed by a trusted identity")
	requireLineage := flags.Bool("require-lineage", false, "reject descendants whose parent was not checked")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return errors.New("usage: arkmesh verify [flags] <capsule-directory>")
	}

	trusted := make([]identity.PublicIdentity, 0, len(trustPaths))
	for _, path := range trustPaths {
		public, err := identity.LoadPublic(path)
		if err != nil {
			return fmt.Errorf("load trusted identity %q: %w", path, err)
		}
		trusted = append(trusted, public)
	}
	verificationOptions := capsule.VerifyOptions{
		Trusted:          trusted,
		RequireSignature: *requireSignature,
		RequireTrusted:   *requireTrusted,
	}
	manifest, authenticity, err := capsule.VerifyAuthenticated(flags.Arg(0), verificationOptions)
	if err != nil {
		return err
	}

	lineage := capsule.DescribeLineage(manifest)
	if manifest.ParentCapsuleID == "" {
		if strings.TrimSpace(*parentPath) != "" {
			return errors.New("root capsule does not declare a parent")
		}
	} else if strings.TrimSpace(*parentPath) == "" {
		if *requireLineage {
			return fmt.Errorf("capsule declares parent %s but no --parent was supplied", manifest.ParentCapsuleID)
		}
	} else {
		parentManifest, parentAuthenticity, err := capsule.VerifyAuthenticated(*parentPath, verificationOptions)
		if err != nil {
			return fmt.Errorf("verify parent capsule: %w", err)
		}
		lineage, err = capsule.CheckLineage(manifest, authenticity, parentManifest, parentAuthenticity)
		if err != nil {
			return err
		}
	}

	fmt.Fprintf(stdout, "verified %s\nassets: %d\nsignature: %s\nlineage: %s\n", manifest.CapsuleID, len(manifest.Assets), authenticity.Status, lineage.Status)
	if authenticity.SignerID != "" {
		fmt.Fprintf(stdout, "signer: %s\n", authenticity.SignerID)
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
  arkmesh pack --name NAME --out DIR --asset role=/path/to/file [--asset ...] [--signing-key PRIVATE_IDENTITY] [--parent PARENT_CAPSULE]
  arkmesh inspect DIR
  arkmesh verify [--trust PUBLIC_IDENTITY] [--require-signature] [--require-trusted] [--parent PARENT_CAPSULE] [--require-lineage] DIR

Integrity verification remains available for unsigned capsules. Descendants require the same signer as their parent. Trust is local and explicit.`)
}
