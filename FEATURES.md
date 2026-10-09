# Features — securitymd

_Audit: 2026-10-08, post-rebuild (module `github.com/LarsArtmann/securitymd`). Evidence verified against code this date._

## Core

| # | Feature                      | Status           | Notes                                                                                                                 | Evidence                                       |
| - | ---------------------------- | ---------------- | --------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------- |
| 1 | Policy validation (10 rules) | FULLY_FUNCTIONAL | Section rules + content-quality rules, stable kebab rule IDs                                                          | `pkg/policy/validate.go`, `validate_test.go`   |
| 2 | Policy generation            | FULLY_FUNCTIONAL | Embedded template, never-overwrite, dry-run, atomic idempotent write                                                  | `pkg/policy/generate.go`, `generate_test.go`   |
| 3 | Detection                    | FULLY_FUNCTIONAL | Candidate locations, missing-file finding with rendered AfterCode                                                     | `pkg/policy/detect.go`                         |
| 4 | Git identity detection       | FULLY_FUNCTIONAL | https/ssh/gitlab-nested remote parsing, latest-tag best effort                                                        | `pkg/policy/project.go`, `project_test.go`     |
| 5 | BuildFlow provider           | FULLY_FUNCTIONAL | toolsdk Spec, Detect + Repair, `contact-email` option; verified live in BuildFlow (detect → repair → re-detect clean) | `pkg/provider/provider.go`, `provider_test.go` |
| 6 | CLI (validate/setup/status)  | FULLY_FUNCTIONAL | Exit 0/1 by error findings; covered via acceptance tests                                                              | `cmd/securitymd/`                              |
| 7 | Dogfood invariant            | FULLY_FUNCTIONAL | Generated template must pass its own validator — pinned by test                                                       | `pkg/policy/validate_test.go`                  |
| 8 | BDD acceptance tests         | FULLY_FUNCTIONAL | 7 Ginkgo specs                                                                                                        | `test/acceptance/validation_test.go`           |

## Partial / gaps

| #  | Feature                          | Status               | Gap                                                                                   |
| -- | -------------------------------- | -------------------- | ------------------------------------------------------------------------------------- |
| 9  | Multi-format output (JSON/SARIF) | PARTIALLY_FUNCTIONAL | Implemented in CLI; no golden-file test for SARIF shape (TODO_LIST)                   |
| 10 | Trigger breadth                  | PARTIALLY_FUNCTIONAL | Docs-only repos without dependency manifests don't activate the provider              |
| 11 | CLI unit tests                   | PARTIALLY_FUNCTIONAL | `cmd/` untested directly; behavior covered via acceptance/provider tests              |
| 12 | Nix package output               | MISSING              | flake exposes devshells only (no `packages.<system>.securitymd` build)                |
| 13 | Published module                 | MISSING              | Not on GitHub yet — `go install …@latest` impossible (publish checklist in TODO_LIST) |
