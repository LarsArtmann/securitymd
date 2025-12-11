package internal

import (
	"fmt"
	"os"
	"strings"
)

// SecurityValidator validates security policy files
type SecurityValidator struct{}

// NewSecurityValidator creates a new security validator
func NewSecurityValidator() *SecurityValidator {
	return &SecurityValidator{}
}

// SecurityValidationResult represents the result of security validation
type SecurityValidationResult struct {
	Valid    bool
	Errors   []string
	Warnings []string
	File     string
}

// ValidateSECURITYMd validates a SECURITY.md file
func (sv *SecurityValidator) ValidateSECURITYMd(filePath string) (*SecurityValidationResult, error) {
	result := &SecurityValidationResult{
		File: filePath,
	}

	// Read file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")

	// Check for required sections
	requiredSections := []struct {
		title   string
		pattern string
		level   string // "error" or "warning"
		message string
	}{
		{
			title:   "Security Contact Information",
			pattern: "# Security Policy",
			level:   "error",
			message: "Missing Security Policy header",
		},
		{
			title:   "Reporting a Vulnerability",
			pattern: "## Reporting a Vulnerability",
			level:   "error",
			message: "Missing 'Reporting a Vulnerability' section",
		},
		{
			title:   "Supported Versions",
			pattern: "## Supported Versions",
			level:   "error",
			message: "Missing 'Supported Versions' section",
		},
		{
			title:   "Security Practices",
			pattern: "## Security Practices",
			level:   "error",
			message: "Missing 'Security Practices' section",
		},
		{
			title:   "Security Contact",
			pattern: "@",
			level:   "error",
			message: "Missing contact email address",
		},
		{
			title:   "Response Time",
			pattern: "response time|Response Time|within",
			level:   "warning",
			message: "Should specify response time for vulnerability reports",
		},
	}

	// Check each section
	for _, section := range requiredSections {
		found := false
		for _, line := range lines {
			if strings.Contains(strings.ToLower(line), strings.ToLower(section.pattern)) {
				found = true
				break
			}
		}

		if !found {
			if section.level == "error" {
				result.Errors = append(result.Errors, section.message)
			} else {
				result.Warnings = append(result.Warnings, section.message)
			}
		}
	}

	// Check content quality
	sv.validateContentQuality(contentStr, result)

	// Determine if valid
	result.Valid = len(result.Errors) == 0

	return result, nil
}

// validateContentQuality checks the quality of the content
func (sv *SecurityValidator) validateContentQuality(content string, result *SecurityValidationResult) {
	lines := strings.Split(content, "\n")

	// Minimum content check
	if len(lines) < 20 {
		result.Warnings = append(result.Warnings, "SECURITY.md seems too short (< 20 lines)")
	}

	// Check for template variables that weren't replaced
	for _, line := range lines {
		if strings.Contains(line, "{{") && strings.Contains(line, "}}") {
			result.Errors = append(result.Errors, fmt.Sprintf("Unresolved template variable: %s", strings.TrimSpace(line)))
		}
	}

	// Check for actual content (not just placeholders)
	hasActualContent := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > 20 && !strings.HasPrefix(line, "#") && !strings.Contains(line, "example") {
			hasActualContent = true
			break
		}
	}

	if !hasActualContent {
		result.Errors = append(result.Errors, "SECURITY.md lacks substantive content")
	}

	// Check for version information
	hasVersion := false
	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), "version") &&
			(strings.Contains(line, "v") || strings.Contains(line, ".")) {
			hasVersion = true
			break
		}
	}

	if !hasVersion {
		result.Warnings = append(result.Warnings, "No version information found")
	}
}

// PrintResults prints validation results in a user-friendly format
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
