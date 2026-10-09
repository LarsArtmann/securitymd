# Project Split Executive Report: template-SECURITY

> **Never executed — superseded** — no split ever happened; the 2026-10-08 securitymd rebuild (`1c55390`) unified everything into one focused repo (`pkg/policy` + `pkg/provider` + `cmd/securitymd`) and deleted the scripts this plan would have moved. Each proposed project below carries its verdict.

## Current State Analysis

The `template-SECURITY` project currently encompasses a Go-based command-line interface (CLI) for security validation and setup, alongside various automation scripts and extensive documentation. Its structure suggests a monolithic approach where core logic, CLI, and operational scripts reside within a single repository.

## Proposed Project Split

To achieve a higher degree of focus, reusability, and independent development, the `template-SECURITY` project is proposed to be split into the following very focused projects:

### ~~1. `go-security-lib` (Core Library)~~ NOT-DO — core lives in-repo as pkg/policy; never split

- **Purpose**: To provide a highly reusable and modular Go library encapsulating core security analysis, project detection, and common security utility functions. This library will serve as the foundational backend for other security tools.
- **Key Components**:
  - `internal/project_detector.go`
  - `internal/security_tool.go`
  - `internal/security_validator.go`
  - `internal/types/types.go`
  - All associated unit tests.
- **Benefits**: Promotes reusability across multiple projects, enables independent testing and versioning of the core logic, and reduces cognitive load by separating concerns.

### ~~2. `template-security-cli` (CLI Application)~~ NOT-DO — CLI rebuilt in-repo as cmd/securitymd

- **Purpose**: To serve as the primary command-line interface for users to interact with the security tooling. It will consume the `go-security-lib` to perform its functions.
- **Key Components**:
  - `cmd/template-security/main.go`
  - `cmd/template-security/setup.go`
  - `cmd/template-security/status.go`
  - `cmd/template-security/validate.go`
- **Dependencies**: Will depend on `go-security-lib`.
- **Benefits**: Decouples the user interface from the core logic, allowing for independent development cycles and potentially alternative interfaces in the future (e.g., a web UI).

### ~~3. `security-automation-scripts` (Automation Scripts)~~ NOT-DO — scripts deleted at 1c55390, never migrated

- **Purpose**: A dedicated repository for operational scripts related to security, compliance, build processes, and metrics generation. These scripts are often platform-specific or context-dependent.
- **Key Components**:
  - `scripts/build.sh`
  - `scripts/compliance-check.sh`
  - `scripts/generate-metrics.sh`
  - `scripts/security-setup.sh`
  - `scripts/validate-policies.sh`
- **Benefits**: Centralizes operational scripts, making them easier to manage, update, and audit. Reduces clutter in the main code repositories.

### ~~4. `security-docs` (Documentation & Reporting)~~ NOT-DO — docs stay in-repo (docs/ + living root docs)

- **Purpose**: A centralized, dedicated repository for all documentation, status updates, planning documents, and executive reports related to the security initiative.
- **Key Components**:
  - `docs/status/`
  - `docs/planning/`
  - `IMPROVEMENT_PLAN.md`
- **Benefits**: Provides a single source of truth for all security-related documentation, enhances discoverability, and allows for specialized tooling (e.g., static site generators) for documentation.

## Configuration and Workflows

Configuration files (e.g., `.golangci.yml`, `.go-arch-lint.yml`) and GitHub Actions workflows (`.github/workflows`) will reside within the respective code repositories (`go-security-lib` and `template-security-cli`) where they are most relevant, ensuring that each project has its own specific linting, building, and testing configurations.

## Conclusion

This proposed split aims to enhance modularity, foster independent development, improve reusability, and streamline the focus of each component within the overall security initiative. It aligns with best practices for larger projects, promoting a more maintainable and scalable architecture.
