package internal

import (
	"context"
	"fmt"
	"os"
	"strings"

	finding "github.com/larsartmann/go-finding"
	"github.com/larsartmann/go-finding/pipeline"
)

const toolName = "template-security"

// SecurityValidator validates security policy files.
// It implements pipeline.Detector for integration with the go-finding ecosystem.
type SecurityValidator struct {
	filePath string
}

// NewSecurityValidator creates a new security validator for the given file.
// The filePath is used by Detect() to know which file to validate.
func NewSecurityValidator() *SecurityValidator {
	return &SecurityValidator{
		filePath: "",
	}
}

// WithFile sets the file path to validate. Returns a new validator.
func (sv *SecurityValidator) WithFile(path string) *SecurityValidator {
	return &SecurityValidator{filePath: path}
}

// Name implements pipeline.Detector.
func (sv *SecurityValidator) Name() string {
	return toolName
}

// Detect implements pipeline.Detector.
func (sv *SecurityValidator) Detect(ctx context.Context) ([]finding.Finding, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("security validation cancelled: %w", ctx.Err())
	default:
	}

	report, err := sv.ValidateSECURITYMd(sv.filePath)
	if err != nil {
		return nil, err
	}

	return report.FindingsSnapshot(), nil
}

// lineChecker is a helper to find lines matching any of the patterns.
func lineChecker(lines []string, patterns ...string) bool {
	for _, line := range lines {
		for _, pattern := range patterns {
			if strings.Contains(strings.ToLower(line), strings.ToLower(pattern)) {
				return true
			}
		}
	}

	return false
}

// buildFinding creates a Finding using the Builder API with consistent defaults.
func buildFinding(
	rule, message string,
	severity finding.Severity,
	file string,
	line int,
) (finding.Finding, error) {
	result, err := finding.NewBuilder(finding.RuleName(rule), toolName, message, severity, finding.Pos(finding.FilePath(file), line, 0)).
		WithCategory(finding.CategorySecurity).
		WithTags(finding.TagSecurity).
		WithFixStrategy(finding.FixStrategySuggest).
		WithConfidence(1.0).
		Build()
	if err != nil {
		return result, fmt.Errorf(
			"failed to build finding {rule:%q file:%q line:%d message:%q severity:%v}: %w",
			rule,
			file,
			line,
			message,
			severity,
			err,
		)
	}

	return result, nil
}

// ValidateSECURITYMd validates a SECURITY.md file and returns a Report containing findings.
func (sv *SecurityValidator) ValidateSECURITYMd(
	filePath string,
) (*finding.Report, error) {
	report := finding.NewReport(finding.ToolInfo{Name: toolName, Version: "dev"})

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %q: %w", filePath, err)
	}

	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")

	sv.checkRequiredSections(filePath, lines, report)
	sv.validateContentQuality(filePath, contentStr, lines, report)

	report.ComputeSummary()

	return report, nil
}

func (sv *SecurityValidator) checkRequiredSections(
	filePath string,
	lines []string,
	report *finding.Report,
) {
	requiredSections := []struct {
		pattern string
		level   string
		rule    string
		message string
	}{
		{"# Security Policy", "error", "missing-header", "Missing Security Policy header"},
		{
			"## Reporting a Vulnerability",
			"error",
			"missing-reporting",
			"Missing 'Reporting a Vulnerability' section",
		},
		{
			"## Supported Versions",
			"error",
			"missing-versions",
			"Missing 'Supported Versions' section",
		},
		{
			"## Security Practices",
			"error",
			"missing-practices",
			"Missing 'Security Practices' section",
		},
		{"@", "error", "missing-contact", "Missing contact email address"},
		{
			"response time",
			"warning",
			"missing-response-time",
			"Should specify response time for vulnerability reports",
		},
	}

	for _, section := range requiredSections {
		if !lineChecker(lines, section.pattern) {
			severity := finding.SeverityWarning
			if section.level == "error" {
				severity = finding.SeverityError
			}

			f, err := buildFinding(section.rule, section.message, severity, filePath, 0)
			if err != nil {
				continue
			}

			report.AddFinding(f)
		}
	}
}

func (sv *SecurityValidator) addFinding(
	rule, message string,
	severity finding.Severity,
	file string,
	line int,
	report *finding.Report,
) {
	f, err := buildFinding(rule, message, severity, file, line)
	if err == nil {
		report.AddFinding(f)
	}
}

func (sv *SecurityValidator) validateContentQuality(
	filePath string,
	_ string,
	lines []string,
	report *finding.Report,
) {
	const (
		minLines      = 20
		minLineLength = 20
	)

	if len(lines) < minLines {
		sv.addFinding("too-short", "SECURITY.md seems too short (< 20 lines)",
			finding.SeverityWarning, filePath, 0, report)
	}

	for i, line := range lines {
		if strings.Contains(line, "{{") && strings.Contains(line, "}}") {
			sv.addFinding(
				"unresolved-template",
				"Unresolved template variable: "+strings.TrimSpace(line),
				finding.SeverityError,
				filePath,
				i+1,
				report,
			)
		}
	}

	if !hasActualContent(lines, minLineLength) {
		sv.addFinding("no-content", "SECURITY.md lacks substantive content",
			finding.SeverityError, filePath, 0, report)
	}

	if !hasVersionInformation(lines) {
		sv.addFinding("no-version-info", "No version information found",
			finding.SeverityWarning, filePath, 0, report)
	}
}

func hasActualContent(lines []string, minLength int) bool {
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) > minLength && !strings.HasPrefix(trimmed, "#") &&
			!strings.Contains(trimmed, "example") {
			return true
		}
	}

	return false
}

func hasVersionInformation(lines []string) bool {
	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), "version") &&
			(strings.Contains(line, "v") || strings.Contains(line, ".")) {
			return true
		}
	}

	return false
}

// ReportIsValid returns true if the report has no error-severity findings.
func ReportIsValid(report *finding.Report) bool {
	return len(finding.Filter(report.FindingsSnapshot(), finding.BySeverity(finding.SeverityError))) == 0
}

// DetectFile is a convenience function that creates a validator for a single file,
// runs Detect, and returns the findings as a Report.
func DetectFile(ctx context.Context, filePath string) (*finding.Report, error) {
	detector := &SecurityValidator{filePath: filePath}

	findings, err := detector.Detect(ctx)
	if err != nil {
		return nil, err
	}

	report := finding.NewReport(finding.ToolInfo{Name: toolName, Version: "dev"})
	report.AddFindings(findings)
	report.ComputeSummary()

	return report, nil
}

// PrintResults prints validation results in a user-friendly format.
func (sv *SecurityValidator) PrintResults(reports []*finding.Report) {
	totalFiles := len(reports)
	validFiles := 0

	for _, report := range reports {
		if ReportIsValid(report) {
			validFiles++
		}
	}

	fmt.Printf("🔒 Security Validation Results\n")
	fmt.Printf("================================\n")
	fmt.Printf("Files checked: %d\n", totalFiles)
	fmt.Printf("Valid files: %d\n", validFiles)
	fmt.Printf("Invalid files: %d\n\n", totalFiles-validFiles)

	for _, report := range reports {
		file := "<unknown>"
		if len(report.FindingsSnapshot()) > 0 {
			file = string(report.FindingsSnapshot()[0].Position.File)
		}

		valid := ReportIsValid(report)

		status := "✅"
		if !valid {
			status = "❌"
		}

		fmt.Printf("%s %s\n", status, file)

		for _, f := range report.FindingsSnapshot() {
			if f.Severity == finding.SeverityError {
				fmt.Printf("  ❌ Error: %s\n", f.Message)
			} else {
				fmt.Printf("  ⚠️  Warning: %s\n", f.Message)
			}
		}

		if len(report.FindingsSnapshot()) > 0 {
			fmt.Println()
		}
	}
}

// Compile-time interface check.
var _ pipeline.Detector = (*SecurityValidator)(nil)
