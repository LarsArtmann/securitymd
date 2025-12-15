# GitHub Security Policy Validation

This GitHub Action automatically validates your security policy files to ensure they meet GitHub best practices and include all required information.

## What it Checks

### Security Policy Requirements ✅

- **SECURITY.md** file exists and contains required sections
- **Contact Information** with valid email addresses
- **Version Support** information
- **Vulnerability Reporting** process
- **Security Practices** documentation

### Quality Checks 🔍

- Unresolved template variables (`{{variable}}`)
- Placeholder text (`example.com`)
- Minimum content length (20+ lines)
- Required sections presence

## Usage

### Automatic Trigger

The workflow runs automatically on:

- **Push** to main/master/develop branches
- **Pull Requests** to main/master branches
- When security policy files change

### Manual Run

```bash
# Local validation
./bin/template-security validate --file SECURITY.md

# Validate all security files
./bin/template-security validate
```

## Files Monitored

- `SECURITY.md` - Main security policy
- `security-policy.md` - Enterprise security policy
- `templates/` - Template files
- `internal/security_*.go` - Validation logic

## Error Handling

If validation fails:

1. **PR Comment**: Automated comment with fix suggestions
2. **Status Check**: Failed status on PR
3. **GitHub Summary**: Detailed validation report
4. **Build Failure**: Prevents merge

## Required Sections

### SECURITY.md must include:

- ✅ **Security Policy** header
- ✅ **Reporting a Vulnerability** section
- ✅ **Supported Versions** information
- ✅ **Security Practices** details
- ✅ **Contact Information** with email

### Optional but Recommended:

- 📋 **What to Expect** timeline
- 🔐 **Safe Harbor** statement
- 🏆 **Security Acknowledgments**
- 💰 **Bug Bounty** information

## Fixing Common Issues

### Template Variables

```diff
- {{CONTACT_EMAIL}}
+ security@yourcompany.com
```

### Missing Sections

```diff
+ ## Reporting a Vulnerability
+ If you discover a security vulnerability, please report it to us privately...
```

### Placeholder Text

```diff
- example.com
+ yourcompany.com
```

## Examples

### Good Security Policy

- ✅ Clear contact information
- ✅ Detailed reporting process
- ✅ Version support information
- ✅ Security practices
- ✅ Response time expectations

### Needs Improvement

- ❌ Missing email contact
- ❌ No version information
- ❌ Incomplete reporting process
- ❌ Unresolved template variables

## Integration with Other Tools

This validation works well with:

- **Dependabot**: Automated dependency updates
- **CodeQL**: Code security scanning
- **Secret Scanning**: Credential detection
- **GitHub Security Advisories**: Vulnerability disclosure

## Configuration

You can customize validation by modifying:

- `.github/workflows/security-validation.yml`
- Required sections list
- Validation rules
- Error messages

## Support

- **Issues**: [GitHub Issues](https://github.com/LarsArtmann/template-SECURITY/issues)
- **Documentation**: [README.md](https://github.com/LarsArtmann/template-SECURITY)
- **Security**: Report vulnerabilities via `security@larsartmann.com`

---

_This action is part of the template-SECURITY project for automated security policy management._
