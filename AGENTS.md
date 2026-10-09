# AGENTS.md — securitymd

SECURITY.md policy linter and generator. Go, single binary, self-registering BuildFlow provider via go-finding/toolsdk.

- **Module:** `github.com/LarsArtmann/securitymd` — **directory is still `template-SECURITY`** (rename to `securitymd` pending; do not be confused, the code is new)
- **Not yet published on GitHub** — `go install …@latest` impossible until rename+push+tag; BuildFlow consumes it via local replace + flake input
- Single detector core (`pkg/policy`) behind both the CLI and the toolsdk provider — never two implementations

## Commands

```bash
GOWORK=off go build ./...          # build (GOWORK=off: parent fleet workspace interferes)
GOWORK=off go test ./...           # all tests (unit + provider + Ginkgo BDD)
GOWORK=off go mod tidy             # always GOWORK=off for module commands
gofumpt -w ./pkg ./cmd ./test      # formatter is gofumpt, not gofmt
golangci-lint run --timeout 5m     # 90+ linters; must stay at 0 issues
go build -o bin/securitymd ./cmd/securitymd   # binary
```

Quality gate is BuildFlow (`buildflow` / `buildflow --fix`); do not hand-write lint/format scripts. `flake.nix` provides devshells only (no package output yet).

## Architecture

```
pkg/policy/            # core: validate.go, generate.go, detect.go, project.go, template.md (go:embed)
pkg/provider/          # toolsdk self-registration (Spec with Detect + Repair, contact-email option)
cmd/securitymd/        # cobra CLI: validate, setup, status
test/acceptance/       # Ginkgo BDD specs
```

- **Detect**: WorkingDir from ctx → first existing of `SECURITY.md`, `.github/SECURITY.md`, `docs/SECURITY.md` → missing-file finding (severity error, FixStrategyDirect with rendered AfterCode when git identity derivable) or `Validate()`
- **Repair**: `Generate()` — create-only, **never overwrites a human-written policy**; skip-with-reason when org/repo not derivable from `git remote origin`; dry-run aware
- **Template**: embedded via `go:embed` (`pkg/policy/template.md`). Dogfood invariant is pinned by a test: the rendered template must pass its own validator
- **Contact default**: GitHub advisory link always; email only when explicitly provided (CLI `--email` / provider `contact-email` option). Deliberate — no fabricated `security@domain` addresses. Changing this fleet-wide is a Lars decision, not a code change
- Stable kebab rule IDs (`missing-file`, `missing-header`, `unresolved-template`, …) — suppressions and configs key on them
- Errors wrapped via `finding.NewValidationError`/`NewIOError`/`NewParseError`, never `fmt.Errorf`

## Docs layout

- `docs/status/` — post-rebuild status reports ONLY (the two 2026-10 reports)
- `docs/archive/pre-rebuild/` — ALL pre-rebuild history (status/planning/modularization + the five old root reports), annotated with inline ~~verdicts~~ and a manifest README; closed history, do not open without a concrete need
- Living docs (README/FEATURES/TODO_LIST/ROADMAP/CHANGELOG/DOMAIN_LANGUAGE) refreshed by the 2026-10-09 docs-health pass

## BuildFlow integration (the reason this repo exists)

Wired in `/home/lars/projects/BuildFlow` (its AGENTS.md is authoritative for that repo):

- Blank import in `tools/providers/sdk_imports.go` + guard test `TestSecuritymdProviderRegistered`
- `require` + `replace` in BOTH BuildFlow `go.mod` (indirect) and `tools/go.mod` (direct — separate workspace module)
- Vendoring in BuildFlow: `go work vendor` (NOT `go mod vendor` — fails in workspace mode)
- Nix FOD: securitymd is a BuildFlow flake input + preparedSrc dep (local `replace` directives are IGNORED by the nix prepared-source build — this cost a failed FOD cycle; see nix-private-go-repos skill). After publishing, drop replaces + flake input and use the tag

## Testing patterns

- Unit: `pkg/**/*_test.go`, testify, table-driven, `t.Parallel()`
- Provider tests: full detect → repair → verify loop in a real temp git repo (`testhelpers_test.go` isolates git env)
- BDD: Ginkgo `Describe/It/By` in `test/acceptance/`
- `cmd/` has no tests (excluded from lint paths too); behavior is covered via provider + acceptance tests

## Conventions

- Emit `finding.Finding` only — no converter glue to other formats (family doctrine)
- Lint config `.golangci.yml`: `pkg/` is fully linted; `cmd/` and `test/` excluded via paths. exhaustruct_v5 ignore-patterns must be fully anchored `^module/path.Struct$` — path-style patterns are silent no-ops
- Suppressions in code need a reason comment (3 documented `nolint`/`#nosec` sites exist)
- Deletions via `trash`, never `rm`
- Config file was removed entirely on purpose (old `.template-security.yaml` never worked); knobs are CLI flags + the toolsdk option

## Known state (2026-10-09)

- All tests green, build green, golangci-lint 0 issues, `nix flake check` green (re-verified 2026-10-09 docs-health pass), BuildFlow integration verified live (detect+repair+re-detect clean)
- Pre-rebuild docs fully archived + annotated (see docs/archive/pre-rebuild/README.md manifest)
- Publishing checklist (needs Lars): rename GitHub repo → push → tag v1.0.0 → drop BuildFlow replaces/flake input → re-vendor → `nix run .#update-vendor-hash`
