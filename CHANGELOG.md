# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [Unreleased]

### Added

- Nothing yet.

### Changed

- Nothing yet.

### Fixed

- Nothing yet.

## [1.1.0] - 2026-10-09

### Added

- Docs-health AUDIT pass (2026-10-09 evening): post-rebuild status/planning reports fully annotated (599 inline verdicts citing carrying commits) and archived under per-directory `archived/` homes with bulk-archive manifests; the docs-gate now scans all three archive homes (`docs/archive/pre-rebuild/`, `docs/status/archived/`, `docs/planning/archived/`); FEATURES row 11 table repaired (unescaped pipes had truncated the row); TODO_LIST synced with the local directory rename and BuildFlow's vendoring retirement; ROADMAP gained the hardening next-list harvest (fuzz severity parser, property fuzz, CLI JSON golden, coverage ratchet, golangci CI leg, sweep `--json`, gosec overlap audit, resolver-column refactor, re-survey cadence)

### Changed

- CLI framework migrated from cobra to cmdguard v4 (fleet-standard CLI plumbing: typed flag structs, fang rendering, signal handling, panic recovery); commands, flags, and the README-pinned exit-code contract are unchanged — the 16-scenario subprocess contract test pins the surface

### Fixed

- Piped output is now ANSI-free even when the ambient environment forces color (`CLICOLOR_FORCE`/`TTY_FORCE`): the guard unsets the forcing variables when stdout is not a TTY, so machine-readable output consumed by CI and the BuildFlow provider can no longer leak escape sequences (pinned by a hostile-env contract scenario)

## [1.0.0] - 2026-10-09

### Added

- Suppression comments (`securitymd:ignore(rule) reason`): matched findings stay visible but exit-neutral; unknown rules or missing reasons are errors
- Per-rule severity overrides: `validate --set-severity rule=level` and the provider tool option `severity-overrides` (`missing-file=warning` is the incremental-fleet-adoption unblocker)
- `setup --force`: regenerates an existing policy in place after writing a timestamped `SECURITY.md.<timestamp>.bak` (refuse-by-default unchanged without it)
- `--location root|.github|docs`: canonical write target for `setup` and preferred detection order for `validate` and `status` (candidate-order override)
- Interactive `setup`: prompts for org/repo when no git remote is derivable and stdin is a TTY; never prompts in CI or piped contexts
- `README.md` as a provider trigger: docs-only repositories (no dependency manifests) now activate the provider
- SARIF + JSON golden-file tests pinning the exact CLI output contract (`pkg/policy/golden_test.go`, refresh with `-update`); mutation/discrimination proof that the rule table is non-vacuous (3 sabotage targets, each caught by the suite)
- Native fuzz target for remote-URL parsing (`pkg/policy/project_fuzz_test.go`; 10.9M execs green) with parser hardening: scheme allowlist (https/ssh/git@), trailing-slash and CRLF tolerance, clean-coordinates invariant, minimum host/org/repo shape
- Version-cell caching: the latest git tag is memoized per directory per process (one `git describe` per run)
- flake `packages.<system>.securitymd` binary output (sandboxed Go 1.27 build, pinned vendorHash, git fixture tests pass in the FOD)
- `apps.docs-gate` / `checks.docs-gate`: one-command docs-integrity gate (un-annotated archive detection + dangling `docs/` reference alarm), wired into `nix flake check`
- `pkg/provider`: self-registering BuildFlow provider (go-finding/toolsdk Spec) with Detect + Repair, `contact-email` tool option, dry-run support — consumed by BuildFlow via blank import (`-s securitymd`, `--fix` generates missing policies)
- `missing-file` rule: severity error, direct fix strategy with rendered policy preview when org/repo are derivable from the git remote
- Embedded canonical template (`go:embed`) — the binary no longer depends on a `templates/` directory relative to CWD
- Never-overwrite generation, dry-run, atomic idempotent writes, git-remote identity detection (https/ssh/gitlab-nested), latest-tag best effort
- Multi-signal contact detection (GitHub advisory link and `security.txt` count as contact), line-precise `unresolved-template` findings
- Stable kebab-case rule IDs for suppressions and configurations
- Provider contract tests including a full detect → repair → verify loop in a real temp git repo, plus a dogfood test pinning that the rendered template passes its own validator
- Docs: all pre-rebuild status/planning/modularization reports annotated with inline resolutions and archived under `docs/archive/pre-rebuild/` (manifest inside); `docs/status/` now holds only post-rebuild reports
- CLI subprocess contract test: 12 scenarios re-exec the real binary and pin README's exit-code table (clean/findings/operational/suppression escape hatches/severity escalation)
- README drift guard: the documented rule tables must match the code's rule and severity tables (`sectionRules`, `contentRuleSeverities`, `missingFileSeverity`) — rule drift now fails the build instead of the docs
- Single-source policy fixtures: `pkg/policy/testdata/policy/{compliant,flawed}.md` replace four hand-rolled string copies across unit, golden, suppression, severity, and acceptance tests
- Suppression × severity-override combination tests: escalation of a suppressed finding keeps the evidence, stays exit-neutral, and (unsuppressed) trips the gate
- CI hardening in `security-validation.yml`: 80% coverage floor (actual 87.6%), 30s `FuzzParseGitRemote` smoke, govulncheck (action v1.1.0, SHA-pinned), and the nix docs-gate as its own job (install-nix-action v31.11.1, SHA-pinned)
- Fleet rollout prep: `scripts/fleet-securitymd-sweep.sh` (preview/`--fix` sweep across repo lists) and the announcement draft `docs/planning/2026-10-09_fleet-gate-announcement-draft.md` (execution Lars-gated)
- Suppressed-finding goldens: JSON keeps the suppression evidence, SARIF export drops the finding — the asymmetry is now pinned byte-for-byte
- Suppression expiry: `securitymd:ignore(rule) until YYYY-MM-DD reason` scopes a suppression to a calendar day (suppressed through the end of the given UTC day, active again the next); a malformed date is inert; expired directives attach nothing anywhere so every consumer agrees the debt is due
- All-suppressed golden case: the escape hatch at full extent pins JSON evidence (reason + expiry) and an empty SARIF result set
- Drift-guard red run: the README guard demonstrably fails on a renamed rule ID (end-to-end against a mutated README) and on a severity flip
- Provider emits ACTIVE findings only: suppressed findings are stripped at the toolsdk boundary because BuildFlow's findings gate counts the provider result and does not honor suppression metadata
- CLI contract table grown to 16 scenarios: setup refusal and honest identity skip, suppression expiry end-to-end (unexpired neutral, expired trips), and ANSI-free piped output

### Changed

- **Module renamed**: `github.com/LarsArtmann/template-SECURITY` → `github.com/LarsArtmann/securitymd`; binary is now `securitymd`
- Rebuilt on go-finding + toolsdk + linter-autoconfigure-sdk + go-atomic-write; dropped viper (≈15 indirect deps) and go-finding/pipeline
- `validate` now exits non-zero when error-severity findings remain (previously exited 0 for a missing SECURITY.md)

### Fixed

- **Exit-code contract now matches the docs**: findings exit 1, operational failures (bad flags, IO errors, refusals) exit 2 — previously every error exited 1, making README's `2` unreachable
- **Critical-severity findings now trip the exit gate**: the gate used go-finding's `BySeverity` (equality), so a finding escalated to `critical` via `--set-severity` silently passed CI; the gate is now threshold-based (`BySeverityAtLeast(error)`)
- `status` renders escalated (`critical`) findings as errors instead of warnings, counts suppressed findings separately instead of as errors, and marks them in the listing
- **`setup` no longer crashes in CI**: the interactive-prompt gate used a `ModeCharDevice` check, which classifies `/dev/null` as a terminal — headless runs (CI, scripts, subprocesses) died with a prompt-EOF exit 2 instead of skipping honestly; TTY detection now uses `go-isatty`
- `setup` no longer prompts for identity when a policy already exists — the prompt was a guaranteed no-op (generation refuses by design), pure interrogation noise
- `parseGitRemote` no longer invents identities from unsupported remote forms (`git://`, `file://`, local paths, `host:org` without repo) — these now degrade to an honest incomplete identity instead of a plausible-looking wrong advisory link
- `setup --location .github|docs` creates the target directory before the atomic write (previously only the root target existed)
- Template lookup no longer breaks when the binary runs outside the source tree
- CI workflow: Go version now follows `go.mod` (was pinned to 1.21); removed non-existent `make` targets; paths updated to the new layout

### Removed

- **Breaking**: `.template-security.yaml` config file support (the config never worked — snake_case keys vs camelCase tags); knobs are now CLI flags and the `contact-email` provider option
- **Breaking**: policy-type system and domain types (`internal/types`), `setup` interactive template picker
- `internal/` packages, `cmd/template-security/`, `templates/`, and the legacy `scripts/` (~1,300 lines of shell duplicating Go functionality)

## [0.1.0] - 2026-01-01

### Added

- Initial release
