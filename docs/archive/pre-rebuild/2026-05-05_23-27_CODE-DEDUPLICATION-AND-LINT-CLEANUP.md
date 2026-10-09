# Comprehensive Project Status — template-SECURITY

> **Superseded 2026-10-08** — describes the pre-rebuild tree (`internal/`, `cmd/template-security/`, `scripts/`), deleted in the securitymd rebuild (`1c55390`). Successor report: `docs/status/2026-10-08_23-40_securitymd-rebuild-buildflow-provider-status.md`. All open items below carry inline verdicts.

**Date:** 2026-05-05 23:27 | **Branch:** master | **Head:** f29259c | **Ahead:** 4 commits from origin

---

## Executive Summary

Session 3 today: completed code deduplication (art-dupl → 0 clones), linting configuration overhaul (.golangci.yml), and CLI helper cleanup. All 22 tests pass (9 acceptance + 13 unit). Build and vet clean. 52 lint warnings remain (down from earlier sessions). Total coverage 22.6% — cmd package at 0%. The project has a solid validation/generation core but significant gaps in test coverage, documentation freshness, and dead code.

**Overall Health: 🟡 FUNCTIONAL BUT NEEDS WORK**

| Metric               | Value       | Status |
| -------------------- | ----------- | ------ |
| Build                | Clean       | ✅     |
| `go vet`             | Clean       | ✅     |
| Tests                | 22/22 pass  | ✅     |
| Coverage             | 22.6% total | 🔴     |
| Lint warnings        | 52          | 🟡     |
| Code duplication     | 0 clones    | ✅     |
| `cmd/` coverage      | 0.0%        | 🔴     |
| `internal/` coverage | 32.6%       | 🟡     |
| Uncommitted changes  | 6 files     | 🔴     |
| Stale docs           | Multiple    | 🟡     |
| Dead code            | Present     | 🟡     |

---

## A) FULLY DONE ✅

_Historically accurate for the pre-rebuild tree (deleted `1c55390`). The addFinding/Builder/closure patterns below died with it; the rebuild re-implements the same concerns in `pkg/policy`._

### 1. Core Validation Engine

- `internal/security_validator.go` — Full SECURITY.md validation with:
  - Required section checks (header, reporting, versions, practices, contact)
  - Content quality checks (min lines, template variables, substantive content, version info)
  - `pipeline.Detector` interface implementation with context cancellation
  - `addFinding` helper (extracted this session — eliminates duplication)
  - `buildFinding` Builder API with consistent category/tags/confidence
- `internal/security_validator_test.go` — 13 unit tests covering all validation paths
- `test/acceptance/` — 9 Ginkgo BDD acceptance tests (100% coverage of acceptance suite)

### 2. go-finding Integration (Sessions 1–2)

| Integration Point                        | Status      |
| ---------------------------------------- | ----------- |
| `finding.Report` as data model           | ✅ Complete |
| `finding.Builder` for all findings       | ✅ Complete |
| `finding.FindingError` for error types   | ✅ Complete |
| `pipeline.Detector` interface            | ✅ Complete |
| `WriteJSON()` / `WriteSARIFFiltered()`   | ✅ Complete |
| Severity filtering with `-severity` flag | ✅ Complete |
| Context-aware `Detect(ctx)`              | ✅ Complete |
| `finding.Filter` predicates              | ✅ Complete |

### 3. CLI Commands

- `template-security validate` — Validate SECURITY.md files (text/JSON/SARIF output)
- `template-security setup` — Generate SECURITY.md from templates
- `template-security status` — Show security policy status

### 4. Code Deduplication (This Session)

- Extracted `addFinding` method on `SecurityValidator` — eliminated 4× repeated `buildFinding` + error check pattern
- Inlined `runSetup` and `runValidate` as closures — eliminated semantic clone detection on cobra function signatures
- **Result: 0 clone groups** (art-dupl --semantic -t 15)

### 5. Lint Configuration Overhaul (This Session)

- `.golangci.yml` restructured: disabled `depguard` (was noise), tuned `exhaustruct` excludes, `varnamelen` min-length 3, `funlen` reduced to 60 lines, added `gosec`, `gochecknoinits`, `paralleltest` config
- Renamed `printInfo` → `printInfof`, `printWarning` → `printWarningf` (goprintffuncname)
- Removed unused `printSuccess`, `printError` dead code
- Added doc comments to global vars in `setup.go` (gochecknoglobals)

### 6. Template Generation

- `internal/security_tool.go` — Template loading, processing via `text/template`, config loading via viper
- `templates/SECURITY.md` — Base template with variables
- `.template-security.yaml` — Config file support

---

## B) PARTIALLY DONE 🟡

### 1. Lint Warnings (52 total) — done

0 issues since the 2026-10-08 rebuild (config carried over and fixed in code, not suppression).

| Category           | Count | Examples                                                   |
| ------------------ | ----- | ---------------------------------------------------------- |
| `depguard`         | ~8    | Imports not in allow-list (disabled but still warning)     |
| `gochecknoglobals` | ~7    | CLI flag vars (orgName, contactEmail, etc.)                |
| `exhaustruct`      | ~3    | cobra.Command, finding.ToolInfo, SecurityValidator         |
| `forbidigo`        | ~5    | fmt.Printf in PrintResults, printTextReport                |
| `varnamelen`       | ~3    | `f` for findings, `i` for index                            |
| `cyclop`           | 2     | `runSetup` (19), `validateContentQuality` (11) — both > 10 |
| `wrapcheck`        | ~3    | Unwrapped errors from external packages                    |
| `err113`           | 2     | Dynamic errors via fmt.Errorf                              |
| `unused`           | 2     | `printSuccess`, `printError` — removed in working tree     |
| Other              | ~17   | revive, goconst, goprintffuncname, etc.                    |

### 2. CI/CD Pipeline — done at `72085c2`

Workflow rewritten: go-version-file, build+test+self-validate, no make.

- GitHub Actions workflow exists but references Go 1.21 (project uses 1.26.2)
- Uses `make build || just build || go build` fallback — no makefile exists, justfile deprecated
- No golangci-lint step in CI
- No test running in CI
- No coverage reporting

### 3. Domain Types — done (resolved by deletion at `1c55390`)

The split brain was killed by deleting `internal/types/` wholesale; `finding.Finding` is the only model in the rebuild.

- `internal/types/types.go` — Rich domain model exists (SecurityPolicy, ValidationResult, Template, Contact, Version, etc.) but **most types are unused**
- `internal/types/ids.go` — ID types defined but only `ID` and `PolicyID` are used
- Validation uses `finding.Finding` instead of custom `ValidationResult` type
- This is a "split brain" — two parallel type systems

### 4. Documentation — done at `72085c2`/`4a8987a` (successor docs pass)

README/CHANGELOG/AGENTS rewritten; legacy docs archived.

- `README.md` — Accurate for features but missing go-finding integration details
- `IMPROVEMENT_PLAN.md` — Stale (references old architecture, many items done or irrelevant)
- `BDD_TESTS_REVIEW.md` — Exists but not reviewed for freshness
- `CHANGELOG.md` — Not updated for recent changes
- `PARTS.md`, `PROJECT_SPLIT_EXECUTIVE_REPORT.md`, `PUBLIC_OR_PRIVATE.md` — Legacy docs from project split, potentially stale

---

## C) NOT STARTED 🔴

### Critical Gaps

~~1. **`cmd/` test coverage at 0.0%** — No tests for validate, setup, status commands~~ NOT-DO — old CLI deleted at 1c55390; successor cmd/ is test-free by design (provider + acceptance coverage)
~~2. **No integration tests** — No test that runs the full CLI binary end-to-end~~ done — rebuild ships provider detect→repair→verify loop tests in a real git repo
~~3. **No `FEATURES.md`** — No feature inventory document~~ done at 4a8987a
~~4. **No `TODO_LIST.md`** — No tracked TODO list~~ done at 4a8987a
~~5. **No flake.nix** — AGENTS.md mandates nix flakes, project has justfile + shell scripts~~ done at 9d3094a (2026-06-17)
~~6. **`scripts/` directory** — 5 shell scripts (build.sh, compliance-check.sh, etc.) exist but:~~ done at 1c55390 — scripts/ deleted (~1,300 lines)
   - Not integrated with any build system
   - `build.sh` duplicates functionality of `go build`
   - `compliance-check.sh` (8.8KB) and `security-setup.sh` (16KB) — legacy pre-Go shell scripts
~~7. **`report/` and `reports/` directories** — Empty directories, purpose unclear~~ done — directories/artifacts gone (reports/ trashed 2026-10-08, .gitignore covers it)
~~8. **No `AGENTS.md` at project level** — Only global ~/.config/crush/AGENTS.md exists~~ done — AGENTS.md rewritten for the rebuild at 72085c2
~~9. **SARIF output not tested** — No test for `-format sarif`~~ routed — TODO_LIST keeps the SARIF golden-file test for the rebuilt CLI
~~10. **JSON output not tested** — No test for `-format json`~~ routed — same TODO_LIST golden-file item covers output-format testing

### Missing Quality Infrastructure

~~11. **No golangci-lint in CI**~~ NOT-DO — CI is build+test+self-validate; the 0-issue lint bar is enforced locally/BuildFlow
~~12. **No code coverage enforcement**~~ NOT-DO — no numeric coverage gate by design; three test suites gate behavior
~~13. **No release automation** (goreleaser, etc.)~~ NOT-DO — publish via git tag (TODO_LIST checklist); goreleaser never adopted
~~14. **No version tagging strategy**~~ routed — TODO_LIST publish checklist (tag v1.0.0)
~~15. **No CONTRIBUTING.md**~~ done at 8864cb7 (added 2026-05-28)

---

## D) TOTALLY FUCKED UP 💥

### 1. Split Brain: Two Type Systems — done at `1c55390`

Deleted. The rebuild emits `finding.Finding` only; no parallel domain model exists.

The project has **two parallel type systems** that don't talk to each other:

- **`internal/types/types.go`**: Rich domain model — `SecurityPolicy`, `ValidationResult`, `Template`, `Contact`, `Version`, `Project`, `Error`, `Warning`, etc.
- **`go-finding`**: `finding.Finding`, `finding.Report`, `finding.Severity`, `finding.Builder`

The validation engine uses `finding.*` exclusively. The domain types are **dead code** — defined but never instantiated or used by any business logic. This creates confusion about which is the "real" model.

### 2. Cyclomatic Complexity — moot

`runSetup` and `validateContentQuality` were deleted with the old tree.

- `runSetup`: complexity 19 (max 10) — deeply nested config/flag detection logic
- `validateContentQuality`: complexity 11 (max 10) — partially improved by `addFinding` extraction

### 3. Dead Shell Scripts — done at `1c55390`

`scripts/` deleted in the rebuild.

`scripts/` contains 5 shell scripts totaling ~33KB that predate the Go rewrite. They duplicate Go functionality and are not referenced by any build system. `compliance-check.sh` and `security-setup.sh` are particularly large and likely unmaintained.

### 4. Empty Directories — done

`report/`/`reports/` are gone (reports/ artifacts trashed 2026-10-08).

- `report/` — Empty
- `reports/` — Empty

These suggest abandoned features or leftover from shell-to-Go migration.

### 5. Global Mutable State in CLI — moot

The rebuilt CLI (`cmd/securitymd`) has no package-level flag globals.

`setup.go` uses 5 package-level `var` for CLI flags (orgName, contactEmail, policyType, outputDir, quickMode). Cobra best practice is struct-based command state. These globals make testing the setup command impossible in isolation.

### 6. `status.go` Still Uses `runStatus` Pattern — moot

status.go was rewritten in the rebuild.

While `setup.go` and `validate.go` were inlined to closures, `status.go` still passes `runStatus` function to `newCommand` — inconsistent pattern.

---

## E) WHAT WE SHOULD IMPROVE 📈

### Architecture

~~1. **Resolve the type split brain** — Decide: use `internal/types/` as the domain model with `go-finding` as output format, OR delete the unused types. The current state is the worst of both worlds.~~ done — resolved by deletion: internal/types removed at 1c55390; finding.Finding is the only model
~~2. **Extract CLI command state into structs** — Replace global vars with `type setupOptions struct { orgName, contactEmail, ... }` passed through cobra command context.~~ moot — CLI rebuilt at 1c55390 without package-level flag globals
~~3. **Reduce cyclomatic complexity** — Break `runSetup` (19 → <10) and `validateContentQuality` (11 → <10) into smaller focused functions.~~ moot — runSetup/validateContentQuality deleted with the old tree

### Testing

~~4. **Target 60%+ coverage** — Currently 22.6%. `cmd/` at 0% is unacceptable.~~ NOT-DO — no numeric target in the successor; unit + provider + BDD suites gate behavior
~~5. **Test all output formats** — JSON and SARIF output paths have zero test coverage.~~ routed — TODO_LIST SARIF golden-file item
~~6. **Add CLI integration tests** — Use `cobra.test` or exec-based tests against the built binary.~~ done — provider loop tests + BDD acceptance in the rebuild

### Quality

~~7. **Fix all 52 lint warnings** — Many are easy wins (varnamelen, wrapcheck, err113).~~ done — 0 issues since the rebuild (2026-10-08 report item 9)
~~8. **Delete dead code** — Unused types in `internal/types/`, empty directories, legacy shell scripts.~~ done at 1c55390 — dead types, sentinels, and scripts all deleted
~~9. **Update CI to Go 1.26** — Currently references Go 1.21, which doesn't support the language features used.~~ done at 72085c2 — go-version-file: go.mod
~~10. **Add golangci-lint to CI** — Linting should be automated, not manual.~~ NOT-DO — lint gate stays local/BuildFlow

### Documentation

~~11. **Create `FEATURES.md`** — Feature inventory with status indicators.~~ done at 4a8987a
~~12. **Create `TODO_LIST.md`** — Tracked from codebase analysis.~~ done at 4a8987a
~~13. **Update `IMPROVEMENT_PLAN.md`** — Currently references pre-go-finding architecture.~~ done at 72085c2 — archived to docs/archive/pre-rebuild/
~~14. **Create project-level `AGENTS.md`** — Domain knowledge, build commands, conventions.~~ done at 72085c2
~~15. **Update `CHANGELOG.md`** — Missing entries for go-finding migration and deduplication work.~~ done at 72085c2

### Build System

~~16. **Migrate to `flake.nix`** — Per AGENTS.md mandate. Replace justfile + shell scripts.~~ done at 9d3094a
~~17. **Add `goreleaser`** — Proper cross-platform builds and releases.~~ NOT-DO

---

## F) TOP 25 THINGS TO DO NEXT

| #  | Priority | Task                                                                                    | Impact | Effort |
| -- | -------- | --------------------------------------------------------------------------------------- | ------ | ------ |
~~| 1  | 🔴 P0    | Resolve type split brain — delete unused `internal/types/` or integrate with go-finding | High   | 2h     |~~ done — deleted at 1c55390
~~| 2  | 🔴 P0    | Add tests for `cmd/` — validate, setup, status commands                                 | High   | 3h     |~~ NOT-DO — provider + acceptance coverage in the successor
~~| 3  | 🔴 P0    | Fix CI workflow — Go 1.26, add test step, add golangci-lint step                        | High   | 1h     |~~ done at 72085c2 (build+test+self-validate; lint stays local)
~~| 4  | 🔴 P0    | Delete dead shell scripts in `scripts/` (or migrate to nix)                             | Medium | 30m    |~~ done at 1c55390
~~| 5  | 🔴 P0    | Delete empty `report/` and `reports/` directories                                       | Low    | 1m     |~~ done — gone
~~| 6  | 🟠 P1    | Reduce `runSetup` complexity from 19 to <10                                             | High   | 1h     |~~ moot — deleted
~~| 7  | 🟠 P1    | Test JSON output format (`-format json`)                                                | Medium | 30m    |~~ routed — TODO_LIST golden-file item
~~| 8  | 🟠 P1    | Test SARIF output format (`-format sarif`)                                              | Medium | 30m    |~~ routed — TODO_LIST golden-file item
~~| 9  | 🟠 P1    | Fix all lint warnings (52 → 0)                                                          | Medium | 2h     |~~ done — 0 issues since rebuild
~~| 10 | 🟠 P1    | Create project-level `AGENTS.md` with build/test/lint commands                          | Medium | 30m    |~~ done at 72085c2
~~| 11 | 🟠 P1    | Create `FEATURES.md` feature inventory                                                  | Medium | 1h     |~~ done at 4a8987a
~~| 12 | 🟠 P1    | Create `TODO_LIST.md` from codebase analysis                                            | Medium | 1h     |~~ done at 4a8987a
~~| 13 | 🟡 P2    | Convert CLI globals to struct-based command state                                       | Medium | 1h     |~~ moot — rebuilt CLI has no flag globals
~~| 14 | 🟡 P2    | Inline `runStatus` in `status.go` (consistent with setup/validate pattern)              | Low    | 10m    |~~ moot — status.go rewritten
~~| 15 | 🟡 P2    | Update `README.md` with go-finding integration details                                  | Medium | 30m    |~~ done at 519916d
~~| 16 | 🟡 P2    | Update `CHANGELOG.md` for sessions 1–3                                                  | Low    | 30m    |~~ done at 72085c2
~~| 17 | 🟡 P2    | Update `IMPROVEMENT_PLAN.md` to reflect current state                                   | Medium | 30m    |~~ done — archived at 72085c2
~~| 18 | 🟡 P2    | Add `CONTRIBUTING.md`                                                                   | Medium | 30m    |~~ done at 8864cb7
~~| 19 | 🟡 P2    | Migrate justfile → flake.nix (per AGENTS.md mandate)                                    | High   | 3h     |~~ done at 9d3094a
~~| 20 | 🟢 P3    | Add goreleaser for release automation                                                   | Medium | 2h     |~~ NOT-DO
~~| 21 | 🟢 P3    | Add coverage enforcement (minimum threshold in CI)                                      | Medium | 30m    |~~ NOT-DO
~~| 22 | 🟢 P3    | Review `BDD_TESTS_REVIEW.md` for freshness                                              | Low    | 15m    |~~ done — archived at 72085c2
~~| 23 | 🟢 P3    | Review and clean legacy docs (PARTS.md, PROJECT_SPLIT_EXECUTIVE_REPORT.md)              | Low    | 30m    |~~ done 2026-10-09 (this pass)
~~| 24 | 🟢 P3    | Add end-to-end CLI tests (exec-based, test the binary)                                  | High   | 2h     |~~ done — provider loop tests
~~| 25 | 🟢 P3    | Add version tagging and release workflow                                                | Medium | 1h     |~~ routed — TODO_LIST publish checklist

---

## G) TOP #1 QUESTION I CANNOT FIGURE OUT MYSELF 🤔

**~~Should `internal/types/` be deleted entirely, or should the domain types be used as the internal model with `go-finding` serving purely as the output/reporting layer?~~** — RESOLVED 2026-10-08: deleted entirely. `finding.Finding` is the single model (family doctrine: no converter glue); the dogfood test pins the template↔validator contract instead of a parallel type system.

Arguments for keeping `internal/types/`:

- Richer domain model (Version with dates/status, Contact with response times, Template with variables)
- Could serve as the "source of truth" that gets converted to `finding.Finding` for output
- Future expansion (template management, version tracking) would need these types anyway

Arguments for deleting:

- Currently 100% dead code — nothing instantiates any of these types
- `go-finding` already provides a complete model (Finding, Report, Severity, Category, Tags)
- Maintaining two parallel type systems is a maintenance burden
- YAGNI — no current feature needs these types

The decision here fundamentally shapes the architecture going forward. It's your call.

---

## File Change Summary (This Session — Uncommitted)

| File                                   | Changes                                                                                    |
| -------------------------------------- | ------------------------------------------------------------------------------------------ |
| `.golangci.yml`                        | Restructured: disabled depguard, added exhaustruct excludes, tuned linters                 |
| `cmd/template-security/cmd_helpers.go` | Renamed printInfo→printInfof, printWarning→printWarningf; removed printSuccess, printError |
| `cmd/template-security/setup.go`       | Inlined runSetup as closure; added doc comments to globals                                 |
| `cmd/template-security/status.go`      | Updated printInfo/printWarning calls to new names                                          |
| `cmd/template-security/validate.go`    | Inlined runValidate as closure                                                             |
| `internal/security_validator.go`       | Extracted addFinding helper; reformatted unresolved-template call                          |

---

## Codebase Metrics

| Metric                      | Value                                              |
| --------------------------- | -------------------------------------------------- |
| Total Go LOC                | 2,330                                              |
| Go source files             | 12                                                 |
| Packages                    | 4 (cmd, internal, internal/types, test/acceptance) |
| Tests                       | 22                                                 |
| Test pass rate              | 100%                                               |
| Total coverage              | 22.6%                                              |
| `internal/` coverage        | 32.6%                                              |
| `cmd/` coverage             | 0.0%                                               |
| `test/acceptance/` coverage | 100%                                               |
| Lint warnings               | 52                                                 |
| Clone groups                | 0                                                  |
| Dependencies (direct)       | 7                                                  |
| Go version                  | 1.26.2                                             |

## Resolution (2026-10-09)

Every open item in B/C/D/E/F/G carries an inline verdict. Reference hashes: `1c55390` (old tree deleted), `72085c2`/`4a8987a`/`519916d` (docs pass), `9d3094a` (flake.nix), `8864cb7` (CONTRIBUTING.md). Surviving intents (SARIF/JSON golden test, version tagging) are routed to TODO_LIST. Archivable: no open items remain.
