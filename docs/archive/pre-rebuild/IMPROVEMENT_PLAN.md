# Template Security - Architecture Improvement Plan

## Analysis Summary

### Current Issues

1. **Duplicate Type Definitions**: `PolicyType` defined in multiple places
2. **Minimal Domain Model**: No clear separation of concerns
3. **No Tests**: Critical validation logic lacks test coverage
4. **Manual Template Processing**: Using string replacement instead of proper templating
5. **No Configuration Management**: Everything hardcoded or flag-based
6. **Poor Error Handling**: Basic error handling without structure
7. **Recent Refactor Mess**: Many deleted files and inconsistent state

## Execution Plan (Sorted by Impact vs Effort)

### High Impact, Low Effort (Immediate Wins)

1. **Fix Duplicate Type Definitions** (15 min) - Remove duplication, create single source of truth
2. **Add Basic Test Structure** (30 min) - Set up testing framework and add initial tests
3. **Implement Proper Templating** (45 min) - Use Go's `text/template` instead of string replacement
4. **Add Configuration File Support** (30 min) - Support for `.template-security.yaml` config file
5. **Improve Error Types** (30 min) - Create proper error types and handling

### High Impact, Medium Effort (Next Priority)

6. **Redesign Domain Model** (2 hours) - Create proper domain entities and value objects
7. **Add Repository Pattern** (1.5 hours) - Abstract file system operations
8. **Implement Service Layer** (2 hours) - Separate business logic from application logic
9. **Add Comprehensive Test Suite** (3 hours) - Unit tests for all major components
10. **Add Validation Rules Engine** (2 hours) - Make validation rules configurable

### Medium Impact, Low Effort (Quality of Life)

11. **Add Logging** (30 min) - Structured logging with different levels
12. **Improve CLI Help** (20 min) - Better help text and examples
13. **Add Auto-detection Improvements** (45 min) - Better project detection logic
14. **Add Output Formats** (30 min) - Support JSON output for CI/CD

### Medium Impact, Medium Effort (Future Enhancements)

15. **Add Plugin System** (3 hours) - Allow custom validation rules
16. **Add Template Repository** (2 hours) - Support multiple template sources
17. **Add Integration Tests** (2 hours) - End-to-end testing
18. **Add Performance Metrics** (1.5 hours) - Measure validation performance

### Low Impact, High Effort (Long-term)

19. **Add Web Interface** (5 hours) - Simple web UI for validation
20. **Add Database Persistence** (4 hours) - Store validation history
21. **Add API Server** (6 hours) - REST API for validation services

## Domain Model Improvements

### Current Types

```go
type PolicyType string
type PolicyConfig struct { ... }
type SecurityValidationResult struct { ... }
```

### Proposed Domain Model

```go
// Domain Entities
type SecurityPolicy struct {
    ID          PolicyID
    Project     Project
    Versions    []Version
    Contacts    []Contact
    Content     string
    ValidatedAt time.Time
}

type Project struct {
    Name         string
    Organization string
    Domain       string
    Repository   string
}

type Version struct {
    Number          string
    SupportedUntil  time.Time
    Status          VersionStatus
    IsLatest        bool
}

type Contact struct {
    Type         ContactType
    Value        string
    ResponseTime time.Duration
}

// Value Objects
type PolicyID string
type VersionStatus int
type ContactType int

// Repository Interfaces
type PolicyRepository interface {
    Find(path string) (*SecurityPolicy, error)
    Save(policy *SecurityPolicy) error
}

type ValidationRepository interface {
    Save(result ValidationResult) error
    FindByPolicyID(id PolicyID) ([]ValidationResult, error)
}

// Service Interfaces
type ValidationService interface {
    Validate(policy *SecurityPolicy) (*ValidationResult, error)
    GetValidationRules() []ValidationRule
}

type TemplateService interface {
    Generate(config TemplateConfig) (*SecurityPolicy, error)
    ListTemplates() ([]Template, error)
}
```

## Technology Improvements

### Current Dependencies

- `github.com/spf13/cobra` - CLI framework
- `github.com/fatih/color` - Terminal colors

### Recommended Additions

- `gopkg.in/yaml.v3` - Configuration file support
- `github.com/stretchr/testify` - Testing framework
- `github.com/sirupsen/logrus` - Structured logging
- `github.com/spf13/viper` - Configuration management
- `github.com/go-playground/validator/v10` - Input validation
- `github.com/fatih/structs` - Structure manipulation

### External Libraries to Consider

- `github.com/Masterminds/sprig` - Template function library
- `github.com/mattn/go-zglob` - File globbing
- `github.com/bmatcuk/doublestar` - Glob patterns
- `github.com/mitchellh/mapstructure` - Map to struct conversion
- `github.com/carlmjohnson/flowmatic` - Flow control

## Next Steps

1. Commit current changes
2. Execute high impact, low effort items
3. Set up proper CI/CD with testing
4. Gradually implement medium effort items
5. Plan for long-term architectural improvements
