package internal

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/LarsArtmann/template-CLI/pkg/sdk/fileops"
	"github.com/LarsArtmann/template-CLI/pkg/sdk/vfs"
)

// PolicyValidator validates security policies against best practices
type PolicyValidator struct {
	fileSystem vfs.FileSystem
	fileOps    *fileops.Service
}

// NewPolicyValidator creates new policy validator
func NewPolicyValidator() *PolicyValidator {
	vfsImpl := vfs.NewOSVFS()
	fileOps := fileops.NewService(fileops.ServiceOptions{
		BaseDirectory: ".",
		ValidatePaths: true,
	})

	return &PolicyValidator{
		fileSystem: vfsImpl,
		fileOps:    fileOps,
	}
}

// ValidationResult represents validation result
type ValidationResult struct {
	Valid           bool              `json:"valid"`
	Issues          []ValidationIssue `json:"issues"`
	Score           float64           `json:"score"`
	Recommendations []string          `json:"recommendations"`
}

// ValidationIssue represents a specific validation issue
type ValidationIssue struct {
	Type        string `json:"type"`     // "error", "warning", "info"
	Category    string `json:"category"` // "content", "structure", "security"
	Description string `json:"description"`
	Suggestion  string `json:"suggestion"`
	Line        int    `json:"line,omitempty"`
}

// ValidateSecurityPolicy validates security policy content
func (pv *PolicyValidator) ValidateSecurityPolicy(content string) ValidationResult {
	result := ValidationResult{
		Valid:  true,
		Issues: make([]ValidationIssue, 0),
		Score:  100.0,
	}

	// Check for essential sections
	sections := []string{
		"purpose",
		"policy statement",
		"scope",
		"roles and responsibilities",
		"security controls",
		"compliance",
	}

	for _, section := range sections {
		if !strings.Contains(strings.ToLower(content), section) {
			result.Issues = append(result.Issues, ValidationIssue{
				Type:        "warning",
				Category:    "content",
				Description: fmt.Sprintf("Missing section: %s", section),
				Suggestion:  fmt.Sprintf("Add section about %s", section),
			})
			result.Score -= 10
		}
	}

	// Check for contact information
	if !strings.Contains(strings.ToLower(content), "contact") &&
		!strings.Contains(strings.ToLower(content), "email") {
		result.Issues = append(result.Issues, ValidationIssue{
			Type:        "error",
			Category:    "content",
			Description: "Missing contact information",
			Suggestion:  "Add security contact email or phone",
		})
		result.Score -= 20
	}

	// Check for version information
	if !strings.Contains(strings.ToLower(content), "version") {
		result.Issues = append(result.Issues, ValidationIssue{
			Type:        "info",
			Category:    "content",
			Description: "No version information found",
			Suggestion:  "Add version or date information",
		})
		result.Score -= 5
	}

	// Check for review schedule
	if !strings.Contains(strings.ToLower(content), "review") &&
		!strings.Contains(strings.ToLower(content), "update") {
		result.Issues = append(result.Issues, ValidationIssue{
			Type:        "warning",
			Category:    "content",
			Description: "No review schedule found",
			Suggestion:  "Add policy review schedule",
		})
		result.Score -= 10
	}

	if result.Score < 0 {
		result.Score = 0
	}

	result.Valid = result.Score >= 70

	return result
}

// ValidateIncidentResponse validates incident response plan
func (pv *PolicyValidator) ValidateIncidentResponse(content string) ValidationResult {
	result := ValidationResult{
		Valid:  true,
		Issues: make([]ValidationIssue, 0),
		Score:  100.0,
	}

	// Check for incident classification
	if !strings.Contains(strings.ToLower(content), "classification") &&
		!strings.Contains(strings.ToLower(content), "severity") {
		result.Issues = append(result.Issues, ValidationIssue{
			Type:        "error",
			Category:    "content",
			Description: "Missing incident classification",
			Suggestion:  "Add incident severity classification system",
		})
		result.Score -= 20
	}

	// Check for response procedures
	responseSteps := []string{
		"detection",
		"containment",
		"eradication",
		"recovery",
		"lessons learned",
	}

	for _, step := range responseSteps {
		if !strings.Contains(strings.ToLower(content), step) {
			result.Issues = append(result.Issues, ValidationIssue{
				Type:        "warning",
				Category:    "content",
				Description: fmt.Sprintf("Missing response step: %s", step),
				Suggestion:  fmt.Sprintf("Add procedures for %s", step),
			})
			result.Score -= 10
		}
	}

	// Check for contact information
	if !pv.hasContactInfo(content) {
		result.Issues = append(result.Issues, ValidationIssue{
			Type:        "error",
			Category:    "content",
			Description: "Missing incident response contacts",
			Suggestion:  "Add 24/7 incident response contact information",
		})
		result.Score -= 25
	}

	// Check for escalation procedures
	if !strings.Contains(strings.ToLower(content), "escalation") {
		result.Issues = append(result.Issues, ValidationIssue{
			Type:        "warning",
			Category:    "content",
			Description: "Missing escalation procedures",
			Suggestion:  "Add incident escalation procedures",
		})
		result.Score -= 15
	}

	if result.Score < 0 {
		result.Score = 0
	}

	result.Valid = result.Score >= 75

	return result
}

// ValidatePrivacyPolicy validates privacy policy content
func (pv *PolicyValidator) ValidatePrivacyPolicy(content string) ValidationResult {
	result := ValidationResult{
		Valid:  true,
		Issues: make([]ValidationIssue, 0),
		Score:  100.0,
	}

	// Check for GDPR essential sections
	gdprSections := []string{
		"data controller",
		"data processor",
		"legal basis",
		"data subject rights",
		"data retention",
		"data breach",
		"international transfers",
	}

	for _, section := range gdprSections {
		if !strings.Contains(strings.ToLower(content), section) {
			result.Issues = append(result.Issues, ValidationIssue{
				Type:        "error",
				Category:    "content",
				Description: fmt.Sprintf("Missing GDPR section: %s", section),
				Suggestion:  fmt.Sprintf("Add section about %s", section),
			})
			result.Score -= 15
		}
	}

	// Check for contact information (DPO)
	if !pv.hasContactInfo(content) {
		result.Issues = append(result.Issues, ValidationIssue{
			Type:        "error",
			Category:    "content",
			Description: "Missing Data Protection Officer contact",
			Suggestion:  "Add DPO contact information",
		})
		result.Score -= 25
	}

	// Check for cookie policy
	if !strings.Contains(strings.ToLower(content), "cookie") {
		result.Issues = append(result.Issues, ValidationIssue{
			Type:        "warning",
			Category:    "content",
			Description: "Missing cookie policy",
			Suggestion:  "Add cookie usage and consent information",
		})
		result.Score -= 10
	}

	// Check for data subject rights
	rights := []string{
		"right to access",
		"right to rectification",
		"right to erasure",
		"right to portability",
		"right to object",
	}

	rightsFound := 0
	for _, right := range rights {
		if strings.Contains(strings.ToLower(content), right) {
			rightsFound++
		}
	}

	if rightsFound < len(rights) {
		result.Issues = append(result.Issues, ValidationIssue{
			Type:        "warning",
			Category:    "content",
			Description: fmt.Sprintf("Missing data subject rights (%d/%d)", rightsFound, len(rights)),
			Suggestion:  "Add comprehensive data subject rights information",
		})
		result.Score -= float64((len(rights) - rightsFound) * 5)
	}

	if result.Score < 0 {
		result.Score = 0
	}

	result.Valid = result.Score >= 80

	return result
}

// ValidateBugBountyPolicy validates bug bounty program policy
func (pv *PolicyValidator) ValidateBugBountyPolicy(content string) ValidationResult {
	result := ValidationResult{
		Valid:  true,
		Issues: make([]ValidationIssue, 0),
		Score:  100.0,
	}

	// Check for essential bug bounty sections
	sections := []string{
		"scope",
		"out of scope",
		"rewards",
		"reporting",
		"safe harbor",
		"disclosure",
	}

	for _, section := range sections {
		if !strings.Contains(strings.ToLower(content), section) {
			result.Issues = append(result.Issues, ValidationIssue{
				Type:        "warning",
				Category:    "content",
				Description: fmt.Sprintf("Missing section: %s", section),
				Suggestion:  fmt.Sprintf("Add section about %s", section),
			})
			result.Score -= 15
		}
	}

	// Check for reward structure
	if !strings.Contains(strings.ToLower(content), "reward") &&
		!strings.Contains(strings.ToLower(content), "bounty") {
		result.Issues = append(result.Issues, ValidationIssue{
			Type:        "error",
			Category:    "content",
			Description: "Missing reward structure",
			Suggestion:  "Add clear reward/bounty amounts and criteria",
		})
		result.Score -= 25
	}

	// Check for submission guidelines
	if !strings.Contains(strings.ToLower(content), "submit") &&
		!strings.Contains(strings.ToLower(content), "report") {
		result.Issues = append(result.Issues, ValidationIssue{
			Type:        "error",
			Category:    "content",
			Description: "Missing submission guidelines",
			Suggestion:  "Add clear vulnerability submission instructions",
		})
		result.Score -= 20
	}

	// Check for contact information
	if !pv.hasContactInfo(content) {
		result.Issues = append(result.Issues, ValidationIssue{
			Type:        "error",
			Category:    "content",
			Description: "Missing bug bounty contact information",
			Suggestion:  "Add security contact for vulnerability reports",
		})
		result.Score -= 20
	}

	// Check for legal safe harbor
	if !strings.Contains(strings.ToLower(content), "safe harbor") &&
		!strings.Contains(strings.ToLower(content), "legal protection") {
		result.Issues = append(result.Issues, ValidationIssue{
			Type:        "warning",
			Category:    "content",
			Description: "Missing legal safe harbor",
			Suggestion:  "Add legal protection for authorized research",
		})
		result.Score -= 15
	}

	if result.Score < 0 {
		result.Score = 0
	}

	result.Valid = result.Score >= 75

	return result
}

// ValidateGitHubSecurity validates GitHub SECURITY.md file
func (pv *PolicyValidator) ValidateGitHubSecurity(content string) ValidationResult {
	result := ValidationResult{
		Valid:  true,
		Issues: make([]ValidationIssue, 0),
		Score:  100.0,
	}

	// Check for essential GitHub security sections
	sections := []string{
		"supported versions",
		"reporting a vulnerability",
		"security practices",
		"acknowledgments",
	}

	for _, section := range sections {
		if !strings.Contains(strings.ToLower(content), section) {
			result.Issues = append(result.Issues, ValidationIssue{
				Type:        "warning",
				Category:    "content",
				Description: fmt.Sprintf("Missing section: %s", section),
				Suggestion:  fmt.Sprintf("Add section about %s", section),
			})
			result.Score -= 15
		}
	}

	// Check for vulnerability reporting
	if !strings.Contains(strings.ToLower(content), "vulnerability") {
		result.Issues = append(result.Issues, ValidationIssue{
			Type:        "error",
			Category:    "content",
			Description: "Missing vulnerability reporting information",
			Suggestion:  "Add clear vulnerability reporting instructions",
		})
		result.Score -= 30
	}

	// Check for security contact
	if !pv.hasContactInfo(content) {
		result.Issues = append(result.Issues, ValidationIssue{
			Type:        "error",
			Category:    "content",
			Description: "Missing security contact information",
			Suggestion:  "Add security contact email or form",
		})
		result.Score -= 25
	}

	// Check for supported versions
	if !strings.Contains(strings.ToLower(content), "supported versions") {
		result.Issues = append(result.Issues, ValidationIssue{
			Type:        "warning",
			Category:    "content",
			Description: "Missing supported versions information",
			Suggestion:  "Add list of supported versions and security update policy",
		})
		result.Score -= 15
	}

	// Check for response time
	if !strings.Contains(strings.ToLower(content), "response time") &&
		!strings.Contains(strings.ToLower(content), "sla") {
		result.Issues = append(result.Issues, ValidationIssue{
			Type:        "info",
			Category:    "content",
			Description: "Missing response time commitment",
			Suggestion:  "Add security vulnerability response time commitment",
		})
		result.Score -= 10
	}

	if result.Score < 0 {
		result.Score = 0
	}

	result.Valid = result.Score >= 80

	return result
}

// Helper methods

func (pv *PolicyValidator) hasContactInfo(content string) bool {
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	return emailRegex.MatchString(content) ||
		strings.Contains(strings.ToLower(content), "contact")
}

// ValidateFile validates a specific policy file
func (pv *PolicyValidator) ValidateFile(filename string) ValidationResult {
	readResult := pv.fileOps.ReadFile(filename)
	if readResult.IsError() {
		return ValidationResult{
			Valid: false,
			Issues: []ValidationIssue{
				{
					Type:        "error",
					Category:    "file",
					Description: fmt.Sprintf("Cannot read file: %s", filename),
					Suggestion:  "Check file permissions and path",
				},
			},
			Score: 0,
		}
	}

	content, _ := readResult.Get()

	switch {
	case filename == "SECURITY.md" || filename == "github-security.md":
		return pv.ValidateGitHubSecurity(content)
	case filename == "security-policy.md" || filename == "enterprise-policy.md":
		return pv.ValidateSecurityPolicy(content)
	case filename == "incident-response.md":
		return pv.ValidateIncidentResponse(content)
	case filename == "privacy-policy.md":
		return pv.ValidatePrivacyPolicy(content)
	case filename == "bug-bounty-policy.md" || filename == "bug-bounty-program.md":
		return pv.ValidateBugBountyPolicy(content)
	default:
		return ValidationResult{
			Valid: false,
			Issues: []ValidationIssue{
				{
					Type:        "error",
					Category:    "file",
					Description: fmt.Sprintf("Unknown policy file type: %s", filename),
					Suggestion:  "Use known policy file types",
				},
			},
			Score: 0,
		}
	}
}
