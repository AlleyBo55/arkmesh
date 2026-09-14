package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"
	"time"

	"arkmesh/internal/capsule"
)

func runAudit(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return errors.New("usage: arkmesh audit <sample|history>")
	}
	switch args[0] {
	case "sample":
		flags := flag.NewFlagSet("audit sample", flag.ContinueOnError)
		flags.SetOutput(stderr)
		treePath := flags.String("tree", "", "authenticated chunk tree file")
		tolerance := flags.Float64("tolerance", 0.01, "fraction of the object you refuse to lose unnoticed")
		confidence := flags.Float64("confidence", 0.99, "required statistical confidence")
		samples := flags.Int("samples", 0, "explicit sample count, overriding tolerance")
		logPath := flags.String("log", "", "append the result to a retention log")
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		if flags.NArg() != 1 || strings.TrimSpace(*treePath) == "" {
			return errors.New("usage: arkmesh audit sample --tree FILE [--tolerance F] [--confidence F] [--samples N] [--log FILE] CAPSULE")
		}
		asset, tree, objectPath, err := loadRepairInputs(flags.Arg(0), *treePath)
		if err != nil {
			return err
		}
		result, err := capsule.AuditAsset(objectPath, asset, tree, capsule.AuditPolicy{
			Tolerance:  *tolerance,
			Confidence: *confidence,
			Samples:    *samples,
		}, nil)
		if err != nil {
			return err
		}
		fmt.Fprintf(stdout, "audited asset\nasset: %s\nchunks: %d\nsampled: %d\nfailed: %d\n",
			result.AssetSHA256, result.ChunkCount, result.Sampled, len(result.Failed))
		if len(result.Failed) == 0 {
			if result.Exhaustive {
				fmt.Fprintln(stdout, "statement: every chunk was checked and matched the signed tree")
			} else {
				fmt.Fprintf(stdout, "statement: at least %.2f%% of this object is present, with %.0f%% confidence\n",
					result.MinIntactFraction*100, result.Confidence*100)
				fmt.Fprintf(stdout, "bound: at most %d of %d chunks could still be missing\n", result.MaxBadChunks, result.ChunkCount)
			}
		} else {
			for _, index := range result.Failed {
				fmt.Fprintf(stdout, "failed chunk %d\n", index)
			}
			fmt.Fprintln(stdout, "statement: damage observed, so no retrievability bound is claimed")
		}
		if strings.TrimSpace(*logPath) != "" {
			record := capsule.NewAuditRecord(result, time.Now())
			if err := capsule.AppendAuditRecord(*logPath, record); err != nil {
				return err
			}
			fmt.Fprintf(stdout, "recorded: %s\n", *logPath)
		}
		if len(result.Failed) > 0 {
			return fmt.Errorf("%d of %d sampled chunks failed", len(result.Failed), result.Sampled)
		}
		return nil
	case "history":
		flags := flag.NewFlagSet("audit history", flag.ContinueOnError)
		flags.SetOutput(stderr)
		logPath := flags.String("log", "", "retention log to summarize")
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		if flags.NArg() != 0 || strings.TrimSpace(*logPath) == "" {
			return errors.New("usage: arkmesh audit history --log FILE")
		}
		records, err := capsule.LoadAuditRecords(*logPath)
		if err != nil {
			return err
		}
		summary := capsule.SummarizeAuditRecords(records)
		fmt.Fprintf(stdout, "retention log: %s\naudits: %d\nclean: %d\ndamaged: %d\nchunks sampled: %d\n",
			*logPath, summary.Records, summary.Clean, summary.Damaged, summary.TotalSampled)
		if summary.Records > 0 {
			fmt.Fprintf(stdout, "first checked: %s\nlast checked: %s\n", summary.FirstChecked, summary.LastChecked)
		}
		fmt.Fprintln(stdout, "note: this log records observed checks only. It cannot prove that checks were not omitted.")
		return nil
	default:
		return fmt.Errorf("unknown audit command %q", args[0])
	}
}
