# Features — securitymd

_Audit: 2026-10-09, post-rebuild (module `github.com/LarsArtmann/securitymd`). Evidence verified against code this date (docs-health pass)._

## Core

| # | Feature                      | Status           | Notes                                                                                                                 | Evidence                                                  |
| - | ---------------------------- | ---------------- | --------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------- |
| 1 | Policy validation (11 rules) | FULLY_FUNCTIONAL | 10 content rules (`pkg/policy/validate.go`) + `missing-file` (`pkg/policy/detect.go`); stable kebab rule IDs          | `pkg/policy/validate.go`, `detect.go`, `validate_test.go` |
| 2 | Policy generation            | FULLY_FUNCTIONAL | Embedded template, never-overwrite, dry-run, atomic idempotent write                                                  | `pkg/policy/generate.go`, `generate_test.go`              |
| 3 | Detection                    | FULLY_FUNCTIONAL | Candidate locations, missing-file finding with rendered AfterCode                                                     | `pkg/policy/detect.go`                                    |
| 4 | Git identity detection       | FULLY_FUNCTIONAL | https/ssh/gitlab-nested remote parsing, latest-tag best effort                                                        | `pkg/policy/project.go`, `project_test.go`                |
| 5 | BuildFlow provider           | FULLY_FUNCTIONAL | toolsdk Spec, Detect + Repair, `contact-email` option; verified live in BuildFlow (detect → repair → re-detect clean) | `pkg/provider/provider.go`, `provider_test.go`            |
| 6 | CLI (validate/setup/status)  | FULLY_FUNCTIONAL | Exit 0/1 by error findings; covered via acceptance tests                                                              | `cmd/securitymd/`                                         |
| 7 | Dogfood invariant            | FULLY_FUNCTIONAL | Generated template must pass its own validator — pinned by test                                                       | `pkg/policy/validate_test.go`                             |
| 8 | BDD acceptance tests         | FULLY_FUNCTIONAL | 7 Ginkgo specs                                                                                                        | `test/acceptance/validation_test.go`                      |

**Mutation proof (2026-10-09):** the rule table is proven non-vacuous. Sabotaging `missing-contact`'s patterns (7 tests red across unit + golden + acceptance), neutering the `too-short` threshold 20→0 (1 test red), and shifting the `unresolved-template` line off by one (2 tests red: line-precision unit test + golden) each failed the suite; reverting restored green. The SARIF/JSON goldens act as an independent net on every mutation.

## Partial / gaps

| #  | Feature                          | Status               | Gap                                                                                                          |
| -- | -------------------------------- | -------------------- | ------------------------------------------------------------------------------------------------------------ |
| 9  | Multi-format output (JSON/SARIF) | FULLY_FUNCTIONAL     | none — output contract pinned by golden-file tests (JSON + SARIF 2.1.0; clean, mixed, missing-file fixtures) |
| 10 | Trigger breadth                  | FULLY_FUNCTIONAL     | none — README.md activates docs-only repos (coverage test pins the trigger)                                  |
| 11 | CLI unit tests                   | PARTIALLY_FUNCTIONAL | `cmd/` untested directly; behavior covered via acceptance/provider tests                                     |
| 12 | Nix package output               | FULLY_FUNCTIONAL     | none — `packages.<system>.securitymd` builds sandboxed (Go 1.27 pinned, git fixture tests pass in FOD)       |
| 13 | Published module                 | MISSING              | Not on GitHub yet — `go install …@latest` impossible (publish checklist in TODO_LIST)                        |
