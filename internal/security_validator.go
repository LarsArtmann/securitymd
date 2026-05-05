package internal

import (
	"fmt"
	"os"
	"strings"

	finding "github.com/larsartmann/go-finding"
)

const toolName = "template-security"

// SecurityValidator validates security policy files.
type SecurityValidator struct{}

// NewSecurityValidator creates a new security validator.
func NewSecurityValidator() *SecurityValidator {
	return &SecurityValidator{}
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

// makeFinding creates a Finding for a validation issue at the given file and line.
func makeFinding(rule, message string, severity finding.Severity, file string, line int) finding.Finding {
	return finding.NewFinding(
		rule,
		toolName,
		message,
		severity,
		finding.Pos(file, line, 0),
		1.0,
	)
}

// ValidateSECURITYMd validates a SECURITY.md file and returns a Report containing findings.
func (sv *SecurityValidator) ValidateSECURITYMd(
	filePath string,
) (*finding.Report, error) {
	report := finding.NewReport(finding.ToolInfo{Name: toolName})

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")

	requiredSections := []struct {
		pattern  string
		level    string
		rule     string
		message  string
	}{
		{"# Security Policy", "error", "missing-header", "Missing Security Policy header"},
		{"## Reporting a Vulnerability", "error", "missing-reporting", "Missing 'Reporting a Vulnerability' section"},
		{"## Supported Versions", "error", "missing-versions", "Missing 'Supported Versions' section"},
		{"## Security Practices", "error", "missing-practices", "Missing 'Security Practices' section"},
		{"@", "error", "missing-contact", "Missing contact email address"},
		{"response time", "warning", "missing-response-time", "Should specify response time for vulnerability reports"},
	}

	for _, section := range requiredSections {
		if !lineChecker(lines, section.pattern) {
			severity := finding.SeverityWarning
			if section.level == "error" {
				severity = finding.SeverityError
			}

			f := makeFinding(section.rule, section.message, severity, filePath, 0)
			f.Category = finding.CategorySecurity

			report.AddFinding(f)
		}
	}

	sv.validateContentQuality(filePath, contentStr, report)

	return report, nil
}

// validateContentQuality checks the quality of the content and adds findings.
func (sv *SecurityValidator) validateContentQuality(
	filePath string,
	content string,
	report *finding.Report,
) {
	const (
		minLines      = 20
		minLineLength = 20
	)

	lines := strings.Split(content, "\n")

	if len(lines) < minLines {
		f := makeFinding("too-short", "SECURITY.md seems too short (< 20 lines)", finding.SeverityWarning, filePath, 0)
		f.Category = finding.CategorySecurity
		report.AddFinding(f)
	}

	for i, line := range lines {
		if strings.Contains(line, "{{") && strings.Contains(line, "}}") {
			f := makeFinding(
				"unresolved-template",
				"Unresolved template variable: "+strings.TrimSpace(line),
				finding.SeverityError,
				filePath,
				i+1,
			)
			f.Category = finding.CategorySecurity
			report.AddFinding(f)
		}
	}

	if !hasActualContent(lines, minLineLength) {
		f := makeFinding("no-content", "SECURITY.md lacks substantive content", finding.SeverityError, filePath, 0)
		f.Category = finding.CategorySecurity
		report.AddFinding(f)
	}

	if !hasVersionInformation(lines) {
		f := makeFinding("no-version-info", "No version information found", finding.SeverityWarning, filePath, 0)
		f.Category = finding.CategorySecurity
		report.AddFinding(f)
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
	for _, f := range report.Findings {
		if f.Severity == finding.SeverityError {
			return false
		}
	}

	return true
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
		if len(report.Findings) > 0 {
			file = report.Findings[0].Position.File
		}

		valid := ReportIsValid(report)
		status := "✅"
		if !valid {
			status = "❌"
		}

		fmt.Printf("%s %s\n", status, file)

		for _, f := range report.Findings {
			if f.Severity == finding.SeverityError {
				fmt.Printf("  ❌ Error: %s\n", f.Message)
			} else {
				fmt.Printf("  ⚠️  Warning: %s\n", f.Message)
			}
		}

		if len(report.Findings) > 0 {
			fmt.Println()
		}
	}
}
