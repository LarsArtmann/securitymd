# securitymd 🔒

> **Ensure your GitHub repository has a compliant SECURITY.md**

A CLI tool and [BuildFlow](https://github.com/LarsArtmann/BuildFlow) provider that validates and generates `SECURITY.md` files. Detects a missing policy, checks an existing one for the sections GitHub expects, and generates a compliant policy from an embedded template — never overwriting a file you wrote by hand.

## Features

- **Validate** existing SECURITY.md files against GitHub's expected sections
- **Generate** new SECURITY.md files from an embedded canonical template
- **Never overwrites** an existing policy — generation is additive only
- **BuildFlow provider** via [go-finding toolsdk](https://github.com/larsartmann/go-finding): detect → repair → verify as a first-class DAG tool
- **go-finding findings** with stable rule IDs, SARIF/JSON output, fix strategies

## 🚀 Quick Start

```bash
# Install (once the repo is renamed/published under LarsArtmann/securitymd)
go install github.com/LarsArtmann/securitymd/cmd/securitymd@latest

# Validate the current repository's policy
securitymd validate

# Generate a SECURITY.md if none exists
securitymd setup

# Generate with explicit identity and contact email
securitymd setup --organization AcmeCorp --repository widget --email security@acme.com

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

Exit codes: `0` clean · `1` error-severity findings · `2` operational failure.

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
