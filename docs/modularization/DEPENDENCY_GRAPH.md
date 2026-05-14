# Dependency Graph — template-SECURITY

_Generated: 2026-05-14 | Phase: Current State Analysis_

---

## Current State: Single Module (Monolith)

```
github.com/LarsArtmann/template-SECURITY (single go.mod)
├── cmd/template-security    → imports: internal, internal/types, cobra, fatih/color, go-finding
├── internal                 → imports: internal/types, go-finding, go-finding/pipeline, viper
├── internal/types           → imports: (stdlib only — time)
└── test/acceptance          → imports: internal, internal/types, go-finding, ginkgo, gomega
```

## Internal Package Dependency Graph

```
cmd/template-security
├── internal (SecurityTool, SecurityValidator, ProjectDetector, Config, PolicyConfig, etc.)
│   └── internal/types (PolicyType, Version, ID, PolicyID, VersionStatus, Contact, etc.)
└── (no direct import of internal/types from cmd — goes through internal)

test/acceptance
├── internal
└── (no direct import of internal/types — goes through internal)
```

### Observations

1. `internal` is a **flat catch-all package** — validation, policy generation, config loading, project detection, and output formatting all live here
2. `internal/types` has **zero importers** for most of its types — only `PolicyType` and `Version` are used
3. `cmd/template-security` directly imports `internal/types` in `setup.go` (for `PolicyConfig.Type` cast)
4. `test/acceptance` imports `internal` heavily but never `internal/types` directly (except via `policy_generation_test.go`)

## External Dependencies (Direct)

| Package | Used By | Purpose | Production/Test |
|---|---|---|---|
| `github.com/spf13/cobra` | `cmd/template-security` | CLI framework | Production |
| `github.com/spf13/viper` | `internal` | Config loading | Production |
| `github.com/fatih/color` | `cmd/template-security`, `internal` | Colored output | Production |
| `github.com/larsartmann/go-finding` | `internal`, `cmd/template-security` | Finding/report model | Production |
| `github.com/larsartmann/go-finding/pipeline` | `internal` | Detector interface | Production |
| `github.com/onsi/ginkgo/v2` | `test/acceptance` | BDD test framework | Test only |
| `github.com/onsi/gomega` | `test/acceptance` | BDD matchers | Test only |
| `github.com/stretchr/testify` | `internal` (unit tests) | Test assertions | Test only |

## Coupling Analysis

### Concern Clusters within `internal`

| Concern | Files | External Deps | Could Be Module? |
|---|---|---|---|
| Validation | `security_validator.go`, `security_validator_test.go` | `go-finding`, `go-finding/pipeline` | Yes — **validator** |
| Policy Generation | `security_tool.go`, `security_tool_test.go` | `viper`, `go-finding` | Yes — **generator** |
| Project Detection | `project_detector.go` | (stdlib only) | Yes — **detector** |
| Domain Types | `types/types.go`, `types/ids.go` | (stdlib only — time) | Yes — **types (core)** |

### Coupling Hotspots

1. **`internal` package is overloaded** — 4 distinct concerns in one package. Each has its own test file and its own external dependency profile.
2. **`viper` only used by `security_tool.go`** — config loading concern isolated to one file, but forces entire `internal` package to carry the dependency.
3. **`go-finding/pipeline` only used by `security_validator.go`** — the `pipeline.Detector` interface is only relevant to validation.
4. **Dead domain types** — `SecurityPolicy`, `ValidationResult`, `Template`, `Contact`, `Project`, `ValidationMessage`, `TemplateVariable` are defined but never used anywhere. This inflates the types package unnecessarily.

### Import Flow (Simplified)

```
┌──────────────────────────┐
│ cmd/template-security    │  (CLI entry point)
│  main.go, setup.go,      │
│  validate.go, status.go  │
└────────┬─────────────────┘
         │ imports
         ▼
┌──────────────────────────┐
│ internal                 │  (catch-all: validation + generation + detection + config)
│  security_validator.go   │──── go-finding, go-finding/pipeline
│  security_tool.go        │──── viper, go-finding
│  project_detector.go     │──── (stdlib only)
└────────┬─────────────────┘
         │ imports
         ▼
┌──────────────────────────┐
│ internal/types           │  (domain types — 90% unused)
│  types.go, ids.go        │──── (time only)
└──────────────────────────┘
```

## Proposed Module Dependency Graph (After Modularization)

```
                    ┌─────────────┐
                    │   /core     │  (domain types: PolicyType, Version, ID)
                    │   types     │  Zero external deps
                    └──────┬──────┘
                           │
              ┌────────────┼────────────────┐
              │            │                │
    ┌─────────▼──────┐  ┌──▼───────────┐  ┌─▼──────────────┐
    │  /validator    │  │  /generator  │  │  /detector     │
    │  validation    │  │  policy gen  │  │  project info  │
    │  go-finding    │  │  viper       │  │  (stdlib)      │
    │  pipeline      │  │  go-finding  │  │                │
    └────────────────┘  └──────────────┘  └────────────────┘
              │            │                │
              └────────────┼────────────────┘
                           │
                ┌──────────▼──────────┐
                │  /cmd               │  (CLI: cobra, color, go-finding)
                │  template-security  │
                └─────────────────────┘
```

### DAG Verification

```
core ← validator  ← cmd
core ← generator  ← cmd
core ← detector   ← cmd
```

- **No cycles** — all arrows point downward
- **Core has zero internal deps** — pure domain types
- **No lateral deps** — validator doesn't import generator, detector doesn't import validator
- **cmd imports all leaf modules** — orchestration happens at the CLI layer
