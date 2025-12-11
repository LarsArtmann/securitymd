package domain

import (
	"time"

	"github.com/LarsArtmann/template-SECURITY/v2/internal/types"
	"github.com/LarsArtmann/template-SECURITY/v2/internal/errors"
)

// SecurityPolicy represents a complete security policy domain model
type SecurityPolicy struct {
	ID           PolicyID                     `json:"id"`
	Type         types.PolicyType             `json:"type"`
	Organization Organization                `json:"organization"`
	Contact      Contact                      `json:"contact"`
	Versions     []Version                    `json:"versions"`
	Validation   Validation                   `json:"validation"`
	Metadata     PolicyMetadata               `json:"metadata"`
	Content      string                       `json:"content"`
	GeneratedAt  time.Time                    `json:"generated_at"`
	UpdatedAt    time.Time                    `json:"updated_at"`
}

// PolicyID represents a unique policy identifier
type PolicyID struct {
	Value string `json:"value"`
}

// IsValid checks if PolicyID is valid
func (pid PolicyID) IsValid() bool {
	return len(pid.Value) > 0 && len(pid.Value) <= 100
}

// Organization represents organization/domain information
type Organization struct {
	Name         string          `json:"name"`
	Domain       string          `json:"domain"`
	Email        string          `json:"email"`
	Website      string          `json:"website"`
	Type         OrganizationType `json:"type"`
	Size         OrgSize         `json:"size"`
	Industry     string          `json:"industry"`
	GitHub       GitHubInfo      `json:"github"`
}

// OrganizationType represents organization classification
type OrganizationType string

const (
	OrgTypeStartup    OrganizationType = "startup"
	OrgTypeSME        OrganizationType = "sme"
	OrgTypeEnterprise  OrganizationType = "enterprise"
	OrgTypeOpenSource  OrganizationType = "opensource"
	OrgTypePersonal    OrganizationType = "personal"
)

// IsValid checks if organization type is valid
func (ot OrganizationType) IsValid() bool {
	switch ot {
	case OrgTypeStartup, OrgTypeSME, OrgTypeEnterprise,
		 OrgTypeOpenSource, OrgTypePersonal:
		return true
	default:
		return false
	}
}

// OrgSize represents organization size
type OrgSize uint16

const (
	SizeTiny    OrgSize = 1-10
	SizeSmall   OrgSize = 11-50
	SizeMedium  OrgSize = 51-200
	SizeLarge   OrgSize = 201-1000
	SizeHuge    OrgSize = 1001-10000
	SizeMassive OrgSize = 10001+
)

// String implements fmt.Stringer
func (os OrgSize) String() string {
	switch {
	case os <= 10:
		return "1-10"
	case os <= 50:
		return "11-50"
	case os <= 200:
		return "51-200"
	case os <= 1000:
		return "201-1000"
	case os <= 10000:
		return "1001-10000"
	default:
		return "10000+"
	}
}

// GitHubInfo represents GitHub-specific organization information
type GitHubInfo struct {
	Owner           string                `json:"owner"`
	Repository      string                `json:"repository"`
	FullName        string                `json:"full_name"`
	URL             string                `json:"url"`
	Host            string                `json:"host"`
	DefaultBranch   string                `json:"default_branch"`
	Protocol        types.GitProtocol      `json:"protocol"`
	Status          types.RepositoryStatus `json:"status"`
	IsGitHub        bool                  `json:"is_github"`
	URLs            GitHubURLs            `json:"urls"`
}

// GitHubURLs represents various GitHub URLs
type GitHubURLs struct {
	Repository    string `json:"repository"`
	Issues        string `json:"issues"`
	Security      string `json:"security"`
	Advisories    string `json:"advisories"`
	Pulse         string `json:"pulse"`
	Network       string `json:"network"`
	Settings      string `json:"settings"`
	Actions       string `json:"actions"`
	Releases      string `json:"releases"`
	Wiki          string `json:"wiki"`
	API           string `json:"api"`
}

// Contact represents contact information for security reports
type Contact struct {
	Primary         ContactMethod `json:"primary"`
	Alternative     ContactMethod `json:"alternative"`
	SecurityTeam    ContactMethod `json:"security_team"`
	Responsibilities Responsibilities `json:"responsibilities"`
}

// ContactMethod represents a contact method
type ContactMethod struct {
	Type        ContactType `json:"type"`
	Value       string      `json:"value"`
	Description string      `json:"description"`
	Encrypted   bool        `json:"encrypted"`
}

// ContactType represents different contact types
type ContactType string

const (
	ContactTypeEmail ContactType = "email"
	ContactTypeWeb   ContactType = "web"
	ContactTypePGP   ContactType = "pgp"
	ContactTypeSlack ContactType = "slack"
)

// IsValid checks if contact type is valid
func (ct ContactType) IsValid() bool {
	switch ct {
	case ContactTypeEmail, ContactTypeWeb, ContactTypePGP, ContactTypeSlack:
		return true
	default:
		return false
	}
}

// Responsibilities represents security response responsibilities
type Responsibilities struct {
	ResponseTime     time.Duration `json:"response_time"`
	InvestigationTime time.Duration `json:"investigation_time"`
	PatchTime        time.Duration `json:"patch_time"`
	DisclosureTime   time.Duration `json:"disclosure_time"`
	SafeHarbor       bool          `json:"safe_harbor"`
	CoordinateCVE    bool          `json:"coordinate_cve"`
}

// Version represents software version support information
type Version struct {
	Name            string    `json:"name"`
	SemanticVersion  string    `json:"semantic_version"`
	SupportedUntil   time.Time `json:"supported_until"`
	Status          VersionStatus `json:"status"`
	IsLatest        bool      `json:"is_latest"`
	IsPrevious      bool      `json:"is_previous"`
}

// VersionStatus represents version status
type VersionStatus string

const (
	StatusSupported    VersionStatus = "supported"
	StatusDeprecated VersionStatus = "deprecated"
	StatusEOL        VersionStatus = "eol"
	StatusSecurity    VersionStatus = "security"
)

// IsValid checks if version status is valid
func (vs VersionStatus) IsValid() bool {
	switch vs {
	case StatusSupported, StatusDeprecated, StatusEOL, StatusSecurity:
		return true
	default:
		return false
	}
}

// Validation represents validation results
type Validation struct {
	Status       types.ValidationStatus `json:"status"`
	Score        uint8                `json:"score"`  // 0-100
	Issues       []ValidationIssue     `json:"issues"`
	Warnings     []ValidationWarning   `json:"warnings"`
	CheckedAt    time.Time            `json:"checked_at"`
	Quality      QualityMetrics        `json:"quality"`
}

// ValidationIssue represents a validation error
type ValidationIssue struct {
	ID          string              `json:"id"`
	Field       string              `json:"field"`
	Category    ValidationCategory  `json:"category"`
	Level       types.ValidationLevel `json:"level"`
	Message     string              `json:"message"`
	Suggestion  string              `json:"suggestion"`
	Line        uint32              `json:"line,omitempty"`
	Column      uint32              `json:"column,omitempty"`
}

// ValidationCategory represents validation issue categories
type ValidationCategory string

const (
	CategoryStructure ValidationCategory = "structure"
	CategoryContent ValidationCategory = "content"
	CategorySecurity ValidationCategory = "security"
	CategoryContact ValidationCategory = "contact"
	CategoryLegal ValidationCategory = "legal"
)

// IsValid checks if validation category is valid
func (vc ValidationCategory) IsValid() bool {
	switch vc {
	case CategoryStructure, CategoryContent, CategorySecurity,
		 CategoryContact, CategoryLegal:
		return true
	default:
		return false
	}
}

// ValidationWarning represents a validation warning
type ValidationWarning struct {
	ID          string              `json:"id"`
	Field       string              `json:"field"`
	Category    ValidationCategory  `json:"category"`
	Level       types.ValidationLevel `json:"level"`
	Message     string              `json:"message"`
	Suggestion  string              `json:"suggestion"`
	Line        uint32              `json:"line,omitempty"`
	Column      uint32              `json:"column,omitempty"`
}

// QualityMetrics represents content quality metrics
type QualityMetrics struct {
	Words         uint16 `json:"words"`
	Lines         uint16 `json:"lines"`
	Characters    uint16 `json:"characters"`
	Sections      uint8  `json:"sections"`
	Readability   uint8  `json:"readability"`    // 0-100
	Completeness  uint8  `json:"completeness"`  // 0-100
	BestPractices uint8  `json:"best_practices"` // 0-100
}

// PolicyMetadata represents policy metadata
type PolicyMetadata struct {
	TemplateVersion  string      `json:"template_version"`
	GeneratorVersion  string      `json:"generator_version"`
	VariableCount    uint16      `json:"variable_count"`
	ResolvedCount    uint16      `json:"resolved_count"`
	Variables        []Variable  `json:"variables"`
	Project          ProjectInfo `json:"project"`
}

// Variable represents a template variable
type Variable struct {
	Placeholder string            `json:"placeholder"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Value       string             `json:"value"`
	Priority    types.Priority     `json:"priority"`
	Category    types.VariableCategory `json:"category"`
	Detected    bool               `json:"detected"`
	Overridden  bool               `json:"overridden"`
	Required    bool               `json:"required"`
}

// ProjectInfo represents project information
type ProjectInfo struct {
	Name         string               `json:"name"`
	Type         types.ProjectType    `json:"type"`
	TechStack    []types.TechStack   `json:"tech_stack"`
	Languages    []string             `json:"languages"`
	Frameworks   []string             `json:"frameworks"`
	Dependencies  []string             `json:"dependencies"`
	Size         ProjectSize          `json:"size"`
	Complexity   ProjectComplexity    `json:"complexity"`
}

// ProjectSize represents project size classification
type ProjectSize string

const (
	SizeMicro ProjectSize = "micro"
	SizeSmall ProjectSize = "small"
	SizeMedium ProjectSize = "medium"
	SizeLarge ProjectSize = "large"
	SizeXLarge ProjectSize = "xlarge"
)

// IsValid checks if project size is valid
func (ps ProjectSize) IsValid() bool {
	switch ps {
	case SizeMicro, SizeSmall, SizeMedium, SizeLarge, SizeXLarge:
		return true
	default:
		return false
	}
}

// ProjectComplexity represents project complexity classification
type ProjectComplexity string

const (
	ComplexitySimple    ProjectComplexity = "simple"
	ComplexityModerate  ProjectComplexity = "moderate"
	ComplexityComplex   ProjectComplexity = "complex"
	ComplexityVeryComplex ProjectComplexity = "very_complex"
)

// IsValid checks if project complexity is valid
func (pc ProjectComplexity) IsValid() bool {
	switch pc {
	case ComplexitySimple, ComplexityModerate, ComplexityComplex, ComplexityVeryComplex:
		return true
	default:
		return false
	}
}

// BountyProgram represents bounty program information
type BountyProgram struct {
	Enabled    bool                    `json:"enabled"`
	URL        string                  `json:"url"`
	Platform   BountyPlatform          `json:"platform"`
	Rewards    map[types.BountySeverity]Reward `json:"rewards"`
	Process    BountyProcess           `json:"process"`
	Policy     BountyPolicy            `json:"policy"`
	HallOfFame []Researcher           `json:"hall_of_fame"`
}

// BountyPlatform represents bounty platform types
type BountyPlatform string

const (
	PlatformGitHub BountyPlatform = "github"
	PlatformHackerOne BountyPlatform = "hackerone"
	PlatformBugcrowd BountyPlatform = "bugcrowd"
	PlatformIntigriti BountyPlatform = "intigriti"
	PlatformYesWeHack BountyPlatform = "yeswehack"
	PlatformSelfHosted BountyPlatform = "self_hosted"
)

// IsValid checks if bounty platform is valid
func (bp BountyPlatform) IsValid() bool {
	switch bp {
	case PlatformGitHub, PlatformHackerOne, PlatformBugcrowd,
		 PlatformIntigriti, PlatformYesWeHack, PlatformSelfHosted:
		return true
	default:
		return false
	}
}

// Reward represents bounty reward information
type Reward struct {
	Amount    uint32 `json:"amount"`
	Currency  string `json:"currency"`
	MinAmount uint32 `json:"min_amount,omitempty"`
	MaxAmount uint32 `json:"max_amount,omitempty"`
}

// BountyProcess represents bounty submission process
type BountyProcess struct {
	SubmissionURL   string        `json:"submission_url"`
	TriageTime      time.Duration `json:"triage_time"`
	ResponseTime    time.Duration `json:"response_time"`
	PaymentTime     time.Duration `json:"payment_time"`
	Requirements    []string      `json:"requirements"`
	Exclusions      []string      `json:"exclusions"`
}

// BountyPolicy represents bounty policy terms
type BountyPolicy struct {
	Scope           []string `json:"scope"`
	OutOfScope      []string `json:"out_of_scope"`
	SafeHarbor      bool     `json:"safe_harbor"`
	DisclosurePolicy string   `json:"disclosure_policy"`
	PublicDisclosure bool     `json:"public_disclosure"`
}

// Researcher represents security researcher information
type Researcher struct {
	Name         string    `json:"name"`
	Username     string    `json:"username"`
	Platform     string    `json:"platform"`
	Findings     uint16    `json:"findings"`
	HighImpact   uint16    `json:"high_impact"`
	FirstFinding time.Time `json:"first_finding"`
	LastFinding  time.Time `json:"last_finding"`
}

// IsValid checks if SecurityPolicy is valid
func (sp SecurityPolicy) IsValid() bool {
	return sp.ID.IsValid() &&
		sp.Type.IsValid() &&
		len(sp.Organization.Name) > 0 &&
		len(sp.Contact.Primary.Value) > 0 &&
		len(sp.Versions) > 0
}