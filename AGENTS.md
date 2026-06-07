# AGENTS.md — template-SECURITY

CLI tool for validating and generating SECURITY.md files. Written in Go, using Cobra for CLI, go-finding for structured findings/reports, and Ginkgo for BDD acceptance tests.

## Commands

```bash
just build              # Build binary to bin/template-security (also: ./scripts/build.sh)
just test-unit          # go test ./internal/... -v
just test-bdd           # go test ./test/acceptance/... -v
just test-all           # Both unit + BDD
just validate           # Build + run validate subcommand
just status             # Build + run status subcommand

go build -o bin/template-security ./cmd/template-security              # Direct build
go test ./internal/... -v                                              # Unit tests only
go test ./test/acceptance/... -v                                       # BDD tests only
```

Build script injects version info via ldflags: `main.version`, `main.commit`, `main.date`.

## Architecture

```
cmd/template-security/    # CLI entrypoint (Cobra commands: setup, validate, status)
internal/                  # Core business logic (package `internal`)
  security_validator.go    # Validates SECURITY.md files, implements pipeline.Detector
  security_tool.go         # Generates SECURITY.md from templates, loads config
  project_detector.go      # Auto-detects project name/org/domain from git/package.json/go.mod/etc.
  types/                   # Domain types (PolicyType, Version, Config structs, IDs)
templates/SECURITY.md      # Go text/template used for policy generation
test/acceptance/            # Ginkgo BDD acceptance tests
```

### Control Flow

- **Validate**: CLI → `SecurityValidator.ValidateSECURITYMd()` → checks required sections + content quality → returns `finding.Report` with `[]finding.Finding`
- **Setup/Generate**: CLI → `SecurityTool.GeneratePolicy()` → reads `templates/SECURITY.md` → Go `text/template` rendering → writes SECURITY.md
- **Config discovery**: `SecurityTool.FindConfigFile()` walks up from CWD looking for `.template-security.yaml`

### Key Dependencies

- `github.com/larsartmann/go-finding` — finding/report data model, SARIF/JSON output, pipeline.Detector interface
- `github.com/spf13/cobra` — CLI framework
- `github.com/spf13/viper` — config file loading (`.template-security.yaml`)
- `github.com/onsi/ginkgo/v2` + `github.com/onsi/gomega` — BDD acceptance tests
- `github.com/stretchr/testify` — unit test assertions
- `github.com/fatih/color` — colored CLI output

## Testing Patterns

- **Unit tests**: `internal/*_test.go` — package `internal` (white-box), uses `testify`. Table-driven tests with `t.Parallel()`. Temp files via `t.TempDir()`.
- **BDD acceptance tests**: `test/acceptance/` — Ginkgo/Gomega with `ginkgo.Describe`/`ginkgo.It`/`ginkgo.By`. Helper `withTempFile()` in `test_helpers.go`.
- **No tests for**: `cmd/template-security/` (0% coverage), `internal/project_detector.go` (354 lines, no dedicated tests).

## Known Issues

- **Compile error in `cmd/template-security/validate.go:195`**: `report.WriteSARIFFiltered` call is missing `context.Context` argument. The go-finding API signature changed to require `(context.Context, io.Writer, finding.Severity)` but the call still uses `(*os.File, finding.Severity)`.
- **CI workflow is broken**: `.github/workflows/security-validation.yml` uses Go 1.21 but go.mod requires 1.26.3. It also references `make` which doesn't exist.
- **Dead code in `internal/types/`**: Rich domain types (`SecurityPolicy`, `ValidationResult`, `Template`, etc.) are defined but mostly unused. Only `PolicyType` and `Version` are actually consumed.
- **Error sentinels in `security_tool.go`** (`ErrConfigNotFound`, etc.) are defined but never returned — functions create new errors via `finding.New*Error()` instead.
- **Shell scripts in `scripts/`** (~1,323 lines) are legacy, duplicating Go functionality.

## Linting

`.golangci.yml` is extremely strict (90+ linters). Key settings:

- `exhaustruct` excluded for `cmd/` and key internal files
- `forbidigo` allows `fmt.Print*` in cmd and specific internal files
- `funlen`: 80 lines / 70 statements
- `gocognit`: max complexity 40
- `cyclop`: max complexity 15
- `varnamelen`: min 3 chars
- `ireturn`: allows error, empty, anon, stdlib, generic

`.go-arch-lint.yml` defines DDD-style layers (domain, application, infrastructure, config, main) with dependency rules, but the current codebase doesn't actually follow this structure — it's a template/aspiration config.

## Conventions

- Main branch: `master` (per `git-town.toml`)
- Go module: `github.com/LarsArtmann/template-SECURITY`
- Config file: `.template-security.yaml` (auto-discovered walking up from CWD)
- Template variables use Go `text/template` syntax: `{{.Organization}}`, `{{.ContactEmail}}`, etc.
- CLI uses `color.Cyan`/`color.Red`/`color.Green`/`color.Yellow` for output via `fatih/color`
- `finding.NewBuilder()` is the standard way to construct findings with rule, tool name, message, severity, position
- All errors wrapped via `finding.NewValidationError`, `finding.NewIOError`, `finding.NewParseError` (not `fmt.Errorf`)
