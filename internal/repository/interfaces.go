package repository

import (
	"context"

	"github.com/LarsArtmann/template-SECURITY/v2/internal/domain"
	"github.com/LarsArtmann/template-SECURITY/v2/internal/types"
)

// ProjectRepository defines operations for project data persistence
type ProjectRepository interface {
	// Save stores project information
	Save(ctx context.Context, project domain.ProjectInfo) error
	
	// Get retrieves project information
	Get(ctx context.Context, id domain.PolicyID) (domain.ProjectInfo, error)
	
	// Delete removes project information
	Delete(ctx context.Context, id domain.PolicyID) error
	
	// List returns all projects
	List(ctx context.Context) ([]domain.ProjectInfo, error)
	
	// Search finds projects by criteria
	Search(ctx context.Context, criteria ProjectSearchCriteria) ([]domain.ProjectInfo, error)
}

// PolicyRepository defines operations for security policy persistence
type PolicyRepository interface {
	// Save stores a security policy
	Save(ctx context.Context, policy domain.SecurityPolicy) error
	
	// Get retrieves a security policy by ID
	Get(ctx context.Context, id domain.PolicyID) (domain.SecurityPolicy, error)
	
	// GetByType retrieves policies by type
	GetByType(ctx context.Context, policyType types.PolicyType) ([]domain.SecurityPolicy, error)
	
	// Update modifies an existing policy
	Update(ctx context.Context, policy domain.SecurityPolicy) error
	
	// Delete removes a security policy
	Delete(ctx context.Context, id domain.PolicyID) error
	
	// List returns all policies
	List(ctx context.Context) ([]domain.SecurityPolicy, error)
	
	// Validate checks policy validity
	Validate(ctx context.Context, policy domain.SecurityPolicy) domain.Validation
}

// VariableRepository defines operations for template variable persistence
type VariableRepository interface {
	// Save stores template variables
	Save(ctx context.Context, variables map[string]domain.TemplateVariable) error
	
	// Get retrieves template variables by policy ID
	Get(ctx context.Context, policyID domain.PolicyID) (map[string]domain.TemplateVariable, error)
	
	// Update modifies existing variables
	Update(ctx context.Context, policyID domain.PolicyID, variables map[string]domain.TemplateVariable) error
	
	// Delete removes variables for a policy
	Delete(ctx context.Context, policyID domain.PolicyID) error
	
	// FindByCategory returns variables by category
	FindByCategory(ctx context.Context, policyID domain.PolicyID, category types.VariableCategory) ([]domain.TemplateVariable, error)
	
	// FindUnresolved returns unresolved required variables
	FindUnresolved(ctx context.Context, policyID domain.PolicyID) ([]domain.TemplateVariable, error)
}

// ValidationRepository defines operations for validation results persistence
type ValidationRepository interface {
	// Save stores validation results
	Save(ctx context.Context, result domain.Validation) error
	
	// Get retrieves validation results by policy ID
	Get(ctx context.Context, policyID domain.PolicyID, version string) (domain.Validation, error)
	
	// GetLatest returns the most recent validation
	GetLatest(ctx context.Context, policyID domain.PolicyID) (domain.Validation, error)
	
	// GetHistory returns validation history
	GetHistory(ctx context.Context, policyID domain.PolicyID, limit int) ([]domain.Validation, error)
	
	// Search finds validations by criteria
	Search(ctx context.Context, criteria ValidationSearchCriteria) ([]domain.Validation, error)
}

// GitHubRepository defines operations for GitHub integration
type GitHubRepository interface {
	// GetRepositoryInfo retrieves GitHub repository information
	GetRepositoryInfo(ctx context.Context, url string) (domain.GitHubInfo, error)
	
	// GetOwnerInfo retrieves GitHub owner information
	GetOwnerInfo(ctx context.Context, owner string) (domain.Organization, error)
	
	// GetSecurityAdvisories retrieves security advisories
	GetSecurityAdvisories(ctx context.Context, owner, repo string) ([]SecurityAdvisory, error)
	
	// GetContributors retrieves repository contributors
	GetContributors(ctx context.Context, owner, repo string) ([]GitHubContributor, error)
	
	// IsGitHubRepository checks if URL is a GitHub repository
	IsGitHubRepository(ctx context.Context, url string) bool
}

// BountyRepository defines operations for bounty program data
type BountyRepository interface {
	// Save stores bounty program information
	Save(ctx context.Context, bounty domain.BountyProgram) error
	
	// Get retrieves bounty program by organization
	Get(ctx context.Context, organization string) (domain.BountyProgram, error)
	
	// Update modifies bounty program
	Update(ctx context.Context, bounty domain.BountyProgram) error
	
	// GetLeaderboard returns bounty leaderboard
	GetLeaderboard(ctx context.Context, organization string) ([]domain.Researcher, error)
	
	// AddResearcher adds researcher to hall of fame
	AddResearcher(ctx context.Context, organization string, researcher domain.Researcher) error
}

// Search criteria types

// ProjectSearchCriteria defines project search parameters
type ProjectSearchCriteria struct {
	Name         string              `json:"name"`
	Type         types.ProjectType   `json:"type"`
	TechStack    []types.TechStack  `json:"tech_stack"`
	Languages    []string            `json:"languages"`
	Limit        int                 `json:"limit"`
	Offset       int                 `json:"offset"`
	SortBy       ProjectSortBy       `json:"sort_by"`
	SortOrder    SortOrder           `json:"sort_order"`
}

// ValidationSearchCriteria defines validation search parameters
type ValidationSearchCriteria struct {
	PolicyID     *domain.PolicyID           `json:"policy_id,omitempty"`
	Status       types.ValidationStatus      `json:"status"`
	ScoreFrom    *uint8                    `json:"score_from,omitempty"`
	ScoreTo      *uint8                    `json:"score_to,omitempty"`
	DateFrom     *time.Time                 `json:"date_from,omitempty"`
	DateTo       *time.Time                 `json:"date_to,omitempty"`
	Limit        int                        `json:"limit"`
	Offset       int                        `json:"offset"`
	SortBy       ValidationSortBy           `json:"sort_by"`
	SortOrder    SortOrder                  `json:"sort_order"`
}

// Sort and pagination types

// ProjectSortBy defines project sort options
type ProjectSortBy string

const (
	SortProjectByName        ProjectSortBy = "name"
	SortProjectByType       ProjectSortBy = "type"
	SortProjectByCreated    ProjectSortBy = "created_at"
	SortProjectByUpdated    ProjectSortBy = "updated_at"
	SortProjectByTechStack  ProjectSortBy = "tech_stack"
)

// ValidationSortBy defines validation sort options
type ValidationSortBy string

const (
	SortValidationByScore     ValidationSortBy = "score"
	SortValidationByDate      ValidationSortBy = "checked_at"
	SortValidationByStatus    ValidationSortBy = "status"
	SortValidationByIssues     ValidationSortBy = "issues_count"
)

// SortOrder defines sort direction
type SortOrder string

const (
	SortOrderAscending  SortOrder = "asc"
	SortOrderDescending SortOrder = "desc"
)

// GitHub-specific types

// SecurityAdvisory represents a GitHub security advisory
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

// CVE represents a CVE entry
type CVE struct {
	ID        string `json:"id"`
	Severity  string `json:"severity"`
	CWEs      []CWE  `json:"cwes"`
	Published string `json:"published"`
	Updated   string `json:"updated"`
}

// CWE represents a Common Weakness Enumeration
type CWE struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Vulnerability represents a vulnerability
type Vulnerability struct {
	Package       string   `json:"package"`
	Ecosystem     string   `json:"ecosystem"`
	Severity      string   `json:"severity"`
	Requirements  []string `json:"requirements"`
	FixedVersions []string `json:"fixed_versions"`
}

// GitHubContributor represents a GitHub contributor
type GitHubContributor struct {
	Login         string    `json:"login"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Type          string    `json:"type"`
	Contributions  uint32    `json:"contributions"`
	FirstContribution time.Time `json:"first_contribution"`
	LastContribution  time.Time `json:"last_contribution"`
}

// RepositoryConfig represents repository configuration
type RepositoryConfig struct {
	Type        string            `json:"type"`
	Connection  string            `json:"connection"`
	Options     map[string]interface{} `json:"options"`
	Cache       CacheConfig       `json:"cache"`
}

// CacheConfig represents cache configuration
type CacheConfig struct {
	Enabled     bool          `json:"enabled"`
	TTL         time.Duration `json:"ttl"`
	MaxSize     int64         `json:"max_size"`
	EvictionPolicy string      `json:"eviction_policy"`
}

// UnitOfWork represents transactional operations
type UnitOfWork interface {
	// Begin starts a transaction
	Begin(ctx context.Context) error
	
	// Commit commits the transaction
	Commit() error
	
	// Rollback rolls back the transaction
	Rollback() error
	
	// GetRepositories returns repositories within the unit of work
	GetRepositories() Repositories
}

// Repositories groups all repositories for unit of work
type Repositories struct {
	Project   ProjectRepository
	Policy    PolicyRepository
	Variable  VariableRepository
	Validation ValidationRepository
	GitHub    GitHubRepository
	Bounty    BountyRepository
}