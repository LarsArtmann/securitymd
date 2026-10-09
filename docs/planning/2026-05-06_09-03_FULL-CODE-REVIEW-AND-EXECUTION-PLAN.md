# Full Code Review & Execution Plan

> **Superseded 2026-10-08** — every file reviewed below (`cmd/template-security/`, `internal/`, `templates/`, `scripts/`) was deleted in the securitymd rebuild (`1c55390`). The findings are an accurate historical record of the dead tree; the successor lives in `pkg/policy`, `pkg/provider`, `cmd/securitymd`. All open items below carry inline verdicts.

_Date: 2026-05-06 | Commit: 03dcf51 | Reviewer: Senior Staff Architect_

---

## Executive Summary

**Build:** ✅ Clean | **Lint:** ✅ 0 issues | **Tests:** ✅ 24/24 passing (after fixing 2 bugs)

The project is a focused CLI tool with clean separation between CLI (`cmd/`) and logic (`internal/`). The go-finding integration is solid. The main issues are: dead domain types, no CLI tests, a broken CI workflow, and several stale documentation files.

---

## File-by-File Review Findings

### `cmd/template-security/main.go` (34 lines)

- ✅ Clean entry point, proper error handling
- ⚠️ `version`/`commit`/`date` globals — acceptable for ldflags pattern

### `cmd/template-security/validate.go` (237 lines)

- ⚠️ **`outputFormat`/`minimumSeverity` as package globals** — prevents concurrent test execution
- ⚠️ **`_ = stat`** at line 95 — dead code, file existence checked but unused
- ⚠️ **`printSummary` ignores first param** — takes `[]*finding.Report` but uses `_`
- ✅ Good error wrapping pattern
- ✅ Proper severity filtering with `parseSeverity()`

### `cmd/template-security/setup.go` (145 lines)

- ⚠️ **6 package-level global variables** (`orgName`, `contactEmail`, etc.) — prevents concurrent test execution, makes testing harder
- ⚠️ **`ValidatePolicyType` is exported but in `main` package** — unreachable from tests
- ⚠️ **Two `ProjectDetector` instances created** (lines 95, 109) — should be one
- ✅ Good config file fallback chain

### `cmd/template-security/status.go` (75 lines)

- ✅ Clean, focused
- ⚠️ Hardcoded `"templates"` path — should use config or constant
- ⚠️ Magic string `.md` length check — fragile

### `cmd/template-security/cmd_helpers.go` (30 lines)

- ✅ Clean, minimal

### `internal/security_validator.go` (311 lines)

- ✅ **`pipeline.Detector` implementation** is correct with context cancellation
- ✅ **`buildFinding()`** with Builder API — consistent, well-structured
- ⚠️ **`filePath` field** set in constructor but never used — `WithFile()` creates new instances, `ValidateSECURITYMd()` takes param instead
- ⚠️ **`lineChecker()` is package-level** — should be a method or live in a helper package
- ⚠️ **`level` field as string** in `checkRequiredSections` — should be `finding.Severity`
- ⚠️ **`hasActualContent()` excludes lines with "example"** — this is a smell; legitimate content might contain "example"
- ⚠️ **`PrintResults()` uses `fmt.Printf`** — should use `color` for consistency with other commands
- ✅ Good compile-time interface check

### `internal/security_tool.go` (250 lines)

- ⚠️ **Error sentinels defined but never returned** — `ErrConfigNotFound`, `ErrInvalidConfig`, `ErrTemplateNotFound`, etc. are dead code
- ⚠️ **`readTemplate()` ignores `PolicyType` parameter** — enterprise vs github has no effect
- ⚠️ **`prepareTemplateData()` tries `time.ParseDuration("1y")`** — this doesn't work; `time.Duration` doesn't support "y" unit
- ⚠️ **Hardcoded `"templates/SECURITY.md"`** path — breaks when run from different directory
- ⚠️ **`TemplateData.AdditionalFields` is `map[string]any`** but only strings ever put in
- ✅ Good config defaults chain

### `internal/project_detector.go` (354 lines)

- ⚠️ **LARGEST FILE — over 350 line limit** from architect checklist
- ⚠️ **Zero dedicated tests** — complex URL parsing, domain detection, multi-source fallback
- ⚠️ **`detectFromPackageJSON()` uses hand-rolled string parsing** — comment says "use JSON parsing" but doesn't
- ⚠️ **`detectOrgFromPackageJSON()` has potential index out of bounds** at `contentStr[start:start+500]` if file is shorter
- ⚠️ **`DetectDomain()` produces nonsensical domains** like `template-security.com` from project name
- ⚠️ **Each detect method calls `detectFromGitRemote()` independently** — no caching, runs `git config` multiple times
- ✅ Good use of `strings.CutPrefix` and `strings.SplitSeq`

### `internal/types/types.go` (122 lines)

- ❌ **CRITICAL: Dead type system** — Only `PolicyType`, `Version`, and `VersionStatus` are used. The following are NEVER used anywhere:
  - `Contact` / `ContactType`
  - `SecurityPolicy`
  - `Project`
  - `ValidationResult`
  - `ValidationMessage` / `Error` / `Warning`
  - `Template` / `TemplateVariable`
- ✅ Used types (`PolicyType`, `Version`) are well-defined

### `internal/types/ids.go` (8 lines)

- ❌ **`ID` and `PolicyID` types never used** — dead code

### `test/acceptance/test_helpers.go` (28 lines)

- ✅ **Fixed: temp file race condition** — was removing file before writing
- ✅ Now uses `defer os.Remove()` and proper filename pattern

### `test/acceptance/validation_test.go` (188 lines)

- ✅ Good BDD structure with table-driven test cases
- ✅ Proper severity-based assertions

### `test/acceptance/policy_generation_test.go` (161 lines)

- ✅ Good test isolation with temp directories
- ⚠️ Test for "custom variables" doesn't actually verify the variable appears in output

---

## Architectural Assessment

### Strengths

1. **Clean cmd/internal split** — CLI and logic properly separated
2. **go-finding integration** — proper `pipeline.Detector` implementation
3. **BDD acceptance tests** — ginkgo-based, user-focused
4. **Zero lint issues** — golangci-lint passes clean

### Weaknesses

~~1. **Type split brain** — `internal/types/` has 8 unused types vs `go-finding` used types~~ done — resolved by deletion at 1c55390; finding.Finding is the only model
~~2. **No CLI tests** — 0% coverage on `cmd/`~~ NOT-DO — successor covers behavior via provider + acceptance tests
~~3. **Package globals** — prevent concurrent execution and testing~~ moot — rebuilt CLI has no flag globals
~~4. **Dead error sentinels** — 6 `Err*` vars never returned~~ done at 1c55390 — file deleted wholesale
~~5. **Broken CI** — wrong Go version, missing build tool~~ done at 72085c2 — workflow rewritten
~~6. **`time.ParseDuration("1y")` bug** — silently fails, always uses default~~ moot — code deleted

---

## Pareto Plan (80/20 Breakdown)

### 1% → 51% Impact (Do First)

| # | Task                                            | Impact                    | Effort | File(s)              |
| - | ----------------------------------------------- | ------------------------- | ------ | -------------------- |
~~| 1 | Delete dead types in `internal/types/`          | Removes split brain       | 15min  | `types.go`, `ids.go` |~~ done at 1c55390
~~| 2 | Fix CI workflow (Go version, remove make)       | Enables CI                | 15min  | `.github/workflows/` |~~ done at 72085c2
~~| 3 | Delete dead error sentinels in security_tool.go | Removes confusion         | 10min  | `security_tool.go`   |~~ done at 1c55390 — security_tool.go deleted wholesale
~~| 4 | Delete `.go-arch-lint.yml`                      | Removes misleading config | 5min   | `.go-arch-lint.yml`  |~~ done at 1c55390
~~| 5 | Delete empty dirs (`report/`, `reports/`)       | Clean tree                | 2min   | dirs                 |~~ done — gone

### 4% → 64% Impact (Do Second)

| #  | Task                                                                 | Impact             | Effort | File(s)                  |
| -- | -------------------------------------------------------------------- | ------------------ | ------ | ------------------------ |
~~| 6  | Add CLI tests for validate command                                   | 0→50% cmd coverage | 60min  | `cmd/template-security/` |~~ NOT-DO — successor cmd/ test-free by design (provider + acceptance coverage)
~~| 7  | Remove `t.Parallel()` from `detectFromPackageJSON` test file         | Already done       | 0min   | —                        |~~ historical — noted as already done when written
~~| 8  | Fix `readTemplate()` to use config `TemplateDir`                     | Correctness        | 20min  | `security_tool.go`       |~~ moot — template now embedded (go:embed), no path lookup exists
~~| 9  | Cache `detectFromGitRemote()` result                                 | Performance        | 15min  | `project_detector.go`    |~~ moot — one git call per detect in the rebuild
~~| 10 | Convert `level` string → `finding.Severity` in checkRequiredSections | Type safety        | 10min  | `security_validator.go`  |~~ moot — function deleted

### 20% → 80% Impact (Do Third)

| #  | Task                                                                        | Impact                 | Effort | File(s)                    |
| -- | --------------------------------------------------------------------------- | ---------------------- | ------ | -------------------------- |
~~| 11 | Refactor package globals → struct-based command state                       | Testability            | 45min  | `validate.go`, `setup.go`  |~~ moot — rebuilt CLI has no flag globals
~~| 12 | Add ProjectDetector tests                                                   | Coverage               | 60min  | `project_detector_test.go` |~~ done — pkg/policy/project_test.go
~~| 13 | Fix `time.ParseDuration("1y")` — use `time.ParseDuration` or just `AddDate` | Correctness            | 10min  | `security_tool.go`         |~~ moot — code deleted
~~| 14 | Fix `withTempFile` pattern in test helpers                                  | Already done           | 0min   | —                          |~~ historical — already done when written
~~| 15 | Update `CHANGELOG.md`                                                       | Documentation          | 30min  | `CHANGELOG.md`             |~~ done at 72085c2
~~| 16 | Clean up stale docs (archive 2025-12-11 status reports)                     | Reduces confusion      | 15min  | `docs/status/`             |~~ done 2026-10-09 — this docs-health pass annotates and archives them
~~| 17 | Delete stale `docs/planning/go-composable-business-types-usage.md`          | Removes dead reference | 2min   | docs                       |~~ done 2026-10-09 — archived by this pass

### Remaining (80%+ Impact Items)

| #  | Task                                                        | Impact                | Effort | File(s)                 |
| -- | ----------------------------------------------------------- | --------------------- | ------ | ----------------------- |
~~| 18 | Replace hand-rolled JSON parsing in `detectFromPackageJSON` | Robustness            | 30min  | `project_detector.go`   |~~ moot — manifest parsing removed; identity comes from git remote
~~| 19 | Add SARIF output test                                       | Coverage              | 30min  | `validate_test.go`      |~~ routed — TODO_LIST golden-file item
~~| 20 | Add JSON output test                                        | Coverage              | 30min  | `validate_test.go`      |~~ routed — same TODO_LIST item
~~| 21 | Fix `hasActualContent()` "example" exclusion                | Correctness           | 10min  | `security_validator.go` |~~ moot — validator rewritten (hasSubstantiveContent in the rebuild)
~~| 22 | Create project-level `AGENTS.md`                            | Documentation         | 30min  | `AGENTS.md`             |~~ done at 72085c2
~~| 23 | Migrate justfile → flake.nix                                | Per AGENTS.md mandate | 120min | `flake.nix`             |~~ done at 9d3094a
~~| 24 | Add goreleaser                                              | Release automation    | 60min  | `.goreleaser.yml`       |~~ NOT-DO — publish via git tag
~~| 25 | Delete legacy shell scripts                                 | Clean tree            | 10min  | `scripts/`              |~~ done at 1c55390
~~| 26 | Wire `pipeline.New()` + `Run()` for full pipeline           | Feature               | 60min  | `internal/`             |~~ NOT-DO — BuildFlow's DAG orchestrates detect/repair now
~~| 27 | Add enterprise template                                     | Feature completeness  | 45min  | `templates/`            |~~ NOT-DO — one canonical embedded template by design

## Resolution (2026-10-09)

Every Pareto item and every weakness carries an inline verdict. Reference hashes: `1c55390` (old tree deleted), `72085c2`/`4a8987a` (docs pass), `9d3094a` (flake.nix). Surviving intents (SARIF/JSON golden test) are routed to TODO_LIST. Archivable: no open items remain.
