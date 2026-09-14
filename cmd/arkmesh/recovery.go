package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strconv"
	"strings"

	"arkmesh/internal/capsule"
	"arkmesh/internal/identity"
)

func runRecovery(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return errors.New("usage: arkmesh recovery <policy|approve|assemble>")
	}
	switch args[0] {
	case "policy":
		return runRecoveryPolicy(args[1:], stdout, stderr)
	case "approve":
		flags := flag.NewFlagSet("recovery approve", flag.ContinueOnError)
		flags.SetOutput(stderr)
		policyPath := flags.String("policy", "", "recovery policy file")
		signingKey := flags.String("signing-key", "", "private recovery identity")
		parentPath := flags.String("parent", "", "checkpointed parent capsule")
		output := flags.String("out", "", "new detached approval file")
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		if flags.NArg() != 1 || strings.TrimSpace(*policyPath) == "" || strings.TrimSpace(*signingKey) == "" || strings.TrimSpace(*parentPath) == "" || strings.TrimSpace(*output) == "" {
			return errors.New("usage: arkmesh recovery approve --policy POLICY --signing-key PRIVATE_IDENTITY --parent PARENT --out FILE CHILD")
		}
		policy, parent, parentAuth, child, childAuth, err := loadRecoveryEndpoints(*policyPath, *parentPath, flags.Arg(0), identity.RevocationSet{})
		if err != nil {
			return err
		}
		privateKey, err := identity.LoadPrivate(*signingKey)
		if err != nil {
			return err
		}
		approval, err := capsule.ApproveRecovery(*output, policy, parent, parentAuth, child, childAuth, privateKey)
		if err != nil {
			return err
		}
		fmt.Fprintf(stdout, "created recovery approval\npath: %s\npolicy: %s\napprover: %s\nchild: %s\n", *output, approval.PolicyID, approval.ApproverID, approval.ChildCapsuleID)
		return nil
	case "assemble":
		flags := flag.NewFlagSet("recovery assemble", flag.ContinueOnError)
		flags.SetOutput(stderr)
		policyPath := flags.String("policy", "", "recovery policy file")
		parentPath := flags.String("parent", "", "checkpointed parent capsule")
		var approvalPaths repeatedFlags
		flags.Var(&approvalPaths, "approval", "detached recovery approval (repeatable)")
		var revocationPaths repeatedFlags
		flags.Var(&revocationPaths, "revocations", "local revocation policy file (repeatable)")
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		if flags.NArg() != 1 || strings.TrimSpace(*policyPath) == "" || strings.TrimSpace(*parentPath) == "" || len(approvalPaths) == 0 {
			return errors.New("usage: arkmesh recovery assemble --policy POLICY --parent PARENT --approval FILE [--approval ...] [--revocations FILE] CHILD")
		}
		revocations, err := loadRevocations(revocationPaths)
		if err != nil {
			return err
		}
		policy, parent, parentAuth, child, childAuth, err := loadRecoveryEndpoints(*policyPath, *parentPath, flags.Arg(0), revocations)
		if err != nil {
			return err
		}
		approvals := make([]capsule.RecoveryApproval, 0, len(approvalPaths))
		for _, path := range approvalPaths {
			approval, err := capsule.LoadRecoveryApproval(path)
			if err != nil {
				return fmt.Errorf("load recovery approval %q: %w", path, err)
			}
			approvals = append(approvals, approval)
		}
		transition, err := capsule.AssembleRecovery(flags.Arg(0), policy, parent, parentAuth, child, childAuth, approvals, revocations)
		if err != nil {
			return err
		}
		fmt.Fprintf(stdout, "assembled threshold recovery\npolicy: %s\napprovals: %d\nchild: %s\n", transition.PolicyID, len(transition.Approvals), transition.ChildCapsuleID)
		return nil
	default:
		return fmt.Errorf("unknown recovery command %q", args[0])
	}
}

func runRecoveryPolicy(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 || args[0] != "create" {
		return errors.New("usage: arkmesh recovery policy create --threshold N --identity PUBLIC_IDENTITY [--identity ...] --out FILE")
	}
	flags := flag.NewFlagSet("recovery policy create", flag.ContinueOnError)
	flags.SetOutput(stderr)
	thresholdText := flags.String("threshold", "", "required number of distinct approvals")
	output := flags.String("out", "", "new recovery policy file")
	var identityPaths repeatedFlags
	flags.Var(&identityPaths, "identity", "public recovery identity (repeatable)")
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}
	if flags.NArg() != 0 || strings.TrimSpace(*thresholdText) == "" || strings.TrimSpace(*output) == "" || len(identityPaths) < 2 {
		return errors.New("usage: arkmesh recovery policy create --threshold N --identity PUBLIC_IDENTITY [--identity ...] --out FILE")
	}
	threshold, err := strconv.Atoi(*thresholdText)
	if err != nil {
		return fmt.Errorf("invalid recovery threshold %q", *thresholdText)
	}
	members := make([]identity.PublicIdentity, 0, len(identityPaths))
	for _, path := range identityPaths {
		member, err := identity.LoadPublic(path)
		if err != nil {
			return fmt.Errorf("load recovery identity %q: %w", path, err)
		}
		members = append(members, member)
	}
	policy, err := capsule.NewRecoveryPolicy(threshold, members)
	if err != nil {
		return err
	}
	if err := capsule.WriteRecoveryPolicy(*output, policy); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "created recovery policy\npath: %s\npolicy: %s\nthreshold: %d of %d\n", *output, policy.PolicyID, policy.Threshold, len(policy.Members))
	return nil
}

func loadRecoveryEndpoints(policyPath, parentPath, childPath string, revocations identity.RevocationSet) (capsule.RecoveryPolicy, capsule.Manifest, capsule.Authenticity, capsule.Manifest, capsule.Authenticity, error) {
	policy, err := capsule.LoadRecoveryPolicy(policyPath)
	if err != nil {
		return capsule.RecoveryPolicy{}, capsule.Manifest{}, capsule.Authenticity{}, capsule.Manifest{}, capsule.Authenticity{}, err
	}
	parent, parentAuth, err := capsule.VerifyAuthenticated(parentPath, capsule.VerifyOptions{Revocations: revocations, RequireSignature: true})
	if err != nil {
		return capsule.RecoveryPolicy{}, capsule.Manifest{}, capsule.Authenticity{}, capsule.Manifest{}, capsule.Authenticity{}, fmt.Errorf("verify recovery parent: %w", err)
	}
	child, childAuth, err := capsule.VerifyAuthenticated(childPath, capsule.VerifyOptions{Revocations: revocations, RequireSignature: true})
	if err != nil {
		return capsule.RecoveryPolicy{}, capsule.Manifest{}, capsule.Authenticity{}, capsule.Manifest{}, capsule.Authenticity{}, fmt.Errorf("verify recovery child: %w", err)
	}
	return policy, parent, parentAuth, child, childAuth, nil
}
