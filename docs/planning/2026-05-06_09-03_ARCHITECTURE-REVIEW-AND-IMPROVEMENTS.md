# Architecture Review & Improvement Opportunities

_Date: 2026-05-06 | Commit: 03dcf51_

---

## Current Architecture Assessment

### Module Dependency Graph

```
cmd/template-security
├── internal/security_validator  (pipeline.Detector)
├── internal/security_tool       (policy generation)
├── internal/project_detector    (auto-detection)
├── internal/types               (shared types)
├── github.com/larsartmann/go-finding  (finding model)
├── github.com/spf13/cobra       (CLI framework)
├── github.com/spf13/viper       (config loading)
└── github.com/fatih/color       (terminal colors)
```

### Scalability: 6/10

**Positive:**
- CLI + internal split allows independent evolution
- `pipeline.Detector` interface enables adding new validators as plugins
- go-finding ecosystem provides standardized output format

**Limiting:**
- Package globals in `validate.go`/`setup.go` prevent concurrent command execution
- Hardcoded template path (`"templates/SECURITY.md"`) breaks when binary is installed globally
- No plugin/extension system for custom validation rules
- Single-threaded validation — no file-level parallelism

### Modularity: 7/10

**Positive:**
- Clean `cmd/` → `internal/` boundary
- `SecurityValidator` properly implements `pipeline.Detector`
- Types isolated in `internal/types/`
- Template rendering separated from business logic

**Limiting:**
- `SecurityTool` is a god object (config loading + template rendering + policy generation + file discovery)
- `ProjectDetector` mixes I/O (git, file reads) with parsing logic — hard to test
- No interfaces between `cmd/` and `internal/` — direct struct coupling

### Service Orientation: 3/10

- Single binary, no service decomposition
- No HTTP/gRPC API — CLI only
- Not composable as a library (internal package blocks external imports)
- No streaming/event-based processing

### Composability: 5/10

**Positive:**
- `pipeline.Detector` enables composition in go-finding pipelines
- Config via viper supports env vars + YAML + CLI flags
- Multi-format output (text/JSON/SARIF)

**Limiting:**
- `internal/` package prevents reuse as a library
- No public API surface
- Template rendering tightly coupled to file system

---

## Deepening Opportunities

### Candidate 1: Extract Validator from I/O

**Files:** `internal/security_validator.go` (311 lines)
**Problem:** `ValidateSECURITYMd()` does file I/O + validation + report building in one function. `Detect()` wraps it. Testing requires temp files.
**Solution:** Split into:
- `ValidateContent(content string) *finding.Report` — pure function, zero I/O
- `ValidateFile(path string) *finding.Report` — reads file, calls `ValidateContent`
**Benefits:** Content validation testable without file I/O. `Detect()` becomes a thin wrapper. Deletion test confirms: removing the I/O layer concentrates complexity in the pure validator.

### Candidate 2: Extract Template Renderer from SecurityTool

**Files:** `internal/security_tool.go` (250 lines)
**Problem:** `SecurityTool` handles config loading, file discovery, template reading, template rendering, and file writing. It's 5 responsibilities in one struct.
**Solution:** Extract `TemplateRenderer` with interface:
- `Render(template string, data TemplateData) (string, error)`
- `RenderFile(templatePath string, data TemplateData) (string, error)`
**Benefits:** Testable without temp files. Template rendering becomes mockable. `SecurityTool` becomes a thin orchestrator.

### Candidate 3: Extract ConfigProvider Interface

**Files:** `internal/security_tool.go`, `cmd/template-security/setup.go`
**Problem:** Setup command directly calls `SecurityTool.LoadConfig()` + `FindConfigFile()`. Config resolution logic is scattered across cmd and internal.
**Solution:** Create `ConfigProvider` interface:
```go
type ConfigProvider interface {
    Resolve() (*Config, error)
}
```
**Benefits:** Setup command becomes testable. Config resolution logic centralized. Alternative config sources (env-only, flags-only) become trivial.

### Candidate 4: ProjectDetector I/O Separation

**Files:** `internal/project_detector.go` (354 lines)
**Problem:** Every detection method directly calls `os.ReadFile()` or `exec.Command()`. No way to test without real files/git repos.
**Solution:** Inject `FileSystem` interface:
```go
type FileSystem interface {
    ReadFile(name string) ([]byte, error)
    Stat(name string) (os.FileInfo, error)
    Getwd() (string, error)
    RunGit(args ...string) (string, error)
}
```
**Benefits:** Full testability. Can mock git responses. Can test edge cases (malformed package.json, etc.).

### Candidate 5: Dead Type System Cleanup

**Files:** `internal/types/types.go` (122 lines), `internal/types/ids.go` (8 lines)
**Problem:** 8 types defined, only 3 used. Creates confusion about what the "real" domain model is.
**Solution:** Delete unused types (`Contact`, `SecurityPolicy`, `Project`, `ValidationResult`, `Template`, `ValidationMessage`, `Error`, `Warning`, `ID`, `PolicyID`). Keep only `PolicyType`, `Version`, `VersionStatus`.
**Benefits:** Eliminates split brain. Developers see only what's real.

---

## Recommended Architecture (Target)

```
cmd/template-security/
├── main.go              (entry point)
├── validate.go          (CLI wiring, delegates to Validator)
├── setup.go             (CLI wiring, delegates to Generator)
└── status.go            (CLI wiring, delegates to StatusChecker)

internal/
├── validator/
│   ├── validator.go     (SecurityValidator — pure content validation)
│   └── validator_test.go
├── generator/
│   ├── generator.go     (TemplateRenderer + Generator)
│   └── generator_test.go
├── detector/
│   ├── detector.go      (ProjectDetector — with I/O interface)
│   └── detector_test.go
├── config/
│   ├── config.go        (Config + ConfigProvider)
│   └── config_test.go
└── types/
    └── types.go         (ONLY PolicyType, Version, VersionStatus)
```

### Key Changes from Current → Target

1. **Split `internal/` into sub-packages** — each with single responsibility
2. **Delete dead types** — keep only what's used
3. **Extract interfaces for I/O** — FileSystem, ConfigProvider
4. **Pure functions for business logic** — `ValidateContent()` takes string, not file path
5. **No package globals** — command state in structs

---

## Concrete Recommendations

### Make the repo more Service-Oriented
1. Move core logic to `pkg/` (public) instead of `internal/` (private)
2. Expose `ValidateContent(content string) *finding.Report` as a public function
3. Add HTTP API wrapper for remote validation (optional)

### Make the repo more Composable
1. Define `Validator` interface in `pkg/` — allows custom validators
2. Support `io.Reader`/`io.Writer` instead of file paths
3. Add middleware/hooks for pre/post validation
4. Support pipeline composition via go-finding's `pipeline.New()`
