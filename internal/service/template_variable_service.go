package service

import (
	"context"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/LarsArtmann/template-SECURITY/v2/internal/domain"
	"github.com/LarsArtmann/template-SECURITY/v2/internal/errors"
	"github.com/LarsArtmann/template-SECURITY/v2/internal/types"
)

// TemplateVariableService handles template variable detection and management
type TemplateVariableService struct {
	gitService     GitService
	projectService ProjectDetectionService
	githubService  GitHubService
	registry       domain.TemplateVariableRegistry
	errorHandler   errors.ErrorHandler
}

// NewTemplateVariableService creates a new template variable service
func NewTemplateVariableService() *TemplateVariableService {
	registry := domain.NewTemplateVariableRegistry()
	errorHandler := errors.NewErrorHandler(true)

	return &TemplateVariableService{
		gitService:     NewGitService(),
		projectService: NewProjectDetectionService(),
		githubService:  NewGitHubService(),
		registry:       registry,
		errorHandler:   errorHandler,
	}
}

// DetectAllVariables detects all available template variables
func (tvs *TemplateVariableService) DetectAllVariables(ctx context.Context) (*domain.TemplateVariableRegistry, error) {
	// Clear existing registry
	tvs.registry.Clear()

	// Detect project variables
	projectVars, err := tvs.detectProjectVariables(ctx)
	if err != nil {
		tvs.errorHandler.Handle(err)
	}
	for _, variable := range projectVars {
		tvs.registry.Register(variable)
	}

	// Detect organization variables
	orgVars, err := tvs.detectOrganizationVariables(ctx)
	if err != nil {
		tvs.errorHandler.Handle(err)
	}
	for _, variable := range orgVars {
		tvs.registry.Register(variable)
	}

	// Detect contact variables
	contactVars, err := tvs.detectContactVariables(ctx)
	if err != nil {
		tvs.errorHandler.Handle(err)
	}
	for _, variable := range contactVars {
		tvs.registry.Register(variable)
	}

	// Detect security variables
	securityVars, err := tvs.detectSecurityVariables(ctx)
	if err != nil {
		tvs.errorHandler.Handle(err)
	}
	for _, variable := range securityVars {
		tvs.registry.Register(variable)
	}

	// Detect legal variables
	legalVars, err := tvs.detectLegalVariables(ctx)
	if err != nil {
		tvs.errorHandler.Handle(err)
	}
	for _, variable := range legalVars {
		tvs.registry.Register(variable)
	}

	// Detect temporal variables
	temporalVars := tvs.detectTemporalVariables(ctx)
	for _, variable := range temporalVars {
		tvs.registry.Register(variable)
	}

	// Validate registry
	validation := tvs.registry.Validate()
	if validation.Status == types.StatusInvalid {
		return nil, errors.NewValidationError(
			"VAL004",
			"Variable registry validation failed",
			"registry",
			types.ValidationLevelError,
		)
	}

	return &tvs.registry, nil
}

// detectProjectVariables detects project-specific variables
func (tvs *TemplateVariableService) detectProjectVariables(ctx context.Context) ([]domain.TemplateVariable, error) {
	var variables []domain.TemplateVariable

	// Project name
	projectName, err := tvs.projectService.DetectProjectName(ctx)
	if err != nil {
		return nil, errors.NewProjectDetectionError(
			"DET002",
			"Failed to detect project name",
			"project_detection",
			"DetectProjectName",
		).WithContext("error", err.Error())
	}

	if projectName != "MyProject" {
		variables = append(variables, domain.TemplateVariable{
			Placeholder: "{{PROJECT_NAME}}",
			Name:        "Project Name",
			Description: "Project or repository name",
			Value:       projectName,
			Priority:    types.PriorityCritical,
			Category:    types.CategoryProject,
			Detected:    true,
			Required:    false,
			Source:      domain.SourceProject,
		})
	}

	// Repository name
	repoName, err := tvs.gitService.GetRepositoryName(ctx)
	if err == nil && repoName != "" {
		variables = append(variables, domain.TemplateVariable{
			Placeholder: "{{REPOSITORY_NAME}}",
			Name:        "Repository Name",
			Description: "Repository name from git",
			Value:       repoName,
			Priority:    types.PriorityCritical,
			Category:    types.CategoryProject,
			Detected:    true,
			Required:    false,
			Source:      domain.SourceGitConfig,
		})
	}

	// Project type
	projectType := tvs.detectProjectType()
	if projectType != "" {
		variables = append(variables, domain.TemplateVariable{
			Placeholder: "{{PROJECT_TYPE}}",
			Name:        "Project Type",
			Description: "Type of project (web, api, library, etc.)",
			Value:       projectType,
			Priority:    types.PriorityMedium,
			Category:    types.CategoryProject,
			Detected:    true,
			Required:    false,
			Source:      domain.SourceProject,
		})
	}

	// Technology stack
	techStack := tvs.detectTechnologyStack()
	for tech, info := range techStack {
		variables = append(variables, domain.TemplateVariable{
			Placeholder: "{{TECH_STACK_" + strings.ToUpper(tech) + "}}",
			Name:        "Technology: " + tech,
			Description: "Technology in stack: " + tech,
			Value:       info,
			Priority:    types.PriorityMedium,
			Category:    types.CategoryProject,
			Detected:    true,
			Required:    false,
			Source:      domain.SourceProject,
		})
	}

	return variables, nil
}

// detectOrganizationVariables detects organization-specific variables
func (tvs *TemplateVariableService) detectOrganizationVariables(ctx context.Context) ([]domain.TemplateVariable, error) {
	var variables []domain.TemplateVariable

	// Organization name
	org, err := tvs.projectService.DetectOrganization(ctx)
	if err != nil {
		return nil, err
	}

	if org != "MyProject" {
		variables = append(variables, domain.TemplateVariable{
			Placeholder: "{{ORGANIZATION}}",
			Name:        "Organization",
			Description: "Organization or company name",
			Value:       org,
			Priority:    types.PriorityCritical,
			Category:    types.CategoryOrganization,
			Detected:    true,
			Required:    true,
			Source:      domain.SourceProject,
		})
	}

	// Domain
	domain, err := tvs.projectService.DetectDomain(ctx)
	if err != nil {
		return nil, err
	}

	if domain != "example.com" {
		variables = append(variables, domain.TemplateVariable{
			Placeholder: "{{DOMAIN}}",
			Name:        "Domain",
			Description: "Organization domain",
			Value:       domain,
			Priority:    types.PriorityHigh,
			Category:    types.CategoryOrganization,
			Detected:    true,
			Required:    true,
			Source:      domain.SourceProject,
		})
	}

	// GitHub-specific variables
	githubInfo, err := tvs.githubService.GetGitHubInfo(ctx)
	if err == nil && githubInfo.IsGitHub {
		variables = append(variables, domain.TemplateVariable{
			Placeholder: "{{GITHUB_OWNER}}",
			Name:        "GitHub Owner",
			Description: "GitHub repository owner",
			Value:       githubInfo.Owner,
			Priority:    types.PriorityHigh,
			Category:    types.CategoryOrganization,
			Detected:    true,
			Required:    false,
			Source:      domain.SourceGitHub,
		})

		variables = append(variables, domain.TemplateVariable{
			Placeholder: "{{GITHUB_HOST}}",
			Name:        "GitHub Host",
			Description: "GitHub host (github.com or enterprise)",
			Value:       githubInfo.Host,
			Priority:    types.PriorityMedium,
			Category:    types.CategoryOrganization,
			Detected:    true,
			Required:    false,
			Source:      domain.SourceGitHub,
		})
	}

	return variables, nil
}

// detectContactVariables detects contact-related variables
func (tvs *TemplateVariableService) detectContactVariables(ctx context.Context) ([]domain.TemplateVariable, error) {
	var variables []domain.TemplateVariable

	// Detect email from git config
	email, err := tvs.gitService.GetUserEmail(ctx)
	if err == nil && tvs.isValidEmail(email) {
		variables = append(variables, domain.TemplateVariable{
			Placeholder: "{{CONTACT_EMAIL}}",
			Name:        "Contact Email",
			Description: "Contact email for security reports",
			Value:       email,
			Priority:    types.PriorityCritical,
			Category:    types.CategoryContact,
			Detected:    true,
			Required:    true,
			Source:      domain.SourceGitConfig,
		})
	}

	// Generate organization email
	org, _ := tvs.projectService.DetectOrganization(ctx)
	domain, _ := tvs.projectService.DetectDomain(ctx)
	if org != "" && domain != "" {
		orgEmail := "security@" + strings.ToLower(domain)
		variables = append(variables, domain.TemplateVariable{
			Placeholder: "{{ORGANIZATION_EMAIL}}",
			Name:        "Organization Email",
			Description: "Default organization security email",
			Value:       orgEmail,
			Priority:    types.PriorityHigh,
			Category:    types.CategoryContact,
			Detected:    true,
			Required:    false,
			Source:      domain.SourceAutoDetected,
		})
	}

	return variables, nil
}

// detectSecurityVariables detects security-specific variables
func (tvs *TemplateVariableService) detectSecurityVariables(ctx context.Context) ([]domain.TemplateVariable, error) {
	var variables []domain.TemplateVariable

	// Bounty rewards (based on organization size)
	org, _ := tvs.projectService.DetectOrganization(ctx)
	bountyType := tvs.detectBountyType(org)

	rewards := map[types.BountySeverity]Reward{
		types.SeverityCritical: tvs.detectBountyAmount("critical", bountyType),
		types.SeverityHigh:     tvs.detectBountyAmount("high", bountyType),
		types.SeverityMedium:   tvs.detectBountyAmount("medium", bountyType),
		types.SeverityLow:      tvs.detectBountyAmount("low", bountyType),
	}

	for severity, reward := range rewards {
		severityStr := string(severity)
		variables = append(variables, domain.TemplateVariable{
			Placeholder: "{{" + strings.ToUpper(severityStr) + "_REWARD}}",
			Name:        strings.Title(severityStr) + " Severity Reward",
			Description: severityStr + " severity vulnerability reward",
			Value:       string(rune(reward.Amount)),
			Priority:    types.PriorityMedium,
			Category:    types.CategorySecurity,
			Detected:    true,
			Required:    false,
			Source:      domain.SourceAutoDetected,
		})
	}

	// GitHub-specific security URLs
	githubInfo, err := tvs.githubService.GetGitHubInfo(ctx)
	if err == nil && githubInfo.IsGitHub {
		urls := tvs.githubService.GetGitHubURLs(githubInfo)

		variables = append(variables, domain.TemplateVariable{
			Placeholder: "{{SECURITY_ADVISORIES_URL}}",
			Name:        "Security Advisories URL",
			Description: "GitHub security advisories URL",
			Value:       urls["advisories"],
			Priority:    types.PriorityHigh,
			Category:    types.CategorySecurity,
			Detected:    true,
			Required:    false,
			Source:      domain.SourceGitHub,
		})

		variables = append(variables, domain.TemplateVariable{
			Placeholder: "{{BOUNTY_PROGRAM_URL}}",
			Name:        "Bounty Program URL",
			Description: "Bug bounty program URL",
			Value:       urls["security"] + "/policy",
			Priority:    types.PriorityMedium,
			Category:    types.CategorySecurity,
			Detected:    true,
			Required:    false,
			Source:      domain.SourceGitHub,
		})
	}

	return variables, nil
}

// detectLegalVariables detects legal-related variables
func (tvs *TemplateVariableService) detectLegalVariables(ctx context.Context) ([]domain.TemplateVariable, error) {
	var variables []domain.TemplateVariable

	domain, _ := tvs.projectService.DetectDomain(ctx)

	variables = append(variables, domain.TemplateVariable{
		Placeholder: "{{TERMS_URL}}",
		Name:        "Terms of Service URL",
		Description: "Terms of service URL",
		Value:       "https://" + domain + "/terms",
		Priority:    types.PriorityLow,
		Category:    types.CategoryLegal,
		Detected:    domain != "example.com",
		Required:    false,
		Source:      domain.SourceAutoDetected,
	})

	variables = append(variables, domain.TemplateVariable{
		Placeholder: "{{PRIVACY_POLICY_URL}}",
		Name:        "Privacy Policy URL",
		Description: "Privacy policy URL",
		Value:       "https://" + domain + "/privacy",
		Priority:    types.PriorityLow,
		Category:    types.CategoryLegal,
		Detected:    domain != "example.com",
		Required:    false,
		Source:      domain.SourceAutoDetected,
	})

	variables = append(variables, domain.TemplateVariable{
		Placeholder: "{{LICENSE_URL}}",
		Name:        "License URL",
		Description: "License URL",
		Value:       "https://creativecommons.org/licenses/by-sa/4.0/",
		Priority:    types.PriorityLow,
		Category:    types.CategoryLegal,
		Detected:    true,
		Required:    false,
		Source:      domain.SourceDefault,
	})

	return variables, nil
}

// detectTemporalVariables detects time-related variables
func (tvs *TemplateVariableService) detectTemporalVariables(ctx context.Context) []domain.TemplateVariable {
	var variables []domain.TemplateVariable

	now := time.Now()

	variables = append(variables, domain.TemplateVariable{
		Placeholder: "{{CURRENT_DATE}}",
		Name:        "Current Date",
		Description: "Current date in YYYY-MM-DD format",
		Value:       now.Format("2006-01-02"),
		Priority:    types.PriorityHigh,
		Category:    types.CategoryTemporal,
		Detected:    true,
		Required:    false,
		Source:      domain.SourceAutoDetected,
	})

	variables = append(variables, domain.TemplateVariable{
		Placeholder: "{{CURRENT_YEAR}}",
		Name:        "Current Year",
		Description: "Current year",
		Value:       now.Format("2006"),
		Priority:    types.PriorityHigh,
		Category:    types.CategoryTemporal,
		Detected:    true,
		Required:    false,
		Source:      domain.SourceAutoDetected,
	})

	variables = append(variables, domain.TemplateVariable{
		Placeholder: "{{SUPPORT_END_DATE}}",
		Name:        "Support End Date",
		Description: "End of support date for current version",
		Value:       now.AddDate(1, 0, 0).Format("2006-01-02"),
		Priority:    types.PriorityMedium,
		Category:    types.CategoryTemporal,
		Detected:    true,
		Required:    false,
		Source:      domain.SourceAutoDetected,
	})

	variables = append(variables, domain.TemplateVariable{
		Placeholder: "{{LAST_UPDATED}}",
		Name:        "Last Updated Date",
		Description: "Last updated date",
		Value:       now.Format("2006-01-02"),
		Priority:    types.PriorityMedium,
		Category:    types.CategoryTemporal,
		Detected:    true,
		Required:    false,
		Source:      domain.SourceAutoDetected,
	})

	return variables
}

// Helper methods

func (tvs *TemplateVariableService) detectProjectType() string {
	if tvs.fileExists("package.json") {
		return "web"
	}
	if tvs.fileExists("go.mod") {
		return "api"
	}
	if tvs.fileExists("Cargo.toml") {
		return "cli"
	}
	if tvs.fileExists("pyproject.toml") {
		return "python"
	}
	return "unknown"
}

func (tvs *TemplateVariableService) detectTechnologyStack() map[string]string {
	stack := make(map[string]string)

	if tvs.fileExists("package.json") {
		stack["javascript"] = "Node.js"
		if tvs.fileExists("next.config.js") {
			stack["framework"] = "Next.js"
		} else if tvs.fileExists("vite.config.js") {
			stack["framework"] = "Vite"
		}
	}

	if tvs.fileExists("go.mod") {
		stack["go"] = "Go"
	}

	if tvs.fileExists("pyproject.toml") || tvs.fileExists("requirements.txt") {
		stack["python"] = "Python"
	}

	if tvs.fileExists("Cargo.toml") {
		stack["rust"] = "Rust"
	}

	if tvs.fileExists("Dockerfile") {
		stack["containerization"] = "Docker"
	}

	return stack
}

func (tvs *TemplateVariableService) detectBountyType(org string) string {
	if strings.Contains(strings.ToLower(org), "corp") ||
		strings.Contains(strings.ToLower(org), "company") ||
		strings.Contains(strings.ToLower(org), "inc") {
		return "enterprise"
	}
	return "opensource"
}

func (tvs *TemplateVariableService) detectBountyAmount(severity, bountyType string) uint32 {
	if bountyType == "enterprise" {
		return map[string]uint32{
			"critical": 5000,
			"high":     2000,
			"medium":   500,
			"low":      100,
		}[severity]
	}

	return map[string]uint32{
		"critical": 1000,
		"high":     500,
		"medium":   200,
		"low":      50,
	}[severity]
}

func (tvs *TemplateVariableService) isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func (tvs *TemplateVariableService) fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// GetRegistry returns the current variable registry
func (tvs *TemplateVariableService) GetRegistry() *domain.TemplateVariableRegistry {
	return &tvs.registry
}

// ExportMap exports registry as simple map for template processing
func (tvs *TemplateVariableService) ExportMap() map[string]string {
	return tvs.registry.ExportMap()
}

// ValidateRegistry validates the current registry
func (tvs *TemplateVariableService) ValidateRegistry() domain.Validation {
	return tvs.registry.Validate()
}

// Reward represents bounty reward information
type Reward struct {
	Amount   uint32
	Currency string
}
