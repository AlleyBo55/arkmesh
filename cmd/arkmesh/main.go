package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"arkmesh/internal/capsule"
)

type assetFlags []string

func (values *assetFlags) String() string {
	return strings.Join(*values, ",")
}

func (values *assetFlags) Set(value string) error {
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
	case "pack":
		err = runPack(args[1:], stdout, stderr)
	case "inspect":
		err = runInspect(args[1:], stdout)
	case "verify":
		err = runVerify(args[1:], stdout)
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

func runPack(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("pack", flag.ContinueOnError)
	flags.SetOutput(stderr)
	name := flags.String("name", "", "human-readable capsule name")
	output := flags.String("out", "", "output capsule directory")
	var assets assetFlags
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

	sources := make([]capsule.AssetSource, 0, len(assets))
	for _, value := range assets {
		role, path, ok := strings.Cut(value, "=")
		if !ok || strings.TrimSpace(role) == "" || strings.TrimSpace(path) == "" {
			return fmt.Errorf("invalid --asset %q; expected role=/path/to/file", value)
		}
		sources = append(sources, capsule.AssetSource{Role: role, Path: path})
	}

	manifest, err := capsule.Pack(*name, *output, sources, time.Now())
	if err != nil {
		return err
	}
	fmt.Fprintf(stdout, "packed %s\nassets: %d\npath: %s\n", manifest.CapsuleID, len(manifest.Assets), *output)
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

func runVerify(args []string, stdout io.Writer) error {
	if len(args) != 1 {
		return errors.New("usage: arkmesh verify <capsule-directory>")
	}
	manifest, err := capsule.Verify(args[0])
	if err != nil {
		return err
	}
	fmt.Fprintf(stdout, "verified %s\nassets: %d\n", manifest.CapsuleID, len(manifest.Assets))
	return nil
}

func printUsage(writer io.Writer) {
	fmt.Fprintln(writer, `ArkMesh v0alpha1

Usage:
  arkmesh pack --name NAME --out DIR --asset role=/path/to/file [--asset ...]
  arkmesh inspect DIR
  arkmesh verify DIR

This version proves local content integrity only. It does not authenticate authors.`)
}
