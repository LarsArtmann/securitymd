package internal

import (
	"context"
	"fmt"
	"strings"

	"github.com/LarsArtmann/template-CLI/pkg/sdk/fileops"
	"github.com/LarsArtmann/template-CLI/pkg/sdk/vfs"
	"github.com/spf13/afero"

	"github.com/LarsArtmann/template-SECURITY/v2/internal/domain"
	"github.com/LarsArtmann/template-SECURITY/v2/internal/service"
	"github.com/LarsArtmann/template-SECURITY/v2/internal/types"
)

// SecurityTool represents main security policy tool with clean architecture
type SecurityTool struct {
	templateManager   *vfs.TemplateManager
	fileSystem        vfs.FileSystem
	fileOps           *fileops.Service
	processor         *vfs.TemplateProcessor
	security          *vfs.TemplateSecurity
	variableService   *service.TemplateVariableService
	validationService *service.ValidationService
	githubService     *service.GitHubService
	projectService    *service.ProjectDetectionService
	errorHandler      *service.ErrorHandler
}

// NewSecurityTool creates a new security tool instance
func NewSecurityTool() *SecurityTool {
	vfsImpl := vfs.NewOSVFS()
	fileOps := fileops.NewService(fileops.ServiceOptions{
		BaseDirectory: ".",
		ValidatePaths: true,
	})

	return &SecurityTool{
		templateManager:   vfs.NewTemplateManager(vfsImpl),
		fileSystem:        vfsImpl,
		fileOps:           fileOps,
		processor:         vfs.NewTemplateProcessor(),
		security:          vfs.NewTemplateSecurity(),
		variableService:   service.NewTemplateVariableService(),
		validationService: service.NewValidationService(),
		githubService:     service.NewGitHubService(),
		projectService:    service.NewProjectDetectionService(),
		errorHandler:      service.NewErrorHandler(true),
	}
}

// PolicyConfig represents security policy configuration
type PolicyConfig struct {
	Type         types.PolicyType  `json:"type"`
	Organization string            `json:"organization"`
	Email        string            `json:"email"`
	Variables    map[string]string `json:"variables"`
	QuickMode    bool              `json:"quick_mode"`
	Options      PolicyOptions     `json:"options"`
}

// PolicyOptions represents additional policy options
type PolicyOptions struct {
	Validation bool   `json:"validation"`
	AutoDetect bool   `json:"auto_detect"`
	OutputFile string `json:"output_file"`
	Overwrite  bool   `json:"overwrite"`
	Template   string `json:"template"`
}

// GeneratePolicy generates a security policy
func (st *SecurityTool) GeneratePolicy(ctx context.Context, config PolicyConfig) (domain.SecurityPolicy, error) {
	// Validate input
	if err := st.validateConfig(config); err != nil {
		return domain.SecurityPolicy{}, fmt.Errorf("invalid configuration: %w", err)
	}

	// Detect variables if auto-detect is enabled
	var registry *domain.TemplateVariableRegistry
	if config.Options.AutoDetect {
		detected, err := st.variableService.DetectAllVariables(ctx)
		if err != nil {
			st.errorHandler.Handle(err)
			return domain.SecurityPolicy{}, fmt.Errorf("failed to detect variables: %w", err)
		}
		registry = detected
	} else {
		registry = st.variableService.GetRegistry()
	}

	// Add user-specified variables (override detected)
	for placeholder, value := range config.Variables {
		variable := domain.TemplateVariable{
			Placeholder: placeholder,
			Value:       value,
			Source:      domain.SourceUserInput,
			Priority:    types.PriorityCritical,
			Detected:    true,
			Overridden:  true,
		}
		registry.Register(variable)
	}

	// Get template content
	template, err := st.getTemplate(ctx, config.Type, config.Options.Template)
	if err != nil {
		return domain.SecurityPolicy{}, fmt.Errorf("failed to get template: %w", err)
	}

	// Process template
	content := st.processTemplate(template, registry.ExportMap())

	// Create policy ID
	policyID := domain.PolicyID{
		Value: st.generatePolicyID(config.Type, config.Organization),
	}

	// Build domain model
	policy := domain.SecurityPolicy{
		ID:           policyID,
		Type:         config.Type,
		Organization: st.buildOrganization(ctx, config),
		Contact:      st.buildContact(ctx, config),
		Versions:     st.buildVersions(),
		Validation:   st.buildValidation(ctx, content),
		Metadata:     st.buildMetadata(ctx, registry, content),
		Content:      content,
		GeneratedAt:  ctx.Value("now").(time.Time), // Will be set in commands
	}

	// Save policy if output file specified
	if config.Options.OutputFile != "" {
		if err := st.savePolicy(content, config.Options.OutputFile); err != nil {
			return domain.SecurityPolicy{}, fmt.Errorf("failed to save policy: %w", err)
		}
	}

	return policy, nil
}

// ValidatePolicy validates security policies
func (st *SecurityTool) ValidatePolicy(ctx context.Context, filePath string) (domain.Validation, error) {
	return st.validationService.ValidateSecurityPolicy(ctx, filePath)
}

// ValidateAllPolicies validates all policies in directory
func (st *SecurityTool) ValidateAllPolicies(ctx context.Context) ([]domain.Validation, error) {
	return st.validationService.ValidateAllPolicies(ctx)
}

// GetPolicy retrieves a policy by ID
func (st *SecurityTool) GetPolicy(ctx context.Context, id domain.PolicyID) (domain.SecurityPolicy, error) {
	// In this file-based implementation, read from file
	filePath := st.getPolicyFilePath(id)
	return st.readPolicyFile(filePath)
}

// Private helper methods

func (st *SecurityTool) validateConfig(config PolicyConfig) error {
	if !config.Type.IsValid() {
		return fmt.Errorf("invalid policy type: %s", config.Type)
	}

	if config.Organization == "" {
		return fmt.Errorf("organization name is required")
	}

	if config.Email == "" {
		return fmt.Errorf("contact email is required")
	}

	return nil
}

func (st *SecurityTool) getTemplate(ctx context.Context, policyType types.PolicyType, customTemplate string) (string, error) {
	if customTemplate != "" {
		// Load custom template
		data, err := st.fileSystem.ReadFile(customTemplate)
		if err != nil {
			return "", fmt.Errorf("failed to read custom template: %w", err)
		}
		return string(data), nil
	}

	// Load built-in template
	templatePath := st.getTemplatePath(policyType)
	data, err := st.fileSystem.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template: %w", err)
	}

	return string(data), nil
}

func (st *SecurityTool) getTemplatePath(policyType types.PolicyType) string {
	switch policyType {
	case types.PolicyTypeGitHub:
		return "templates/github-security.md"
	case types.PolicyTypeEnterprise:
		return "templates/enterprise-policy.md"
	case types.PolicyTypeBugBounty:
		return "templates/bug-bounty-program.md"
	case types.PolicyTypeIncidentResponse:
		return "templates/incident-response-plan.md"
	case types.PolicyTypePrivacyPolicy:
		return "templates/privacy-policy.md"
	default:
		return "templates/github-security.md"
	}
}

func (st *SecurityTool) processTemplate(template string, variables map[string]string) string {
	result := template

	// Multi-pass variable resolution for nested variables
	maxPasses := 3
	for pass := 0; pass < maxPasses; pass++ {
		changed := false
		for placeholder, value := range variables {
			oldResult := result
			result = strings.ReplaceAll(result, placeholder, value)
			if oldResult != result {
				changed = true
			}
		}
		if !changed {
			break
		}
	}

	return result
}

func (st *SecurityTool) generatePolicyID(policyType types.PolicyType, organization string) string {
	timestamp := time.Now().Format("20060102-150405")
	return fmt.Sprintf("%s-%s-%s", policyType, strings.ToLower(organization), timestamp)
}

func (st *SecurityTool) buildOrganization(ctx context.Context, config PolicyConfig) domain.Organization {
	// Get GitHub info if available
	githubInfo, _ := st.githubService.GetGitHubInfo(ctx)

	org := domain.Organization{
		Name:    config.Organization,
		Email:   config.Email,
		Website: "https://" + st.detectDomain(config.Organization),
		Type:    domain.OrgTypeOpenSource,
		Size:    domain.SizeSmall,
		GitHub:  githubInfo,
	}

	return org
}

func (st *SecurityTool) buildContact(ctx context.Context, config PolicyConfig) domain.Contact {
	return domain.Contact{
		Primary: domain.ContactMethod{
			Type:        domain.ContactTypeEmail,
			Value:       config.Email,
			Description: "Primary security contact",
			Encrypted:   false,
		},
		Alternative: domain.ContactMethod{
			Type:        domain.ContactTypeWeb,
			Value:       "GitHub Security",
			Description: "GitHub security advisories",
			Encrypted:   false,
		},
		Responsibilities: domain.Responsibilities{
			ResponseTime:      24 * time.Hour,
			InvestigationTime: 72 * time.Hour,
			PatchTime:         14 * 24 * time.Hour,
			DisclosureTime:    90 * 24 * time.Hour,
			SafeHarbor:        true,
			CoordinateCVE:     true,
		},
	}
}

func (st *SecurityTool) buildVersions() []domain.Version {
	now := time.Now()

	return []domain.Version{
		{
			Name:            "v2.x",
			SemanticVersion: "2.0.0",
			SupportedUntil:  now.AddDate(1, 0, 0),
			Status:          domain.StatusSupported,
			IsLatest:        true,
		},
		{
			Name:            "v1.x",
			SemanticVersion: "1.0.0",
			SupportedUntil:  now,
			Status:          domain.StatusDeprecated,
			IsLatest:        false,
			IsPrevious:      true,
		},
	}
}

func (st *SecurityTool) buildValidation(ctx context.Context, content string) domain.Validation {
	// Simple validation
	score := uint8(100)
	if len(content) < 500 {
		score = 50
	}

	status := types.StatusValid
	if score < 70 {
		status = types.StatusWarning
	}

	return domain.Validation{
		Status:    status,
		Score:     score,
		Issues:    []domain.ValidationIssue{},
		Warnings:  []domain.ValidationWarning{},
		CheckedAt: time.Now(),
		Quality: domain.QualityMetrics{
			Words:         uint16(len(strings.Fields(content))),
			Lines:         uint16(len(strings.Split(content, "\n"))),
			Characters:    uint16(len(content)),
			Sections:      10, // Approximate
			Readability:   score,
			Completeness:  score,
			BestPractices: score,
		},
	}
}

func (st *SecurityTool) buildMetadata(ctx context.Context, registry *domain.TemplateVariableRegistry, content string) domain.PolicyMetadata {
	stats := registry.Stats()

	return domain.PolicyMetadata{
		TemplateVersion:  "2.0",
		GeneratorVersion: "v2.0.0",
		VariableCount:    stats.TotalVariables,
		ResolvedCount:    stats.TotalVariables - stats.UnresolvedCount,
		Variables:        st.convertRegistryToVariables(registry.GetAll()),
		Project:          st.buildProjectInfo(ctx),
	}
}

func (st *SecurityTool) convertRegistryToVariables(variables map[string]domain.TemplateVariable) []domain.Variable {
	var result []domain.Variable

	for _, variable := range variables {
		result = append(result, domain.Variable{
			Placeholder: variable.Placeholder,
			Name:        variable.Name,
			Description: variable.Description,
			Value:       variable.Value,
			Priority:    variable.Priority,
			Category:    variable.Category,
			Detected:    variable.Detected,
			Overridden:  variable.Overridden,
			Required:    variable.Required,
		})
	}

	return result
}

func (st *SecurityTool) buildProjectInfo(ctx context.Context) domain.ProjectInfo {
	projectInfo, _ := st.projectService.AnalyzeProject(ctx)
	return projectInfo
}

func (st *SecurityTool) detectDomain(organization string) string {
	// Simple domain detection
	if strings.Contains(organization, ".") {
		return organization
	}
	return strings.ToLower(organization) + ".com"
}

func (st *SecurityTool) savePolicy(content, outputPath string) error {
	return st.fileOps.WriteFile(outputPath, []byte(content))
}

func (st *SecurityTool) getPolicyFilePath(id domain.PolicyID) string {
	return fmt.Sprintf("policies/%s.json", id.Value)
}

func (st *SecurityTool) readPolicyFile(filePath string) (domain.SecurityPolicy, error) {
	// Implementation would read from file system
	return domain.SecurityPolicy{}, fmt.Errorf("not implemented")
}

// SetErrorHandler allows dependency injection for testing
func (st *SecurityTool) SetErrorHandler(handler service.ErrorHandler) {
	st.errorHandler = handler
}
