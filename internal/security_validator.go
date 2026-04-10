package internal

import (
	"fmt"
	"os"
	"strings"
)

// SecurityValidator validates security policy files.
type SecurityValidator struct{}

// NewSecurityValidator creates a new security validator.
func NewSecurityValidator() *SecurityValidator {
	return &SecurityValidator{}
}

// SecurityValidationResult represents the result of security validation.
type SecurityValidationResult struct {
	Valid    bool
	Errors   []string
	Warnings []string
	File     string
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

// firstMatchingLine returns the first line matching any of the patterns.
func firstMatchingLine(lines []string, patterns ...string) string {
	for _, line := range lines {
		for _, pattern := range patterns {
			if strings.Contains(strings.ToLower(line), strings.ToLower(pattern)) {
				return line
			}
		}
	}
	return ""
}

// ValidateSECURITYMd validates a SECURITY.md file.
func (sv *SecurityValidator) ValidateSECURITYMd(
	filePath string,
) (*SecurityValidationResult, error) {
	result := &SecurityValidationResult{
		File: filePath,
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")

	requiredSections := []struct {
		pattern string
		level   string
		message string
	}{
		{"# Security Policy", "error", "Missing Security Policy header"},
		{"## Reporting a Vulnerability", "error", "Missing 'Reporting a Vulnerability' section"},
		{"## Supported Versions", "error", "Missing 'Supported Versions' section"},
		{"## Security Practices", "error", "Missing 'Security Practices' section"},
		{"@", "error", "Missing contact email address"},
		{"response time", "warning", "Should specify response time for vulnerability reports"},
	}

	for _, section := range requiredSections {
		if !lineChecker(lines, section.pattern) {
			if section.level == "error" {
				result.Errors = append(result.Errors, section.message)
			} else {
				result.Warnings = append(result.Warnings, section.message)
			}
		}
	}

	sv.validateContentQuality(contentStr, result)

	result.Valid = len(result.Errors) == 0

	return result, nil
}

// validateContentQuality checks the quality of the content.
func (sv *SecurityValidator) validateContentQuality(
	content string,
	result *SecurityValidationResult,
) {
	const minLines = 20
	const minLineLength = 20

	lines := strings.Split(content, "\n")

	if len(lines) < minLines {
		result.Warnings = append(result.Warnings, "SECURITY.md seems too short (< 20 lines)")
	}

	for _, line := range lines {
		if strings.Contains(line, "{{") && strings.Contains(line, "}}") {
			result.Errors = append(
				result.Errors,
				"Unresolved template variable: "+strings.TrimSpace(line),
			)
		}
	}

	hasActualContent := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) > minLineLength && !strings.HasPrefix(trimmed, "#") && !strings.Contains(trimmed, "example") {
			hasActualContent = true
			break
		}
	}

	if !hasActualContent {
		result.Errors = append(result.Errors, "SECURITY.md lacks substantive content")
	}

	if !hasVersionInformation(lines) {
		result.Warnings = append(result.Warnings, "No version information found")
	}
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

// PrintResults prints validation results in a user-friendly format.
func (sv *SecurityValidator) PrintResults(results []*SecurityValidationResult) {
	totalFiles := len(results)
	validFiles := 0

	for _, result := range results {
		if result.Valid {
			validFiles++
		}
	}

	fmt.Printf("🔒 Security Validation Results\n")
	fmt.Printf("================================\n")
	fmt.Printf("Files checked: %d\n", totalFiles)
	fmt.Printf("Valid files: %d\n", validFiles)
	fmt.Printf("Invalid files: %d\n\n", totalFiles-validFiles)

	for _, result := range results {
		status := "✅"
		if !result.Valid {
			status = "❌"
		}

		fmt.Printf("%s %s\n", status, result.File)

		// Print errors
		for _, err := range result.Errors {
			fmt.Printf("  ❌ Error: %s\n", err)
		}

		// Print warnings
		for _, warning := range result.Warnings {
			fmt.Printf("  ⚠️  Warning: %s\n", warning)
		}

		if len(result.Errors) > 0 || len(result.Warnings) > 0 {
			fmt.Println()
		}
	}
}
