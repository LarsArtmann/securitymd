package domain

import (
	"time"

	"github.com/LarsArtmann/template-SECURITY/v2/internal/types"
	"github.com/LarsArtmann/template-SECURITY/v2/internal/errors"
)

// TemplateVariable represents a template variable with type safety
type TemplateVariable struct {
	Placeholder string                    `json:"placeholder"`
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	Value       string                    `json:"value"`
	Priority    types.Priority           `json:"priority"`
	Category    types.VariableCategory  `json:"category"`
	Detected    bool                      `json:"detected"`
	Overridden  bool                      `json:"overridden"`
	Required    bool                      `json:"required"`
	Source      VariableSource            `json:"source"`
	CreatedAt   time.Time                  `json:"created_at"`
	UpdatedAt   time.Time                  `json:"updated_at"`
}

// VariableSource represents where a variable comes from
type VariableSource string

const (
	SourceAutoDetected VariableSource = "auto_detected"
	SourceGitConfig    VariableSource = "git_config"
	SourceProject      VariableSource = "project_files"
	SourceGitHub       VariableSource = "github_api"
	SourceUserInput    VariableSource = "user_input"
	SourceDefault      VariableSource = "default_value"
)

// IsValid checks if variable source is valid
func (vs VariableSource) IsValid() bool {
	switch vs {
	case SourceAutoDetected, SourceGitConfig, SourceProject,
		 SourceGitHub, SourceUserInput, SourceDefault:
		return true
	default:
		return false
	}
}

// IsValid checks if TemplateVariable is valid
func (tv TemplateVariable) IsValid() bool {
	return len(tv.Placeholder) > 0 &&
		len(tv.Name) > 0 &&
		tv.Category.IsValid() &&
		tv.Source.IsValid() &&
		(tv.Placeholder[0] == '{' && tv.Placeholder[1] == '{' &&
		 tv.Placeholder[len(tv.Placeholder)-2] == '}' && tv.Placeholder[len(tv.Placeholder)-1] == '}')
}

// TemplateVariableRegistry manages template variables with type safety
type TemplateVariableRegistry struct {
	variables map[string]TemplateVariable
	sources   map[VariableSource]uint16
	updated   time.Time
}

// NewTemplateVariableRegistry creates a new registry
func NewTemplateVariableRegistry() *TemplateVariableRegistry {
	return &TemplateVariableRegistry{
		variables: make(map[string]TemplateVariable),
		sources:   make(map[VariableSource]uint16),
		updated:   time.Now(),
	}
}

// Register adds a variable to the registry
func (tvr *TemplateVariableRegistry) Register(variable TemplateVariable) error {
	if !variable.IsValid() {
		return ValidationError{
			Field:   "TemplateVariable",
			Value:   variable.Placeholder,
			Message: "variable is not valid",
			Level:   types.ValidationLevelError,
		}
	}

	// Check if variable already exists
	if existing, exists := tvr.variables[variable.Placeholder]; exists {
		// Override only if new variable has higher priority
		if variable.Priority <= existing.Priority {
			// Mark as overridden but keep existing
			existing.Overridden = true
			existing.UpdatedAt = time.Now()
			tvr.variables[variable.Placeholder] = existing
		} else {
			// Replace with higher priority variable
			variable.Overridden = existing.Overridden
			tvr.variables[variable.Placeholder] = variable
		}
	} else {
		// Add new variable
		variable.CreatedAt = time.Now()
		variable.UpdatedAt = time.Now()
		tvr.variables[variable.Placeholder] = variable
	}

	// Update source count
	tvr.sources[variable.Source]++
	tvr.updated = time.Now()

	return nil
}

// Get retrieves a variable by placeholder
func (tvr *TemplateVariableRegistry) Get(placeholder string) (TemplateVariable, bool) {
	variable, exists := tvr.variables[placeholder]
	return variable, exists
}

// GetAll returns all variables
func (tvr *TemplateVariableRegistry) GetAll() map[string]TemplateVariable {
	result := make(map[string]TemplateVariable)
	for k, v := range tvr.variables {
		result[k] = v
	}
	return result
}

// GetByCategory returns variables by category
func (tvr *TemplateVariableRegistry) GetByCategory(category types.VariableCategory) []TemplateVariable {
	var result []TemplateVariable
	for _, variable := range tvr.variables {
		if variable.Category == category {
			result = append(result, variable)
		}
	}
	return result
}

// GetByPriority returns variables by priority level
func (tvr *TemplateVariableRegistry) GetByPriority(priority types.Priority) []TemplateVariable {
	var result []TemplateVariable
	for _, variable := range tvr.variables {
		if variable.Priority == priority {
			result = append(result, variable)
		}
	}
	return result
}

// GetRequired returns required variables
func (tvr *TemplateVariableRegistry) GetRequired() []TemplateVariable {
	var result []TemplateVariable
	for _, variable := range tvr.variables {
		if variable.Required {
			result = append(result, variable)
		}
	}
	return result
}

// GetUnresolved returns variables with empty values
func (tvr *TemplateVariableRegistry) GetUnresolved() []TemplateVariable {
	var result []TemplateVariable
	for _, variable := range tvr.variables {
		if len(variable.Value) == 0 && variable.Required {
			result = append(result, variable)
		}
	}
	return result
}

// Validate checks all variables for validity and completeness
func (tvr *TemplateVariableRegistry) Validate() Validation {
	var issues []ValidationIssue
	var warnings []ValidationWarning

	// Check for unresolved required variables
	for _, variable := range tvr.GetUnresolved() {
		issues = append(issues, ValidationIssue{
			ID:          "unresolved_required_variable",
			Field:       variable.Placeholder,
			Category:    CategoryContent,
			Level:       types.ValidationLevelError,
			Message:     "Required variable is unresolved",
			Suggestion:  "Provide a value for " + variable.Description,
		})
	}

	// Check for overridden variables
	for _, variable := range tvr.variables {
		if variable.Overridden {
			warnings = append(warnings, ValidationWarning{
				ID:          "overridden_variable",
				Field:       variable.Placeholder,
				Category:    CategoryContent,
				Level:       types.ValidationLevelWarning,
				Message:     "Variable was overridden by higher priority source",
				Suggestion:  "Review variable priorities if unexpected",
			})
		}
	}

	score := uint8(100)
	if len(issues) > 0 {
		score -= uint8(len(issues)) * 20
		if score < 0 {
			score = 0
		}
	}

	status := StatusValid
	if len(issues) > 0 {
		status = StatusInvalid
	} else if len(warnings) > 0 {
		status = StatusWarning
	}

	return Validation{
		Status:    status,
		Score:     score,
		Issues:    issues,
		Warnings:  warnings,
		CheckedAt: time.Now(),
		Quality: QualityMetrics{
			Sections:     uint8(len(tvr.variables)),
			Completeness: score,
		},
	}
}

// Merge merges another registry into this one
func (tvr *TemplateVariableRegistry) Merge(other *TemplateVariableRegistry) {
	for placeholder, variable := range other.variables {
		tvr.Register(variable)
	}
}

// ExportMap exports registry as simple map for template processing
func (tvr *TemplateVariableRegistry) ExportMap() map[string]string {
	result := make(map[string]string)
	for placeholder, variable := range tvr.variables {
		result[placeholder] = variable.Value
	}
	return result
}

// Stats returns registry statistics
func (tvr *TemplateVariableRegistry) Stats() RegistryStats {
	return RegistryStats{
		TotalVariables: uint16(len(tvr.variables)),
		TotalSources:   uint16(len(tvr.sources)),
		RequiredCount:  uint16(len(tvr.GetRequired())),
		UnresolvedCount: uint16(len(tvr.GetUnresolved())),
		CategoryCount:  uint16(len(tvr.GetCategories())),
		LastUpdated:    tvr.updated,
	}
}

// GetCategories returns all variable categories present
func (tvr *TemplateVariableRegistry) GetCategories() []types.VariableCategory {
	categories := make(map[types.VariableCategory]bool)
	for _, variable := range tvr.variables {
		categories[variable.Category] = true
	}

	var result []types.VariableCategory
	for category := range categories {
		result = append(result, category)
	}
	return result
}

// RegistryStats represents registry statistics
type RegistryStats struct {
	TotalVariables    uint16                `json:"total_variables"`
	TotalSources      uint16                `json:"total_sources"`
	RequiredCount     uint16                `json:"required_count"`
	UnresolvedCount  uint16                `json:"unresolved_count"`
	CategoryCount     uint16                `json:"category_count"`
	LastUpdated       time.Time             `json:"last_updated"`
	Sources          map[VariableSource]uint16 `json:"sources"`
}

// Clear removes all variables from registry
func (tvr *TemplateVariableRegistry) Clear() {
	tvr.variables = make(map[string]TemplateVariable)
	tvr.sources = make(map[VariableSource]uint16)
	tvr.updated = time.Now()
}

// Size returns the number of variables in the registry
func (tvr *TemplateVariableRegistry) Size() int {
	return len(tvr.variables)
}