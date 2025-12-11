package internal

import (
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// TemplateVariable represents a template variable with metadata
type TemplateVariable struct {
	Name        string
	Description string
	Value       string
	Priority    int    // 1=high, 2=medium, 3=low
	Category    string // "project", "organization", "contact", "security", "legal"
	Detected    bool
}

// TemplateVariableDetector detects and manages template variables
type TemplateVariableDetector struct {
	detector          *ProjectDetector
	githubIntegration *GitHubIntegration
}

// NewTemplateVariableDetector creates a new template variable detector
func NewTemplateVariableDetector() *TemplateVariableDetector {
	return &TemplateVariableDetector{
		detector:          NewProjectDetector(),
		githubIntegration: NewGitHubIntegration(),
	}
}

// DetectAllVariables detects all available template variables
func (tvd *TemplateVariableDetector) DetectAllVariables() map[string]TemplateVariable {
	variables := make(map[string]TemplateVariable)

	// Detect project variables
	projectVars := tvd.detectProjectVariables()
	for k, v := range projectVars {
		variables[k] = v
	}

	// Detect organization variables
	orgVars := tvd.detectOrganizationVariables()
	for k, v := range orgVars {
		variables[k] = v
	}

	// Detect contact variables
	contactVars := tvd.detectContactVariables()
	for k, v := range contactVars {
		variables[k] = v
	}

	// Detect security variables
	securityVars := tvd.detectSecurityVariables()
	for k, v := range securityVars {
		variables[k] = v
	}

	// Detect legal variables
	legalVars := tvd.detectLegalVariables()
	for k, v := range legalVars {
		variables[k] = v
	}

	// Detect temporal variables
	temporalVars := tvd.detectTemporalVariables()
	for k, v := range temporalVars {
		variables[k] = v
	}

	return variables
}

// detectProjectVariables detects project-specific variables
func (tvd *TemplateVariableDetector) detectProjectVariables() map[string]TemplateVariable {
	projectName := tvd.detector.DetectProjectName()

	variables := map[string]TemplateVariable{
		"{{PROJECT_NAME}}": {
			Name:        "{{PROJECT_NAME}}",
			Description: "Project or repository name",
			Value:       projectName,
			Priority:    1,
			Category:    "project",
			Detected:    projectName != "MyProject",
		},
		"{{REPOSITORY_NAME}}": {
			Name:        "{{REPOSITORY_NAME}}",
			Description: "Repository name from git",
			Value:       tvd.getRepositoryName(),
			Priority:    1,
			Category:    "project",
			Detected:    tvd.getRepositoryName() != "",
		},
	}

	// Detect project type
	projectType := tvd.detectProjectType()
	if projectType != "" {
		variables["{{PROJECT_TYPE}}"] = TemplateVariable{
			Name:        "{{PROJECT_TYPE}}",
			Description: "Type of project (web, api, library, etc.)",
			Value:       projectType,
			Priority:    2,
			Category:    "project",
			Detected:    true,
		}
	}

	// Detect technology stack
	techStack := tvd.detectTechnologyStack()
	for tech, info := range techStack {
		placeholder := "{{TECH_STACK_" + strings.ToUpper(tech) + "}}"
		variables[placeholder] = TemplateVariable{
			Name:        placeholder,
			Description: "Technology in stack: " + tech,
			Value:       info,
			Priority:    2,
			Category:    "project",
			Detected:    true,
		}
	}

	return variables
}

// detectOrganizationVariables detects organization-specific variables
func (tvd *TemplateVariableDetector) detectOrganizationVariables() map[string]TemplateVariable {
	org := tvd.detector.DetectOrganization()
	domain := tvd.detector.DetectDomain()

	variables := map[string]TemplateVariable{
		"{{ORGANIZATION}}": {
			Name:        "{{ORGANIZATION}}",
			Description: "Organization or company name",
			Value:       org,
			Priority:    1,
			Category:    "organization",
			Detected:    org != "MyProject",
		},
		"{{DOMAIN}}": {
			Name:        "{{DOMAIN}}",
			Description: "Organization domain",
			Value:       domain,
			Priority:    1,
			Category:    "organization",
			Detected:    domain != "example.com",
		},
	}

	// Add GitHub-specific organization variables
	githubInfo := tvd.githubIntegration.GetGitHubInfo()
	if githubInfo.IsGitHub {
		variables["{{GITHUB_OWNER}}"] = TemplateVariable{
			Name:        "{{GITHUB_OWNER}}",
			Description: "GitHub repository owner",
			Value:       githubInfo.Owner,
			Priority:    1,
			Category:    "organization",
			Detected:    true,
		}
		variables["{{GITHUB_HOST}}"] = TemplateVariable{
			Name:        "{{GITHUB_HOST}}",
			Description: "GitHub host (github.com or enterprise)",
			Value:       githubInfo.Host,
			Priority:    2,
			Category:    "organization",
			Detected:    true,
		}
	}

	return variables
}

// detectContactVariables detects contact-related variables
func (tvd *TemplateVariableDetector) detectContactVariables() map[string]TemplateVariable {
	variables := map[string]TemplateVariable{}

	// Detect email from git config
	email := tvd.detectGitEmail()
	if email != "" {
		variables["{{CONTACT_EMAIL}}"] = TemplateVariable{
			Name:        "{{CONTACT_EMAIL}}",
			Description: "Contact email for security reports",
			Value:       email,
			Priority:    1,
			Category:    "contact",
			Detected:    true,
		}
	}

	// Generate organization email
	org := tvd.detector.DetectOrganization()
	domain := tvd.detector.DetectDomain()
	if org != "" && domain != "" {
		orgEmail := "security@" + strings.ToLower(domain)
		variables["{{ORGANIZATION_EMAIL}}"] = TemplateVariable{
			Name:        "{{ORGANIZATION_EMAIL}}",
			Description: "Default organization security email",
			Value:       orgEmail,
			Priority:    2,
			Category:    "contact",
			Detected:    true,
		}
	}

	return variables
}

// detectSecurityVariables detects security-specific variables
func (tvd *TemplateVariableDetector) detectSecurityVariables() map[string]TemplateVariable {
	variables := map[string]TemplateVariable{
		// Bounty rewards (based on organization size)
		"{{CRITICAL_REWARD}}": {
			Name:        "{{CRITICAL_REWARD}}",
			Description: "Critical vulnerability reward",
			Value:       tvd.detectBountyAmount("critical"),
			Priority:    2,
			Category:    "security",
			Detected:    true,
		},
		"{{HIGH_REWARD}}": {
			Name:        "{{HIGH_REWARD}}",
			Description: "High severity vulnerability reward",
			Value:       tvd.detectBountyAmount("high"),
			Priority:    2,
			Category:    "security",
			Detected:    true,
		},
		"{{MEDIUM_REWARD}}": {
			Name:        "{{MEDIUM_REWARD}}",
			Description: "Medium severity vulnerability reward",
			Value:       tvd.detectBountyAmount("medium"),
			Priority:    2,
			Category:    "security",
			Detected:    true,
		},
		"{{LOW_REWARD}}": {
			Name:        "{{LOW_REWARD}}",
			Description: "Low severity vulnerability reward",
			Value:       tvd.detectBountyAmount("low"),
			Priority:    2,
			Category:    "security",
			Detected:    true,
		},
	}

	// Add GitHub-specific security URLs
	githubInfo := tvd.githubIntegration.GetGitHubInfo()
	if githubInfo.IsGitHub {
		urls := tvd.githubIntegration.GetGitHubURLs(githubInfo)
		variables["{{SECURITY_ADVISORIES_URL}}"] = TemplateVariable{
			Name:        "{{SECURITY_ADVISORIES_URL}}",
			Description: "GitHub security advisories URL",
			Value:       urls["advisories"],
			Priority:    1,
			Category:    "security",
			Detected:    true,
		}
		variables["{{BOUNTY_PROGRAM_URL}}"] = TemplateVariable{
			Name:        "{{BOUNTY_PROGRAM_URL}}",
			Description: "Bug bounty program URL",
			Value:       urls["security"] + "/policy",
			Priority:    2,
			Category:    "security",
			Detected:    true,
		}
	}

	return variables
}

// detectLegalVariables detects legal-related variables
func (tvd *TemplateVariableDetector) detectLegalVariables() map[string]TemplateVariable {
	domain := tvd.detector.DetectDomain()

	return map[string]TemplateVariable{
		"{{TERMS_URL}}": {
			Name:        "{{TERMS_URL}}",
			Description: "Terms of service URL",
			Value:       "https://" + domain + "/terms",
			Priority:    3,
			Category:    "legal",
			Detected:    domain != "example.com",
		},
		"{{PRIVACY_POLICY_URL}}": {
			Name:        "{{PRIVACY_POLICY_URL}}",
			Description: "Privacy policy URL",
			Value:       "https://" + domain + "/privacy",
			Priority:    3,
			Category:    "legal",
			Detected:    domain != "example.com",
		},
		"{{LICENSE_URL}}": {
			Name:        "{{LICENSE_URL}}",
			Description: "License URL",
			Value:       "https://creativecommons.org/licenses/by-sa/4.0/",
			Priority:    3,
			Category:    "legal",
			Detected:    true,
		},
	}
}

// detectTemporalVariables detects time-related variables
func (tvd *TemplateVariableDetector) detectTemporalVariables() map[string]TemplateVariable {
	now := time.Now()

	return map[string]TemplateVariable{
		"{{CURRENT_DATE}}": {
			Name:        "{{CURRENT_DATE}}",
			Description: "Current date in YYYY-MM-DD format",
			Value:       now.Format("2006-01-02"),
			Priority:    1,
			Category:    "temporal",
			Detected:    true,
		},
		"{{CURRENT_YEAR}}": {
			Name:        "{{CURRENT_YEAR}}",
			Description: "Current year",
			Value:       now.Format("2006"),
			Priority:    1,
			Category:    "temporal",
			Detected:    true,
		},
		"{{SUPPORT_END_DATE}}": {
			Name:        "{{SUPPORT_END_DATE}}",
			Description: "End of support date for current version",
			Value:       now.AddDate(1, 0, 0).Format("2006-01-02"),
			Priority:    2,
			Category:    "temporal",
			Detected:    true,
		},
		"{{PREVIOUS_SUPPORT_END_DATE}}": {
			Name:        "{{PREVIOUS_SUPPORT_END_DATE}}",
			Description: "End of support date for previous version",
			Value:       now.Format("2006-01-02"),
			Priority:    2,
			Category:    "temporal",
			Detected:    true,
		},
		"{{LAST_UPDATED}}": {
			Name:        "{{LAST_UPDATED}}",
			Description: "Last updated date",
			Value:       now.Format("2006-01-02"),
			Priority:    1,
			Category:    "temporal",
			Detected:    true,
		},
	}
}

// Helper methods

func (tvd *TemplateVariableDetector) getRepositoryName() string {
	githubInfo := tvd.githubIntegration.GetGitHubInfo()
	if githubInfo.IsGitHub {
		return githubInfo.Repository
	}

	// Fallback to project name
	return tvd.detector.DetectProjectName()
}

func (tvd *TemplateVariableDetector) detectProjectType() string {
	// Detect based on files and structure
	if tvd.fileExists("package.json") {
		return "web"
	}
	if tvd.fileExists("go.mod") {
		return "api"
	}
	if tvd.fileExists("Cargo.toml") {
		return "cli"
	}
	if tvd.fileExists("pyproject.toml") {
		return "python"
	}
	return "unknown"
}

func (tvd *TemplateVariableDetector) detectTechnologyStack() map[string]string {
	stack := make(map[string]string)

	// Detect JavaScript/Node.js
	if tvd.fileExists("package.json") {
		stack["javascript"] = "Node.js"
		if tvd.fileExists("next.config.js") {
			stack["framework"] = "Next.js"
		} else if tvd.fileExists("vite.config.js") {
			stack["framework"] = "Vite"
		}
	}

	// Detect Go
	if tvd.fileExists("go.mod") {
		stack["go"] = "Go"
	}

	// Detect Python
	if tvd.fileExists("pyproject.toml") || tvd.fileExists("requirements.txt") {
		stack["python"] = "Python"
	}

	// Detect Rust
	if tvd.fileExists("Cargo.toml") {
		stack["rust"] = "Rust"
	}

	// Detect Docker
	if tvd.fileExists("Dockerfile") {
		stack["containerization"] = "Docker"
	}

	return stack
}

func (tvd *TemplateVariableDetector) detectBountyAmount(severity string) string {
	// Simple heuristic based on organization/project size
	// In a real implementation, this could be configurable
	org := tvd.detector.DetectOrganization()

	// Enterprise organizations get higher rewards
	if strings.Contains(strings.ToLower(org), "corp") ||
		strings.Contains(strings.ToLower(org), "company") ||
		strings.Contains(strings.ToLower(org), "inc") {
		return map[string]string{
			"critical": "5000",
			"high":     "2000",
			"medium":   "500",
			"low":      "100",
		}[severity]
	}

	// Default open source rewards
	return map[string]string{
		"critical": "1000",
		"high":     "500",
		"medium":   "200",
		"low":      "50",
	}[severity]
}

func (tvd *TemplateVariableDetector) detectGitEmail() string {
	cmd := exec.Command("git", "config", "--get", "user.email")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	email := strings.TrimSpace(string(output))
	// Validate email format
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return ""
	}

	return email
}

func (tvd *TemplateVariableDetector) fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
