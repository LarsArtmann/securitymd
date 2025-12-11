# 🔥 Quick Commands (1% effort → 51% results)

# 🚀 Quick GitHub security setup
github:
    @echo "🚀 Quick GitHub SECURITY.md setup..."
    @just build && ./bin/template-security setup --type github --organization "$(git config --get remote.origin.url 2>/dev/null | xargs basename -s .git 2>/dev/null || echo 'MyProject')" --email "security@$(git config --get remote.origin.url 2>/dev/null | sed 's/.*\/\([^/]*\)\/.*/\1/' 2>/dev/null || echo 'example').com" --quick

# 🏢 Quick enterprise security setup  
enterprise:
    @echo "🏢 Quick enterprise security setup..."
    @just build && ./bin/template-security setup --type enterprise --organization "$(git config --get remote.origin.url 2>/dev/null | xargs basename -s .git 2>/dev/null || echo 'MyCompany')" --email "security@$(git config --get remote.origin.url 2>/dev/null | sed 's/.*\/\([^/]*\)\/.*/\1/' 2>/dev/null || echo 'example').com" --quick

# ✅ Quick validate
check:
    @echo "✅ Quick security validation..."
    @just build && ./bin/template-security validate

# 🧪 Quick test
test-quick:
    @echo "🧪 Quick functionality test..."
    @just build
    @just clean-generated
    @./bin/template-security setup --type github --organization "TestProject" --email "security@test.com" --quick
    @test -f SECURITY.md && echo "✅ SECURITY.md generated successfully"
    @./bin/template-security validate
    @echo "✅ All tests passed!"

# Template Security v2 - Go Implementation
# Enterprise-grade security policy generation with template-CLI SDK

default:
    @just --list

# 🔧 Build the Go application
build:
    @echo "🔧 Building template-security..."
    @./scripts/build.sh

# 🧹 Clean build artifacts and cache
clean:
    @echo "🧹 Cleaning..."
    @rm -rf bin/
    @rm -f security-*.md
    @rm -f *.tmp
    @echo "✅ Clean complete"

# 🚀 Setup security policies - interactive wizard
setup:
    @echo "🚀 Setting up security policies..."
    @./bin/template-security setup

# 📋 Generate GitHub SECURITY.md - quick setup for open source
github-security:
    @echo "📋 Generating GitHub SECURITY.md..."
    @./bin/template-security setup --type github --organization "MyProject" --email "security@example.com" --quick

# 🏢 Generate enterprise security policy
enterprise-policy:
    @echo "🏢 Generating enterprise security policy..."
    @./bin/template-security setup --type enterprise --organization "MyProject" --email "security@example.com" --quick

# 🐛 Generate bug bounty program
bug-bounty:
    @echo "🐛 Generating bug bounty program..."
    @./bin/template-security setup --type bug-bounty --organization "MyProject" --email "security@example.com" --quick

# 📋 Generate incident response plan
incident-response:
    @echo "📋 Generating incident response plan..."
    @./bin/template-security setup --type incident-response --organization "MyProject" --email "security@example.com" --quick

# 🔒 Generate privacy policy
privacy-policy:
    @echo "🔒 Generating privacy policy..."
    @./bin/template-security setup --type privacy-policy --organization "MyProject" --email "security@example.com" --quick

# ✅ Validate security policies
validate:
    @echo "✅ Validating security policies..."
    @./bin/template-security validate

# 📊 Generate compliance report
compliance-report FRAMEWORK:
    @echo "📊 Generating compliance report for: {{FRAMEWORK}}"
    @./bin/template-security compliance --framework {{FRAMEWORK}}

# 📊 Generate comprehensive compliance dashboard
compliance-dashboard:
    @echo "📊 Generating comprehensive compliance dashboard..."
    @./bin/template-security compliance --report

# 📈 Generate security metrics
metrics:
    @echo "📈 Generating security metrics..."
    @./bin/template-security metrics

# 📊 Generate executive security report
executive-report:
    @echo "📊 Generating executive security report..."
    @./bin/template-security metrics --executive

# 📋 Show security policy status
status:
    @echo "📋 Security Policy Status"
    @./bin/template-security status

# 🛠 Development mode - build and watch
dev:
    @echo "🛠 Development mode - building and testing..."
    @just build
    @echo "📝 Running quick tests..."
    @./bin/template-security --version
    @./bin/template-security status
    @echo "✅ Development ready!"

# 🧪 Run integration tests
test:
    @echo "🧪 Running integration tests..."
    @./bin/template-security setup --type github --organization "TestCorp" --email "security@testcorp.com" --quick
    @test -f SECURITY.md || (echo "❌ SECURITY.md not generated" && exit 1)
    @./bin/template-security setup --type enterprise --organization "TestCorp" --email "security@testcorp.com" --quick
    @test -f security-policy.md || (echo "❌ security-policy.md not generated" && exit 1)
    @echo "✅ All tests passed!"

# 🚀 Quick start - build and setup
quick-start:
    @echo "🚀 Template Security Quick Start"
    @echo "=============================="
    @just build
    @just status
    @echo ""
    @echo "✅ Template Security v2 is ready!"
    @echo "Next steps:"
    @echo "1. Run 'just setup' for interactive wizard"
    @echo "2. Run 'just github' for quick GitHub setup"
    @echo "3. Run 'just check' to validate policies"

# 📦 Install dependencies
deps:
    @echo "📦 Installing dependencies..."
    @go mod download
    @go mod tidy
    @echo "✅ Dependencies ready"

# 🔍 List available templates
list-templates:
    @echo "🔍 Available Security Policy Templates:"
    @echo "======================================"
    @ls -la templates/ 2>/dev/null || echo "No templates directory found"

# 🧪 Clean generated policies and test
clean-generated:
    @echo "🧹 Cleaning generated policies..."
    @rm -f SECURITY.md security-policy.md bug-bounty-policy.md incident-response.md privacy-policy.md
    @echo "✅ Clean complete"

# 🔍 Help - show available commands
help:
    @echo "🔒 Template Security v2 - Go Implementation"
    @echo "=========================================="
    @echo ""
    @echo "🔥 Quick Commands (80/20 optimized):"
    @echo "  just github           - Auto-detect repo & generate SECURITY.md"
    @echo "  just enterprise       - Auto-detect org & generate enterprise policy"
    @echo "  just check            - Quick validation of security policies"
    @echo "  just test-quick       - Quick functionality test"
    @echo ""
    @echo "Standard Commands:"
    @echo "  just setup            - Interactive security policy wizard"
    @echo "  just github-security   - Generate GitHub SECURITY.md"
    @echo "  just enterprise-policy - Generate enterprise policy"
    @echo "  just validate         - Validate all policies"
    @echo "  just status           - Show policy status"
    @echo ""
    @echo "Policy Generation:"
    @echo "  just bug-bounty       - Generate bug bounty program"
    @echo "  just incident-response - Generate incident response plan"
    @echo "  just privacy-policy   - Generate privacy policy"
    @echo ""
    @echo "Development:"
    @echo "  just build            - Build the application"
    @echo "  just test             - Run integration tests"
    @echo "  just dev              - Development mode"
    @echo "  just clean            - Clean build artifacts"
    @echo ""
    @echo "For detailed help, run: ./bin/template-security --help"