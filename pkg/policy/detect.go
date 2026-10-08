package policy

import (
	"context"
	"fmt"
	"strings"

	"github.com/larsartmann/go-finding"
	autoconfigure "github.com/larsartmann/linter-autoconfigure-sdk"
)

// RuleMissingFile is the rule ID reported when no SECURITY.md exists in any
// candidate location. Its fix strategy is direct: the provider's Repair
// generates the file.
const RuleMissingFile = finding.RuleName("missing-file")

// Detect locates the repository's SECURITY.md (root, .github/, docs/) and
// validates it. A missing file is itself an error finding with a direct fix
// strategy: the provider's Repair generates one.
func Detect(ctx context.Context) ([]finding.Finding, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("securitymd detection cancelled: %w", ctx.Err())
	default:
	}

	dir := autoconfigure.WorkingDir(ctx)

	path, found := autoconfigure.FirstExisting(dir, CandidateLocations...)
	if !found {
		f, err := missingFileFinding(dir)
		if err != nil {
			return nil, err
		}

		return []finding.Finding{f}, nil
	}

	findings, err := Validate(path)
	if err != nil {
		return nil, fmt.Errorf("validate %s: %w", path, err)
	}

	return findings, nil
}

// DetectNamed is intentionally absent: the provider builds its detector via
// finding.NamedDetectorFunc directly to keep the interface return out of this
// package's API.
// Report wraps Detect's findings in a finding.Report for the CLI's JSON and
// SARIF outputs.
func Report(ctx context.Context) (*finding.Report, error) {
	findings, err := Detect(ctx)
	if err != nil {
		return nil, err
	}

	report := finding.NewReport(finding.ToolInfo{Name: string(ToolName), Version: "dev"})
	report.AddFindings(findings)
	report.ComputeSummary()

	return report, nil
}

func missingFileFinding(dir string) (finding.Finding, error) {
	builder := finding.NewBuilder(
		RuleMissingFile,
		ToolName,
		"No SECURITY.md found (looked in "+dir+" for "+strings.Join(CandidateLocations, ", ")+")",
		finding.SeverityError,
		finding.FilePos(finding.FilePath(CandidateLocations[0])),
	).
		WithCategory(finding.CategorySecurity).
		WithTags(finding.TagSecurity).
		WithSuggestion("Run `securitymd setup` (or `buildflow --fix` with the securitymd provider) to generate one")

	// A direct fix needs its content: when the repo identity is derivable the
	// finding carries exactly what Generate would write. Otherwise the fix
	// stays suggest-only (the user must supply --organization/--repository).
	if preview, ok := renderPolicyPreview(dir); ok {
		builder = builder.WithFixStrategy(finding.FixStrategyDirect).WithAfterCode(preview)
	} else {
		builder = builder.WithFixStrategy(finding.FixStrategySuggest)
	}

	f, err := builder.Build()
	if err != nil {
		return f, fmt.Errorf("build missing-file finding: %w", err)
	}

	return f, nil
}

// renderPolicyPreview best-effort renders what Generate would write into dir;
// ok=false when the repo identity cannot be derived or rendering fails.
func renderPolicyPreview(dir string) (string, bool) {
	identity := DetectRepoIdentity(dir)
	if !identity.IsComplete() {
		return "", false
	}

	content, err := renderForIdentity(dir, identity, "")
	if err != nil {
		return "", false
	}

	return content, true
}
