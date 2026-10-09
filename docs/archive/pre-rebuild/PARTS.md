# PARTS.md — Component Analysis & Extraction Potential

> **Superseded** — analyzes the deleted `internal/` tree; no extraction ever happened. Verdicts on the Action Items below. Fittingly, this doc's Phase-2 idea "consider publishing `securitymd`" predicted the 2026-10-08 rebuild's name.

> Analysis of template-SECURITY for reusable library/SDK extraction opportunities
>
> **Last Updated:** 2026-03-03
> **Version:** 1.0

---

## Executive Summary

This document analyzes the `template-SECURITY` project to identify components suitable for extraction as standalone reusable libraries or SDKs. Each component is evaluated against existing alternatives and assessed for unique value proposition.

**Key Findings:**

| Component         | Extraction Value | Recommendation                        |
| ----------------- | ---------------- | ------------------------------------- |
| ProjectDetector   | **High**         | Extract as `projectmeta`              |
| SecurityValidator | **Medium**       | Keep internal, consider open-sourcing |
| SecurityTool      | **Low**          | Keep internal (thin wrapper)          |
| Domain Types      | **Medium**       | Extract with ProjectDetector          |

---

## Component Analysis

### 1. ProjectDetector (`internal/project_detector.go`)

**What it does:**

Detects project metadata (name, organization, domain) from multiple sources with a priority-based fallback chain:

1. Git remote origin URL (HTTPS and SSH formats)
2. `package.json` (Node.js)
3. `go.mod` (Go)
4. `Cargo.toml` (Rust)
5. `pyproject.toml` (Python)
6. Directory name (fallback)

**Current Implementation:**

```go
type ProjectDetector struct{}

func (pd *ProjectDetector) DetectProjectName() string
func (pd *ProjectDetector) DetectOrganization() string
func (pd *ProjectDetector) DetectDomain() string
```

**Alternatives:**

| Library                                                               | Stars | Focus                      | Gap                                   |
| --------------------------------------------------------------------- | ----- | -------------------------- | ------------------------------------- |
| [go-git/go-git](https://github.com/go-git/go-git)                     | 6k+   | Pure Go git implementation | Only git, no multi-language           |
| [boostsecurityio/poutine](https://github.com/boostsecurityio/poutine) | -     | Security analysis          | Has `ParseRepoAndOrg` but specialized |
| [dbinky/Pommel](https://github.com/dbinky/Pommel)                     | -     | Subproject detection       | Different scope (monorepos)           |

**Unique Value Proposition:**

1. **Polyglot Detection** — Single interface for 5+ language ecosystems
2. **Fallback Chain** — Graceful degradation from most to least reliable source
3. **Zero Dependencies** — Pure stdlib, no external git library required
4. **Context-Aware** — Extracts organization/domain from git URLs, not just project name
5. **Structured Output** — Returns typed metadata, not just strings

**Recommended Library Name:** `projectmeta`

**Proposed API:**

```go
package projectmeta

type Metadata struct {
    Name         string
    Organization string
    Domain       string
    Repository   string
    Source       DetectionSource
}

type DetectionSource string

const (
    SourceGit       DetectionSource = "git"
    SourcePackageJSON DetectionSource = "package.json"
    SourceGoMod     DetectionSource = "go.mod"
    SourceCargo     DetectionSource = "cargo.toml"
    SourcePyProject DetectionSource = "pyproject.toml"
    SourceDirectory DetectionSource = "directory"
)

type Detector struct {
    // Configuration options
}

func NewDetector(opts ...Option) *Detector
func (d *Detector) Detect(ctx context.Context) (*Metadata, error)
func (d *Detector) DetectFromPath(ctx context.Context, path string) (*Metadata, error)
```

**Use Cases Beyond template-SECURITY:**

- CI/CD pipelines needing project context
- Code generators requiring project metadata
- Documentation generators
- Security scanners (SBOM, dependency analysis)
- Monorepo tooling

**Extraction Effort:** Medium (2-3 days)
**Value:** High

---

### 2. SecurityValidator (`internal/security_validator.go`)

**What it does:**

Validates `SECURITY.md` files against industry best practices:

- Required sections (Security Policy header, Reporting Vulnerability, Supported Versions, Security Practices)
- Contact information presence
- Response time commitments
- Content quality (minimum length, no template variables, substantive content)

**Current Implementation:**

```go
type SecurityValidator struct{}

type SecurityValidationResult struct {
    Valid    bool
    Errors   []string
    Warnings []string
    File     string
}

func (sv *SecurityValidator) ValidateSECURITYMd(filePath string) (*SecurityValidationResult, error)
func (sv *SecurityValidator) PrintResults(results []*SecurityValidationResult)
```

**Alternatives:**

| Library    | Focus | Gap                                       |
| ---------- | ----- | ----------------------------------------- |
| None found | -     | No dedicated SECURITY.md validator exists |

Related but different:

- Markdown linters (mdl, remark-lint) — Generic, not security-policy aware
- GitHub's security policy guidance — Documentation only, no tooling

**Unique Value Proposition:**

1. **Opinionated Standards** — Based on GitHub's official security policy guidelines
2. **Security-Specific Rules** — Checks for contact info, response times, safe harbor
3. **Quality Scoring** — Not just pass/fail, but actionable warnings
4. **CI/CD Ready** — Designed for automation with exit codes

**Recommended Library Name:** `securitymd` or keep internal

**Proposed API (if extracted):**

```go
package securitymd

type Validator struct {
    rules []Rule
}

type ValidationResult struct {
    Valid    bool
    Score    int           // 0-100
    Errors   []Violation
    Warnings []Violation
}

type Violation struct {
    Code     string
    Message  string
    Line     int
    Severity Severity
}

func NewValidator(opts ...Option) *Validator
func (v *Validator) Validate(content string) (*ValidationResult, error)
func (v *Validator) ValidateFile(path string) (*ValidationResult, error)
```

**Use Cases Beyond template-SECURITY:**

- Pre-commit hooks
- CI/CD quality gates
- GitHub Actions marketplace action
- Security compliance audits

**Extraction Effort:** Low (1-2 days)
**Value:** Medium (niche but useful)

---

### 3. SecurityTool (`internal/security_tool.go`)

**What it does:**

- Loads configuration from YAML files (using viper)
- Processes Go templates for SECURITY.md generation
- Manages template data with defaults

**Current Implementation:**

```go
type SecurityTool struct{}

func (st *SecurityTool) LoadConfig(configPath string) (*Config, error)
func (st *SecurityTool) FindConfigFile() string
func (st *SecurityTool) GeneratePolicy(ctx context.Context, config PolicyConfig) error
```

**Assessment:**

This is a **thin orchestration layer** combining:

1. Config loading (viper — per HOW_TO_GOLANG.md, should use koanf)
2. Template processing (text/template stdlib)
3. File I/O

**Recommendation:** Keep internal. Not suitable for extraction because:

1. Too coupled to template-SECURITY's domain
2. Better alternatives exist for each concern:
   - Config: koanf (per HOW_TO_GOLANG.md)
   - Templates: templ (type-safe), text/template (stdlib)
3. No unique value over composing existing libraries

**Refactoring Needed:**

Per HOW_TO_GOLANG.md, replace:

- `spf13/viper` → `knadh/koanf/v2`

---

### 4. Domain Types (`internal/types/types.go`)

**What it contains:**

Well-defined domain types for security policies:

```go
type PolicyType string      // github, enterprise
type VersionStatus string   // supported, deprecated, end-of-life
type ContactType string     // email, web, api

type Version struct { ... }
type Contact struct { ... }
type SecurityPolicy struct { ... }
type Project struct { ... }
type ValidationResult struct { ... }
type Template struct { ... }
```

**Assessment:**

These types are **security-policy specific**. They could be valuable if extracted alongside:

1. **ProjectDetector** — `Project` type is generic enough
2. **SecurityValidator** — Validation types would be part of that library

**Recommendation:**

- Extract `Project` type with ProjectDetector
- Keep security-specific types internal unless SecurityValidator is extracted

---

## Recommended Extraction Strategy

### Phase 1: Extract `projectmeta` (High Value)

```
projectmeta/
├── detector.go       # Core detection logic
├── sources/
│   ├── git.go        # Git remote detection
│   ├── package_json.go
│   ├── gomod.go
│   ├── cargo.go
│   └── pyproject.go
├── types.go          # Metadata, DetectionSource
├── options.go        # Configuration options
└── detector_test.go
```

**Benefits:**

- Reusable across all LarsArtmann projects
- Clean separation of concerns
- Independent versioning
- Can be used by other Go tools

### Phase 2: Consider `securitymd` (Medium Value)

Only if there's demand from the community or internal use cases.

---

## Comparison to PROJECT_SPLIT_EXECUTIVE_REPORT.md

The existing executive report proposes:

| Their Proposal                        | My Assessment                                     |
| ------------------------------------- | ------------------------------------------------- |
| `go-security-lib` (all internal code) | **Too broad**. Extract ProjectDetector separately |
| `template-security-cli`               | **Correct**. CLI consumes libraries               |
| `security-automation-scripts`         | **Correct**. Scripts are operational              |
| `security-docs`                       | **Questionable**. Docs should live with code      |

**Key Difference:**

The executive report treats all internal code as one library. I recommend extracting **ProjectDetector** as its own focused library because:

1. It has clear boundaries
2. It solves a general problem (project metadata detection)
3. It has no dependencies on security-specific logic
4. It's useful beyond security tooling

---

## Action Items

~~1. [ ] Create `projectmeta` library repository~~ NOT-DO — never happened; detection lives in pkg/policy/project.go (rebuild 1c55390)
~~2. [ ] Extract ProjectDetector code with refined API~~ NOT-DO — same
~~3. [ ] Add comprehensive tests with ginkgo/gomega~~ done differently — pkg/policy/project_test.go (testify) covers identity parsing
~~4. [ ] Update template-SECURITY to use `projectmeta`~~ NOT-DO — moot; tool renamed securitymd, detection in-repo
~~5. [ ] Replace viper with koanf in template-SECURITY~~ done at 1c55390 — viper dropped wholesale (no replacement needed)
~~6. [ ] Consider publishing `securitymd` if community interest exists~~ routed — publish checklist in TODO_LIST; the name securitymd WAS adopted by the 2026-10-08 rebuild

---

## References

- [HOW_TO_GOLANG.md](/Users/larsartmann/projects/library-policy/HOW_TO_GOLANG.md) — Library standards and patterns
- [PROJECT_SPLIT_EXECUTIVE_REPORT.md](PROJECT_SPLIT_EXECUTIVE_REPORT.md) — Previous analysis
- [GitHub Security Policy Guidelines](https://docs.github.com/en/code-security/getting-started/adding-a-security-policy-to-your-repository)
