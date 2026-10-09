# Modularization Proposal — template-SECURITY

> **Never executed — superseded** — the 2026-10-08 securitymd rebuild (`1c55390`) solved the same coupling by deletion: two focused packages (`pkg/policy`, `pkg/provider`), viper and pipeline dropped, dead types gone. No `go.work`, no sub-modules. All findings and migration steps below carry inline verdicts.

_Generated: 2026-05-14 | ~~Status: Self-Reviewed (Phase 4 complete)~~ never executed_

---

## 1. Executive Summary

template-SECURITY is a Go CLI tool that validates and generates `SECURITY.md` files. The project currently has a single `go.mod` with a flat `internal` package mixing 4 distinct concerns (validation, policy generation, project detection, domain types). This proposal splits the project into **4 semi-independent sub-modules** with clear boundaries, enabling:

- **Independent dependency management** — `viper` won't pollute the validator module, `go-finding/pipeline` won't leak into the generator
- **Faster CI** — only rebuild changed modules
- **Cleaner public API** — each module exposes a focused interface
- **Reusability** — the validator and detector can be consumed independently as libraries

### What Changes

| Before                   | After                                 |
| ------------------------ | ------------------------------------- |
| 1 `go.mod`               | 4 `go.mod` files + 1 `go.work`        |
| Flat `internal` package  | 4 focused modules with clear DAG      |
| All deps in one `go.mod` | Each module carries only its own deps |

### Expected Benefits

- **Smaller dependency surface** per module (e.g., detector has zero third-party deps)
- **Compile-time enforced boundaries** between validation, generation, and detection
- **Independent versioning** possible if modules are published separately
- **Dead code isolation** — unused domain types become visible in the core module's zero-consumption

---

## 2. Current State Analysis

### Module Landscape

| Module                                     | Path       | Internal Deps       | External Deps (Direct)                                   | Replace Directives | State    |
| ------------------------------------------ | ---------- | ------------------- | -------------------------------------------------------- | ------------------ | -------- |
| `github.com/LarsArtmann/template-SECURITY` | `/` (root) | N/A (single module) | cobra, viper, color, go-finding, ginkgo, gomega, testify | None               | Monolith |

### Package Dependency Graph

See [DEPENDENCY_GRAPH.md](./DEPENDENCY_GRAPH.md) for full analysis.

### Key Findings

~~1. **`internal` is a catch-all** — 4 files, 4 concerns, 4 different external dependency profiles all crammed into one package~~ moot — internal/ deleted at 1c55390
~~2. **`internal/types` has 90% dead code** — `SecurityPolicy`, `ValidationResult`, `Template`, `Contact`, `Project` are defined but never referenced outside the types package itself~~ moot — types deleted entirely at 1c55390
~~3. **`viper` is only used by `security_tool.go`** — but the entire project carries its transitive dependency tree~~ moot — viper dropped wholesale in the rebuild
~~4. **`go-finding/pipeline` is only used by `security_validator.go`** — the `pipeline.Detector` interface is irrelevant to policy generation~~ moot — pipeline dropped wholesale in the rebuild
~~5. **`testify` is used in unit tests** — but it's in the production `go.mod` `require` block (acceptable since Go handles this, but worth noting)~~ historical — testify still used in the rebuild's unit tests (accepted)
~~6. **No circular dependencies** — the current structure is already a DAG, making modularization straightforward~~ historical — DAG property inherited by the rebuild

### Coupling Hotspots

| Hotspot                                        | Impact                                        | Mitigation                                        |
| ---------------------------------------------- | --------------------------------------------- | ------------------------------------------------- |
| `internal` flat package                        | All concerns share one namespace              | Split into separate modules                       |
| `viper` in `internal`                          | Forces config dep on all consumers            | Isolate to generator module                       |
| Dead types in `internal/types`                 | Misleading API surface                        | Remove dead types, keep only used ones            |
| `SecurityValidator.PrintResults` in `internal` | Output formatting mixed with validation logic | Move to CLI layer or separate presentation module |

---

## 3. Proposed Module Structure

### Module Definitions

| # | Name & Path  | Purpose                                                                | Deps (Prod)                                                                  | Deps (Test)        | Public API                                                                                           |
| - | ------------ | ---------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ------------------ | ---------------------------------------------------------------------------------------------------- |
| 1 | `/core`      | Domain types shared across all modules                                 | None                                                                         | None               | `PolicyType`, `Version`, `VersionStatus`, `ID`, `PolicyID`, `Config`, `PolicyConfig`, `TemplateData` |
| 2 | `/validator` | SECURITY.md validation and finding generation                          | `core`, `go-finding`, `go-finding/pipeline`                                  | `testify`          | `SecurityValidator`, `DetectFile`, `ReportIsValid`                                                   |
| 3 | `/generator` | SECURITY.md template generation and config loading                     | `core`, `go-finding`, `viper`                                                | `testify`          | `SecurityTool`, `NewSecurityTool`                                                                    |
| 4 | `/detector`  | Project name/org/domain detection from git, package.json, go.mod, etc. | None                                                                         | (minimal)          | `ProjectDetector`, `NewProjectDetector`                                                              |
| 5 | `/` (root)   | CLI entry point (`cmd/template-security`)                              | `core`, `validator`, `generator`, `detector`, `cobra`, `color`, `go-finding` | `ginkgo`, `gomega` | CLI binary                                                                                           |

### DAG Verification

```
core ← validator ← (root/cmd)
core ← generator ← (root/cmd)
core ← detector  ← (root/cmd)

No cycles. No lateral dependencies.
```

### Dependency Diagram

```
                  ┌─────────────┐
                  │    core     │
                  │  (types)    │
                  └──────┬──────┘
                         │
       ┌─────────────────┼──────────────────┐
       │                 │                  │
┌──────▼──────┐  ┌───────▼───────┐  ┌──────▼──────┐
│  validator  │  │   generator   │  │  detector   │
│  go-finding │  │  viper        │  │  (stdlib)   │
│  pipeline   │  │  go-finding   │  │             │
└──────┬──────┘  └───────┬───────┘  └──────┬──────┘
       │                 │                  │
       └─────────┬───────┴──────────────────┘
                 │
       ┌─────────▼──────────┐
       │  cmd               │
       │  cobra, color      │
       │  go-finding        │
       └────────────────────┘
```

---

## 4. Replace / Workspace Strategy

**Chosen: `go.work` at repo root**

| Strategy             | Decision    | Rationale                                           |
| -------------------- | ----------- | --------------------------------------------------- |
| `go.work` file       | ✅ Selected | 4 modules in same repo, not published independently |
| `replace` directives | ❌ Not used | Per-module replace is messy with 4+ modules         |
| Versioned imports    | ❌ Not used | Not published to Go proxy                           |

### Rules

1. `go.work` at repo root lists all 5 modules
2. No `replace` directives in any `go.mod`
3. `go.work` is committed (not in `.gitignore`) since all modules live in the same repo
4. `go mod tidy` must work both with and without the workspace

### go.work Structure

```go
go 1.26.2

use (
    .
    ./core
    ./validator
    ./generator
    ./detector
)
```

---

## 5. Test Dependency Isolation

| Module       | Production Deps                                                              | Test-Only Deps     |
| ------------ | ---------------------------------------------------------------------------- | ------------------ |
| `core`       | None                                                                         | None               |
| `validator`  | `core`, `go-finding`, `go-finding/pipeline`                                  | `testify`          |
| `generator`  | `core`, `go-finding`, `viper`                                                | `testify`          |
| `detector`   | None                                                                         | None               |
| root (`cmd`) | `core`, `validator`, `generator`, `detector`, `cobra`, `color`, `go-finding` | `ginkgo`, `gomega` |

### Testhelpers Strategy

- **No shared testhelpers module** — the existing `test/acceptance/test_helpers.go` is small (27 lines) and stays at the root
- Each module has its own `_test.go` files within the module
- Acceptance tests in `test/acceptance/` remain at root level, importing modules as needed

---

## 6. Interface Extraction

No interface/implementation split is needed for this project. The modules are already thin enough that each exposes concrete types. However:

### Moves Required

| From                                | To                   | What                                                             |
| ----------------------------------- | -------------------- | ---------------------------------------------------------------- |
| `internal/types/types.go`           | `core/`              | `PolicyType`, `Version`, `VersionStatus`                         |
| `internal/types/ids.go`             | `core/`              | `ID`, `PolicyID`                                                 |
| `internal/security_tool.go` (types) | `core/`              | `Config`, `PolicyConfig`, `TemplateData`                         |
| `internal/security_validator.go`    | `validator/`         | `SecurityValidator`, `DetectFile`, `ReportIsValid`, helper funcs |
| `internal/security_tool.go` (logic) | `generator/`         | `SecurityTool`, `NewSecurityTool`, error vars, all methods       |
| `internal/project_detector.go`      | `detector/`          | `ProjectDetector`                                                |
| `cmd/template-security/*`           | stays in root `cmd/` | No change                                                        |

### Types to Remove (Dead Code)

The following types in `internal/types/types.go` are defined but never used. They should be **removed** before or during modularization:

- `SecurityPolicy`
- `ValidationResult`
- `ValidationMessage` (`Error`, `Warning` aliases)
- `Template`
- `TemplateVariable`
- `Contact`
- `ContactType` and its constants (`ContactTypeEmail`, `ContactTypeWeb`, `ContactTypeAPI`)
- `Project`

---

## 7. Versioning Strategy

| Strategy             | Decision    | Rationale                                |
| -------------------- | ----------- | ---------------------------------------- |
| Shared version       | ✅ Selected | Single team, tight coupling, single repo |
| Independent semver   | ❌          | Overkill for a small CLI tool            |
| Root-only versioning | ❌          | Same as shared version in practice       |

All modules share a single git tag (`v1.2.3`). If a module is ever extracted to its own repo, it can be re-versioned at that point.

---

## 8. Migration Strategy

Ordered steps, each independently executable:

~~1. **Remove dead types** from `internal/types/` — pure deletion, no import changes~~ done at 1c55390 — all types deleted
~~2. **Extract `core` module** — move `internal/types/` → `core/`, create `go.mod`, move used types only~~ NOT-DO — modularization never executed; superseded by the rebuild (two packages: pkg/policy, pkg/provider)
~~3. **Extract `detector` module** — move `internal/project_detector.go` → `detector/`, create `go.mod`~~ NOT-DO — superseded: detection lives in pkg/policy/project.go
~~4. **Extract `validator` module** — move `internal/security_validator.go` + tests → `validator/`, create `go.mod`~~ NOT-DO — superseded: validation lives in pkg/policy/validate.go
~~5. **Extract `generator` module** — move `internal/security_tool.go` + tests → `generator/`, create `go.mod`, carry `viper`~~ NOT-DO — superseded: generation lives in pkg/policy/generate.go
~~6. **Create `go.work`** — wire all modules together~~ NOT-DO — single module in the rebuild
~~7. **Update root `go.mod`** — remove moved code, depend on sub-modules~~ moot — go.mod rewritten by the rebuild
~~8. **Update `cmd/template-security/` imports** — point to new module paths~~ moot — CLI rebuilt as cmd/securitymd
~~9. **Update `test/acceptance/` imports** — point to new module paths~~ moot — acceptance tests rewritten
~~10. **Verify build, test, lint** — full green suite~~ done — the rebuild's own gate (tests green, lint 0, flake check green)
~~11. **Update documentation** — README, AGENTS.md, flake.nix if applicable~~ done at 72085c2/4a8987a/519916d — docs pass

---

## 9. Risk Assessment

| Risk                                         | Likelihood | Impact | Mitigation                                               |
| -------------------------------------------- | ---------- | ------ | -------------------------------------------------------- |
| Import path breaks across modules            | Medium     | High   | Update all imports systematically; build after each step |
| `viper` config loading breaks in generator   | Low        | Medium | Test config loading after extraction                     |
| Test file package declarations need updating | Medium     | Low    | Tests move with their module                             |
| `go.work` causes issues in CI                | Low        | High   | Test `go mod tidy` without workspace                     |
| Acceptance tests break due to import changes | Medium     | Medium | Update acceptance test imports last                      |
| Existing CI workflow breaks                  | Medium     | Low    | CI workflow is already broken (Go 1.21 vs 1.26.2)        |

---

## 10. Build System Impact

| Component            | Change Required                                             |
| -------------------- | ----------------------------------------------------------- |
| `go.mod` (root)      | Slim down to only CLI deps + sub-module requires            |
| `go.work` (new)      | Create with all 5 modules                                   |
| `justfile`           | Update build/test commands if they reference specific paths |
| `.github/workflows/` | Update to use `go.work`, fix Go version                     |
| `.golangci.yml`      | May need path updates for lint targets                      |
| `flake.nix`          | Does not exist yet — create after modularization            |
| `scripts/`           | Legacy shell scripts — out of scope for this modularization |

---

## 11. Self-Review (Phase 4)

### Findings Incorporated

1. **`Config`, `PolicyConfig`, `TemplateData` moved to `core`** — these are shared between generator and CLI. Placing them in `core` avoids the generator module being a dependency of the CLI just for type definitions.

2. **`ReportIsValid` stays in `validator`** — it's a validation concern used by both the CLI and acceptance tests. Consumers import `validator` to use it. No split brain.

3. **Banned dependencies noted but not in scope:**
   - `viper` → should be `koanf` (banned per how-to-golang)
   - `testify` → should be `ginkgo/v2 + gomega` (banned per how-to-golang)
   - These migrations are **orthogonal** to modularization and should be addressed separately to avoid scope explosion.

4. **Module granularity is correct** — the detector has zero deps and is genuinely independent. Splitting further (e.g., separating `Config` loading from `SecurityTool`) would create modules too thin to justify their own `go.mod`.

5. **No split brains introduced** — every type is defined in exactly one module. The DAG is clean with no lateral dependencies.

6. **Existing CI is already broken** — Go 1.21 in CI vs 1.26.2 in go.mod. Modularization doesn't make this worse.

7. **`go.work` is the correct strategy** — verified no mixing with `replace` directives. `go.work` is committed to the repo since all modules are co-located.
