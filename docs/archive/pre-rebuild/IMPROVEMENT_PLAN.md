# Template Security - Architecture Improvement Plan

> **Superseded** — plan against the deleted tree; items either shipped in some later era or are NOT-DO by the rebuild's design. All numbered items carry inline verdicts.

## Analysis Summary

### Current Issues

~~1. **Duplicate Type Definitions**: `PolicyType` defined in multiple places~~ moot — old tree deleted at 1c55390; single model now (finding.Finding)
~~2. **Minimal Domain Model**: No clear separation of concerns~~ moot — deliberate: no parallel domain model in the rebuild
~~3. **No Tests**: Critical validation logic lacks test coverage~~ done — three suites (unit, provider loop, BDD), all green
~~4. **Manual Template Processing**: Using string replacement instead of proper templating~~ moot — Go text/template now (embedded, dogfood-tested)
~~5. **No Configuration Management**: Everything hardcoded or flag-based~~ NOT-DO — config system removed on purpose
~~6. **Poor Error Handling**: Basic error handling without structure~~ done — finding.FindingError family
~~7. **Recent Refactor Mess**: Many deleted files and inconsistent state~~ moot — clean tree since the rebuild

## Execution Plan (Sorted by Impact vs Effort)

### High Impact, Low Effort (Immediate Wins)

~~1. **Fix Duplicate Type Definitions** (15 min) - Remove duplication, create single source of truth~~ moot — types deleted at 1c55390
~~2. **Add Basic Test Structure** (30 min) - Set up testing framework and add initial tests~~ done — pkg tests + acceptance suite
~~3. **Implement Proper Templating** (45 min) - Use Go's `text/template` instead of string replacement~~ done — embedded text/template
~~4. **Add Configuration File Support** (30 min) - Support for `.template-security.yaml` config file~~ NOT-DO — config removed entirely (breaking change, documented)
~~5. **Improve Error Types** (30 min) - Create proper error types and handling~~ done — finding error family

### High Impact, Medium Effort (Next Priority)

~~6. **Redesign Domain Model** (2 hours) - Create proper domain entities and value objects~~ NOT-DO — finding.Finding is the model
~~7. **Add Repository Pattern** (1.5 hours) - Abstract file system operations~~ NOT-DO — direct file I/O, tested against real temp repos
~~8. **Implement Service Layer** (2 hours) - Separate business logic from application logic~~ NOT-DO — purged in Dec 2025, never returned
~~9. **Add Comprehensive Test Suite** (3 hours) - Unit tests for all major components~~ done — unit + provider + BDD
~~10. **Add Validation Rules Engine** (2 hours) - Make validation rules configurable~~ NOT-DO — stable rule IDs; per-rule severity is the TODO_LIST knob

### Medium Impact, Low Effort (Quality of Life)

~~11. **Add Logging** (30 min) - Structured logging with different levels~~ NOT-DO — CLI output is the interface
~~12. **Improve CLI Help** (20 min) - Better help text and examples~~ moot — CLI rebuilt with cobra help
~~13. **Add Auto-detection Improvements** (45 min) - Better project detection logic~~ done — git-remote identity parsing
~~14. **Add Output Formats** (30 min) - Support JSON output for CI/CD~~ done — JSON + SARIF in the rebuilt CLI

### Medium Impact, Medium Effort (Future Enhancements)

~~15. **Add Plugin System** (3 hours) - Allow custom validation rules~~ NOT-DO — stable rule IDs instead
~~16. **Add Template Repository** (2 hours) - Support multiple template sources~~ NOT-DO — one canonical embedded template
~~17. **Add Integration Tests** (2 hours) - End-to-end testing~~ done — provider detect→repair→verify loop
~~18. **Add Performance Metrics** (1.5 hours) - Measure validation performance~~ Won't implement — no perf need

### Low Impact, High Effort (Long-term)

~~19. **Add Web Interface** (5 hours) - Simple web UI for validation~~ NOT-DO
~~20. **Add Database Persistence** (4 hours) - Store validation history~~ NOT-DO — stateless tool
~~21. **Add API Server** (6 hours) - REST API for validation services~~ NOT-DO

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

~~1. Commit current changes~~ moot — historical; the tree itself was later deleted
~~2. Execute high impact, low effort items~~ moot — resolved variously (see verdicts above)
~~3. Set up proper CI/CD with testing~~ done at 72085c2
~~4. Gradually implement medium effort items~~ superseded — the rebuild re-scoped everything
~~5. Plan for long-term architectural improvements~~ superseded — ROADMAP owns long-term now
