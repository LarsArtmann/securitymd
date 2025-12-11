package errors

import (
	"fmt"
	"time"

	"github.com/LarsArtmann/template-SECURITY/v2/internal/types"
)

// BaseError represents the base error structure
type BaseError struct {
	Code      string                 `json:"code"`
	Message   string                 `json:"message"`
	Field     string                 `json:"field,omitempty"`
	Value     string                 `json:"value,omitempty"`
	Component string                 `json:"component"`
	Timestamp time.Time              `json:"timestamp"`
	Context   map[string]interface{} `json:"context,omitempty"`
}

// Error implements error interface
func (be BaseError) Error() string {
	return fmt.Sprintf("[%s] %s", be.Code, be.Message)
}

// WithContext adds context to error
func (be BaseError) WithContext(key string, value interface{}) BaseError {
	if be.Context == nil {
		be.Context = make(map[string]interface{})
	}
	be.Context[key] = value
	return be
}

// NewBaseError creates a new base error
func NewBaseError(code, message, component string) BaseError {
	return BaseError{
		Code:      code,
		Message:   message,
		Component: component,
		Timestamp: time.Now(),
	}
}

// ValidationError represents validation errors
type ValidationError struct {
	BaseError
	Level  types.ValidationLevel `json:"level"`
	Line   uint32                `json:"line,omitempty"`
	Column uint32                `json:"column,omitempty"`
	File   string                `json:"file,omitempty"`
}

// NewValidationError creates a new validation error
func NewValidationError(code, message, field string, level types.ValidationLevel) ValidationError {
	return ValidationError{
		BaseError: NewBaseError(code, message, "validator"),
		Level:     level,
		Field:     field,
	}
}

// WithLocation adds file location to validation error
func (ve ValidationError) WithLocation(file string, line, column uint32) ValidationError {
	ve.File = file
	ve.Line = line
	ve.Column = column
	return ve
}

// ConfigurationError represents configuration errors
type ConfigurationError struct {
	BaseError
	Section string `json:"section"`
	Key     string `json:"key"`
}

// NewConfigurationError creates a new configuration error
func NewConfigurationError(code, message, section, key string) ConfigurationError {
	return ConfigurationError{
		BaseError: NewBaseError(code, message, "config"),
		Section:   section,
		Key:       key,
	}
}

// ProjectDetectionError represents project detection errors
type ProjectDetectionError struct {
	BaseError
	Source string `json:"source"`
	Method string `json:"method"`
}

// NewProjectDetectionError creates a new project detection error
func NewProjectDetectionError(code, message, source, method string) ProjectDetectionError {
	return ProjectDetectionError{
		BaseError: NewBaseError(code, message, "detector"),
		Source:    source,
		Method:    method,
	}
}

// TemplateError represents template-related errors
type TemplateError struct {
	BaseError
	Template string `json:"template"`
	Variable string `json:"variable,omitempty"`
}

// NewTemplateError creates a new template error
func NewTemplateError(code, message, template string) TemplateError {
	return TemplateError{
		BaseError: NewBaseError(code, message, "template"),
		Template:  template,
	}
}

// WithVariable adds variable info to template error
func (te TemplateError) WithVariable(variable string) TemplateError {
	te.Variable = variable
	return te
}

// GitHubError represents GitHub integration errors
type GitHubError struct {
	BaseError
	Repository string `json:"repository"`
	Endpoint   string `json:"endpoint"`
	Status     int    `json:"status,omitempty"`
}

// NewGitHubError creates a new GitHub error
func NewGitHubError(code, message, repository, endpoint string) GitHubError {
	return GitHubError{
		BaseError:  NewBaseError(code, message, "github"),
		Repository: repository,
		Endpoint:   endpoint,
	}
}

// WithStatus adds HTTP status to GitHub error
func (ghe GitHubError) WithStatus(status int) GitHubError {
	ghe.Status = status
	return ghe
}

// FileOperationError represents file operation errors
type FileOperationError struct {
	BaseError
	Operation  string `json:"operation"`
	FilePath   string `json:"file_path"`
	Permission bool   `json:"permission,omitempty"`
}

// NewFileOperationError creates a new file operation error
func NewFileOperationError(code, message, operation, filePath string) FileOperationError {
	return FileOperationError{
		BaseError: NewBaseError(code, message, "file_ops"),
		Operation: operation,
		FilePath:  filePath,
	}
}

// WithPermission adds permission info to file error
func (foe FileOperationError) WithPermission(permission bool) FileOperationError {
	foe.Permission = permission
	return foe
}

// RegistryError represents template variable registry errors
type RegistryError struct {
	BaseError
	Variable   string             `json:"variable,omitempty"`
	Registry   string             `json:"registry"`
	Constraint RegistryConstraint `json:"constraint"`
}

// RegistryConstraint represents registry constraint types
type RegistryConstraint string

const (
	ConstraintPlaceholderInvalid RegistryConstraint = "placeholder_invalid"
	ConstraintCategoryInvalid    RegistryConstraint = "category_invalid"
	ConstraintPriorityInvalid    RegistryConstraint = "priority_invalid"
	ConstraintSourceInvalid      RegistryConstraint = "source_invalid"
	ConstraintDuplicate          RegistryConstraint = "duplicate"
	ConstraintRequired           RegistryConstraint = "required"
)

// NewRegistryError creates a new registry error
func NewRegistryError(code, message, registry string, constraint RegistryConstraint) RegistryError {
	return RegistryError{
		BaseError:  NewBaseError(code, message, "registry"),
		Registry:   registry,
		Constraint: constraint,
	}
}

// WithVariable adds variable info to registry error
func (re RegistryError) WithVariable(variable string) RegistryError {
	re.Variable = variable
	return re
}

// CompilationError represents compilation/build errors
type CompilationError struct {
	BaseError
	Target     string `json:"target"`
	Step       string `json:"step"`
	Diagnostic string `json:"diagnostic"`
}

// NewCompilationError creates a new compilation error
func NewCompilationError(code, message, target, step string) CompilationError {
	return CompilationError{
		BaseError: NewBaseError(code, message, "compiler"),
		Target:    target,
		Step:      step,
	}
}

// WithDiagnostic adds diagnostic info to compilation error
func (ce CompilationError) WithDiagnostic(diagnostic string) CompilationError {
	ce.Diagnostic = diagnostic
	return ce
}

// Error codes for consistent error handling
const (
	// Validation error codes
	ErrCodeRequiredSection     = "VAL001"
	ErrCodeMissingContact      = "VAL002"
	ErrCodeInvalidEmail        = "VAL003"
	ErrCodeUnresolvedVariable  = "VAL004"
	ErrCodeInsufficientContent = "VAL005"

	// Configuration error codes
	ErrCodeInvalidPolicyType   = "CFG001"
	ErrCodeMissingOrganization = "CFG002"
	ErrCodeInvalidOutputFormat = "CFG003"

	// Project detection error codes
	ErrCodeGitConfigAccess = "DET001"
	ErrCodeFileReadAccess  = "DET002"
	ErrCodeInvalidGitURL   = "DET003"

	// Template error codes
	ErrCodeTemplateNotFound  = "TPL001"
	ErrCodeTemplateCorrupted = "TPL002"
	ErrCodeVariableCycle     = "TPL003"

	// GitHub error codes
	ErrCodeGitHubAPIFailed    = "GH001"
	ErrCodeRepositoryNotFound = "GH002"
	ErrCodeInvalidToken       = "GH003"

	// File operation error codes
	ErrCodeFileNotFound     = "FILE001"
	ErrCodePermissionDenied = "FILE002"
	ErrCodeDiskFull         = "FILE003"

	// Registry error codes
	ErrCodeInvalidPlaceholder = "REG001"
	ErrCodeDuplicateVariable  = "REG002"
	ErrCodeInvalidCategory    = "REG003"

	// Compilation error codes
	ErrCodeBuildFailed       = "BUILD001"
	ErrCodeDependencyMissing = "BUILD002"
	ErrCodeTestFailed        = "BUILD003"
)

// ErrorHandler provides centralized error handling
type ErrorHandler struct {
	logErrors     bool
	errorChannels map[string]chan error
}

// NewErrorHandler creates a new error handler
func NewErrorHandler(logErrors bool) *ErrorHandler {
	return &ErrorHandler{
		logErrors:     logErrors,
		errorChannels: make(map[string]chan error),
	}
}

// Handle processes an error
func (eh *ErrorHandler) Handle(err error) {
	if err == nil {
		return
	}

	// Log error if enabled
	if eh.logErrors {
		// In a real implementation, use proper logging
		fmt.Printf("ERROR: %s\n", err.Error())
	}

	// Send to appropriate channel if registered
	if typedErr, ok := err.(interface{ Component() string }); ok {
		component := typedErr.Component()
		if channel, exists := eh.errorChannels[component]; exists {
			select {
			case channel <- err:
			default:
				// Channel full, handle overflow
			}
		}
	}
}

// RegisterChannel registers an error channel for a component
func (eh *ErrorHandler) RegisterChannel(component string, ch chan error) {
	eh.errorChannels[component] = ch
}

// Component method for typed errors
func (ve ValidationError) Component() string        { return ve.Component }
func (ce ConfigurationError) Component() string     { return ce.Component }
func (pde ProjectDetectionError) Component() string { return pde.Component }
func (te TemplateError) Component() string          { return te.Component }
func (ghe GitHubError) Component() string           { return ghe.Component }
func (foe FileOperationError) Component() string    { return foe.Component }
func (re RegistryError) Component() string          { return re.Component }
func (ce CompilationError) Component() string       { return ce.Component }

// ErrorCollection represents multiple errors
type ErrorCollection struct {
	Errors []error `json:"errors"`
}

// Error implements error interface
func (ec ErrorCollection) Error() string {
	if len(ec.Errors) == 0 {
		return "no errors"
	}
	if len(ec.Errors) == 1 {
		return ec.Errors[0].Error()
	}
	return fmt.Sprintf("%d errors: %s", len(ec.Errors), ec.Errors[0].Error())
}

// Add adds an error to the collection
func (ec *ErrorCollection) Add(err error) {
	if err != nil {
		ec.Errors = append(ec.Errors, err)
	}
}

// HasErrors returns true if collection has errors
func (ec ErrorCollection) HasErrors() bool {
	return len(ec.Errors) > 0
}

// Count returns number of errors
func (ec ErrorCollection) Count() int {
	return len(ec.Errors)
}

// Clear removes all errors
func (ec *ErrorCollection) Clear() {
	ec.Errors = nil
}
