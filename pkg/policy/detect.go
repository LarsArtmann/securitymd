package policy

import (
	"context"
	"errors"
	"fmt"

	"github.com/larsartmann/go-finding"
	"github.com/larsartmann/linter-autoconfigure-sdk/autoconfigure"
)

// ErrMissingFile is not an error to return — it is the rule name Detect uses
// when no policy file exists. Kept as a documented constant.
const ruleMissingFile = finding.RuleName("missing-file")

// Detect locates the repository's SECURITY.md (root, .github/, docs/) and
// validates it. A missing file is itself an error finding with a direct fix
// strategy: the provider's Repair generates one. Implements the
// finding.Detector contract via finding.NamedDetectorFunc.
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

func missingFileFinding(dir string) (finding.Finding, error) {
	f, err := finding.NewBuilder(
		ruleMissingFile,
		ToolName,
		"No SECURITY.md found (looked in "+dir+" at "+joinCandidates()+")",
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

// DetectNamed returns the detector under the tool's canonical name, ready for
// toolsdk.Spec.Detect and pipeline wiring.
func DetectNamed() finding.Detector {
	return finding.NamedDetectorFunc(string(ToolName), Detect)
}

// Report wraps Detect's findings in a finding.Report for the CLI's JSON and
// SARIF outputs.
func Report(ctx context.Context) (*finding.Report, error) {
	findings, err := Detect(ctx)
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}

		return nil, err
	}

	report := finding.NewReport(finding.ToolInfo{Name: ToolName, Version: "dev"})
	report.AddFindings(findings)
	report.ComputeSummary()

	return report, nil
}
