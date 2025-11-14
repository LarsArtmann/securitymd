# Template Security Policy Management
# Provides setup, validation, and maintenance for security policies and compliance

default:
    @just --list

# 🚀 Setup security policies - interactive wizard for policy generation
setup:
    @echo "🚀 Setting up security policies..."
    @if [ ! -f "./security-setup.sh" ]; then \
        echo "❌ security-setup.sh not found. Please ensure you're in the template-SECURITY directory."; \
        exit 1; \
    fi
    @./security-setup.sh --interactive

# ✅ Validate security policies - check completeness and compliance
validate:
    @echo "✅ Validating security policies..."
    @if [ ! -d "./scripts" ]; then mkdir -p scripts; fi
    @./scripts/validate-policies.sh

# 🔍 Validate specific policy file
validate-file FILE:
    @echo "🔍 Validating security policy: {{FILE}}"
    @if [ ! -f "{{FILE}}" ]; then \
        echo "❌ File not found: {{FILE}}"; \
        exit 1; \
    fi
    @./scripts/validate-single-policy.sh "{{FILE}}"

# 📋 Generate GitHub SECURITY.md - quick setup for open source projects
github-security:
    @echo "📋 Generating GitHub SECURITY.md..."
    @./security-setup.sh --type github --quick

# 🏢 Generate enterprise security policy - comprehensive organizational policies
enterprise-policy:
    @echo "🏢 Generating enterprise security policy..."
    @./security-setup.sh --type enterprise --interactive

# 🐛 Generate bug bounty program - responsible disclosure setup
bug-bounty:
    @echo "🐛 Generating bug bounty program..."
    @./security-setup.sh --type bug-bounty --interactive

# 📋 Generate incident response plan - comprehensive incident management
incident-response:
    @echo "📋 Generating incident response plan..."
    @./security-setup.sh --type incident-response --interactive

# 🔒 Generate privacy policy - GDPR/CCPA compliance
privacy-policy:
    @echo "🔒 Generating privacy policy..."
    @./security-setup.sh --type privacy-policy --interactive

# 📊 Generate compliance report - assess policy coverage
compliance-report FRAMEWORK:
    @echo "📊 Generating compliance report for: {{FRAMEWORK}}"
    @./scripts/compliance-check.sh --framework {{FRAMEWORK}}

# 📊 Generate full compliance dashboard - all frameworks
compliance-dashboard:
    @echo "📊 Generating comprehensive compliance dashboard..."
    @./scripts/compliance-check.sh --report

# 🧪 Test security contact channels - verify contact information works
test-contacts:
    @echo "🧪 Testing security contact channels..."
    @./scripts/test-security-contacts.sh

# 🚨 Simulate vulnerability report - test disclosure process
test-vulnerability-report:
    @echo "🚨 Simulating vulnerability report process..."
    @./scripts/simulate-vuln-report.sh

# 🔄 Update policies - refresh with latest standards and requirements
update:
    @echo "🔄 Updating security policies..."
    @./security-setup.sh --update --policies all

# 🔄 Update specific compliance framework
update-compliance FRAMEWORK:
    @echo "🔄 Updating compliance framework: {{FRAMEWORK}}"
    @./scripts/update-compliance.sh --framework {{FRAMEWORK}}

# 📈 Generate security metrics - analyze policy effectiveness
metrics:
    @echo "📈 Generating security metrics..."
    @./scripts/generate-metrics.sh

# 📊 Generate executive security report
executive-report:
    @echo "📊 Generating executive security report..."
    @./scripts/executive-report.sh --format pdf

# 🧹 Clean temporary files - remove generated reports and temporary data
clean:
    @echo "🧹 Cleaning temporary files..."
    @rm -rf tmp/ reports/temp/ *.tmp
    @rm -f security-report-*.pdf compliance-*.json
    @echo "✅ Cleanup complete"

# 📋 Show security policy status - overview of current policies
status:
    @echo "📋 Security Policy Status"
    @echo "========================"
    @echo "Policies directory: $(pwd)/templates"
    @echo "Available templates: $(ls templates/ 2>/dev/null | wc -l | tr -d ' ')"
    @echo "Generated policies: $(find . -name 'SECURITY.md' -o -name '*security*.md' | wc -l | tr -d ' ')"
    @echo ""
    @if [ -f "SECURITY.md" ]; then \
        echo "✅ SECURITY.md exists"; \
    else \
        echo "❌ SECURITY.md not found"; \
    fi
    @if [ -f "incident-response.md" ]; then \
        echo "✅ Incident response plan exists"; \
    else \
        echo "❌ Incident response plan not found"; \
    fi

# 🔍 List available templates - show all security policy templates
list-templates:
    @echo "🔍 Available Security Policy Templates:"
    @echo "======================================"
    @ls -la templates/ 2>/dev/null || echo "No templates directory found"

# 📝 Generate custom template - create new security policy template
create-template NAME:
    @echo "📝 Creating custom security policy template: {{NAME}}"
    @if [ ! -d "templates" ]; then mkdir -p templates; fi
    @cp templates/base-template.md "templates/{{NAME}}.md" 2>/dev/null || \
        echo "# {{NAME}} Security Policy\n\nCustom security policy template.\n" > "templates/{{NAME}}.md"
    @echo "✅ Template created: templates/{{NAME}}.md"
    @echo "📝 Edit the template and then run: just validate-template {{NAME}}"

# ✅ Validate custom template
validate-template NAME:
    @echo "✅ Validating template: {{NAME}}"
    @./scripts/validate-template.sh "templates/{{NAME}}.md"

# 🔧 Install dependencies - setup required tools
install-deps:
    @echo "🔧 Installing dependencies..."
    @echo "Checking for required tools..."
    @command -v curl >/dev/null 2>&1 || { echo "❌ curl is required but not installed"; exit 1; }
    @command -v git >/dev/null 2>&1 || { echo "❌ git is required but not installed"; exit 1; }
    @echo "✅ All dependencies available"

# 📖 Show help - comprehensive usage information
help:
    @echo "🔒 Template Security - Security Policy Management"
    @echo "================================================"
    @echo ""
    @echo "Quick Commands:"
    @echo "  just setup              - Interactive security policy wizard"
    @echo "  just github-security    - Generate GitHub SECURITY.md"
    @echo "  just enterprise-policy  - Generate enterprise policy"
    @echo "  just validate           - Validate all policies"
    @echo "  just status            - Show policy status"
    @echo ""
    @echo "Compliance:"
    @echo "  just compliance-report gdpr     - GDPR compliance report"
    @echo "  just compliance-report soc2     - SOC 2 compliance report"
    @echo "  just compliance-report iso27001 - ISO 27001 compliance report"
    @echo ""
    @echo "Testing:"
    @echo "  just test-contacts              - Test security contact channels"
    @echo "  just test-vulnerability-report  - Test vulnerability disclosure"
    @echo ""
    @echo "Maintenance:"
    @echo "  just update             - Update all policies"
    @echo "  just metrics            - Generate security metrics"
    @echo "  just clean              - Clean temporary files"
    @echo ""
    @echo "For detailed help, see: README.md"

# 🚀 Quick start - complete setup wizard with validation
quick-start:
    @echo "🚀 Security Policy Quick Start"
    @echo "=============================="
    @just install-deps
    @just setup
    @just validate
    @just status
    @echo ""
    @echo "✅ Security policies are ready!"
    @echo "Next steps:"
    @echo "1. Review generated policies"
    @echo "2. Customize for your organization"
    @echo "3. Set up security contact channels"
    @echo "4. Test vulnerability reporting process"

# 🛠 Development mode - watch for policy changes and auto-validate
dev:
    @echo "🛠 Development mode - watching for policy changes..."
    @echo "Press Ctrl+C to stop"
    @while true; do \
        if find . -name "*.md" -newer .last-validation 2>/dev/null | grep -q .; then \
            echo "📝 Policy changes detected, validating..."; \
            just validate; \
            touch .last-validation; \
        fi; \
        sleep 2; \
    done