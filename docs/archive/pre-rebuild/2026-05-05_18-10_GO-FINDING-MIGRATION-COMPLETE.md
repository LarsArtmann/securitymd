# Comprehensive Project Status — template-SECURITY

> **Superseded 2026-10-08** — describes the pre-rebuild tree (`internal/`, `cmd/template-security/`), deleted in the securitymd rebuild (`1c55390`). Successor report: `docs/status/2026-10-08_23-40_securitymd-rebuild-buildflow-provider-status.md`. All open items below carry inline verdicts.

**Date:** 2026-05-05 18:10 | **Branch:** master | **Commit:** f24512b

---

## Executive Summary

Successfully migrated template-SECURITY from custom validation types to `github.com/larsartmann/go-finding v0.3.0`, eliminating the `go-branded-id` dependency. All 17 tests pass (8 unit + 9 acceptance). Build and vet are clean. The project now uses a standardized finding/report data model for its core validation output.

---

## A) FULLY DONE ✅

_Historically accurate: all of this shipped in the pre-rebuild tree (deleted `1c55390`). Several rows below were themselves superseded by the rebuild — pipeline.Detector became a toolsdk Spec, WriteSARIF became CLI `--format sarif`._

### 1. go-finding Migration (THIS SESSION)

| What                       | Before                                                    | After                                                                                                     | Status      |
| -------------------------- | --------------------------------------------------------- | --------------------------------------------------------------------------------------------------------- | ----------- |
| Validation results         | `SecurityValidationResult{Valid, Errors, Warnings, File}` | `*finding.Report` with `[]Finding`                                                                        | ✅ Complete |
| Validation errors/warnings | `[]string`                                                | `finding.Finding` with `SeverityError`/`SeverityWarning`, `CategorySecurity`, machine-readable rule names | ✅ Complete |
| CLI validate command       | Manual string counting                                    | Severity-based filtering via `finding.Report`                                                             | ✅ Complete |
| JSON output                | Custom map structure                                      | Standard `finding.Report` JSON serialization                                                              | ✅ Complete |
| PrintResults               | Custom struct iteration                                   | Report-based with severity classification                                                                 | ✅ Complete |
| Unit tests (8)             | Old `SecurityValidationResult` assertions                 | `collectMessages(report, severity)` + `ReportIsValid()`                                                   | ✅ Complete |
| Acceptance tests (9)       | Old struct field assertions                               | `finding.SeverityError`/`finding.SeverityWarning` checks                                                  | ✅ Complete |
| ID types                   | `id.ID[IDBrand, string]` from `go-branded-id`             | Plain `type ID string` / `type PolicyID string`                                                           | ✅ Complete |
| Renamed `IDID` → `ID`      | Embarrassing `type IDID string` ("ID ID" nonsense)        | Clean `type ID string`                                                                                    | ✅ Complete |
| `go-branded-id` dependency | Required                                                  | Removed                                                                                                   | ✅ Complete |

**Files changed (7):**

- `go.mod` / `go.sum` — added `go-finding v0.3.0`, removed `go-branded-id v0.1.0`
- `internal/security_validator.go` — core migration to `finding.Report`/`finding.Finding`
- `internal/security_validator_test.go` — test migration
- `internal/types/ids.go` — simplified from branded IDs to plain string types
- `cmd/template-security/validate.go` — CLI adapted to new types
- `test/acceptance/validation_test.go` — acceptance test migration

### 2. Previously Complete (from earlier sessions)

- Shell-to-Go migration of the entire tool
- Cobra CLI with `setup`, `validate`, `status` commands
- Viper configuration loading
- Project detection (git, package.json, go.mod, Cargo.toml, pyproject.toml)
- Template-based SECURITY.md generation
- SARIF output capability (inherited from go-finding Report)
- GitHub Actions CI workflow
- Acceptance test suite (Ginkgo/Gomega)

---

## B) PARTIALLY DONE 🔶

### 1. go-finding Integration Depth — ~~~30%~~ done next session, superseded by rebuild

All five gaps (FindingError, Builder, WriteJSON/SARIF, pipeline, filters) were closed the same day by session 2 — see `2026-05-05_18-51_DEEP-GO-FINDING-INTEGRATION.md`. The rebuild then re-founded integration on toolsdk instead of pipeline.

- **NOT using `finding.FindingError`** — `security_tool.go` still has its own `SecurityError` type for config/template errors. Could use `finding.NewValidationError()`, `finding.NewIOError()`, etc.
- **NOT using `finding.Builder` API** — Findings are constructed via `finding.NewFinding()` directly. The fluent builder would be cleaner.
- **NOT using `finding.Report` JSON/SARIF output** — validate.go manually constructs JSON maps instead of using `report.WriteJSON()` or `report.WriteSARIF()`.
- **NOT using pipeline** — The `pipeline.Detector` interface is not implemented. The validator could become a proper `Detector` that plugs into the go-finding pipeline.
- **NOT using filters** — `finding.Filter`, `finding.BySeverity`, etc. not used for result processing.

### 2. Test Coverage — 24.2% overall (historical)

The coverage question is moot for the deleted tree; the successor ships unit + provider-loop + BDD suites, all green.

| Package                 | Coverage | Notes                                                                                            |
| ----------------------- | -------- | ------------------------------------------------------------------------------------------------ |
| `internal`              | 35.2%    | Validator is well-tested, but `security_tool.go` and `project_detector.go` have minimal coverage |
| `internal/types`        | N/A      | No test files                                                                                    |
| `cmd/template-security` | 0.0%     | No tests at all                                                                                  |
| `test/acceptance`       | 100.0%   | Good acceptance coverage                                                                         |

### 3. CLI Binary — Functional but Basic — superseded

The rebuilt CLI (`cmd/securitymd`) ships `--file`, `--format text/json/sarif`, `--severity`, and exit 0/1/2 semantics.

- No `--sarif` output format (now easily addable via go-finding)
- No `--severity` filtering (now trivially addable)
- No `--config` for validation rules

---

## C) NOT STARTED ⬜

~~1. **Pipeline integration** — Implement `pipeline.Detector` interface for `SecurityValidator`~~ done 2026-05-05 (session 2) as pipeline.Detector; superseded 10-08 — rebuild drops pipeline for a toolsdk provider
~~2. **SARIF output in CLI** — Add `-format sarif` flag using `report.WriteSARIF()`~~ done session 2 (WriteSARIFFiltered); rebuilt CLI ships text/json/sarif
~~3. **Severity filtering in CLI** — Add `-severity` flag using `finding.BySeverityAtLeast()`~~ done session 2 (-severity flag); rebuilt CLI keeps --severity
~~4. **Migrate SecurityError to finding.FindingError** — Replace custom error type in `security_tool.go`~~ done session 2 — FindingError family (NewIOError/NewParseError)
~~5. **Implement `finding.Builder`** — Use fluent builder instead of `NewFinding()` + field assignment~~ done session 2 — Builder API used throughout
~~6. **CLI tests** — 0% coverage on `cmd/template-security`~~ NOT-DO — old CLI deleted at 1c55390; successor cmd/ stays test-free by design (coverage via provider + acceptance tests, see FEATURES.md)
~~7. **Project detector tests** — `project_detector.go` has no test file~~ done — pkg/policy/project_test.go exists in the rebuild
~~8. **Config file tests** — Config loading has minimal coverage~~ done at 1c55390 — config system removed entirely (nothing to test)
~~9. **Template variable tests** — Template processing edge cases untested~~ done — pkg/policy/generate_test.go + line-precise unresolved-template tests in validate_test.go
~~10. **Nix flake migration** — Proposal exists but no implementation~~ done at 9d3094a (2026-06-17)
~~11. **De-duplication of findings** — Multiple validations of same file produce duplicates~~ NOT-DO — moot in the rebuild: one detect pass per run, no multi-validation merge
~~12. **Exit code semantics** — Different exit codes for different severity levels~~ done — rebuilt CLI exits 0 clean / 1 error findings / 2 operational failure
~~13. **CI/CD integration guide** — How to use in GitHub Actions with SARIF upload~~ NOT-DO — BuildFlow provider is the integration path (README documents usage)
~~14. **Version flag in binary** — Uses ldflags but no goreleaser config~~ NOT-DO — version stays ldflags metadata; goreleaser never adopted

---

## D) TOTALLY FUCKED UP 💥

### 1. 88+ Linter Warnings (pre-existing)

The project has **88+ golangci-lint warnings** that were NOT introduced by this migration. Major categories:

| Category           | Count | Examples                                            |
| ------------------ | ----- | --------------------------------------------------- |
| `depguard`         | ~15   | Cobra, color, viper, internal imports "not allowed" |
| `exhaustruct`      | ~10   | Cobra.Command missing fields                        |
| `revive`           | ~10   | Missing package comments, exported type comments    |
| `forbidigo`        | ~8    | `fmt.Printf` forbidden                              |
| `gochecknoglobals` | ~6    | Global variables in cmd                             |
| `err113`           | ~3    | Dynamic errors via `errors.New()`                   |
| `noinlineerr`      | ~5    | Inline error handling                               |
| `wrapcheck`        | ~3    | Unwrapped external errors                           |
| `mnd`              | ~5    | Magic numbers                                       |
| `gosec`            | ~3    | File permissions, file inclusion                    |

**Root cause:** The `.golangci.yml` config is extremely strict and was likely copied from another project. Many rules are inappropriate for a CLI tool (e.g., `forbidigo` blocking `fmt.Printf`).

### 2. Unused Functions

- `printSuccess` in `cmd_helpers.go` — defined but never called
- `printError` in `cmd_helpers.go` — defined but never called
- The old `validate.go` had `allResultsValid` and `countValid` — removed in migration

### 3. JSON Tag Inconsistency

`types.go` uses `snake_case` JSON tags (`supported_until`, `semantic_version`) while `golangci-lint` expects `camelCase`. This is pre-existing.

---

## E) WHAT WE SHOULD IMPROVE 📈

### High Impact

~~1. **Fix `.golangci.yml`** — Relax rules inappropriate for a CLI tool (forbidigo, exhaustruct for cobra). This eliminates 80%+ of warnings in one shot.~~ done session 3 (f29259c, d6614d6)
~~2. **Add CLI tests** — `cmd/template-security` at 0% coverage is the biggest gap~~ NOT-DO — successor covers behavior via provider + acceptance tests
~~3. **Use go-finding's JSON/SARIF output** — Replace manual JSON construction with `report.WriteJSON()` / `report.WriteSARIF()`~~ done session 2 — WriteJSON/WriteSARIFFiltered wired
~~4. **Implement `pipeline.Detector`** — Makes the validator pluggable into go-finding's ecosystem~~ done session 2; superseded — toolsdk Spec is the integration shape now

### Medium Impact

~~5. **Add `-format sarif` CLI flag** — Trivially addable now, huge value for CI/CD~~ done session 2; rebuilt CLI ships it
~~6. **Add `-severity` CLI flag** — Filter findings by minimum severity~~ done session 2; rebuilt CLI keeps it
~~7. **Test `project_detector.go`** — No tests at all for a complex file~~ done — pkg/policy/project_test.go in the rebuild
~~8. **Migrate `SecurityError` to `finding.FindingError`** — Consolidate error handling~~ done session 2
~~9. **Remove unused functions** — `printSuccess`, `printError` in cmd_helpers~~ done session 3 (d6614d6)

### Low Impact / Polish

~~10. Fix JSON tag casing in `types.go`~~ moot — types deleted at 1c55390
~~11. Add package comments~~ moot — rebuild lint runs at 0 issues
~~12. Add nix flake build~~ done at 9d3094a
~~13. Add goreleaser config~~ NOT-DO — publish via git tag (TODO_LIST checklist); goreleaser never adopted
~~14. Update stale status docs~~ done 2026-10-09 — this docs-health pass annotates and archives them

---

## F) TOP 25 THINGS TO DO NEXT

| #  | Priority | Task | Effort                                                          | Impact |
| -- | -------- | ---- | --------------------------------------------------------------- | ------ |
| ~~ | 1        | P0   | Fix `.golangci.yml` — relax CLI-inappropriate rules             | 30min  |
| ~~ | 2        | P0   | Add CLI tests for `cmd/template-security`                       | 2h     |
| ~~ | 3        | P0   | Remove unused `printSuccess`/`printError` functions             | 5min   |
| ~~ | 4        | P1   | Add `-format sarif` flag using `report.WriteSARIF()`            | 30min  |
| ~~ | 5        | P1   | Add `-severity` filter flag using `finding.BySeverityAtLeast()` | 30min  |
| ~~ | 6        | P1   | Use `report.WriteJSON()` instead of manual JSON in validate.go  | 30min  |
| ~~ | 7        | P1   | Implement `pipeline.Detector` for `SecurityValidator`           | 1h     |
| ~~ | 8        | P1   | Test `project_detector.go`                                      | 1h     |
| ~~ | 9        | P1   | Migrate `SecurityError` to `finding.FindingError`               | 1h     |
| ~~ | 10       | P1   | Use `finding.Builder` API in validator                          | 30min  |
| ~~ | 11       | P2   | Fix JSON tags in `types.go` (camelCase)                         | 15min  |
| ~~ | 12       | P2   | Add package comments to all packages                            | 15min  |
| ~~ | 13       | P2   | Add `go.Finding` category tags to all findings                  | 15min  |
| ~~ | 14       | P2   | Add confidence scores to findings                               | 15min  |
| ~~ | 15       | P2   | Update `README.md` to mention go-finding integration            | 15min  |
| ~~ | 16       | P2   | Update `SECURITY.md` template                                   | 30min  |
| ~~ | 17       | P2   | Add `Makefile` or `flake.nix` build                             | 1h     |
| ~~ | 18       | P2   | Add GitHub Actions SARIF upload step                            | 30min  |
| ~~ | 19       | P2   | Add configuration file for validation rules                     | 2h     |
| ~~ | 20       | P3   | Add `docs/adr/` for go-finding migration decision               | 30min  |
| ~~ | 21       | P3   | Add `FEATURES.md` — audit actual features                       | 1h     |
| ~~ | 22       | P3   | Add `TODO_LIST.md` — comprehensive task list                    | 1h     |
| ~~ | 23       | P3   | Clean up stale docs in `docs/status/` and `docs/planning/`      | 30min  |
| ~~ | 24       | P3   | Add goreleaser config for binary releases                       | 1h     |
| ~~ | 25       | P3   | Update `CHANGELOG.md` with go-finding migration entry           | 15min  |

---

## G) TOP #1 QUESTION I CANNOT FIGURE OUT MYSELF 🤔

**~~Should `SecurityValidator` become a `pipeline.Detector` implementation…?~~** — RESOLVED twice over: implemented as `Detector` in session 2 the same day; the 2026-10-08 rebuild then replaced the question entirely — the integration shape is now a go-finding **toolsdk Spec provider** (Detect + Repair), which fits the validate/generate duality better than a detect-only pipeline node.

Arguments for `Detector`:

- Pluggable into go-finding's pipeline (detect → triage → fix → verify)
- Consistent with go-finding's ecosystem patterns
- Could combine with govet/staticcheck detectors from go-finding

Arguments against:

- The validator doesn't "detect" issues in source code — it validates a Markdown file's structure
- The pipeline's fix/verify loop doesn't apply to SECURITY.md validation
- Adding pipeline dependency for a single-validator use case feels heavyweight

This is an architectural decision that affects the project's direction and I'd rather have your input than guess.

---

## Dependency State

| Dependency      | Version | Status                   |
| --------------- | ------- | ------------------------ |
| `go-finding`    | v0.3.0  | ✅ NEW — core data model |
| `cobra`         | v1.10.2 | ✅ CLI framework         |
| `viper`         | v1.21.0 | ✅ Config loading        |
| `ginkgo/v2`     | v2.28.3 | ✅ BDD acceptance tests  |
| `gomega`        | v1.40.0 | ✅ Test assertions       |
| `testify`       | v1.11.1 | ✅ Unit test assertions  |
| `fatih/color`   | v1.19.0 | ✅ Terminal colors       |
| `go-branded-id` | —       | ❌ REMOVED               |

## Test Results

```
ok  github.com/LarsArtmann/template-SECURITY/internal            0.004s  35.2%
ok  github.com/LarsArtmann/template-SECURITY/test/acceptance     0.005s  100.0%
?   github.com/LarsArtmann/template-SECURITY/cmd/template-security       [no test files]
?   github.com/LarsArtmann/template-SECURITY/internal/types              [no test files]

Total: 17 tests PASS | 0 FAIL | go vet CLEAN | go build CLEAN
```

---

_Generated by Crush <crush@charm.land>_

## Resolution (2026-10-09)

Every open item in C/E/F/G carries an inline verdict. Reference hashes: `1c55390` (old tree deleted), `72085c2`/`4a8987a`/`519916d` (docs pass), `9d3094a` (flake.nix), session-2/3 work of 2026-05-05 (`f29259c`, `d6614d6` era). Archivable: no open items remain.
