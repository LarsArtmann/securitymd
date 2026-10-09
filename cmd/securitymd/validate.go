package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/LarsArtmann/securitymd/pkg/policy"
	"github.com/fatih/color"
	cmdguard "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
	finding "github.com/larsartmann/go-finding"
)

const (
	outputFormatJSON  = "json"
	outputFormatSarif = "sarif"
	outputFormatText  = "text"
)

// errPolicyFindings is returned when error-severity findings remain: exit
// code 1 means "findings", not "crash".
var errPolicyFindings = errors.New("SECURITY.md validation failed: error-severity findings remain (see above)")

// validateFlags carries everything `securitymd validate` accepts. cmdguard
// parses the struct tags into pflag definitions at construction, so a typo
// in a tag fails at startup instead of at first use.
type validateFlags struct {
	File            string   `flag:"file"                       help:"Validate a specific policy file instead of discovering it"`
	Format          string   `flag:"format" default:"text"      help:"Output format (text, json, sarif)" values:"text,json,sarif"`
	MinimumSeverity string   `flag:"severity" default:"info"    help:"Minimum severity to report (info, warning, error, critical)" values:"info,warning,error,critical"`
	SetSeverity     []string `flag:"set-severity"               help:"Override a rule's severity, rule=level (repeatable, e.g. --set-severity missing-file=warning)"`
	Location        string   `flag:"location" default:"root"    help:"Preferred canonical policy location (root, .github, docs)"`
}

func registerValidateCmd(cli *cmdguard.CLI[appConfig]) error {
	cmd, err := cmdguard.NewCommand(
		"validate",
		validateFlags{},
		runValidate,
		cmdguard.WithShort("Validate security policy"),
		cmdguard.WithLong(`Validate the repository's SECURITY.md for completeness and compliance.

Checks the candidate locations (SECURITY.md, .github/SECURITY.md,
docs/SECURITY.md); a missing file is itself an error finding. Exits 1 when
error-severity findings remain.`),
		cmdguard.WithNoArgs(),
	)
	if err != nil {
		return fmt.Errorf("building validate command: %w", err)
	}

	return cmdguard.AddCommand(cli, cmd)
}

func runValidate(ctx context.Context, _ *appConfig, flags validateFlags) error {
	overrides, err := parseSeverityFlag(flags.SetSeverity)
	if err != nil {
		return err
	}

	if flags.File != "" {
		return validateFile(ctx, flags.File, overrides, flags)
	}

	return validateRepo(ctx, overrides, flags)
}

func parseSeverityFlag(specs []string) (policy.SeverityOverrides, error) {
	overrides, err := policy.ParseSeverityOverrides(strings.Join(specs, ","))
	if err != nil {
		return nil, fmt.Errorf("--set-severity: %w", err)
	}

	return overrides, nil
}

func validateFile(ctx context.Context, filename string, overrides policy.SeverityOverrides, flags validateFlags) error {
	findings, err := policy.Validate(filename)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return outputFindings(ctx, buildReport(findings, overrides), flags)
}

func validateRepo(ctx context.Context, overrides policy.SeverityOverrides, flags validateFlags) error {
	candidates, err := policy.OrderedCandidates(flags.Location)
	if err != nil {
		return fmt.Errorf("invalid --location: %w", err)
	}

	findings, err := policy.DetectIn(ctx, candidates)
	if err != nil {
		return fmt.Errorf("detection failed: %w", err)
	}

	return outputFindings(ctx, buildReport(findings, overrides), flags)
}

// buildReport applies severity overrides and aggregates the findings; the
// summary therefore always reflects the severities the user actually sees.
func buildReport(findings []finding.Finding, overrides policy.SeverityOverrides) *finding.Report {
	report := finding.NewReportFromFindings(
		finding.ToolInfo{Name: string(policy.ToolName), Version: version},
		policy.ApplySeverityOverrides(findings, overrides))

	return report
}

func outputFindings(ctx context.Context, report *finding.Report, flags validateFlags) error {
	findings := report.FindingsSnapshot()

	switch flags.Format {
	case outputFormatJSON:
		if err := report.WriteJSON(os.Stdout); err != nil {
			return fmt.Errorf("write JSON report: %w", err)
		}
	case outputFormatSarif:
		if err := report.WriteSARIFWithOpts(
			ctx,
			os.Stdout,
			finding.WithMinSeverity(parseSeverity(flags.MinimumSeverity)),
		); err != nil {
			return fmt.Errorf("write SARIF report: %w", err)
		}
	default:
		printFindings(findings)
	}

	if hasErrors(findings) {
		return errPolicyFindings
	}

	return nil
}

func printFindings(findings []finding.Finding) {
	if len(findings) == 0 {
		color.Green("✅ SECURITY.md - All checks passed")

		return
	}

	now := time.Now()

	for _, f := range findings {
		location := string(f.Position.File)
		if f.Position.Line > 0 {
			location = fmt.Sprintf("%s:%d", location, f.Position.Line)
		}

		suppressed := f.IsSuppressedAt(now)

		switch {
		case suppressed:
			color.Cyan("🔇 [%s] %s (%s)\n    🔇 suppressed: %s", f.Rule, f.Message, location, f.Suppression.Reason)
		case f.Severity == finding.SeverityError:
			color.Red("❌ [%s] %s (%s)", f.Rule, f.Message, location)
		default:
			color.Yellow("⚠️  [%s] %s (%s)", f.Rule, f.Message, location)
		}

		if !suppressed && f.Suggestion != "" {
			fmt.Printf("    💡 %s\n", f.Suggestion)
		}
	}
}

// hasErrors gates the exit code on ACTIVE findings only: suppressed findings
// stay visible in the output but never fail CI — that is the escape hatch.
// The gate is threshold-based (error or critical): an escalated finding must
// never silently pass.
func hasErrors(findings []finding.Finding) bool {
	active := finding.Filter(findings, func(f finding.Finding) bool {
		return !f.IsSuppressedAt(time.Now())
	})

	return len(finding.Filter(active, finding.BySeverityAtLeast(finding.SeverityError))) > 0
}

func parseSeverity(severity string) finding.Severity {
	parsed := finding.Severity(severity)
	if parsed.IsValid() {
		return parsed
	}

	return finding.SeverityInfo
}
