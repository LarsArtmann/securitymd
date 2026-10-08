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

// DetectNamed returns the detector under the tool's canonical name, ready for
// toolsdk.Spec.Detect.
func DetectNamed() finding.Detector {
	return finding.NamedDetectorFunc(string(ToolName), Detect)
}

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
	f, err := finding.NewBuilder(
		RuleMissingFile,
		ToolName,
		"No SECURITY.md found (looked in "+dir+" for "+strings.Join(CandidateLocations, ", ")+")",
		finding.SeverityError,
		finding.FilePos(finding.FilePath(CandidateLocations[0])),
	).
		WithCategory(finding.CategorySecurity).
		WithTags(finding.TagSecurity).
		WithFixStrategy(finding.FixStrategyDirect).
		WithSuggestion("Run `securitymd setup` (or `buildflow --fix` with the securitymd provider) to generate one").
		Build()
	if err != nil {
		return f, fmt.Errorf("build missing-file finding: %w", err)
	}

	return f, nil
}
