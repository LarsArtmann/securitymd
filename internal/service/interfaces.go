package service

import (
	"context"
	"fmt"
	"time"

	"github.com/LarsArtmann/template-SECURITY/v2/internal/domain"
	"github.com/LarsArtmann/template-SECURITY/v2/internal/types"
)

// SecurityAdvisory represents a GitHub security advisory (moved from repository)
type SecurityAdvisory struct {
	ID             string    `json:"id"`
	GHSAID         string    `json:"ghsa_id"`
	Summary        string    `json:"summary"`
	Description     string    `json:"description"`
	Severity       string    `json:"severity"`
	PublishedAt    time.Time `json:"published_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	CVEs           []CVE     `json:"cves"`
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
}

// CVE represents a CVE entry (moved from repository)
type CVE struct {
	ID        string `json:"id"`
	Severity  string `json:"severity"`
	CWEs      []CWE  `json:"cwes"`
	Published string `json:"published"`
	Updated   string `json:"updated"`
}

// CWE represents a Common Weakness Enumeration (moved from repository)
type CWE struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Vulnerability represents a vulnerability (moved from repository)
type Vulnerability struct {
	Package       string   `json:"package"`
	Ecosystem     string   `json:"ecosystem"`
	Severity      string   `json:"severity"`
	Requirements  []string `json:"requirements"`
	FixedVersions []string `json:"fixed_versions"`
}

// GitService provides git-related operations
type GitService interface {
	// GetRepositoryName gets repository name from git config
	GetRepositoryName(ctx context.Context) (string, error)
	
	// GetUserEmail gets user email from git config
	GetUserEmail(ctx context.Context) (string, error)
	
	// GetRemoteURL gets git remote URL
	GetRemoteURL(ctx context.Context) (string, error)
	
	// GetCurrentBranch gets current branch
	GetCurrentBranch(ctx context.Context) (string, error)
	
	// IsGitRepository checks if current directory is a git repository
	IsGitRepository(ctx context.Context) bool
}

// ProjectDetectionService provides project detection capabilities
type ProjectDetectionService interface {
	// DetectProjectName detects project name from various sources
	DetectProjectName(ctx context.Context) (string, error)
	
	// DetectOrganization detects organization name
	DetectOrganization(ctx context.Context) (string, error)
	
	// DetectDomain detects organization domain
	DetectDomain(ctx context.Context) (string, error)
	
	// DetectTechnologyStack detects technology stack
	DetectTechnologyStack(ctx context.Context) ([]types.TechStack, error)
	
	// DetectProjectType detects project type
	DetectProjectType(ctx context.Context) (types.ProjectType, error)
	
	// AnalyzeProject analyzes project structure and metadata
	AnalyzeProject(ctx context.Context) (domain.ProjectInfo, error)
}

// GitHubService provides GitHub integration capabilities
type GitHubService interface {
	// GetGitHubInfo detects GitHub repository information
	GetGitHubInfo(ctx context.Context) (domain.GitHubInfo, error)
	
	// GetGitHubURLs generates GitHub URLs for repository
	GetGitHubURLs(info domain.GitHubInfo) map[string]string
	
	// IsGitHubRepository checks if repository is on GitHub
	IsGitHubRepository(ctx context.Context) bool
	
	// GetOwnerInfo retrieves GitHub owner information
	GetOwnerInfo(ctx context.Context, owner string) (domain.Organization, error)
	
	// GetSecurityAdvisories retrieves security advisories
	GetSecurityAdvisories(ctx context.Context, owner, repo string) ([]SecurityAdvisory, error)
}

// SecurityPolicyService provides security policy management
type SecurityPolicyService interface {
	// GeneratePolicy generates a security policy
	GeneratePolicy(ctx context.Context, request PolicyGenerationRequest) (domain.SecurityPolicy, error)
	
	// ValidatePolicy validates a security policy
	ValidatePolicy(ctx context.Context, policy domain.SecurityPolicy) domain.Validation
	
	// UpdatePolicy updates an existing policy
	UpdatePolicy(ctx context.Context, policy domain.SecurityPolicy) (domain.SecurityPolicy, error)
	
	// GetPolicy retrieves a policy by ID
	GetPolicy(ctx context.Context, id domain.PolicyID) (domain.SecurityPolicy, error)
	
	// ListPolicies returns all policies
	ListPolicies(ctx context.Context) ([]domain.SecurityPolicy, error)
	
	// DeletePolicy removes a policy
	DeletePolicy(ctx context.Context, id domain.PolicyID) error
}

// TemplateVariableService provides template variable management
type TemplateVariableService interface {
	// DetectAllVariables detects all available template variables
	DetectAllVariables(ctx context.Context) (*domain.TemplateVariableRegistry, error)
	
	// GetRegistry returns current variable registry
	GetRegistry() *domain.TemplateVariableRegistry
	
	// ExportMap exports variables as simple map
	ExportMap() map[string]string
	
	// ValidateRegistry validates current registry
	ValidateRegistry() domain.Validation
	
	// UpdateVariable updates a single variable
	UpdateVariable(ctx context.Context, variable domain.TemplateVariable) error
	
	// GetVariable retrieves a variable by placeholder
	GetVariable(ctx context.Context, placeholder string) (domain.TemplateVariable, bool)
}

// ValidationService provides validation capabilities
type ValidationService interface {
	// ValidateSecurityPolicy validates a security policy file
	ValidateSecurityPolicy(ctx context.Context, filePath string) (domain.Validation, error)
	
	// ValidateAllPolicies validates all policy files in directory
	ValidateAllPolicies(ctx context.Context) ([]domain.Validation, error)
	
	// GetValidationResults retrieves validation history
	GetValidationResults(ctx context.Context, policyID domain.PolicyID) ([]domain.Validation, error)
	
	// GenerateValidationReport generates comprehensive validation report
	GenerateValidationReport(ctx context.Context, results []domain.Validation) (ValidationReport, error)
}

// CI/CD Service interfaces

// GitHubActionService provides GitHub Actions integration
type GitHubActionService interface {
	// GenerateWorkflow generates GitHub Action workflow
	GenerateWorkflow(ctx context.Context, config GitHubActionConfig) (GitHubWorkflow, error)
	
	// ValidateWorkflow validates a workflow configuration
	ValidateWorkflow(ctx context.Context, workflow GitHubWorkflow) error
	
	// GetWorkflows returns available workflow templates
	GetWorkflows(ctx context.Context) []GitHubWorkflowTemplate
}

// NotificationService provides notification capabilities
type NotificationService interface {
	// SendValidationResult sends validation result notification
	SendValidationResult(ctx context.Context, result domain.Validation) error
	
	// SendPolicyUpdate sends policy update notification
	SendPolicyUpdate(ctx context.Context, policy domain.SecurityPolicy) error
	
	// SendSecurityAlert sends security alert notification
	SendSecurityAlert(ctx context.Context, alert SecurityAlert) error
}

// Configuration Service interfaces

// ConfigurationService provides configuration management
type ConfigurationService interface {
	// LoadConfiguration loads configuration from file
	LoadConfiguration(ctx context.Context, filePath string) (Configuration, error)
	
	// SaveConfiguration saves configuration to file
	SaveConfiguration(ctx context.Context, config Configuration, filePath string) error
	
	// GetDefaultConfiguration returns default configuration
	GetDefaultConfiguration(ctx context.Context) Configuration
	
	// ValidateConfiguration validates configuration
	ValidateConfiguration(ctx context.Context, config Configuration) error
}

// Request and response types

// PolicyGenerationRequest represents a policy generation request
type PolicyGenerationRequest struct {
	Type         types.PolicyType             `json:"type"`
	Organization string                      `json:"organization"`
	Email        string                       `json:"email"`
	Variables    map[string]string            `json:"variables"`
	Template     string                       `json:"template"`
	Options      PolicyGenerationOptions      `json:"options"`
}

// PolicyGenerationOptions represents policy generation options
type PolicyGenerationOptions struct {
	QuickMode     bool   `json:"quick_mode"`
	Validation    bool   `json:"validation"`
	AutoDetect    bool   `json:"auto_detect"`
	OutputFile    string `json:"output_file"`
	Overwrite     bool   `json:"overwrite"`
}

// ValidationReport represents a validation report
type ValidationReport struct {
	Summary      ValidationSummary            `json:"summary"`
	Results      []domain.Validation          `json:"results"`
	Metrics      ValidationMetrics           `json:"metrics"`
	GeneratedAt  time.Time                  `json:"generated_at"`
}

// ValidationSummary represents validation summary
type ValidationSummary struct {
	TotalValid     int `json:"total_valid"`
	TotalInvalid   int `json:"total_invalid"`
	TotalWarnings  int `json:"total_warnings"`
	AverageScore   float64 `json:"average_score"`
	HighestScore   uint8 `json:"highest_score"`
	LowestScore    uint8 `json:"lowest_score"`
}

// ValidationMetrics represents validation metrics
type ValidationMetrics struct {
	TotalPoliciesValidated int                `json:"total_policies_validated"`
	ValidationScore        map[types.ValidationLevel]int `json:"validation_score"`
	CommonIssues          []string                 `json:"common_issues"`
	ImprovementSuggestions []string                 `json:"improvement_suggestions"`
}

// GitHubActionConfig represents GitHub Action configuration
type GitHubActionConfig struct {
	Name           string   `json:"name"`
	Description     string   `json:"description"`
	Triggers       []string `json:"triggers"`
	Validation     bool     `json:"validation"`
	Notifications  bool     `json:"notifications"`
	ApprovalFlow   bool     `json:"approval_flow"`
}

// GitHubWorkflow represents a GitHub Action workflow
type GitHubWorkflow struct {
	Name        string                 `json:"name"`
	FileName    string                 `json:"file_name"`
	Content     string                 `json:"content"`
	Metadata    GitHubWorkflowMetadata `json:"metadata"`
	Validated   bool                   `json:"validated"`
}

// GitHubWorkflowMetadata represents workflow metadata
type GitHubWorkflowMetadata struct {
	Version     string    `json:"version"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Author      string    `json:"author"`
	Description string    `json:"description"`
}

// GitHubWorkflowTemplate represents a workflow template
type GitHubWorkflowTemplate struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Category    string                 `json:"category"`
	Content     string                 `json:"content"`
	Options     map[string]interface{} `json:"options"`
}

// SecurityAlert represents a security alert
type SecurityAlert struct {
	ID          string                 `json:"id"`
	Type        AlertType              `json:"type"`
	Severity    AlertSeverity          `json:"severity"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Component   string                 `json:"component"`
	Timestamp   time.Time              `json:"timestamp"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// AlertType represents different alert types
type AlertType string

const (
	AlertTypeValidation    AlertType = "validation"
	AlertTypeSecurity     AlertType = "security"
	AlertTypeCompliance   AlertType = "compliance"
	AlertTypePerformance  AlertType = "performance"
)

// AlertSeverity represents alert severity levels
type AlertSeverity string

const (
	AlertSeverityLow      AlertSeverity = "low"
	AlertSeverityMedium   AlertSeverity = "medium"
	AlertSeverityHigh     AlertSeverity = "high"
	AlertSeverityCritical AlertSeverity = "critical"
)

// Configuration represents application configuration
type Configuration struct {
	Repository    RepositoryConfig    `json:"repository"`
	Services      ServiceConfig      `json:"services"`
	Validation    ValidationConfig    `json:"validation"`
	Templates     TemplateConfig     `json:"templates"`
	Notifications NotificationConfig `json:"notifications"`
}

// RepositoryConfig represents repository configuration
type RepositoryConfig struct {
	Type        string            `json:"type"`
	Connection  string            `json:"connection"`
	Options     map[string]interface{} `json:"options"`
}

// ServiceConfig represents service configuration
type ServiceConfig struct {
	GitHub      GitHubServiceConfig     `json:"github"`
	Validation  ValidationServiceConfig  `json:"validation"`
	Notification NotificationServiceConfig `json:"notification"`
}

// GitHubServiceConfig represents GitHub service configuration
type GitHubServiceConfig struct {
	Enabled     bool   `json:"enabled"`
	Token        string `json:"token"`
	API          string `json:"api"`
	Timeout      string `json:"timeout"`
}

// ValidationServiceConfig represents validation service configuration
type ValidationServiceConfig struct {
	Strict      bool `json:"strict"`
	AutoFix     bool `json:"auto_fix"`
	Level       string `json:"level"`
}

// ValidationServiceConfig represents notification service configuration
type NotificationServiceConfig struct {
	Enabled   bool     `json:"enabled"`
	Channels  []string `json:"channels"`
	Threshold  uint8    `json:"threshold"`
}

// TemplateConfig represents template configuration
type TemplateConfig struct {
	Directory  string            `json:"directory"`
	Variables  map[string]string `json:"variables"`
	Strict     bool              `json:"strict"`
}

// Common service errors

// ServiceError represents a service-level error
type ServiceError struct {
	Service   string `json:"service"`
	Operation string `json:"operation"`
	Code      string `json:"code"`
	Message   string `json:"message"`
}

// Error implements error interface
func (se ServiceError) Error() string {
	return fmt.Sprintf("[%s:%s] %s: %s", se.Service, se.Operation, se.Code, se.Message)
}