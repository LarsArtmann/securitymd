package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/LarsArtmann/template-SECURITY/v2/internal/domain"
	"github.com/LarsArtmann/template-SECURITY/v2/internal/types"
	"github.com/LarsArtmann/template-SECURITY/v2/internal/errors"
)

// validationService implements ValidationService interface
type validationService struct {
	errorHandler ErrorHandler
}

// NewValidationService creates a new validation service
func NewValidationService() ValidationService {
	errorHandler := NewErrorHandler(true)
	return &validationService{
		errorHandler: errorHandler,
	}
}

// ValidateSecurityPolicy validates a security policy file
func (v *validationService) ValidateSecurityPolicy(ctx context.Context, filePath string) (domain.Validation, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return domain.Validation{
			Status:    types.StatusInvalid,
			Score:     0,
			Issues:    []domain.ValidationIssue{},
			Warnings:  []domain.ValidationWarning{},
			CheckedAt: time.Now(),
			Quality: domain.QualityMetrics{},
		}, errors.NewFileOperationError(
			"FILE001",
			"Failed to read security policy file",
			"read",
			filePath,
		)
	}

	return v.validateContent(string(content), filePath)
}

// ValidateAllPolicies validates all policy files in directory
func (v *validationService) ValidateAllPolicies(ctx context.Context) ([]domain.Validation, error) {
	policyFiles := []string{
		"SECURITY.md",
		"security-policy.md",
		"enterprise-policy.md",
		"bug-bounty-program.md",
	}

	var results []domain.Validation

	for _, file := range policyFiles {
		if _, err := os.Stat(file); err == nil {
			validation, err := v.ValidateSecurityPolicy(ctx, file)
			if err != nil {
				v.errorHandler.Handle(err)
				// Add failed validation
				validation = domain.Validation{
					Status:    types.StatusInvalid,
					Score:     0,
					Issues: []domain.ValidationIssue{
						{
							ID:          "VAL001",
							Field:       file,
							Category:    domain.CategoryStructure,
							Level:       types.ValidationLevelError,
							Message:     "Failed to read file",
							Suggestion:  "Check file permissions",
						},
					},
					Warnings:  []domain.ValidationWarning{},
					CheckedAt: time.Now(),
					Quality: domain.QualityMetrics{},
				}
			}
			results = append(results, validation)
		}
	}

	return results, nil
}

// GetValidationResults retrieves validation history
func (v *validationService) GetValidationResults(ctx context.Context, policyID domain.PolicyID) ([]domain.Validation, error) {
	// For now, return empty slice
	// In a real implementation, this would read from persistence layer
	return []domain.Validation{}, nil
}

// GenerateValidationReport generates comprehensive validation report
func (v *validationService) GenerateValidationReport(ctx context.Context, results []domain.Validation) (ValidationReport, error) {
	summary := ValidationSummary{}
	totalScore := 0

	for _, result := range results {
		switch result.Status {
		case types.StatusValid:
			summary.TotalValid++
		case types.StatusInvalid:
			summary.TotalInvalid++
		case types.StatusWarning:
			summary.TotalWarnings++
		}

		totalScore += int(result.Score)
		
		if result.Score > summary.HighestScore {
			summary.HighestScore = result.Score
		}
		if summary.LowestScore == 0 || result.Score < summary.LowestScore {
			summary.LowestScore = result.Score
		}
	}

	if len(results) > 0 {
		summary.AverageScore = float64(totalScore) / float64(len(results))
	}

	metrics := ValidationMetrics{
		TotalPoliciesValidated: len(results),
		ValidationScore: map[types.ValidationLevel]int{},
		CommonIssues:          v.getCommonIssues(results),
		ImprovementSuggestions: v.getImprovementSuggestions(results),
	}

	report := ValidationReport{
		Summary:     summary,
		Results:     results,
		Metrics:     metrics,
		GeneratedAt: time.Now(),
	}

	return report, nil
}

// Private validation methods

func (v *validationService) validateContent(content, filePath string) domain.Validation {
	issues := []domain.ValidationIssue{}
	warnings := []domain.ValidationWarning{}
	score := uint8(100)

	// Basic content checks
	if len(strings.TrimSpace(content)) < 100 {
		issues = append(issues, domain.ValidationIssue{
			ID:          "VAL002",
			Field:       "content",
			Category:    domain.CategoryContent,
			Level:       types.ValidationLevelError,
			Message:     "Content is too short",
			Suggestion:  "Add comprehensive security policy content",
		})
		score -= 30
	}

	// Required sections
	requiredSections := []string{
		"# Security Policy",
		"## Reporting a Vulnerability",
		"## Supported Versions",
		"## Security Contact Information",
	}

	for _, section := range requiredSections {
		if !strings.Contains(content, section) {
			issues = append(issues, domain.ValidationIssue{
				ID:          "VAL003",
				Field:       "sections",
				Category:    domain.CategoryStructure,
				Level:       types.ValidationLevelError,
				Message:     fmt.Sprintf("Missing required section: %s", section),
				Suggestion:  fmt.Sprintf("Add the '%s' section", section),
			})
			score -= 20
		}
	}

	// Contact information checks
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	if !emailRegex.MatchString(content) {
		issues = append(issues, domain.ValidationIssue{
			ID:          "VAL004",
			Field:       "contact",
			Category:    domain.CategoryContact,
			Level:       types.ValidationLevelError,
			Message:     "No valid email address found",
			Suggestion:  "Add a valid email address for security reports",
		})
		score -= 25
	}

	// Template variable checks
	variableRegex := regexp.MustCompile(`\{\{[^}]+\}\}`)
	variables := variableRegex.FindAllString(content, -1)
	
	unresolvedVariables := []string{}
	for _, variable := range variables {
		if v.isCommonUnresolvedVariable(variable) {
			unresolvedVariables = append(unresolvedVariables, variable)
		}
	}

	if len(unresolvedVariables) > 0 {
		warnings = append(warnings, domain.ValidationWarning{
			ID:          "VAL005",
			Field:       "variables",
			Category:    domain.CategoryContent,
			Level:       types.ValidationLevelWarning,
			Message:     fmt.Sprintf("Unresolved template variables: %s", strings.Join(unresolvedVariables, ", ")),
			Suggestion:  "Replace template variables with actual values",
		})
		score -= 10
	}

	// Content quality checks
	lines := strings.Split(content, "\n")
	if len(lines) < 20 {
		warnings = append(warnings, domain.ValidationWarning{
			ID:          "VAL006",
			Field:       "content",
			Category:    domain.CategoryContent,
			Level:       types.ValidationLevelWarning,
			Message:     "Content is shorter than recommended minimum",
			Suggestion:  "Add more detail to your security policy",
		})
		score -= 5
	}

	// Check for placeholder text
	placeholderTexts := []string{"example.com", "yourcompany.com", "replace with", "TODO", "XXX"}
	for _, placeholder := range placeholderTexts {
		if strings.Contains(strings.ToLower(content), placeholder) {
			warnings = append(warnings, domain.ValidationWarning{
				ID:          "VAL007",
				Field:       "content",
				Category:    domain.CategoryContent,
				Level:       types.ValidationLevelWarning,
				Message:     fmt.Sprintf("Found placeholder text: %s", placeholder),
				Suggestion:  "Replace placeholder text with actual values",
			})
			score -= 5
		}
	}

	// Determine final status
	status := types.StatusValid
	if score < 50 {
		status = types.StatusInvalid
	} else if score < 80 || len(warnings) > 0 {
		status = types.StatusWarning
	}

	if score < 0 {
		score = 0
	}

	// Calculate quality metrics
	quality := domain.QualityMetrics{
		Words:         uint16(len(strings.Fields(content))),
		Lines:         uint16(len(lines)),
		Characters:    uint16(len(content)),
		Sections:      uint8(v.countSections(content)),
		Readability:   uint8(v.calculateReadability(content)),
		Completeness:  uint8(v.calculateCompleteness(content)),
		BestPractices: uint8(v.calculateBestPractices(content)),
	}

	return domain.Validation{
		Status:    status,
		Score:     score,
		Issues:    issues,
		Warnings:  warnings,
		CheckedAt: time.Now(),
		Quality:    quality,
	}
}

func (v *validationService) isCommonUnresolvedVariable(variable string) bool {
	commonUnresolved := []string{
		"{{ORGANIZATION}}",
		"{{CONTACT_EMAIL}}",
		"{{DOMAIN}}",
		"{{PROJECT_NAME}}",
		"{{CURRENT_DATE}}",
	}

	variable = strings.TrimSpace(variable)
	for _, unresolved := range commonUnresolved {
		if variable == unresolved {
			return true
		}
	}
	return false
}

func (v *validationService) countSections(content string) int {
	// Count markdown headers (#, ##, ###)
	headerRegex := regexp.MustCompile(`^#{1,6}\s+`)
	lines := strings.Split(content, "\n")
	count := 0

	for _, line := range lines {
		if headerRegex.MatchString(strings.TrimSpace(line)) {
			count++
		}
	}

	return count
}

func (v *validationService) calculateReadability(content string) uint8 {
	// Simple readability score based on content length and structure
	lines := strings.Split(content, "\n")
	score := 50

	// Penalize very short or very long content
	if len(content) < 500 {
		score -= 20
	} else if len(content) > 10000 {
		score -= 10
	}

	// Reward structured content
	sections := v.countSections(content)
	score += sections * 2

	// Penalize extremely long paragraphs
	paragraphs := strings.Split(content, "\n\n")
	for _, paragraph := range paragraphs {
		if len(paragraph) > 1000 {
			score -= 5
		}
	}

	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return uint8(score)
}

func (v *validationService) calculateCompleteness(content string) uint8 {
	score := 0
	totalChecks := 0

	// Check for various sections and content types
	checks := map[string]bool{
		"vulnerability reporting":  strings.Contains(content, "Reporting a Vulnerability"),
		"supported versions":      strings.Contains(content, "Supported Versions"),
		"contact information":     strings.Contains(content, "Contact"),
		"security practices":     strings.Contains(content, "Security Practices"),
		"email contact":          strings.Contains(content, "@"),
		"response time":          strings.Contains(content, "response") || strings.Contains(content, "Response"),
		"supported versions":      strings.Contains(content, "v1.") || strings.Contains(content, "v2."),
	}

	for _, check := range checks {
		totalChecks++
		if check {
			score++
		}
	}

	if totalChecks > 0 {
		return uint8((score * 100) / totalChecks)
	}
	return 0
}

func (v *validationService) calculateBestPractices(content string) uint8 {
	score := 0
	totalChecks := 0

	// Check for best practices
	bpChecks := map[string]bool{
		"private disclosure":     strings.Contains(content, "private") && strings.Contains(content, "disclosure"),
		"safe harbor":           strings.Contains(content, "Safe Harbor") || strings.Contains(content, "safe harbor"),
		"timeline":              strings.Contains(content, "timeline") || strings.Contains(content, "Timeline"),
		"severity":              strings.Contains(content, "severity") || strings.Contains(content, "Severity"),
		"bounty program":        strings.Contains(content, "bounty") || strings.Contains(content, "Bounty"),
		"coordinate disclosure": strings.Contains(content, "coordinate") || strings.Contains(content, "coordinate"),
	}

	for _, check := range bpChecks {
		totalChecks++
		if check {
			score++
		}
	}

	if totalChecks > 0 {
		return uint8((score * 100) / totalChecks)
	}
	return 0
}

func (v *validationService) getCommonIssues(results []domain.Validation) []string {
	issues := make(map[string]int)

	for _, result := range results {
		for _, issue := range result.Issues {
			issues[issue.Message]++
		}
		for _, warning := range result.Warnings {
			issues[warning.Message]++
		}
	}

	var commonIssues []string
	for issue, count := range issues {
		if count > 1 {
			commonIssues = append(commonIssues, fmt.Sprintf("%s (%d occurrences)", issue, count))
		}
	}

	return commonIssues
}

func (v *validationService) getImprovementSuggestions(results []domain.Validation) []string {
	suggestions := []string{}

	// Analyze common patterns across all validations
	hasContactIssues := false
	hasSectionIssues := false
	hasVariableIssues := false

	for _, result := range results {
		for _, issue := range result.Issues {
			switch issue.Category {
			case domain.CategoryContact:
				hasContactIssues = true
			case domain.CategoryStructure:
				hasSectionIssues = true
			case domain.CategoryContent:
				if strings.Contains(issue.Message, "variable") {
					hasVariableIssues = true
				}
			}
		}
	}

	if hasContactIssues {
		suggestions = append(suggestions, "Ensure proper contact information is included in all security policies")
	}

	if hasSectionIssues {
		suggestions = append(suggestions, "Add all required sections following security policy best practices")
	}

	if hasVariableIssues {
		suggestions = append(suggestions, "Replace all template variables with actual values before finalizing policies")
	}

	return suggestions
}