# Security Policy

## Supported Versions

| Version            | Supported Until     |
| ------------------ | ------------------- |
| {{.LatestVersion}} | {{.SupportEndDate}} |

> **Note**: Only the latest major version receives security updates. Please upgrade to the latest version as soon as possible.

## Reporting a Vulnerability

### 🚨 Private Disclosure Process

**Please do NOT report security vulnerabilities through public issues.**

1. **Email us**: {{.ContactEmail}}
2. **Use subject line**: `Security Vulnerability Report - [Brief Description]`
3. **Include details**:
   - Vulnerability type and severity
   - Steps to reproduce
   - Potential impact
   - Proof of concept (if available)

### 📋 What to Expect

| Phase                 | Timeline        | Description                                  |
| --------------------- | --------------- | -------------------------------------------- |
| **Initial Response**  | Within 24 hours | We acknowledge receipt and assess severity   |
| **Investigation**     | Within 72 hours | We reproduce and validate the vulnerability  |
| **Patch Development** | Within 14 days  | We develop and test a security patch         |
| **Disclosure**        | Within 90 days  | Coordinated public disclosure (CVE assigned) |

### 🔐 Safe Harbor

{{.Organization}} commits to:

- **Never pursue legal action** against security researchers who follow this policy
- **Work with researchers** to understand and fix vulnerabilities
- **Credit researchers** in our security advisories (with permission)
- **Maintain communication** throughout the disclosure process

## Security Practices

### 🛡 Development Security

- **Code Review**: All code changes undergo security review
- **Static Analysis**: Automated security scanning on every commit
- **Dependency Scanning**: Regular vulnerability scans of third-party dependencies
- **Security Testing**: Automated security tests in CI/CD pipeline

### 🔒 Infrastructure Security

- **Encryption**: All data transmission uses TLS 1.3+
- **Access Control**: Principle of least privilege with MFA required
- **Monitoring**: 24/7 security monitoring and alerting

---

_This security policy is last updated: {{.LastUpdated}}_
