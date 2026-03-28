# BDD Tests Review

**Project:** template-SECURITY
**Date:** 2026-03-28
**Status:** Critical - No BDD tests exist

---

## Executive Summary

**Verdict: Insufficient BDD coverage. The project lacks behavioral tests entirely.**

| Criterion                 | Status       | Score    |
| ------------------------- | ------------ | -------- |
| Ginkgo Framework          | Not used     | 0/10     |
| End-User Perspective      | Missing      | 0/10     |
| Behavioral Scenarios      | Missing      | 0/10     |
| Given-When-Then Structure | Missing      | 0/10     |
| Business Value Focus      | Missing      | 0/10     |
| **Overall**               | **Critical** | **0/50** |

---

## Current State

### Test Files Found

| File                                  | Type                 | Lines |
| ------------------------------------- | -------------------- | ----- |
| `internal/security_tool_test.go`      | Unit tests (testify) | 172   |
| `internal/security_validator_test.go` | Unit tests (testify) | 270   |

### Framework Usage

```go
// Current approach - NOT BDD
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestSecurityTool_GeneratePolicy(t *testing.T) {
    tests := []struct{...}{...}
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {...})
    }
}
```

### What's Wrong With Current Tests

| Issue             | Current State                     | BDD Requirement                                          |
| ----------------- | --------------------------------- | -------------------------------------------------------- |
| **Naming**        | `TestSecurityTool_GeneratePolicy` | `Describe("Generating a security policy", func() {...})` |
| **Perspective**   | Developer/implementation          | End-user/business                                        |
| **Structure**     | Arrange-Act-Assert                | Given-When-Then                                          |
| **Readability**   | Technical jargon                  | Natural language                                         |
| **Documentation** | None                              | Tests AS documentation                                   |

---

## Gap Analysis

### Missing User Stories

1. **As a developer**, I want to generate a SECURITY.md file so that my project has clear security guidelines
2. **As a security auditor**, I want to validate existing SECURITY.md files so that I can ensure compliance
3. **As a project maintainer**, I want to customize security policies with my organization's details
4. **As a CI/CD pipeline**, I want to fail builds when security documentation is incomplete

### Missing Behavioral Scenarios

#### Policy Generation

```
GIVEN a project without security documentation
WHEN I run the security policy generator
THEN a complete SECURITY.md file is created
AND it contains my organization's contact information
AND it includes supported versions
AND it specifies vulnerability reporting procedures
```

#### Validation

```
GIVEN an existing SECURITY.md file
WHEN I validate it
THEN I receive clear feedback about missing sections
AND warnings about best practices
AND actionable recommendations
```

#### Error Handling

```
GIVEN a template with unresolved variables
WHEN I generate a policy
THEN I receive a clear error message
AND the error identifies the unresolved variable
AND no incomplete file is written
```

---

## Recommendations

### 1. Adopt Ginkgo + Gomega

**Add to go.mod:**

```go
require (
    github.com/onsi/ginkgo/v2 v2.x.x
    github.com/onsi/gomega v1.x.x
)
```

**Why Ginkgo:**

- Industry standard for Go BDD testing
- Expressive DSL: `Describe`, `Context`, `When`, `It`
- Built-in table-driven specs with `DescribeTable`
- Excellent async support with `Eventually`
- Integrates with standard `go test`

### 2. Restructure Test Files

**Current:**

```
internal/
  security_tool_test.go      # Unit tests
  security_validator_test.go # Unit tests
```

**Recommended:**

```
internal/
  security_tool_test.go           # Unit tests (keep for coverage)
test/
  acceptance/
    security_policy_suite_test.go # BDD acceptance tests
    validation_suite_test.go      # BDD validation tests
  integration/
    cli_suite_test.go             # CLI integration tests
```

### 3. Write User-Centric Test Scenarios

**Pattern:**

```go
var _ = Describe("Security Policy Generation", Label("acceptance"), func() {
    var (
        tool     *internal.SecurityTool
        config   internal.PolicyConfig
        outputDir string
    )

    BeforeEach(func() {
        tool = internal.NewSecurityTool()
        outputDir = GinkgoT().TempDir()
        config = internal.PolicyConfig{
            Type:         internal.PolicyTypeGitHub,
            Organization: "Acme Corp",
            ContactEmail: "security@acme.com",
            OutputDir:    outputDir,
        }
    })

    When("generating a complete security policy", func() {
        BeforeEach(func() {
            // Setup: Create template
        })

        It("creates a SECURITY.md file in the output directory", func() {
            Expect(tool.GeneratePolicy(nil, config)).To(Succeed())

            Expect(outputDir + "/SECURITY.md").To(BeAnExistingFile())
        })

        It("includes the organization name", func() {
            Expect(tool.GeneratePolicy(nil, config)).To(Succeed())

            content, err := os.ReadFile(outputDir + "/SECURITY.md")
            Expect(err).NotTo(HaveOccurred())
            Expect(string(content)).To(ContainSubstring("Acme Corp"))
        })

        It("includes the contact email", func() {
            Expect(tool.GeneratePolicy(nil, config)).To(Succeed())

            content, err := os.ReadFile(outputDir + "/SECURITY.md")
            Expect(err).NotTo(HaveOccurred())
            Expect(string(content)).To(ContainSubstring("security@acme.com"))
        })
    })

    When("required fields are missing", func() {
        Context("when organization is empty", func() {
            BeforeEach(func() {
                config.Organization = ""
            })

            It("returns a meaningful error", func() {
                err := tool.GeneratePolicy(nil, config)
                Expect(err).To(MatchError(ContainSubstring("organization")))
            })
        })

        Context("when contact email is invalid", func() {
            BeforeEach(func() {
                config.ContactEmail = "not-an-email"
            })

            It("returns a validation error", func() {
                err := tool.GeneratePolicy(nil, config)
                Expect(err).To(MatchError(ContainSubstring("email")))
            })
        })
    })
})
```

### 4. Add Validation BDD Tests

```go
var _ = Describe("Security Policy Validation", Label("acceptance"), func() {
    var validator *internal.SecurityValidator

    BeforeEach(func() {
        validator = internal.NewSecurityValidator()
    })

    DescribeTable("validating SECURITY.md completeness",
        func(content string, expectValid bool, expectedErrors []string) {
            tmpFile := createTempFile(content)

            result, err := validator.ValidateSECURITYMd(tmpFile)
            Expect(err).NotTo(HaveOccurred())
            Expect(result.Valid).To(Equal(expectValid))

            for _, expected := range expectedErrors {
                Expect(result.Errors).To(ContainElement(ContainSubstring(expected)))
            }
        },
        Entry("complete policy",
            completeSecurityPolicy,
            true,
            []string{}),
        Entry("missing vulnerability reporting section",
            missingReportingSection,
            false,
            []string{"Reporting a Vulnerability"}),
        Entry("missing contact information",
            missingContactInfo,
            false,
            []string{"contact email"}),
        Entry("unresolved template variables",
            unresolvedVariables,
            false,
            []string{"template variable"}),
    )
})
```

---

## Implementation Plan

### Phase 1: Setup (1-2 hours)

- [ ] Add Ginkgo v2 and Gomega to go.mod
- [ ] Create `test/acceptance/` directory structure
- [ ] Create suite file with bootstrap

### Phase 2: Core Scenarios (2-4 hours)

- [ ] Policy generation happy path
- [ ] Policy generation error cases
- [ ] Validation happy path
- [ ] Validation error cases

### Phase 3: Edge Cases (1-2 hours)

- [ ] Table-driven tests for multiple scenarios
- [ ] Error message quality verification
- [ ] File system edge cases

### Phase 4: Integration (2-3 hours)

- [ ] CLI integration tests
- [ ] End-to-end workflow tests
- [ ] CI/CD integration verification

---

## Quality Checklist

Before considering BDD tests complete:

- [ ] Tests read like natural language documentation
- [ ] Each `It` block describes ONE behavior
- [ ] `BeforeEach` establishes clear context
- [ ] `Context`/`When` blocks explain conditions
- [ ] Error messages are meaningful and actionable
- [ ] Tests are independent (no shared mutable state)
- [ ] Tests cover happy path AND error paths
- [ ] Table-driven tests for repetitive scenarios
- [ ] Labels applied for test categorization
- [ ] `go test` passes with full suite

---

## Test Coverage Goals

| Category             | Current | Target |
| -------------------- | ------- | ------ |
| Unit Tests           | ~70%    | 80%    |
| BDD Acceptance Tests | 0%      | 90%    |
| Integration Tests    | 0%      | 70%    |
| E2E Tests            | 0%      | 50%    |

---

## References

- [Ginkgo Documentation](https://onsi.github.io/ginkgo/)
- [Gomega Matchers](https://onsi.github.io/gomega/)
- [BDD Testing Best Practices](https://martinfowler.com/bliki/GivenWhenThen.html)

---

## Conclusion

The current test suite provides **zero** behavioral coverage. Tests are implementation-focused and do not serve as living documentation. Immediate action is required to:

1. **Adopt Ginkgo** - Industry standard for Go BDD
2. **Write user stories** - From end-user perspective
3. **Structure with Given-When-Then** - Clear behavioral intent
4. **Prioritize acceptance tests** - Business value over implementation details

**Priority: Critical**
**Estimated Effort: 1-2 days**
**Impact: High - Tests will serve as executable documentation**
