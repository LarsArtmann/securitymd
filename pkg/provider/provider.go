// Package provider registers securitymd as a BuildFlow DAG tool via the
// go-finding toolsdk provider contract. BuildFlow discovers this Provider
// automatically through toolsdk.All() when a consumer blank-imports this
// package:
//
//	import _ "github.com/LarsArtmann/securitymd/pkg/provider"
//
// Topology: Detect (missing or non-compliant SECURITY.md) → Repair (generate
// the embedded template when — and only when — no policy exists) → Verify
// (BuildFlow re-runs Detect and measures the finding delta).
package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/LarsArtmann/securitymd/pkg/policy"
	"github.com/larsartmann/go-finding"
	"github.com/larsartmann/go-finding/toolsdk"
)

// optionContactEmail lets a repo override the generated policy's contact
// channel via BuildFlow's tool_options (e.g. `tool_options: {securitymd:
// {contact-email: security@example.com}}`). Optional: without it the policy
// points to GitHub private vulnerability reporting.
const optionContactEmail = "contact-email"

// optionSeverityOverrides remaps finding severities per rule via tool_options
// (`tool_options: {securitymd: {severity-overrides: "missing-file=warning"}}`)
// — the fleet adoption unblocker. Format: rule=level[,rule=level...].
const optionSeverityOverrides = "severity-overrides"

// triggerManifests are the project indicators that activate securitymd: any
// repository carrying a dependency manifest, a nix flake, CI workflows, or a
// README deserves a security policy — docs-only repos activate too, since a
// published project without a policy is exactly the one that needs one.
var triggerManifests = []string{ //nolint:gochecknoglobals // declarative trigger table, read-only after init
	"**/go.mod",
	"package.json",
	"Cargo.toml",
	"pyproject.toml",
	"requirements.txt",
	"flake.nix",
	"README.md",
	".github/workflows/*.yml",
	".github/workflows/*.yaml",
}

//nolint:gochecknoglobals // toolsdk contract: specs self-register as package-level vars
var Provider = toolsdk.Register(toolsdk.Spec{
	Name: string(policy.ToolName),
	Description: "Detects a missing or non-compliant SECURITY.md and generates one from the " +
		"embedded template when absent (existing policies are never overwritten)",
	Trigger: toolsdk.AnyLanguage(append([]string{
		"SECURITY.md",
		".github/SECURITY.md",
		"docs/SECURITY.md",
	}, triggerManifests...)...),
	Inputs: append([]string{
		"SECURITY.md",
		".github/SECURITY.md",
		"docs/SECURITY.md",
	}, triggerManifests...),
	Options: []toolsdk.Option{{
		Name:    optionContactEmail,
		Kind:    toolsdk.OptionKindString,
		Default: "",
		Description: "Security contact email baked into a generated SECURITY.md (optional; " +
			"default points to GitHub private vulnerability reporting)",
	}, {
		Name:    optionSeverityOverrides,
		Kind:    toolsdk.OptionKindString,
		Default: "",
		Description: "Per-rule severity overrides as rule=level[,rule=level...] (optional; " +
			"e.g. missing-file=warning for incremental fleet adoption)",
	}},
	Detect: finding.NamedDetectorFunc(string(policy.ToolName), detectWithOptions),
	Repair: toolsdk.RepairerFunc(func(ctx context.Context) (toolsdk.RepairResult, error) {
		result, err := policy.Generate(ctx, policy.GenerateOptions{
			ContactEmail: contactEmailFromContext(ctx),
			DryRun:       toolsdk.DryRunFromContext(ctx),
		})
		if err != nil {
			return toolsdk.RepairResult{}, fmt.Errorf("securitymd repair: %w", err)
		}

		return toolsdk.RepairResult{Description: result.Description}, nil
	}),
})

// detectWithOptions runs the policy detector, applies the repo's severity
// overrides, then strips suppressed findings. The provider result feeds
// consumers' findings gates, and a suppression must gate nobody — in the CLI
// (exit-neutral) just as much as fleet-wide. Evidence stays available through
// the CLI's JSON report, which keeps the full set.
func detectWithOptions(ctx context.Context) ([]finding.Finding, error) {
	findings, err := policy.Detect(ctx)
	if err != nil {
		return nil, fmt.Errorf("securitymd detect: %w", err)
	}

	overrides, err := severityOverridesFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("securitymd %s: %w", optionSeverityOverrides, err)
	}

	return policy.ActiveFindings(policy.ApplySeverityOverrides(findings, overrides), time.Now()), nil
}

// contactEmailFromContext reads the optional contact-email tool option.
func contactEmailFromContext(ctx context.Context) string {
	values, ok := toolsdk.OptionsFromContext(ctx)
	if !ok {
		return ""
	}

	email, _ := values[optionContactEmail].(string)

	return email
}

// severityOverridesFromContext parses the optional severity-overrides tool
// option; an unset or empty option means no overrides.
func severityOverridesFromContext(ctx context.Context) (policy.SeverityOverrides, error) {
	values, ok := toolsdk.OptionsFromContext(ctx)
	if !ok {
		return policy.SeverityOverrides{}, nil
	}

	spec, _ := values[optionSeverityOverrides].(string)

	overrides, err := policy.ParseSeverityOverrides(spec)
	if err != nil {
		return nil, fmt.Errorf("parse option: %w", err)
	}

	return overrides, nil
}
