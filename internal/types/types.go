package types

import "fmt"

// PolicyType represents different security policy types with strong typing
type PolicyType string

const (
	PolicyTypeGitHub           PolicyType = "github"
	PolicyTypeEnterprise       PolicyType = "enterprise" 
	PolicyTypeBugBounty        PolicyType = "bug-bounty"
	PolicyTypeIncidentResponse PolicyType = "incident-response"
	PolicyTypePrivacyPolicy    PolicyType = "privacy-policy"
)

// IsValid checks if policy type is valid
func (pt PolicyType) IsValid() bool {
	switch pt {
	case PolicyTypeGitHub, PolicyTypeEnterprise, PolicyTypeBugBounty,
		 PolicyTypeIncidentResponse, PolicyTypePrivacyPolicy:
		return true
	default:
		return false
	}
}

// String implements fmt.Stringer
func (pt PolicyType) String() string {
	return string(pt)
}

// OutputFormat represents validation output format
type OutputFormat string

const (
	FormatText      OutputFormat = "text"
	FormatJSON      OutputFormat = "json"
	FormatPrometheus OutputFormat = "prometheus"
)

// IsValid checks if output format is valid
func (of OutputFormat) IsValid() bool {
	switch of {
	case FormatText, FormatJSON, FormatPrometheus:
		return true
	default:
		return false
	}
}

// ValidationLevel represents validation severity
type ValidationLevel uint8

const (
	ValidationLevelInfo ValidationLevel = iota
	ValidationLevelWarning
	ValidationLevelError
	ValidationLevelCritical
)

// String implements fmt.Stringer
func (vl ValidationLevel) String() string {
	switch vl {
	case ValidationLevelInfo:
		return "info"
	case ValidationLevelWarning:
		return "warning"
	case ValidationLevelError:
		return "error"
	case ValidationLevelCritical:
		return "critical"
	default:
		return "unknown"
	}
}

// Priority represents feature/variable priority
type Priority uint8

const (
	PriorityLow Priority = iota
	PriorityMedium
	PriorityHigh
	PriorityCritical
)

// String implements fmt.Stringer
func (p Priority) String() string {
	switch p {
	case PriorityLow:
		return "low"
	case PriorityMedium:
		return "medium"
	case PriorityHigh:
		return "high"
	case PriorityCritical:
		return "critical"
	default:
		return "unknown"
	}
}

// VariableCategory represents template variable categories
type VariableCategory string

const (
	CategoryProject      VariableCategory = "project"
	CategoryOrganization VariableCategory = "organization"
	CategoryContact      VariableCategory = "contact"
	CategorySecurity     VariableCategory = "security"
	CategoryLegal        VariableCategory = "legal"
	CategoryTemporal     VariableCategory = "temporal"
)

// IsValid checks if category is valid
func (vc VariableCategory) IsValid() bool {
	switch vc {
	case CategoryProject, CategoryOrganization, CategoryContact,
		 CategorySecurity, CategoryLegal, CategoryTemporal:
		return true
	default:
		return false
	}
}

// ValidationStatus represents validation result status
type ValidationStatus uint8

const (
	StatusValid ValidationStatus = iota
	StatusInvalid
	StatusWarning
	StatusPending
)

// String implements fmt.Stringer
func (vs ValidationStatus) String() string {
	switch vs {
	case StatusValid:
		return "valid"
	case StatusInvalid:
		return "invalid"
	case StatusWarning:
		return "warning"
	case StatusPending:
		return "pending"
	default:
		return "unknown"
	}
}

// BountySeverity represents bounty severity levels
type BountySeverity string

const (
	SeverityCritical BountySeverity = "critical"
	SeverityHigh     BountySeverity = "high"
	SeverityMedium   BountySeverity = "medium"
	SeverityLow      BountySeverity = "low"
)

// IsValid checks if severity is valid
func (bs BountySeverity) IsValid() bool {
	switch bs {
	case SeverityCritical, SeverityHigh, SeverityMedium, SeverityLow:
		return true
	default:
		return false
	}
}

// GitProtocol represents different git remote protocols
type GitProtocol string

const (
	ProtocolHTTPS GitProtocol = "https"
	ProtocolSSH   GitProtocol = "ssh"
	ProtocolGit   GitProtocol = "git"
)

// IsValid checks if protocol is valid
func (gp GitProtocol) IsValid() bool {
	switch gp {
	case ProtocolHTTPS, ProtocolSSH, ProtocolGit:
		return true
	default:
		return false
	}
}

// RepositoryStatus represents repository status
type RepositoryStatus uint8

const (
	StatusPrivate RepositoryStatus = iota
	StatusPublic
	StatusArchived
	StatusFork
)

// String implements fmt.Stringer
func (rs RepositoryStatus) String() string {
	switch rs {
	case StatusPrivate:
		return "private"
	case StatusPublic:
		return "public"
	case StatusArchived:
		return "archived"
	case StatusFork:
		return "fork"
	default:
		return "unknown"
	}
}

// TechStack represents detected technology stacks
type TechStack string

const (
	TechStackGo       TechStack = "go"
	TechStackNode     TechStack = "node"
	TechStackPython   TechStack = "python"
	TechStackRust     TechStack = "rust"
	TechStackDocker   TechStack = "docker"
	TechStackK8s      TechStack = "kubernetes"
	TechStackReact    TechStack = "react"
	TechStackVue      TechStack = "vue"
	TechStackAngular  TechStack = "angular"
)

// IsValid checks if tech stack is valid
func (ts TechStack) IsValid() bool {
	switch ts {
	case TechStackGo, TechStackNode, TechStackPython, TechStackRust,
		 TechStackDocker, TechStackK8s, TechStackReact, TechStackVue, TechStackAngular:
		return true
	default:
		return false
	}
}

// ProjectType represents detected project types
type ProjectType string

const (
	TypeWeb      ProjectType = "web"
	TypeAPI      ProjectType = "api"
	TypeCLI      ProjectType = "cli"
	TypeLibrary  ProjectType = "library"
	TypeMobile   ProjectType = "mobile"
	TypeDesktop  ProjectType = "desktop"
	TypeIoT      ProjectType = "iot"
	TypeSaaS     ProjectType = "saas"
)

// IsValid checks if project type is valid
func (pt ProjectType) IsValid() bool {
	switch pt {
	case TypeWeb, TypeAPI, TypeCLI, TypeLibrary, TypeMobile,
		 TypeDesktop, TypeIoT, TypeSaaS:
		return true
	default:
		return false
	}
}

// Error types for better error handling

// ValidationError represents validation-specific errors
type ValidationError struct {
	Field   string
	Value   string
	Message string
	Level   ValidationLevel
}

// Error implements error interface
func (ve ValidationError) Error() string {
	return fmt.Sprintf("validation error [%s] on field '%s': %s", ve.Level, ve.Field, ve.Message)
}

// ConfigurationError represents configuration-specific errors
type ConfigurationError struct {
	Section string
	Key     string
	Message string
}

// Error implements error interface
func (ce ConfigurationError) Error() string {
	return fmt.Sprintf("configuration error in section '%s', key '%s': %s", ce.Section, ce.Key, ce.Message)
}

// ProjectDetectionError represents project detection errors
type ProjectDetectionError struct {
	Source string
	Method string
	Issue  string
}

// Error implements error interface
func (pde ProjectDetectionError) Error() string {
	return fmt.Sprintf("project detection error from source '%s' (%s): %s", pde.Source, pde.Method, pde.Issue)
}