# AGENTS.md — securitymd

SECURITY.md policy linter and generator. Go, single binary, self-registering BuildFlow provider via go-finding/toolsdk.

- **Module:** `github.com/LarsArtmann/securitymd` — GitHub repo renamed 2026-10-09 (`LarsArtmann/securitymd`) AND the local checkout moved to `~/projects/securitymd` (dir name now matches the module); BuildFlow's local `replace` paths were retargeted in the same change (go.mod + tools/go.mod)
- **PUBLIC + tagged `v1.0.0`** (visibility flipped 2026-10-09 by Lars) — module proxy resolves it (`go list -m github.com/LarsArtmann/securitymd@v1.0.0` verified via proxy.golang.org); `GOPRIVATE`/local replaces no longer needed for THIS repo (BuildFlow keeps its replace+flake-input wiring, now resolvable from the public proxy)
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

Quality gate is BuildFlow (`buildflow` / `buildflow --fix`); do not hand-write lint/format scripts. `flake.nix` also provides `packages.<system>.securitymd` (sandboxed binary, `nix run .#securitymd`) and the docs gate:

```bash
nix run .#docs-gate        # docs integrity: un-annotated archives + dangling docs/ refs (also a flake check)
nix flake check            # includes docs-gate + format + the package build
```

## Architecture

```
pkg/policy/            # core: validate.go, generate.go, detect.go, suppress.go, severity.go, project.go, template.md (go:embed)
pkg/provider/          # toolsdk self-registration (Spec with Detect + Repair, contact-email + severity-overrides options)
cmd/securitymd/        # cmdguard/fang CLI: validate, setup, status (typed flag structs, fang-styled execution)
pkg/policy/testdata/   # SARIF/JSON goldens (refresh: go test ./pkg/policy -run TestReport_golden -update) + policy/*.md canonical fixtures
scripts/               # fleet-securitymd-sweep.sh only (legacy scripts fully removed)
test/acceptance/       # Ginkgo BDD specs
```

- **Detect**: WorkingDir from ctx → first existing of `SECURITY.md`, `.github/SECURITY.md`, `docs/SECURITY.md` → missing-file finding (severity error, FixStrategyDirect with rendered AfterCode when git identity derivable) or `Validate()`
- **Repair**: `Generate()` — create-only, **never overwrites a human-written policy** (`--force` regenerates in place with a timestamped `.bak`); skip-with-reason when org/repo not derivable; dry-run aware; `--location root|.github|docs` picks the canonical write target and `validate --location` reorders candidate detection
- Interactive `setup`: prompts for org/repo ONLY when stdin is a REAL terminal and no remote/flags provide identity — never in CI (pinned by the contract test). GOTCHA: a `ModeCharDevice` check classifies `/dev/null` as a terminal and crashed headless runs with a prompt-EOF exit 2; TTY detection must go through `go-isatty`
- `--location root|.github|docs`: canonical write target for `setup`, preferred detection order for `validate` and `status` (candidate-order override; invalid values are operational failures, exit 2)
- **Exit-code contract (README-pinned)**: 0 clean · 1 error-severity findings (`errPolicyFindings`) · 2 operational failure. The gate is THRESHOLD-based (`BySeverityAtLeast(error)`) — `finding.BySeverity` is EQUALITY and silently passes `critical` findings; never use it for gates. Pinned end-to-end by the subprocess contract test
- **Suppression grammar**: `securitymd:ignore(rule[,rule2]) [until YYYY-MM-DD] reason` — reason required, malformed date inert, `until` grants the whole given UTC day (`ExpiresAt` = next midnight); an expired directive attaches NOTHING anywhere (JSON, SARIF, gate) so every consumer agrees the debt is due
- **Provider emits ACTIVE findings only** (`policy.ActiveFindings`): BuildFlow's findings gate counts the provider result and its `filterFindingsAtOrAbove` ignores suppression metadata (verified at source 2026-10-09), so suppression must be stripped at this repo's boundary — the CLI's JSON/SARIF keep the full evidence set
- **Template**: embedded via `go:embed` (`pkg/policy/template.md`). Dogfood invariant is pinned by a test: the rendered template must pass its own validator
- **Contact default**: GitHub advisory link always; email only when explicitly provided (CLI `--email` / provider `contact-email` option). Deliberate — no fabricated `security@domain` addresses. Changing this fleet-wide is a Lars decision, not a code change
- Stable kebab rule IDs (`missing-file`, `missing-header`, `unresolved-template`, …) — suppressions and configs key on them
- Errors wrapped via `finding.NewValidationError`/`NewIOError`/`NewParseError`, never `fmt.Errorf`
- **CLI layer is cmdguard v4** (2026-10-09, satisfies the library-policy cobra-companions rule): `cmdguard.NewCLI[appConfig]` root + typed flag structs (`validateFlags`/`setupFlags`/`statusFlags` with `flag:`/`values:` enum tags), `fang` styling bundled (ANSI-free when piped, verified by the contract test). Exit mapping stays in `main.go` (`exitCodeFor`); cmdguard prints execution errors exactly once — main must NOT print them again. fang capitalizes error headers, so contract assertions pin the domain message (`unknown policy location "bogus"`), not the frame

## Docs layout

- `docs/status/` — post-rebuild status reports; the NEWEST report is the live snapshot, fully-resolved older ones move to `docs/status/archived/` (annotated inline, manifest README inside)
- `docs/planning/` — live planning docs only (fleet-gate announcement draft); executed plans in `docs/planning/archived/`
- `docs/archive/pre-rebuild/` — ALL pre-rebuild history (status/planning/modularization + the five old root reports), annotated with inline ~~verdicts~~ and a manifest README; closed history, do not open without a concrete need
- The docs-gate scans ALL THREE archive homes for un-annotated files (extended 2026-10-09); annotation convention: `~~item~~ done at <hash>` / `routed — <where>` / `decided` / `NOT-DO` / `superseded`, applied with the docs-health `annotate-status-items.py` tooling
- Living docs (README/FEATURES/TODO_LIST/ROADMAP/CHANGELOG/DOMAIN_LANGUAGE) refreshed by the 2026-10-09 evening docs-health AUDIT pass (post-publish state)

## BuildFlow integration (the reason this repo exists)

Wired in `/home/lars/projects/BuildFlow` (its AGENTS.md is authoritative for that repo):

- Blank import in `tools/providers/sdk_imports.go` + guard test `TestSecuritymdProviderRegistered`
- `require` + `replace` in BOTH BuildFlow `go.mod` (indirect) and `tools/go.mod` (direct — separate workspace module)
- BuildFlow does NOT vendor this repo (or any): `go work vendor` was retired 2026-10-08 (BuildFlow gotcha #229) — use the workspace module cache + `GOPRIVATE` + local replaces
- Nix FOD: securitymd is a BuildFlow flake input + preparedSrc dep; the nix build runs `GOWORK=off` so only the ROOT `go.mod` replace matters (dependency replaces are ignored in that mode) — mkPreparedSource strips the local `/home/...` replace and re-adds `./_local_deps/securitymd` from the `deps` map. Verified 2026-10-09: `nix build .` green, vendorHash invariant, result binary lists securitymd (detect+repair). Flake input URL retargeted to the renamed repo; tracks latest per BuildFlow's 2026-10-05 input policy (flake.lock pins rev 81632cf). Full option channel verified E2E via the built binary: `severity-overrides: missing-file=warning` downgrades the finding, `contact-email` lands in the generated policy, detect→fix→re-detect is clean

## Testing patterns

- Unit: `pkg/**/*_test.go`, Gomega (`NewWithT(t)` in stdlib tests; testify fully removed 2026-10-09 per the library-policy ban), table-driven, `t.Parallel()`; SARIF/JSON goldens in `pkg/policy/testdata/golden/` (path entropy normalized; refresh with `-update` when the output contract intentionally changes)
- **Fixture single source**: ALL compliant/flawed policy texts come from `pkg/policy/testdata/policy/*.md` via `policyFixture`/`compliantPolicy`/`flawedPolicy`/`policyWithoutResponseTime` (fixtures_test.go) — unit, golden, suppression, severity, and acceptance tests share the same bytes; never re-inline a policy string
- **CLI contract test** (`cmd/securitymd/contract_test.go`): re-execs the test binary via `TestMain` + `SECURITYMD_CONTRACT_CHILD=1` so `main()` runs for real and the exit code is what a shell sees; 16 scenarios pin README's exit table (findings vs operational vs escape hatches vs setup skips vs suppression expiry vs ANSI-free piped output). New CLI behavior that changes exit codes MUST get a scenario here
- **README drift guard** (`pkg/policy/readme_drift_test.go`): README's two rule tables must equal the code tables (`sectionRules` + `contentRuleSeverities` + `missingFileSeverity`) in IDs AND severities — change rule code and README together or this fails
- Fuzz: `FuzzParseGitRemote` in `pkg/policy/project_fuzz_test.go` — run `go test ./pkg/policy -run '^$' -fuzz FuzzParseGitRemote -fuzztime 45s`
- Provider tests: full detect → repair → verify loop in a real temp git repo (`testhelpers_test.go` isolates git env; `version_cache_test.go` git fixtures need explicit author env vars to stay hermetic in the nix FOD)
- BDD: Ginkgo `Describe/It/By` in `test/acceptance/`
- `cmd/` prompt helpers are unit-tested (`prompt_test.go`); the cmdguard/fang wiring is covered end-to-end by the contract test (subprocess) + provider + acceptance tests

## Conventions

- Emit `finding.Finding` only — no converter glue to other formats (family doctrine)
- Lint config `.golangci.yml`: `pkg/` is fully linted; `cmd/` and `test/` excluded via paths. exhaustruct_v5 ignore-patterns must be fully anchored `^module/path.Struct$` — path-style patterns are silent no-ops. varnamelen `ignore-names: [g]` is deliberate (gomega's conventional receiver; same exemption cmdguard curates)
- `.buildflow.yml` skips `go-structure-linter` (embedded snapshot ignores `.go-structure-linter.yaml`, which documents why `pkg/` must stay public — BuildFlow imports it cross-repo); standalone `go-structure-linter .` remains usable
- CI coverage gate runs `go test -race` (buildflow's race-detector structural check demanded a configured home for it)
- Suppressions in code need a reason comment (3 documented `nolint`/`#nosec` sites exist)
- Deletions via `trash`, never `rm`
- Config file was removed entirely on purpose (old `.template-security.yaml` never worked); knobs are CLI flags + the toolsdk option

## Gotchas

- **Auto-commit daemon bumps mtimes mid-session**: after any "auto-commit" notice (or unexpectedly clean `git status`), re-View a file before editing it — three "file modified since read" refusals in one day from editing from memory
- `go install` is blocked by the bash tool security layer — use `go build -o <path> ./cmd/securitymd` instead
- The stale LSP diagnostic `project_fuzz_test.go:40:17 unparam` is a cached false positive (actual golangci-lint: 0 issues); Go LSP `documentSymbol` is broken — use view+edit, not symbol replace
- **Stale vendorHash does NOT print a got-hash**: nix short-circuits the vendor FOD to the old cached content (its hash still matches the declared one) and the build dies with `go: inconsistent vendoring` instead — `buildflow -s nix-hash-fix` cannot parse this class. Recovery: set `vendorHash` to the all-A fakeHash, run `nix build` once to capture `got: sha256-…`, paste the real hash (done 2026-10-09 for the cmdguard dep tree; AGAIN 2026-10-09 late evening for the samber/lo addition, hash `sha256-G038DkuE…`)
- **BuildFlow result cache replays findings after their cause is deleted** (e.g. the gitignored `vendor/` finding surviving `trash vendor`): bypass with `BUILDFLOW_NO_RESULT_CACHE=1` for the verdict run
- **Ambient `CLICOLOR_FORCE=1` ANSI-leaks piped CLI output**: fang renders through colorprofile, which honors CLICOLOR_FORCE/TTY_FORCE for pipes and ignores NO_COLOR there (colorprofile env.go: `envNoColor && isatty` gate) — so a color-forcing shell turns the contract test red and would pollute provider JSON. Guarded at the boundary: `main.go` `disableForcedColorWhenPiped` strips both vars when stdout is not a TTY (go-isatty), pinned by `TestCLI_piped_output_is_ansi_free` (now runs with CLICOLOR_FORCE=1)
- **go-auto-upgrade's stdlib2lo ↔ lo2stdlib migrators oscillate**: with `samber/lo` required, stdlib2lo wants loop→lo rewrites while lo2stdlib emits standing "in-place-semantics/no-stdlib-equivalent" warnings at the same sites (10 as of 2026-10-09) — warning-level only, gate stays green; fleet-level policy call for Lars (BuildFlow-side)
- `suggest-lo-reduce` on `registerCommands` is a documented non-fix: the loop propagates errors with early return, `lo.Reduce` has no error channel — the rewrite would be worse code
- **lychee step can fail transiently** (network 404/timeout on the known private-repo + template-placeholder links flip the step from warning-report to exit-2 failure between runs); rerun before diagnosing

## Known state (2026-10-09, fleet-policy compliance pass + evening green-gate recovery)

- All tests green, build green, `nix flake check` green (vendorHash `sha256-G038DkuE…` after the samber/lo dep), golangci-lint 0 findings, gofumpt clean, docs-gate exit 0, full `buildflow` gate PASSED exit 0 (warnings only: link-scan 404s, vulnix CVEs on nixpkgs stdenv, the standing go-auto-upgrade migrator noise below), `library-policy` clean (testify gone; fang+cmdguard present; samber/lo allowed — only samber/do v1 is banned)
- **Dependency refresh 2026-10-09 ~21:13** (auto-committed by BuildFlow's go-mod-update): transitive bumps across the charmbracelet stack (ultraviolet, x/ansi), go-output v0.38.4, koanf v2.3.8, x/term v0.47 — this is what invalidated the old vendorHash and exposed the CLICOLOR_FORCE contract-test red under the user's interactive shell
- **CI**: PROVEN GREEN 2026-10-09 evening — billing was fixed and the hardened workflow ran end-to-end (runs 37985112024, 37985651759, 37986553713: validate + govulncheck + docs-gate all success; the earlier 37957194805 refusal was billing, now resolved). `workflow_dispatch` is wired for manual runs. The coverage step runs `-race`
- Publish state: repo renamed + PUBLIC + master pushed + annotated tags `v1.0.0` (frozen at 5c739ec — proxy + sumdb indexed BEFORE the post-tag hardening; never re-tag) + `v1.1.0` (dc80634, the cmdguard/fang CLI migration + ANSI-free piped-output guard) + GitHub Releases for both (v1.1.0 = latest) + description/topics set + dogfood SECURITY.md regenerated. Proxy + sumdb + clean-dir `go get @v1.1.0` verified 2026-10-09; pkg.go.dev pages render on-demand with lag (TODO_LIST post-release section). Release mechanics: commit → CI green → `git tag -a` → `git push origin <tag>` → `gh release create <tag> --notes-file … --latest`
- Post-tag-version discipline: everything after `v1.1.0` is `v1.2.0`+ material — the CHANGELOG's `[Unreleased]` is empty with placeholders; cut the next section from it
- Fleet rollout prep staged: `scripts/fleet-securitymd-sweep.sh` + announcement draft `docs/planning/2026-10-09_fleet-gate-announcement-draft.md` (execution Lars-gated)
- Pre-rebuild docs fully archived + annotated (see docs/archive/pre-rebuild/README.md manifest)
- BuildFlow-side finding (2026-10-09, verified at source): their `filterFindingsAtOrAbove` ignores `Suppression`; this repo's provider strips suppressed findings at the boundary — see TODO_LIST BuildFlow section
