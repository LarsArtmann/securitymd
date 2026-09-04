# Features Audit — template-SECURITY

_Audit Date: 2026-05-06 | Commit: 03dcf51_

---

## Core Features

| #  | Feature                    | Status              | Description                                                              | Evidence                                                                                |
| -- | -------------------------- | ------------------- | ------------------------------------------------------------------------ | --------------------------------------------------------------------------------------- |
| 1  | **SECURITY.md Validation** | ✅ FULLY_FUNCTIONAL | Validates SECURITY.md files for required sections and content quality    | `internal/security_validator.go` — 6 section checks, 4 quality checks, 22 tests passing |
| 2  | **SECURITY.md Generation** | ✅ FULLY_FUNCTIONAL | Generates SECURITY.md from Go templates with variable substitution       | `internal/security_tool.go` — template rendering, config loading, file output           |
| 3  | **Policy Status Display**  | ✅ FULLY_FUNCTIONAL | Shows current security file status, available templates, next steps      | `cmd/template-security/status.go`                                                       |
| 4  | **CLI Interface**          | ✅ FULLY_FUNCTIONAL | Cobra-based CLI with validate, setup, status subcommands                 | `cmd/template-security/main.go` — version info, help text                               |
| 5  | **Multi-format Output**    | ✅ FULLY_FUNCTIONAL | Text, JSON, SARIF output formats for validation reports                  | `cmd/template-security/validate.go` — `outputReport()`                                  |
| 6  | **Severity Filtering**     | ✅ FULLY_FUNCTIONAL | Filter findings by minimum severity level                                | `cmd/template-security/validate.go` — `parseSeverity()`                                 |
| 7  | **Config File Loading**    | ✅ FULLY_FUNCTIONAL | Auto-discovers and loads `.template-security.yaml` config                | `internal/security_tool.go` — `FindConfigFile()`, `LoadConfig()`                        |
| 8  | **Project Detection**      | ✅ FULLY_FUNCTIONAL | Auto-detects project name, org, domain from git/package.json/go.mod/etc. | `internal/project_detector.go` — 6 detection strategies                                 |
| 9  | **go-finding Integration** | ✅ FULLY_FUNCTIONAL | Implements `pipeline.Detector` interface, uses `finding.Report/Builder`  | `internal/security_validator.go` — `Detect()`, compile-time check                       |
| 10 | **BDD Acceptance Tests**   | ✅ FULLY_FUNCTIONAL | Ginkgo-based acceptance tests for validation and generation              | `test/acceptance/` — 9 specs, all passing                                               |

## Partial / Incomplete Features

| #  | Feature                      | Status                 | Description                                                                | Gap                                                                      |
| -- | ---------------------------- | ---------------------- | -------------------------------------------------------------------------- | ------------------------------------------------------------------------ |
| 11 | **Pipeline Orchestration**   | ⚠️ PARTIALLY_FUNCTIONAL | `pipeline.Detector` implemented but never wired via `pipeline.New()/Run()` | Interface satisfied but no caller uses pipeline mode                     |
| 12 | **Enterprise Policy Type**   | ⚠️ PARTIALLY_FUNCTIONAL | `PolicyTypeEnterprise` constant exists but uses same template as github    | No enterprise-specific template or logic                                 |
| 13 | **CI/CD Integration**        | ⚠️ PARTIALLY_FUNCTIONAL | GitHub Actions workflow exists but broken                                  | Go 1.21 in CI vs 1.26.2 in go.mod; references `make` which doesn't exist |
| 14 | **Configuration Validation** | ⚠️ PARTIALLY_FUNCTIONAL | Config loads but doesn't validate required fields                          | Missing org/email validation on loaded config                            |
| 15 | **Template Variable System** | ⚠️ PARTIALLY_FUNCTIONAL | Variables stored in config but only `SUPPORT_YEARS` is processed           | Other variables available but not used in template                       |

## Non-functional / Dead Code

| #  | Feature                   | Status       | Description                                                                                                 | Evidence                                                      |
| -- | ------------------------- | ------------ | ----------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------- |
| 16 | **Domain Types**          | ❌ DEAD_CODE | `internal/types/` defines rich domain model (`SecurityPolicy`, `ValidationResult`, `Template`) — zero usage | Only `PolicyType` and `Version` are actually used             |
| 17 | **Error Sentinel Values** | ❌ DEAD_CODE | `ErrConfigNotFound`, `ErrInvalidConfig`, etc. in `security_tool.go` — defined but never returned directly   | Functions create new errors via `finding.New*Error()` instead |
| 18 | **Shell Scripts**         | ❌ LEGACY    | 5 shell scripts (~1,323 lines) that duplicate Go functionality                                              | `scripts/compliance-check.sh`, `generate-metrics.sh`, etc.    |

## Missing / Planned

| #  | Feature                   | Status     | Description                                           |
| -- | ------------------------- | ---------- | ----------------------------------------------------- |
| 19 | **CLI Tests**             | ❌ MISSING | Zero tests for `cmd/template-security/` (0% coverage) |
| 20 | **ProjectDetector Tests** | ❌ MISSING | Complex 354-line file with zero dedicated tests       |
| 21 | **Version Subcommand**    | ❌ MISSING | `--version` flag exists but no dedicated subcommand   |
| 22 | **Release Automation**    | ❌ MISSING | No goreleaser, no release tags, no distribution       |
| 23 | **flake.nix Build**       | ❌ MISSING | Project uses justfile (deprecated per AGENTS.md)      |
| 24 | **CHANGELOG**             | ❌ STALE   | Only v0.1.0 entry; months of work undocumented        |
| 25 | **Pre-commit Hook**       | ❌ MISSING | No git hook for SECURITY.md validation                |
