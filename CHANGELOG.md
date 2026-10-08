# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [Unreleased]

### Added

- `pkg/provider`: self-registering BuildFlow provider (go-finding/toolsdk Spec) with Detect + Repair, `contact-email` tool option, dry-run support — consumed by BuildFlow via blank import (`-s securitymd`, `--fix` generates missing policies)
- `missing-file` rule: severity error, direct fix strategy with rendered policy preview when org/repo are derivable from the git remote
- Embedded canonical template (`go:embed`) — the binary no longer depends on a `templates/` directory relative to CWD
- Never-overwrite generation, dry-run, atomic idempotent writes, git-remote identity detection (https/ssh/gitlab-nested), latest-tag best effort
- Multi-signal contact detection (GitHub advisory link and `security.txt` count as contact), line-precise `unresolved-template` findings
- Stable kebab-case rule IDs for suppressions and configurations
- Provider contract tests including a full detect → repair → verify loop in a real temp git repo, plus a dogfood test pinning that the rendered template passes its own validator

### Changed

- **Module renamed**: `github.com/LarsArtmann/template-SECURITY` → `github.com/LarsArtmann/securitymd`; binary is now `securitymd`
- Rebuilt on go-finding + toolsdk + linter-autoconfigure-sdk + go-atomic-write; dropped viper (≈15 indirect deps) and go-finding/pipeline
- `validate` now exits non-zero when error-severity findings remain (previously exited 0 for a missing SECURITY.md)

### Removed

- **Breaking**: `.template-security.yaml` config file support (the config never worked — snake_case keys vs camelCase tags); knobs are now CLI flags and the `contact-email` provider option
- **Breaking**: policy-type system and domain types (`internal/types`), `setup` interactive template picker
- `internal/` packages, `cmd/template-security/`, `templates/`, and the legacy `scripts/` (~1,300 lines of shell duplicating Go functionality)

### Fixed

- Template lookup no longer breaks when the binary runs outside the source tree
- CI workflow: Go version now follows `go.mod` (was pinned to 1.21); removed non-existent `make` targets; paths updated to the new layout

## [0.1.0] - 2026-01-01

### Added

- Initial release
