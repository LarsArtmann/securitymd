package main

import (
	"context"
	"fmt"
	"time"

	cmdguard "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"

	"github.com/LarsArtmann/securitymd/pkg/policy"
	"github.com/fatih/color"
	finding "github.com/larsartmann/go-finding"
)

// statusFlags carries everything `securitymd status` accepts.
type statusFlags struct {
	Location string `flag:"location" default:"root" help:"Preferred canonical policy location (root, .github, docs)"`
}

func registerStatusCmd(cli *cmdguard.CLI[appConfig]) error {
	cmd, err := cmdguard.NewCommand(
		"status",
		statusFlags{},
		runStatus,
		cmdguard.WithShort("Show security policy status"),
		cmdguard.WithLong(`Show where the repository's SECURITY.md lives and how compliant it is.`),
		cmdguard.WithNoArgs(),
	)
	if err != nil {
		return fmt.Errorf("building status command: %w", err)
	}

	return cmdguard.AddCommand(cli, cmd)
}

func runStatus(ctx context.Context, _ *appConfig, flags statusFlags) error {
	candidates, err := policy.OrderedCandidates(flags.Location)
	if err != nil {
		return fmt.Errorf("invalid --location: %w", err)
	}

	report, err := policy.ReportIn(ctx, candidates)
	if err != nil {
		return fmt.Errorf("status failed: %w", err)
	}

	color.White("📋 Security Policy Status")
	color.White("========================")

	findings := report.FindingsSnapshot()

	if len(findings) == 0 {
		color.Green("✅ SECURITY.md found and compliant")

		return nil
	}

	now := time.Now()
	active := finding.Filter(findings, func(f finding.Finding) bool {
		return !f.IsSuppressedAt(now)
	})
	suppressedCount := len(findings) - len(active)
	errors := finding.Filter(active, finding.BySeverityAtLeast(finding.SeverityError))
	warnings := finding.Filter(active, finding.BySeverity(finding.SeverityWarning))

	color.Red("❌ %d error(s), %d warning(s)", len(errors), len(warnings))
	if suppressedCount > 0 {
		color.Cyan("🔇 %d finding(s) suppressed in-file (not counted)", suppressedCount)
	}

	for _, f := range findings {
		switch {
		case f.IsSuppressedAt(now):
			color.Cyan("  🔇 [%s] %s (suppressed)", f.Rule, f.Message)
		case f.Severity.GreaterThanOrEqual(finding.SeverityError):
			color.Red("  • [%s] %s", f.Rule, f.Message)
		default:
			color.Yellow("  • [%s] %s", f.Rule, f.Message)
		}
	}

	fmt.Println()
	color.Cyan("Next steps:")
	color.Cyan("  1. securitymd setup    — generate a policy if missing")
	color.Cyan("  2. securitymd validate — full compliance check")

	return nil
}
