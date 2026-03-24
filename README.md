# Template Security 🔒

> **Ensure your GitHub repository has a compliant SECURITY.md**

A CLI tool that validates and generates `SECURITY.md` files. Checks for required sections (vulnerability reporting, supported versions, security practices, contact info), ensures content quality, and generates compliant templates for GitHub repositories.

## Features

- **Validate** existing SECURITY.md files against GitHub standards
- **Generate** new SECURITY.md files from templates
- **Upsert** improvements to existing policies

## 🚀 Quick Start

```bash
# Install
go install github.com/LarsArtmann/template-SECURITY/cmd/template-security@latest

# Validate existing SECURITY.md
template-security validate

# Generate new SECURITY.md if missing
template-security setup --type github

# Validate specific file
template-security validate --file SECURITY.md
```

## ✅ What We Validate

### Required Sections

- ✅ Security Policy header
- ✅ Reporting a Vulnerability section
- ✅ Supported Versions information
- ✅ Security Practices description
- ✅ Contact email address
- ✅ Response time commitments

### Quality Checks

- ✅ Minimum content length (>20 lines)
- ✅ No unresolved template variables
- ✅ Substantive content (not just placeholders)
- ✅ Version information included
- ✅ Actual email addresses present

### Example Validation Output

```bash
$ template-security validate

🔍 Validating security policy: SECURITY.md

✅ SECURITY.md - All checks passed
  ✅ Security Policy header found
  ✅ Reporting a Vulnerability section present
  ✅ Supported Versions included
  ✅ Security Practices described
  ✅ Contact email found
  ✅ Response time specified

📊 Validation Summary:
✅ All policies passed validation
```

## 🔧 Commands

### `validate`

Check if SECURITY.md exists and meets standards

```bash
# Validate all security files
template-security validate

# Validate specific file
template-security validate --file SECURITY.md
```

### `setup`

Generate new SECURITY.md if missing

```bash
# Generate GitHub-style SECURITY.md
template-security setup --type github

# Generate enterprise security policy
template-security setup --type enterprise
```

### `status`

Show current security file status

```bash
template-security status
```

## 📋 Validation Standards

Based on [GitHub's security policy guidelines](https://docs.github.com/en/code-security/getting-started/adding-a-security-policy-to-your-repository) and industry best practices:

1. **Clear reporting process** with specific contact information
2. **Defined scope** of what's in-bounds for security testing
3. **Response commitments** with specific timelines
4. **Safe harbor provisions** for security researchers
5. **Version support information** for users

## 🔄 CI/CD Integration

Add to your workflow:

```yaml
# .github/workflows/security-validation.yml
name: Security Policy Validation

on:
  pull_request:
    paths:
      - "SECURITY.md"

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Install template-security
        run: go install github.com/LarsArtmann/template-SECURITY/cmd/template-security@latest
      - name: Validate SECURITY.md
        run: template-security validate
```

## 🚨 What We DON'T Do

- ❌ Generate comprehensive enterprise security frameworks
- ❌ Manage bug bounty programs
- ❌ Handle incident response automation
- ❌ Create compliance documentation (GDPR, SOC2, etc.)
- ❌ Provide security scanning or monitoring

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

---

**Made with 🔒 to improve security disclosure practices**
