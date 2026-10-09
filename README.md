# securitymd 🔒

> **Ensure your GitHub repository has a compliant SECURITY.md**

A CLI tool and [BuildFlow](https://github.com/LarsArtmann/BuildFlow) provider that validates and generates `SECURITY.md` files. Detects a missing policy, checks an existing one for the sections GitHub expects, and generates a compliant policy from an embedded template — never overwriting a file you wrote by hand.

## Features

- **Validate** existing SECURITY.md files against GitHub's expected sections
- **Generate** new SECURITY.md files from an embedded canonical template
- **Never overwrites** an existing policy — generation is additive only (`--force` regenerates in place with a timestamped backup)
- **Suppressions and severity overrides** — `securitymd:ignore(rule) [until YYYY-MM-DD] reason` comments and `--set-severity rule=level` keep CI green while findings stay visible; expired suppressions surface the finding again on their own
- **Interactive when human, silent when CI** — `setup` prompts for org/repo only on a real terminal
- **BuildFlow provider** via [go-finding toolsdk](https://github.com/larsartmann/go-finding): detect → repair → verify as a first-class DAG tool
- **go-finding findings** with stable rule IDs, SARIF/JSON output (contract pinned by golden tests), fix strategies

## 🚀 Quick Start

```bash
# Install (works once the repo is PUBLIC — tag v1.0.0 is already cut;
# until the visibility flip, build from source: go build ./cmd/securitymd)
go install github.com/LarsArtmann/securitymd/cmd/securitymd@latest

# Validate the current repository's policy
securitymd validate

# Generate a SECURITY.md if none exists
securitymd setup

# Generate with explicit identity and contact email
securitymd setup --organization AcmeCorp --repository widget --email security@acme.com

# Canonical location (.github/ or docs/) instead of the repo root
securitymd setup --location docs
securitymd validate --location docs
securitymd status --location docs

# Refresh a policy you own (writes SECURITY.md.<timestamp>.bak first)
securitymd setup --force

# Downgrade the adoption blocker for fleets with many unmanaged repos
securitymd validate --set-severity missing-file=warning

# Show compliance status
securitymd status
```

With BuildFlow (blank import already wired in `tools/providers/sdk_imports.go`):

```bash
buildflow -s securitymd              # detect: missing or non-compliant policy
buildflow -s securitymd --fix        # repair: generate the policy when absent
```

Per-repo contact override:

```yaml
# .buildflow.yml
tool_options:
  securitymd:
    contact-email: security@example.com
    severity-overrides: missing-file=warning
```

## ✅ What We Validate

Discovery order: `SECURITY.md`, `.github/SECURITY.md`, `docs/SECURITY.md`.

### Required sections (error)

| Rule                  | Check                                                            |
| --------------------- | ---------------------------------------------------------------- |
| `missing-file`        | No policy found in any candidate location (fix: generate one)    |
| `missing-header`      | Missing `# Security Policy` header                               |
| `missing-reporting`   | Missing "Reporting a Vulnerability" (or equivalent) section      |
| `missing-versions`    | Missing "Supported Versions" section                             |
| `missing-practices`   | Missing "Security Practices" section                             |
| `missing-contact`     | No contact channel: email, GitHub advisory link, or security.txt |
| `unresolved-template` | Leftover `{{.Variable}}` template placeholders (line-precise)    |
| `no-content`          | Placeholder-only content, nothing substantive                    |

### Quality checks (warning)

| Rule                    | Check                                   |
| ----------------------- | --------------------------------------- |
| `missing-response-time` | No response-time commitment for reports |
| `too-short`             | Under 20 lines                          |
| `no-version-info`       | No version information anywhere         |

Exit codes: `0` clean · `1` error-severity findings · `2` operational failure. In a GitHub Actions workflow, gate on the distinction:

```yaml
- name: Validate security policy
  run: securitymd validate
# exit 1 (policy findings) fails the step; exit 2 (operational) surfaces as a tool error
```

Two severity flags, two different jobs: `--severity LEVEL` sets the minimum severity REPORTED in SARIF output, while `--set-severity rule=level` overrides a rule's severity itself (affecting both output and the exit gate).

Colored output disables itself when stdout is not a terminal, when `NO_COLOR` is set, or when `TERM=dumb` — CI logs stay clean.

False positives? Add a suppression comment in the policy (findings stay visible, exit turns neutral):

```markdown
<!-- securitymd:ignore(no-version-info) bootstrap stage, filled next sprint -->
```

Scope a suppression to a calendar day so deferred debt resurfaces on its own (the finding is silenced through the end of the given UTC day, then active again):

```markdown
<!-- securitymd:ignore(no-version-info) until 2026-11-01 bootstrap stage, filled next sprint -->
```

## Generation

The embedded template renders from your git remote (organization/repository), the latest git tag (best effort), and an optional contact email. Without an email the policy points reporters at GitHub private vulnerability reporting (`https://github.com/ORG/REPO/security/advisories/new`) — no fabricated addresses.

Writes are atomic and idempotent (`go-atomic-write`). If any candidate policy already exists, `setup` refuses and says so.

## Architecture

```
cmd/securitymd/        # CLI (cobra): validate, setup, status
pkg/policy/            # Core: detection candidates, content rules, embedded
                       # template, git identity, atomic generation
pkg/provider/          # toolsdk.Spec: BuildFlow self-registration
test/acceptance/       # Ginkgo BDD acceptance tests
```

- Findings: `github.com/larsartmann/go-finding` (no parallel issue model)
- Provider contract: `github.com/larsartmann/go-finding/toolsdk`
- Helpers: `linter-autoconfigure-sdk` (`FirstExisting`, `WorkingDir`)
- Atomic writes: `go-atomic-write`

## Development

```bash
go build ./...          # Build
go test ./...           # Unit + BDD acceptance tests
golangci-lint run       # 90+ linters, zero issues expected
```
