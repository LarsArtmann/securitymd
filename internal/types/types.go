package internal

import "time"

// PolicyType represents different types of security policies
type PolicyType string

const (
	PolicyTypeGitHub     PolicyType = "github"
	PolicyTypeEnterprise PolicyType = "enterprise"
)

// VersionStatus represents the status of a version
type VersionStatus string

const (
	StatusSupported  VersionStatus = "supported"
	StatusDeprecated VersionStatus = "deprecated"
	StatusEOL        VersionStatus = "end-of-life"
)

// ContactType represents different types of contact methods
type ContactType string

const (
	ContactTypeEmail ContactType = "email"
	ContactTypeWeb   ContactType = "web"
	ContactTypeAPI   ContactType = "api"
)

// Version represents a software version with support information
type Version struct {
	Name            string        `json:"name"`
	SemanticVersion string        `json:"semantic_version"`
	SupportedUntil  time.Time     `json:"supported_until"`
	Status          VersionStatus `json:"status"`
	IsLatest        bool          `json:"is_latest"`
	IsPrevious      bool          `json:"is_previous"`
}

// Contact represents security contact information
type Contact struct {
	Type         ContactType `json:"type"`
	Value        string      `json:"value"`
	ResponseTime string      `json:"response_time"`
	Description  string      `json:"description"`
}

// SecurityPolicy represents a complete security policy
type SecurityPolicy struct {
	ID           string     `json:"id"`
	Project      Project    `json:"project"`
	Versions     []Version  `json:"versions"`
	Contacts     []Contact  `json:"contacts"`
	Content      string     `json:"content"`
	Type         PolicyType `json:"type"`
	ValidatedAt  *time.Time `json:"validated_at,omitempty"`
	LastModified time.Time  `json:"last_modified"`
}

// Project represents project information
type Project struct {
	Name         string `json:"name"`
	Organization string `json:"organization"`
	Domain       string `json:"domain"`
	Repository   string `json:"repository"`
	Description  string `json:"description"`
}

// ValidationResult represents the result of validating a security policy
type ValidationResult struct {
	ID          string    `json:"id"`
	PolicyID    string    `json:"policy_id"`
	Valid       bool      `json:"valid"`
	Errors      []Error   `json:"errors"`
	Warnings    []Warning `json:"warnings"`
	Score       int       `json:"score"` // 0-100
	ValidatedAt time.Time `json:"validated_at"`
	Metadata    any       `json:"metadata,omitempty"`
}

// Error represents a validation error
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
	Line    int    `json:"line,omitempty"`
}

// Warning represents a validation warning
type Warning struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
	Line    int    `json:"line,omitempty"`
}

// Template represents a security policy template
type Template struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Type        PolicyType         `json:"type"`
	Content     string             `json:"content"`
	Variables   []TemplateVariable `json:"variables"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

// TemplateVariable represents a variable in a template
type TemplateVariable struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Type         string `json:"type"`
	Required     bool   `json:"required"`
	DefaultValue string `json:"default_value"`
	Example      string `json:"example"`
}
