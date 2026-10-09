# Comprehensive Project Status — template-SECURITY

**Date:** 2026-05-05 18:51 | **Branch:** master | **Last Commit:** 4cfba19

---

## Executive Summary

Two major sessions completed today: (1) migrated to `go-finding v0.3.0` as the data model, and (2) achieved deep integration across every layer — `pipeline.Detector` interface, `Builder` API, `FindingError` errors, native `WriteJSON`/`WriteSARIF` output, and severity filtering. All 17 tests pass. Build and vet are clean. Uncommitted changes from session 2 remain staged.

---

## A) FULLY DONE ✅

### 1. go-finding Migration (Session 1 — commit dd4fb09)

| What                       | Before                                                    | After                                                       |
| -------------------------- | --------------------------------------------------------- | ----------------------------------------------------------- |
| Validation results         | `SecurityValidationResult{Valid, Errors, Warnings, File}` | `*finding.Report` with `[]Finding`                          |
| Validation errors/warnings | `[]string`                                                | `finding.Finding` with severity, category, tags, rule names |
| CLI JSON output            | Manual `json.MarshalIndent` map construction              | `finding.Report` serialization                              |
| ID types                   | `id.ID[IDBrand, string]` via `go-branded-id`              | Plain `type ID string` / `type PolicyID string`             |
| Nonsense `IDID`            | `type IDID string` ("ID ID")                              | Renamed to `type ID string`                                 |
| `go-branded-id`            | Required dependency                                       | Removed                                                     |

### 2. Deep go-finding Integration (Session 2 — uncommitted)

| Integration Point             | Implementation                                                                                         |
| ----------------------------- | ------------------------------------------------------------------------------------------------------ |
| `pipeline.Detector`           | `SecurityValidator` implements `Name()` + `Detect(ctx)` with compile-time check                        |
| `finding.Builder`             | All findings via `NewBuilder().WithCategory().WithTags().WithFixStrategy().WithConfidence().Build()`   |
| `finding.FindingError`        | Replaced custom `SecurityError` — config errors use `NewIOError()`, parse errors use `NewParseError()` |
| `report.WriteJSON()`          | `-format json` uses native go-finding streaming serialization                                          |
| `report.WriteSARIFFiltered()` | `-format sarif` outputs SARIF 2.1.0 for GitHub Code Scanning                                           |
| `-severity` flag              | `finding.BySeverityAtLeast()` predicate for minimum severity filtering                                 |
| `finding.Filter`              | `BySeverity`, `BySeverityAtLeast`, `NotSuppressed` predicates                                          |
| `DetectFile()`                | Convenience: single-file detect → `*finding.Report`                                                    |
| Context awareness             | `Detect(ctx)` respects cancellation                                                                    |

### 3. Previously Complete

- Shell-to-Go migration of entire tool
- Cobra CLI with `setup`, `validate`, `status` commands
- Viper configuration loading
- Project detection (git, package.json, go.mod, Cargo.toml, pyproject.toml)
- Template-based SECURITY.md generation
- GitHub Actions CI workflow
- Acceptance test suite (Ginkgo/Gomega)

---

## B) PARTIALLY DONE 🔶

### 1. go-finding Pipeline Usage — 10%

We implement `Detector` but don't actually **run the pipeline**. The `pipeline.New()` + `Run()` orchestration is not wired. The validator is a standalone detector that could be composed with go-finding's built-in govet/staticcheck detectors, but that composition doesn't exist yet.

### 2. Test Coverage — 23.4% overall

| Package                 | Coverage | Change                                    |
| ----------------------- | -------- | ----------------------------------------- |
| `internal`              | 33.6%    | Down from 35.2% (new untested code added) |
| `internal/types`        | N/A      | No test files                             |
| `cmd/template-security` | 0.0%     | Unchanged — no tests at all               |
| `test/acceptance`       | 100.0%   | Unchanged                                 |

Notable uncovered functions: `DetectFile()` (0.0%), `hasVersionInformation()` (75.0%)

### 3. SARIF Output — 60%

`WriteSARIFFiltered` is wired but:

- File URIs are relative paths, not proper `file://` URIs
- No `--sarif-output` flag for writing to file instead of stdout
- Not tested in any test

---

## C) NOT STARTED ⬜

~~1. **Pipeline composition** — Combine SecurityValidator detector with govet/staticcheck from go-finding~~ NOT-DO — superseded: composition happens in BuildFlow's DAG now (toolsdk provider), not in-process
~~2. **CLI tests** — `cmd/template-security` still at 0% coverage~~ NOT-DO — old CLI deleted at 1c55390; successor cmd/ test-free by design (provider + acceptance coverage)
~~3. **`project_detector.go` tests** — Complex file with zero tests~~ done — pkg/policy/project_test.go exists in the rebuild
~~4. **`security_tool.go` tests** — `NewSecurityError` callers not tested with `FindingError`~~ moot — file deleted at 1c55390; generation covered by pkg/policy/generate_test.go
~~5. **`DetectFile()` tests** — New convenience function, 0% coverage~~ moot — DetectFile deleted with the old tree (detect.go redesigned)
~~6. **SARIF round-trip tests** — No test for `-format sarif` output~~ routed — TODO_LIST keeps the SARIF golden-file test for the rebuilt CLI
~~7. **GitHub Actions SARIF upload** — CI workflow doesn't use SARIF output yet~~ NOT-DO — BuildFlow is the CI integration path
~~8. **`finding.Merge` for multi-file reports** — Could use `finding.Merge()` instead of manual `mergeReports()`~~ moot — rebuild validates one policy per run; no multi-file merge exists
~~9. **`finding.SortBySeverity`** — Not used in text output~~ moot — rebuild output has no multi-severity sort need
~~10. **`finding.GroupByFile`** — Not used in PrintResults~~ moot — single file per validate in the rebuild
~~11. **`finding.ReportFromJSON`** — Not used (could enable re-loading cached results)~~ NOT-DO — no cached-results feature; YAGNI
~~12. **Nix flake migration** — Proposal exists, no implementation~~ done at 9d3094a (2026-06-17)
~~13. **Goreleaser config** — No release automation~~ NOT-DO — publish via git tag (TODO_LIST checklist)
~~14. **`.golangci.yml` fix** — 80%+ of warnings are from inappropriate rules~~ done session 3 of 2026-05-05 (f29259c, d6614d6)
~~15. **Version flag via ldflags** — Exists but untested~~ NOT-DO — cosmetic; version string ships in the rebuilt CLI header
~~16. **Exit codes** — Same exit code regardless of severity level~~ done — rebuilt CLI: 0 clean / 1 error findings / 2 operational failure
~~17. **Configurable validation rules** — Hardcoded required sections~~ routed — TODO_LIST (per-rule severity configuration); config files stay out

---

## D) TOTALLY FUCKED UP 💥

### 1. 88+ Linter Warnings (PRE-EXISTING — NOT FROM OUR WORK)

All warnings are from an overly strict `.golangci.yml`:

| Category           | Count | Root Cause                                         |
| ------------------ | ----- | -------------------------------------------------- |
| `depguard`         | ~15   | Blocks all imports including standard libs         |
| `exhaustruct`      | ~10   | Demands every struct field on Cobra, finding types |
| `forbidigo`        | ~8    | Blocks `fmt.Printf` in a CLI tool that prints      |
| `revive`           | ~10   | Missing package/type comments                      |
| `gochecknoglobals` | ~6    | Cobra flags require globals                        |
| `err113`           | ~3    | Dynamic errors via `errors.New()`                  |
| `noinlineerr`      | ~5    | Inline error handling preference                   |
| `wrapcheck`        | ~3    | Unwrapped external errors                          |
| `mnd`              | ~5    | Magic numbers                                      |
| `gosec`            | ~3    | File permissions                                   |

**None of these were introduced by our work.** The config was copied from another project and needs relaxing.

### 2. Uncommitted Changes

The deep integration session (commit `4cfba19`) did NOT include all files. Six files have uncommitted changes that need committing:

- `internal/security_validator.go` — pipeline.Detector, Builder API
- `internal/security_tool.go` — FindingError migration
- `internal/security_validator_test.go` — updated tests
- `internal/types/types.go` — whitespace from ID rename
- `test/acceptance/validation_test.go` — acceptance test updates
- `cmd/template-security/validate.go` — SARIF, severity flags

### 3. Nil Context in Tests

`security_tool_test.go:77` passes `nil` as context to `GeneratePolicy()` — should use `context.TODO()`.

---

## E) WHAT WE SHOULD IMPROVE 📈

### Critical

~~1. **Commit the uncommitted changes** — Deep integration work is sitting in working tree~~ done — committed by the auto-git daemon the same night (f29259c followed)
~~2. **Fix `.golangci.yml`** — One 30-minute session eliminates 80%+ of warnings~~ done session 3 (f29259c, d6614d6)
~~3. **Add CLI tests** — Biggest coverage gap (0%)~~ NOT-DO — successor covers behavior via provider + acceptance tests

### High Impact

~~4. **Wire `pipeline.New()` + `Run()`** — Actually use the pipeline orchestration~~ NOT-DO — superseded: BuildFlow's DAG orchestrates detect/repair now
~~5. **Compose with go-finding detectors** — Security validation + govet + staticcheck in one pipeline~~ NOT-DO — superseded: same reason, composition lives in BuildFlow
~~6. **Test `DetectFile()`** — New public function, 0% coverage~~ moot — deleted at 1c55390
~~7. **Test `-format sarif`** — No test for SARIF output at all~~ routed — TODO_LIST SARIF golden-file test

### Medium Impact

~~8. **Use `finding.Merge()`** — Replace manual `mergeReports()`~~ moot — no mergeReports left; single report per run
~~9. **Use `finding.SortBySeverity()`** — In PrintResults~~ moot — output redesigned
~~10. **Use `finding.GroupByFile()`** — Group output by file~~ moot — output redesigned
~~11. **Test `project_detector.go`** — Complex file, zero tests~~ done — pkg/policy/project_test.go
~~12. **Fix nil context in tests** — `context.TODO()` instead of `nil`~~ done — tests use context.TODO() (noted in the 2026-07-26 report)

---

## F) TOP 25 THINGS TO DO NEXT

| #  | Priority | Task                                                       | Effort | Impact                        |
| -- | -------- | ---------------------------------------------------------- | ------ | ----------------------------- |
~~| 1  | P0       | **Commit uncommitted deep integration changes**            | 2min   | Uncommitted work at risk      |~~ done — committed by the auto-git daemon the same night
~~| 2  | P0       | Fix `.golangci.yml` — relax CLI-inappropriate rules        | 30min  | Eliminates 80%+ warnings      |~~ done session 3 (f29259c, d6614d6)
~~| 3  | P0       | Remove unused `printSuccess`/`printError` from cmd_helpers | 5min   | Dead code                     |~~ done session 3 (d6614d6)
~~| 4  | P1       | Add CLI tests for `cmd/template-security`                  | 2h     | 0% → 80% coverage             |~~ NOT-DO — provider + acceptance coverage in the successor
~~| 5  | P1       | Wire `pipeline.New()` + `Run()` in validate command        | 1h     | Full pipeline orchestration   |~~ NOT-DO — BuildFlow DAG is the pipeline now
~~| 6  | P1       | Compose with go-finding govet + staticcheck detectors      | 1h     | Multi-tool pipeline           |~~ NOT-DO — same
~~| 7  | P1       | Test `DetectFile()` convenience function                   | 30min  | 0% → 100% coverage            |~~ moot — deleted at 1c55390
~~| 8  | P1       | Add SARIF output test                                      | 30min  | Untested feature              |~~ routed — TODO_LIST golden-file test
~~| 9  | P1       | Use `finding.Merge()` in validate.go                       | 15min  | Use library instead of custom |~~ moot
~~| 10 | P1       | Use `finding.SortBySeverity()` in PrintResults             | 15min  | Better output ordering        |~~ moot
~~| 11 | P2       | Fix nil context in tests → `context.TODO()`                | 5min   | Correctness                   |~~ done
~~| 12 | P2       | Add `--output` flag for file-based SARIF output            | 30min  | CI/CD file output             |~~ NOT-DO — BuildFlow consumes stdout; no file-output flag
~~| 13 | P2       | Test `project_detector.go`                                 | 1h     | Untested complex code         |~~ done — pkg/policy/project_test.go
~~| 14 | P2       | Convert `file://` URIs in SARIF output                     | 30min  | SARIF compliance              |~~ routed — TODO_LIST SARIF golden-file test covers URI shape
~~| 15 | P2       | Update `README.md` with go-finding integration docs        | 30min  | Documentation                 |~~ done at 519916d — README rewritten for securitymd
~~| 16 | P2       | Add `docs/adr/` for go-finding migration decision          | 30min  | Architecture docs             |~~ NOT-DO — decision recorded in CHANGELOG + status reports
~~| 17 | P2       | Add `FEATURES.md` audit                                    | 1h     | Documentation                 |~~ done at 4a8987a
~~| 18 | P2       | Add `TODO_LIST.md`                                         | 1h     | Project management            |~~ done at 4a8987a
~~| 19 | P3       | Add GitHub Actions SARIF upload step                       | 30min  | CI/CD                         |~~ NOT-DO — BuildFlow path
~~| 20 | P3       | Add exit codes by severity                                 | 30min  | CLI usability                 |~~ done — 0/1/2 semantics in the rebuilt CLI
~~| 21 | P3       | Configurable validation rules via YAML                     | 2h     | Extensibility                 |~~ routed — TODO_LIST per-rule severity item
~~| 22 | P3       | Add goreleaser config                                      | 1h     | Distribution                  |~~ NOT-DO
~~| 23 | P3       | Add nix flake build                                        | 1h     | Build automation              |~~ done at 9d3094a
~~| 24 | P3       | Update `CHANGELOG.md`                                      | 15min  | Changelog hygiene             |~~ done at 72085c2
~~| 25 | P3       | Clean up stale docs in `docs/status/` and `docs/planning/` | 30min  | Housekeeping                  |~~ done 2026-10-09 (this pass)

---

## G) TOP #1 QUESTION 🤔

**Should the validate command run a full `pipeline.Run()` with composed detectors (govet + staticcheck + SecurityValidator), or should pipeline usage be a separate `scan` command?**

Running the full pipeline in `validate` would give users a one-command security + quality scan. But it changes the semantics of `validate` from "check SECURITY.md structure" to "run all static analysis." A separate `scan` command would preserve the focused `validate` behavior while adding the composed pipeline as a new capability. This is a UX/product decision.

---

## Build & Test State

```
Build:  CLEAN ✅
Vet:    CLEAN ✅
Tests:  17 PASS / 0 FAIL ✅
Coverage: 23.4% overall
  internal:         33.6%
  cmd:               0.0%
  acceptance:      100.0%
```

## Dependency State

| Dependency      | Version | Role                       |
| --------------- | ------- | -------------------------- |
| `go-finding`    | v0.3.0  | Core data model + pipeline |
| `cobra`         | v1.10.2 | CLI framework              |
| `viper`         | v1.21.0 | Config loading             |
| `ginkgo/v2`     | v2.28.3 | BDD acceptance tests       |
| `gomega`        | v1.40.0 | Test matchers              |
| `testify`       | v1.11.1 | Unit test assertions       |
| `fatih/color`   | v1.19.0 | Terminal colors            |
| `go-branded-id` | REMOVED | —                          |

## Code Stats

```
Total Go lines: 2,346
Files: ~20 .go files
Packages: 4 (cmd, internal, internal/types, test/acceptance)
```

---

_Generated by Crush <crush@charm.land>_
