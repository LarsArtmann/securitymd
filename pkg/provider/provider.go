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

	"github.com/LarsArtmann/securitymd/pkg/policy"
	"github.com/larsartmann/go-finding"
	"github.com/larsartmann/go-finding/toolsdk"
)

// optionContactEmail lets a repo override the generated policy's contact
// channel via BuildFlow's tool_options (e.g. `tool_options: {securitymd:
// {contact-email: security@example.com}}`). Optional: without it the policy
// points to GitHub private vulnerability reporting.
const optionContactEmail = "contact-email"

// triggerManifests are the project indicators that activate securitymd: any
// repository carrying a dependency manifest, a nix flake, or CI workflows
// deserves a security policy.
var triggerManifests = []string{ //nolint:gochecknoglobals // declarative trigger table, read-only after init
	"**/go.mod",
	"package.json",
	"Cargo.toml",
	"pyproject.toml",
	"requirements.txt",
	"flake.nix",
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
	}},
	Detect: finding.NamedDetectorFunc(string(policy.ToolName), policy.Detect),
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

// contactEmailFromContext reads the optional contact-email tool option.
func contactEmailFromContext(ctx context.Context) string {
	values, ok := toolsdk.OptionsFromContext(ctx)
	if !ok {
		return ""
	}

	email, _ := values[optionContactEmail].(string)

	return email
}
