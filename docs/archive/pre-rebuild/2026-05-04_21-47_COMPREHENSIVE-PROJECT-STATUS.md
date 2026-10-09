# template-SECURITY — Comprehensive Status Report

> **Superseded 2026-10-08** — this report describes the pre-rebuild `template-SECURITY` tree (`internal/`, `cmd/template-security/`, `scripts/`, viper config), deleted in the securitymd rebuild (`1c55390`). Successor report: `docs/status/2026-10-08_23-40_securitymd-rebuild-buildflow-provider-status.md`. Every open item below carries an inline verdict; nothing remains open.

**Date:** 2026-05-04 21:47
**Branch:** master
**Last Commit:** `05c5a96` feat(docs): add public/private decision analysis
**Go Version:** 1.26.2
**Test Status:** ALL PASSING (9 BDD acceptance + unit tests)
**Build Status:** CLEAN (`go build`, `go vet`, `go mod tidy` all pass)

---

## Project Overview

A Go CLI tool (`template-security`) that validates and generates `SECURITY.md` files for GitHub repositories. Checks for required sections, content quality, and generates compliant templates.

| Metric             | Value                                                                        |
| ------------------ | ---------------------------------------------------------------------------- |
| Go source lines    | ~2,245                                                                       |
| Shell script lines | ~1,323 (legacy)                                                              |
| Go packages        | 4 (`cmd/template-security`, `internal`, `internal/types`, `test/acceptance`) |
| Dependencies       | 8 direct (cobra, viper, ginkgo, gomega, testify, fatih/color, go-branded-id) |
| BDD test specs     | 9 (all passing)                                                              |
| Unit test files    | 2 (`security_tool_test.go`, `security_validator_test.go`)                    |

---

## a) FULLY DONE

_All items in this section shipped in the pre-rebuild tree and were historically complete when written; that tree was deleted at `1c55390`. Successors live in `pkg/policy`, `pkg/provider`, and `cmd/securitymd` of the rebuild. No open items here._

### 1. Core CLI Application

- **Status:** COMPLETE
- Three commands: `validate`, `setup`, `status`
- `validate` — validates SECURITY.md for required sections, content quality, template variables
- `setup` — generates SECURITY.md from templates with config/flag support
- `status` — shows current security policy status
- JSON output support for CI/CD (`--format json`)
- Colorized terminal output via fatih/color

### 2. Dependency Migration: go-composable-business-types → go-branded-id

- **Status:** COMPLETE (just done this session)
- `internal/types/ids.go` import updated to `github.com/larsartmann/go-branded-id`
- `go.mod` updated with new dependency and `replace` directive pointing to `../go-branded-id`
- Branded ID types (`IDID`, `PolicyID`) using `id.ID[Brand, string]` pattern
- All builds and tests pass cleanly

### 3. BDD Acceptance Test Suite

- **Status:** COMPLETE
- Ginkgo v2 + Gomega framework integrated
- 9 specs covering:
  - Policy validation (complete policies, missing sections, template variables, content quality)
  - Policy generation (file creation, template variable substitution)
  - Configuration loading (find config, load config, error handling)
  - Non-existent file handling
- Test helpers (`withTempFile`, `expectNoError`) extracted

### 4. Domain Types

- **Status:** COMPLETE
- Well-defined types in `internal/types/types.go`: `PolicyType`, `VersionStatus`, `ContactType`, `Version`, `Contact`, `SecurityPolicy`, `Project`, `ValidationResult`, `Template`
- Branded IDs in `internal/types/ids.go`: `IDID`, `PolicyID`
- Types use branded ID types for type-safe identifiers

### 5. Security Validation Engine

- **Status:** COMPLETE
- Validates required sections (Security Policy header, Reporting a Vulnerability, Supported Versions, Security Practices, contact email)
- Content quality checks (minimum length, no template variables, substantive content, version information)
- Structured error/warning output
- User-friendly result printing with colorized output

### 6. Template System

- **Status:** COMPLETE
- Go `text/template` based (not string replacement)
- Template data with defaults (organization, email, version, dates)
- Custom variable support via config
- Template file: `templates/SECURITY.md`

### 7. Configuration Management

- **Status:** COMPLETE (using viper)
- YAML config file support (`.template-security.yaml`)
- Config file discovery (walks up directory tree)
- Environment variable support (`TEMPLATE_SECURITY_` prefix)
- CLI flag overrides

### 8. Project Detection

- **Status:** COMPLETE
- Multi-source detection: git remote, package.json, go.mod, Cargo.toml, pyproject.toml, directory name
- Organization and domain extraction from git URLs (HTTPS + SSH)
- Fallback chain from most to least reliable source

### 9. Git Town Configuration

- **Status:** COMPLETE
- `git-town.toml` configured with `main = "master"` and GitHub API connector

### 10. Linting Configuration

- **Status:** COMPLETE
- `.golangci.yml` with 80+ linters enabled
- `.go-arch-lint.yml` with Clean Architecture / DDD enforcement
- Extremely comprehensive linter setup

### 11. CI/CD Workflow

- **Status:** EXISTS but outdated
- `.github/workflows/security-validation.yml` present
- Validates SECURITY.md on push/PR
- Generates validation report and PR comments

---

## b) PARTIALLY DONE

### 1. Branded ID Integration in Domain Types — NOT-DO

- **Status:** 20% — Types defined in `ids.go` but only `SecurityPolicy.ID` and `ValidationResult.ID`/`ValidationResult.PolicyID` actually use them
- Removed entirely one day later in the go-finding migration (2026-05-05); never reintroduced — the rebuild uses plain types
- `Template.ID` still uses `IDID` type alias which maps to `id.ID[IDBrand, string]` — but `IDBrand` is a generic "ID" brand, not template-specific
- No constructor functions (`NewPolicyID()`, `NewTemplateID()`, etc.)
- No NanoId integration (planned in `go-composable-business-types-usage.md` but not implemented)
- The branded IDs exist but aren't meaningfully differentiated — `IDID` and `PolicyID` are both `id.ID[X, string]` with different brands, which is correct, but there are no factory functions or meaningful usage beyond type definitions

### 2. Error Handling — done 2026-05-05 (session 2)

- **Status:** 60% — Structured `SecurityError` type exists with code, message, field, cause
- Completed the next day: custom `SecurityError` replaced by `finding.FindingError` (see `2026-05-05_18-51_DEEP-GO-FINDING-INTEGRATION.md`)
- Named sentinel errors defined (`ErrConfigNotFound`, `ErrInvalidConfig`, etc.)
- BUT: Not all code paths use `SecurityError` consistently
- `security_validator.go` uses plain `fmt.Errorf` instead of structured errors
- No error wrapping chain in validator

### 3. Documentation — done at `72085c2`/`4a8987a` (successor docs pass)

- **Status:** 50%
- README.md: Good, covers features, commands, validation standards
- CHANGELOG.md: Skeleton only — no actual entries beyond "Initial release"
- SECURITY.md: Contains placeholder values (security@github.com, MyCompany)
- IMPROVEMENT_PLAN.md: Outdated — many items already done, some no longer relevant
- BDD_TESTS_REVIEW.md: Accurate assessment of the state when written, but now outdated (BDD tests exist)
- PARTS.md: Good analysis of extraction potential
- PUBLIC_OR_PRIVATE.md: Good analysis, needs updating now that dependency migration is done
- `docs/planning/go-composable-business-types-usage.md`: Now STALE — references old library name

### 4. CLI Help & Examples — NOT-DO

- **Status:** 70%
- CLI rebuilt as `cmd/securitymd` (validate/setup/status); shell completion is cobra-default, no version subcommand by design
- Root command has good example text
- `validate` has `--file` and `--format` flags
- `setup` has `--type`, `--organization`, `--email`, `--output`, `--quick` flags
- Missing: shell completion, man pages, version command (version vars defined but not a subcommand)

---

## c) NOT STARTED

### 1. Replace viper with koanf — done at `1c55390`

Viper was dropped wholesale in the rebuild (≈15 indirect deps removed); no config system exists anymore, so no koanf was ever needed.

### 2. Extract `projectmeta` Library — NOT-DO

Superseded: project detection lives in `pkg/policy/project.go` of the rebuild; extraction never happened and has no consumer.

### 3. Plugin System for Validation Rules — NOT-DO

Superseded by design: the rebuild exposes stable kebab rule IDs keyed for suppressions/configuration instead of a plugin system.

### 4. Template Repository (Multiple Sources) — NOT-DO

Superseded by design: the rebuild ships exactly one canonical template, embedded via `go:embed` (`pkg/policy/template.md`).

### 5. Integration / E2E Tests — done in the rebuild

The successor ships provider contract tests running a full detect → repair → verify loop in a real temp git repo (`pkg/provider/provider_test.go`) plus 7 Ginkgo BDD specs.

### 6. Logging — Won't implement

CLI tool; colored terminal output via `fatih/color` is the interface. No structured logging in the successor either.

### 7. Performance Metrics / Benchmarking — Won't implement

No benchmarks in the successor; the tool is subprocess-light and no perf need ever materialized.

### 8. Pre-commit Hook — NOT-DO

BuildFlow's findings gate is the enforcement mechanism now; no per-repo hooks.

### 9. Update SECURITY.md with Real Contact Info — routed

Resolved by policy: the rebuild's contact default is the GitHub advisory link (email optional; no fabricated addresses). This repo's own SECURITY.md regeneration sits in the TODO_LIST publish checklist (blocked on the rename).

### 10. Migrate Shell Scripts to Go or Remove — done at `1c55390`

All ~1,300 legacy shell lines were deleted in the rebuild; the two scripts with no Go equivalent (compliance-check, generate-metrics) were dropped as out of scope for a SECURITY.md linter.

### 11. Fix CI/CD Workflow Go Version — done at `72085c2`

CI rewritten: `go-version-file: go.mod`, no `make`, new layout paths.

### 12. Fix `go-arch-lint.yml` to Match Actual Structure — done at `1c55390`

The template config was deleted with the old tree.

### 13. nix flake migration — done at `9d3094a`

`flake.nix` created and justfile removed on 2026-06-17; modernized in the rebuild docs pass (`72085c2`).

---

## d) TOTALLY FUCKED UP

### 1. Stale Documentation Referencing Dead Dependency — done 2026-10-09

`go-composable-business-types-usage.md` archived by the docs-health pass.

### 2. `.template-security.yaml` Has Placeholder Values — done at `1c55390`

The config file was deleted entirely in the rebuild.

### 3. `SECURITY.md` Has Wrong Contact — routed

Still `security@github.com` as of 2026-10-09; the file passes the new validator but predates the rename. Regeneration is tracked in the TODO_LIST publish checklist.

### 4. `go-arch-lint.yml` Is Completely Wrong — done at `1c55390`

Deleted.

### 5. CI Workflow Will Fail — done at `72085c2`

Workflow rewritten (go-version-file, no make).

### 6. CHANGELOG.md Is Empty — done at `72085c2`

Unreleased entry documents the rebuild with all breaking changes.

---

## e) WHAT WE SHOULD IMPROVE

### Critical (Blocks Public Release)

~~1. **Fix SECURITY.md contact information** — Currently says `security@github.com`. Must have real contact.~~ NOT-DO — superseded: contact default is GitHub-advisory-only (no fabricated emails); own-policy regen sits in TODO_LIST publish checklist
~~2. **Fix CI workflow Go version** — 1.21 vs 1.26.2 mismatch guarantees CI failure.~~ done at 72085c2 — CI rewritten to go-version-file: go.mod
~~3. **Fix or remove `.go-arch-lint.yml`** — Currently references non-existent directories. Either customize it or remove it to avoid confusion.~~ done at 1c55390 — .go-arch-lint.yml deleted
~~4. **Update stale `go-composable-business-types-usage.md`** — References dead dependency. Update or remove.~~ done 2026-10-09 — archived by docs-health pass
~~5. **Replace placeholder config values** — `.template-security.yaml` has fake org/email/URLs.~~ done at 1c55390 — .template-security.yaml deleted

### High (Quality & Correctness)

~~6. **Add CHANGELOG entries** — Document all the work that's been done.~~ done at 72085c2 — Unreleased entry covers the rebuild
~~7. **Replace viper with koanf** — Per project standards, viper is not preferred.~~ done at 1c55390 — viper dropped wholesale; no koanf needed
~~8. **Migrate justfile to flake.nix** — justfile is deprecated per AGENTS.md.~~ done at 9d3094a — flake.nix created + justfile removed (2026-06-17)
~~9. **Remove or migrate shell scripts** — 1,323 lines of bash duplicating Go functionality.~~ done at 1c55390 — scripts/ deleted (~1,300 lines)
~~10. **Add CLI integration tests** — Build binary, run commands, verify output.~~ done — successor ships provider detect→repair→verify loop tests + BDD
~~11. **Complete branded ID integration** — Add factory functions, use meaningful brands, integrate NanoId.~~ NOT-DO — branded IDs removed 2026-05-05 (go-finding migration)

### Medium (Polish & Best Practices)

~~12. **Add structured logging** — Replace `fmt.Printf` with slog or similar.~~ Won't implement — CLI tool; colored terminal output is the interface
~~13. **Add version subcommand** — Version vars are defined but inaccessible as a command.~~ NOT-DO — rebuilt CLI ships validate/setup/status only
~~14. **Customize `.golangci.yml`** — 80+ linters is excessive; tune to project needs.~~ done pre-rebuild (f29259c, d6614d6); config survives at 0 issues
~~15. **Add shell completion** — Cobra supports this natively.~~ NOT-DO — cobra ships default completion; nothing to add
~~16. **Update IMPROVEMENT_PLAN.md** — Many items are done or no longer relevant.~~ done at 72085c2 — archived to docs/archive/pre-rebuild/

### Low (Nice to Have)

~~17. **Add pre-commit hook** — Validate SECURITY.md before commits.~~ NOT-DO — BuildFlow findings gate is the mechanism now
~~18. **Add benchmarks** — For validation and template processing.~~ Won't implement — no perf need
~~19. **Add plugin system** — For custom validation rules.~~ NOT-DO — stable kebab rule IDs replace plugin ambitions
~~20. **Extract `projectmeta` library** — As described in PARTS.md.~~ NOT-DO — detection lives in pkg/policy/project.go of the rebuild

---

## f) Top 25 Things We Should Get Done Next

| #  | Priority | Task                                                                                                            | Effort   | Impact                |
| -- | -------- | --------------------------------------------------------------------------------------------------------------- | -------- | --------------------- |
~~| 1  | P0       | Fix SECURITY.md with real contact info                                                                          | 5 min    | Blocks release        |~~ routed — TODO_LIST publish checklist (advisory-only contact default)
~~| 2  | P0       | Fix CI workflow Go version (1.21 → 1.26.2)                                                                      | 10 min   | CI is broken          |~~ done at 72085c2 — CI follows go.mod now
~~| 3  | P0       | Update `go-composable-business-types-usage.md` to reference `go-branded-id`                                     | 30 min   | Misleading docs       |~~ done 2026-10-09 — archived (docs-health pass)
~~| 4  | P0       | Replace placeholder values in `.template-security.yaml`                                                         | 5 min    | Bad defaults          |~~ done at 1c55390 — config deleted
~~| 5  | P1       | Fix or remove `.go-arch-lint.yml`                                                                               | 30 min   | False confidence      |~~ done at 1c55390 — deleted
~~| 6  | P1       | Write CHANGELOG entries for all work since v0.1.0                                                               | 1 hour   | Release readiness     |~~ done at 72085c2
~~| 7  | P1       | Complete branded ID integration (factory functions, meaningful brands)                                          | 2 hours  | Type safety           |~~ NOT-DO — removed 2026-05-05
~~| 8  | P1       | Replace viper with koanf in `security_tool.go`                                                                  | 2 hours  | Project standards     |~~ done at 1c55390 (dropped)
~~| 9  | P1       | Update `PUBLIC_OR_PRIVATE.md` — remove "replace directive" issue (still exists but pointing to correct lib now) | 30 min   | Accurate status       |~~ done at 72085c2 — archived, verdict quoted in manifest
~~| 10 | P1       | Add CLI integration tests (build binary, run validate/setup/status)                                             | 3 hours  | Confidence            |~~ done — provider loop + BDD tests
~~| 11 | P1       | Migrate justfile to flake.nix                                                                                   | 3 hours  | Project standards     |~~ done at 9d3094a
~~| 12 | P1       | Update IMPROVEMENT_PLAN.md to reflect current reality                                                           | 1 hour   | Accurate planning     |~~ done at 72085c2 — archived
~~| 13 | P2       | Remove or refactor `scripts/security-setup.sh` (duplicates Go `setup` command)                                  | 1 hour   | Reduce confusion      |~~ done at 1c55390 (deleted)
~~| 14 | P2       | Remove or refactor `scripts/validate-policies.sh` (duplicates Go `validate` command)                            | 30 min   | Reduce confusion      |~~ done at 1c55390 (deleted)
~~| 15 | P2       | Implement Go equivalents for `scripts/compliance-check.sh` (GDPR/SOC2/ISO27001)                                 | 4 hours  | Feature parity        |~~ NOT-DO — out of scope post-rebuild
~~| 16 | P2       | Implement Go equivalents for `scripts/generate-metrics.sh` (metrics/Prometheus)                                 | 4 hours  | Feature parity        |~~ NOT-DO — out of scope post-rebuild
~~| 17 | P2       | Remove `scripts/build.sh` (replaced by `go build` or flake.nix)                                                 | 5 min    | Cleanup               |~~ done at 1c55390
~~| 18 | P2       | Add structured logging (replace `fmt.Printf` with `slog`)                                                       | 2 hours  | Observability         |~~ Won't implement
~~| 19 | P2       | Add `version` subcommand to CLI                                                                                 | 30 min   | User experience       |~~ NOT-DO
~~| 20 | P2       | Tune `.golangci.yml` — reduce from 80+ linters to project-appropriate set                                       | 1 hour   | Build speed           |~~ done pre-rebuild (f29259c)
~~| 21 | P2       | Update `BDD_TESTS_REVIEW.md` — BDD tests now exist (score should be updated)                                    | 30 min   | Accurate docs         |~~ done at 72085c2 — archived
~~| 22 | P3       | Add pre-commit hook for SECURITY.md validation                                                                  | 1 hour   | Developer experience  |~~ NOT-DO (BuildFlow gate)
~~| 23 | P3       | Extract `projectmeta` library (as described in PARTS.md)                                                        | 2-3 days | Reusability           |~~ NOT-DO
~~| 24 | P3       | Add benchmarks for validation and template processing                                                           | 1 hour   | Performance awareness |~~ Won't implement
~~| 25 | P3       | Decide: make this repo public or private (per PUBLIC_OR_PRIVATE.md analysis)                                    | Decision | Direction             |~~ routed — TODO_LIST publish decision (PUBLIC_OR_PRIVATE: conditionally make public)

---

## g) Top #1 Question I Cannot Figure Out Myself

**~~What is the actual security contact email and organization name for this project?~~** — RESOLVED 2026-10-08: the org is `LarsArtmann`; the contact default is the GitHub advisory link, with email strictly optional (CLI `--email` / provider `contact-email`). The rebuild deliberately never fabricates `security@domain` addresses, which dissolves this question.

Everything hinges on this:

- `SECURITY.md` currently says `security@github.com` and `MyCompany`
- `.template-security.yaml` says `security@mycompany.com` and `MyCompany`
- The safe harbor section says "MyCompany commits to..."
- The template references `hackerone.com/mycompany`

Is this `LarsArtmann` / `larsartmann`? Something else? Without knowing the real org name and contact, all generated policies and the project's own SECURITY.md will remain placeholder garbage.

---

## Architecture Summary

```
template-SECURITY/
├── cmd/template-security/        # CLI entrypoint (cobra)
│   ├── main.go                   # Root command + subcommands
│   ├── validate.go               # validate command + JSON output
│   ├── setup.go                  # setup command with config loading
│   ├── status.go                 # status command
│   └── cmd_helpers.go            # shared CLI helpers
├── internal/
│   ├── security_tool.go          # Template processing, config loading (viper)
│   ├── security_validator.go     # SECURITY.md validation engine
│   ├── project_detector.go       # Multi-source project metadata detection
│   ├── security_tool_test.go     # Unit tests (testify)
│   ├── security_validator_test.go # Unit tests (testify)
│   └── types/
│       ├── types.go              # Domain types (Policy, Version, Contact, etc.)
│       └── ids.go                # Branded ID types (go-branded-id)
├── test/acceptance/              # BDD tests (ginkgo/gomega)
│   ├── security_policy_suite_test.go
│   ├── validation_test.go
│   ├── policy_generation_test.go
│   └── test_helpers.go
├── templates/
│   └── SECURITY.md               # Go text/template for policy generation
├── scripts/                      # Legacy bash scripts (1,323 lines)
├── docs/
│   ├── status/                   # Status reports (this file)
│   └── planning/                 # Planning documents
└── .github/workflows/            # CI (needs Go version fix)
```

### Dependency Graph

```
cmd/template-security
├── internal (security_tool, security_validator, project_detector)
│   ├── internal/types (domain types)
│   │   └── go-branded-id (branded IDs)
│   ├── spf13/viper (config loading) → should be koanf
│   └── spf13/cobra (CLI framework)
├── fatih/color (terminal colors)
└── encoding/json (JSON output)
```

---

## Test Coverage Summary

| Package                 | Tests                | Status   |
| ----------------------- | -------------------- | -------- |
| `internal`              | Unit tests (testify) | PASS     |
| `internal/types`        | No test files        | —        |
| `test/acceptance`       | 9 BDD specs (ginkgo) | ALL PASS |
| `cmd/template-security` | No test files        | —        |

**Gap:** No tests for `cmd/template-security` (CLI layer) or `internal/types` (domain types).

---

## Recent Git History (Last 20 Commits)

```
05c5a96 feat(docs): add public/private decision analysis for template-SECURITY project
5002c9c feat(config): add Git Town workflow configuration
3f113b8 chore: normalize YAML indentation and bump Go toolchain version
26bc523 feat(config): add project metadata.yaml
f230ee9  chore: code style improvements, dependency updates, jscpd report exclusions
5d23774 refactor: improve code clarity and test structure across core validation and test helpers
8258770 refactor: extract error handling helpers in test acceptance suite
ffd94da refactor: extract expectNoError helper, improve code formatting, and remove dead code
957d1c5 refactor: extract command helpers, improve validation logic, and enhance test infrastructure
131a36e refactor: extract lastPartFromEnd helper and update dependencies
a111bc6 chore: remove enterprise-grade Go linting justfile
f5c4318 chore(deps): initialize go module and coverage
8f8c5f1 chore(build): add architecture linting task
40c1c1d Remove .auto-deduplicate.lock file
eb6679c fix(deduplicate): crush refactoring for duplicate dup_1775271425558282000_11
5a6ec62 fix(deduplicate): crush refactoring for duplicate dup_1775271425558285000_12
ee92e72 fix(deduplicate): crush refactoring for duplicate dup_1775271425558275000_10
5a72911 fix(deduplicate): crush refactoring for duplicate dup_1775271425558286000_13
1ae2a85 chore: pre-deduplication commit before crush refactoring
40c518d feat: add ginkgo/gomega BDD testing framework and code cleanup
```

---

_Report generated at 2026-05-04_21-47 by Crush_

## Resolution (2026-10-09)

Every open item in sections b)–g) carries an inline verdict. Reference hashes: `1c55390` (old tree deleted: `internal/`, `cmd/template-security/`, `scripts/`, `templates/`, `.go-arch-lint.yml`, config), `72085c2` + `4a8987a` (docs/CI/archive pass), `519916d` (README), `9d3094a` (flake.nix created 2026-06-17). Items whose intent survived live in the rebuild (`pkg/policy`, `pkg/provider`, `cmd/securitymd`); the publish/public decision is routed to TODO_LIST. Archivable: no open items remain.
