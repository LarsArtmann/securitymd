# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [Unreleased]

### Added

- Suppression comments (`securitymd:ignore(rule) reason`): matched findings stay visible but exit-neutral; unknown rules or missing reasons are errors
- Per-rule severity overrides: `validate --set-severity rule=level` and the provider tool option `severity-overrides` (`missing-file=warning` is the incremental-fleet-adoption unblocker)
- `setup --force`: regenerates an existing policy in place after writing a timestamped `SECURITY.md.<timestamp>.bak` (refuse-by-default unchanged without it)
- `--location root|.github|docs`: canonical write target for `setup` and preferred detection order for `validate` (candidate-order override)
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

### Changed

- **Module renamed**: `github.com/LarsArtmann/template-SECURITY` → `github.com/LarsArtmann/securitymd`; binary is now `securitymd`
- Rebuilt on go-finding + toolsdk + linter-autoconfigure-sdk + go-atomic-write; dropped viper (≈15 indirect deps) and go-finding/pipeline
- `validate` now exits non-zero when error-severity findings remain (previously exited 0 for a missing SECURITY.md)

### Removed

- **Breaking**: `.template-security.yaml` config file support (the config never worked — snake_case keys vs camelCase tags); knobs are now CLI flags and the `contact-email` provider option
- **Breaking**: policy-type system and domain types (`internal/types`), `setup` interactive template picker
- `internal/` packages, `cmd/template-security/`, `templates/`, and the legacy `scripts/` (~1,300 lines of shell duplicating Go functionality)

### Fixed

- `parseGitRemote` no longer invents identities from unsupported remote forms (`git://`, `file://`, local paths, `host:org` without repo) — these now degrade to an honest incomplete identity instead of a plausible-looking wrong advisory link
- `setup --location .github|docs` creates the target directory before the atomic write (previously only the root target existed)
- Template lookup no longer breaks when the binary runs outside the source tree
- CI workflow: Go version now follows `go.mod` (was pinned to 1.21); removed non-existent `make` targets; paths updated to the new layout

## [0.1.0] - 2026-01-01

### Added

- Initial release
