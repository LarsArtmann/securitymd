# Execution Plan — template-SECURITY Modularization

_Generated: 2026-05-14 | Status: Ready for Execution_

---

## Overview

This plan breaks the modularization into 14 ordered tasks. Each task:
- Takes 15–30 minutes
- Leaves the project buildable and testable
- Is independently revertable (single commit)

**Total estimated effort:** 4–6 hours

---

## Pareto Impact Tiers

| Tier | Tasks | Impact |
|---|---|---|
| **1% → 51%** (Foundational) | T1, T2, T3, T4 | Remove dead code, extract core types, create module structure |
| **4% → 64%** (High leverage) | T5, T6, T7, T8, T9 | Extract validator, generator, detector modules |
| **20% → 80%** (Broad value) | T10, T11, T12 | Wire go.work, update CLI and acceptance tests |
| **Remaining** (Polish) | T13, T14 | Verify full suite, update documentation |

---

## Task List

### T1: Remove Dead Domain Types

**Impact:** 1% → 51% | **Effort:** 10 min | **Depends on:** Nothing

Remove unused types from `internal/types/types.go`:
- `SecurityPolicy`, `ValidationResult`, `ValidationMessage` (and aliases `Error`, `Warning`)
- `Template`, `TemplateVariable`
- `Contact`, `ContactType`, `ContactTypeEmail`, `ContactTypeWeb`, `ContactTypeAPI`
- `Project`

Keep only: `PolicyType`, `PolicyTypeGitHub`, `PolicyTypeEnterprise`, `VersionStatus`, `StatusSupported`, `StatusDeprecated`, `StatusEOL`, `Version`

**Verification:**
```bash
go build ./...
go test ./...
go vet ./...
```

**Rollback:** `git revert HEAD` — restores deleted types

---

### T2: Create `core` Module

**Impact:** 1% → 51% | **Effort:** 15 min | **Depends on:** T1

1. Create `core/` directory
2. Create `core/go.mod`:
   ```
   module github.com/LarsArtmann/template-SECURITY/core
   go 1.26.2
   ```
3. Move `internal/types/types.go` → `core/types.go` (rename package to `core`)
4. Move `internal/types/ids.go` → `core/ids.go` (rename package to `core`)
5. Extract `Config`, `PolicyConfig`, `TemplateData` from `internal/security_tool.go` → `core/config.go` (rename package to `core`)
6. Update all imports throughout the codebase:
   - `"github.com/LarsArtmann/template-SECURITY/internal/types"` → `"github.com/LarsArtmann/template-SECURITY/core"`
   - References to `types.PolicyType` → `core.PolicyType`, etc.
7. Remove `internal/types/` directory

**Verification:**
```bash
cd core && go build ./... && go vet ./...
cd /home/lars/projects/template-SECURITY && go build ./... && go test ./...
```

**Rollback:** `git revert HEAD` — restores `internal/types/` and original imports

---

### T3: Create `detector` Module

**Impact:** 4% → 64% | **Effort:** 15 min | **Depends on:** T2

1. Create `detector/` directory
2. Create `detector/go.mod`:
   ```
   module github.com/LarsArtmann/template-SECURITY/detector
   go 1.26.2
   ```
3. Move `internal/project_detector.go` → `detector/project_detector.go`
4. Update package declaration to `detector`
5. Update imports in `cmd/template-security/setup.go`:
   - `"github.com/LarsArtmann/template-SECURITY/internal"` → uses `detector` import
6. Remove `internal/project_detector.go`

**Verification:**
```bash
cd detector && go build ./... && go vet ./...
go build ./... && go test ./...
```

**Rollback:** `git revert HEAD`

---

### T4: Create `validator` Module

**Impact:** 4% → 64% | **Effort:** 20 min | **Depends on:** T2

1. Create `validator/` directory
2. Create `validator/go.mod`:
   ```
   module github.com/LarsArtmann/template-SECURITY/validator
   go 1.26.2

   require (
       github.com/LarsArtmann/template-SECURITY/core v0.0.0
       github.com/larsartmann/go-finding v0.3.0
   )

   replace github.com/LarsArtmann/template-SECURITY/core => ../core
   ```
3. Move `internal/security_validator.go` → `validator/security_validator.go`
4. Move `internal/security_validator_test.go` → `validator/security_validator_test.go`
5. Update package to `validator`
6. Update imports:
   - Remove `"github.com/LarsArtmann/template-SECURITY/internal/types"` references
   - Add `"github.com/LarsArtmann/template-SECURITY/core"` for any shared types (if needed)
   - Local references to `toolName` stay as-is (defined in same package)
7. Run `go mod tidy` in `validator/`
8. Update consumers:
   - `cmd/template-security/validate.go` — import `validator` instead of `internal`
   - `test/acceptance/validation_test.go` — import `validator`

**Verification:**
```bash
cd validator && go build ./... && go test ./... && go vet ./...
go mod tidy && go build ./...
```

**Rollback:** `git revert HEAD`

---

### T5: Create `generator` Module

**Impact:** 4% → 64% | **Effort:** 20 min | **Depends on:** T2

1. Create `generator/` directory
2. Create `generator/go.mod`:
   ```
   module github.com/LarsArtmann/template-SECURITY/generator
   go 1.26.2

   require (
       github.com/LarsArtmann/template-SECURITY/core v0.0.0
       github.com/larsartmann/go-finding v0.3.0
       github.com/spf13/viper v1.21.0
   )

   replace github.com/LarsArtmann/template-SECURITY/core => ../core
   ```
3. Move `internal/security_tool.go` → `generator/security_tool.go`
   - Exclude `Config`, `PolicyConfig`, `TemplateData` (already in `core`)
   - Exclude `ErrConfigNotFound`, `ErrInvalidConfig`, etc. — move these to `generator/` since they're generator-specific
4. Move `internal/security_tool_test.go` → `generator/security_tool_test.go`
5. Update package to `generator`
6. Update imports — reference `core` for `Config`, `PolicyConfig`, `TemplateData`, `PolicyType`, `Version`
7. Run `go mod tidy` in `generator/`
8. Update consumers:
   - `cmd/template-security/setup.go` — import `generator` + `core` instead of `internal`
   - `test/acceptance/policy_generation_test.go` — import `generator` + `core`

**Verification:**
```bash
cd generator && go build ./... && go test ./... && go vet ./...
go mod tidy && go build ./...
```

**Rollback:** `git revert HEAD`

---

### T6: Remove Empty `internal/` Package

**Impact:** 4% → 64% | **Effort:** 5 min | **Depends on:** T3, T4, T5

1. Verify `internal/` is empty (all files moved out)
2. Remove `internal/` directory
3. Run full build and test suite

**Verification:**
```bash
go build ./... && go test ./... && go vet ./...
```

**Rollback:** `git revert HEAD`

---

### T7: Create `go.work`

**Impact:** 20% → 80% | **Effort:** 10 min | **Depends on:** T2, T3, T4, T5

1. Create `go.work` at repo root:
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
2. Remove all `replace` directives from individual `go.mod` files (since `go.work` handles local resolution)
3. Run `go work sync`
4. Run full build: `go build ./...`

**Verification:**
```bash
go work sync
go build ./...
go test ./...
```

**Rollback:** Delete `go.work`, restore `replace` directives

---

### T8: Update Root `go.mod`

**Impact:** 20% → 80% | **Effort:** 15 min | **Depends on:** T6, T7

1. Slim down root `go.mod` to only depend on:
   - Sub-modules: `core`, `validator`, `generator`, `detector`
   - CLI deps: `cobra`, `color`, `go-finding`
   - Test deps: `ginkgo/v2`, `gomega`
2. Remove production deps that moved to sub-modules:
   - `viper` (moved to `generator`)
   - `go-finding/pipeline` (moved to `validator`)
3. Run `go mod tidy`
4. Verify root still builds and tests pass

**Verification:**
```bash
go mod tidy
go build ./...
go test ./...
go vet ./...
```

**Rollback:** `git revert HEAD`

---

### T9: Update CLI Imports

**Impact:** 20% → 80% | **Effort:** 15 min | **Depends on:** T8

Update all files in `cmd/template-security/`:

| File | Old Import | New Import |
|---|---|---|
| `setup.go` | `"github.com/LarsArtmann/template-SECURITY/internal"` | `"github.com/LarsArtmann/template-SECURITY/generator"` + `"github.com/LarsArtmann/template-SECURITY/detector"` |
| `setup.go` | `"github.com/LarsArtmann/template-SECURITY/internal/types"` | `"github.com/LarsArtmann/template-SECURITY/core"` |
| `validate.go` | `"github.com/LarsArtmann/template-SECURITY/internal"` | `"github.com/LarsArtmann/template-SECURITY/validator"` |
| `status.go` | (no internal imports) | No change |
| `main.go` | (no internal imports) | No change |

Update all type references:
- `internal.SecurityValidator` → `validator.SecurityValidator`
- `internal.NewSecurityValidator` → `validator.NewSecurityValidator`
- `internal.ReportIsValid` → `validator.ReportIsValid`
- `internal.SecurityTool` → `generator.SecurityTool`
- `internal.NewSecurityTool` → `generator.NewSecurityTool`
- `internal.Config` → `core.Config`
- `internal.PolicyConfig` → `core.PolicyConfig`
- `internal.NewProjectDetector` → `detector.NewProjectDetector`
- `types.PolicyType(...)` → `core.PolicyType(...)`

**Verification:**
```bash
go build ./...
go test ./...
```

**Rollback:** `git revert HEAD`

---

### T10: Update Acceptance Test Imports

**Impact:** 20% → 80% | **Effort:** 10 min | **Depends on:** T8

Update all files in `test/acceptance/`:

| File | Old Import | New Import |
|---|---|---|
| `policy_generation_test.go` | `"github.com/LarsArtmann/template-SECURITY/internal"` | `"github.com/LarsArtmann/template-SECURITY/generator"` |
| `policy_generation_test.go` | `"github.com/LarsArtmann/template-SECURITY/internal/types"` | `"github.com/LarsArtmann/template-SECURITY/core"` |
| `validation_test.go` | `"github.com/LarsArtmann/template-SECURITY/internal"` | `"github.com/LarsArtmann/template-SECURITY/validator"` |
| `test_helpers.go` | (no internal imports) | No change |
| `security_policy_suite_test.go` | (no internal imports) | No change |

Update type references same as T9.

**Verification:**
```bash
go test ./test/acceptance/...
```

**Rollback:** `git revert HEAD`

---

### T11: Full Verification Pass

**Impact:** Remaining | **Effort:** 10 min | **Depends on:** T9, T10

Run the complete verification suite:

```bash
# Per-module builds
cd core && go build ./... && go vet ./... && cd ..
cd validator && go build ./... && go test ./... && go vet ./... && cd ..
cd generator && go build ./... && go test ./... && go vet ./... && cd ..
cd detector && go build ./... && go vet ./... && cd ..

# Root-level (uses go.work)
go build ./...
go test ./...
go vet ./...
go mod tidy && go mod verify
go work sync

# Verify no replace directives in go.mod files (go.work handles it)
grep -r "replace" */go.mod go.mod
```

All checks must pass. If any fail, fix immediately before proceeding.

**Rollback:** Fix forward — this is a verification step, not a change step.

---

### T12: Update `.golangci.yml`

**Impact:** Remaining | **Effort:** 10 min | **Depends on:** T11

1. Review `.golangci.yml` for any path-specific configurations
2. Update lint target paths if needed (e.g., `internal/...` → new module paths)
3. Run `golangci-lint run ./...` to verify

**Verification:**
```bash
golangci-lint run ./...
```

**Rollback:** `git revert HEAD`

---

### T13: Update Documentation

**Impact:** Remaining | **Effort:** 15 min | **Depends on:** T11

Update the following files to reflect new module structure:

1. **README.md** — Update installation instructions if module paths changed
2. **FEATURES.md** — Update file path references
3. **AGENTS.md** (if exists at project level) — Update build/test commands per module
4. **`.github/workflows/security-validation.yml`** — Fix Go version (1.21 → 1.26.2), update build commands for go.work
5. **`go-arch-lint.yml`** — Update architecture lint rules for new module paths

**Verification:**
```bash
go build ./... && go test ./...
```

**Rollback:** `git revert HEAD`

---

### T14: Final Commit

**Impact:** Remaining | **Effort:** 5 min | **Depends on:** T12, T13

Final commit with comprehensive message documenting the modularization:

```
refactor: modularize template-SECURITY into 4 sub-modules

- Extract core (domain types), validator, generator, detector
- Add go.work for multi-module development
- Remove dead domain types (SecurityPolicy, ValidationResult, etc.)
- Move Config/PolicyConfig/TemplateData to core for sharing
- Update all import paths across CLI and acceptance tests
- Each module has its own go.mod with minimal dependencies

Modules:
  core/       — domain types (zero deps)
  validator/  — SECURITY.md validation (go-finding, pipeline)
  generator/  — policy generation (viper, go-finding)
  detector/   — project detection (stdlib only)
  /           — CLI entry point (cobra, color)

Planning docs: docs/modularization/
```

---

## Task Dependency Graph

```
T1 (remove dead types)
 └── T2 (create core)
      ├── T3 (create detector)
      ├── T4 (create validator)
      └── T5 (create generator)
           └── T6 (remove internal/)
                └── T7 (create go.work)
                     └── T8 (update root go.mod)
                          ├── T9 (update CLI imports)
                          └── T10 (update acceptance test imports)
                               └── T11 (full verification)
                                    ├── T12 (update golangci)
                                    └── T13 (update docs)
                                         └── T14 (final commit)
```

## Parallel Execution Opportunities

- T3, T4, T5 can run in **parallel** (all depend only on T2)
- T9, T10 can run in **parallel** (both depend on T8)
- T12, T13 can run in **parallel** (both depend on T11)
