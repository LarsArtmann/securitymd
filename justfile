# Template Security v2 - SECURITY.md Validation Tool
# Focused on validating and generating SECURITY.md files

default:
    @just --list

# 🔧 Build the Go application
build:
    @echo "🔧 Building template-security..."
    @./scripts/build.sh

# ✅ Validate SECURITY.md exists and meets standards
validate:
    @echo "✅ Validating SECURITY.md..."
    @just build && ./bin/template-security validate

# 🚀 Quick GitHub SECURITY.md setup
github:
    @echo "🚀 Quick GitHub SECURITY.md setup..."
    @just build && ./bin/template-security setup --type github --organization "$(git config --get remote.origin.url 2>/dev/null | xargs basename -s .git 2>/dev/null || echo 'MyProject')" --email "security@$(git config --get remote.origin.url 2>/dev/null | sed 's/.*\/\([^/]*\)\/.*/\1/' 2>/dev/null || echo 'example').com"

# 📋 Show current SECURITY.md status
status:
    @echo "📋 SECURITY.md Status"
    @just build && ./bin/template-security status

# 🧪 Test validation functionality
test:
    @echo "🧪 Testing SECURITY.md validation..."
    @just build
    @./bin/template-security validate
    @echo "✅ Validation test complete"

# 🧪 Run all unit tests
test-unit:
    @echo "🧪 Running unit tests..."
    @go test ./internal/... -v

# 🧪 Run BDD acceptance tests
test-bdd:
    @echo "🧪 Running BDD acceptance tests..."
    @go test ./test/acceptance/... -v

# 🧪 Run all tests (unit + BDD)
test-all: test-unit test-bdd
    @echo "✅ All tests complete"

# 🧹 Clean build artifacts
clean:
    @echo "🧹 Cleaning..."
    @rm -rf bin/
    @rm -f SECURITY.md
    @echo "✅ Clean complete"

# 📦 Install dependencies
deps:
    @echo "📦 Installing dependencies..."
    @go mod download
    @go mod tidy
    @echo "✅ Dependencies ready"

# 🔍 Help - show available commands
help:
    @echo "🔒 Template Security - SECURITY.md Validation"
    @echo "=============================================="
    @echo ""
    @echo "Core Commands:"
    @echo "  just validate    - Validate SECURITY.md exists and meets standards"
    @echo "  just github      - Generate GitHub-style SECURITY.md"
    @echo "  just status      - Show current SECURITY.md status"
    @echo ""
    @echo "Development:"
    @echo "  just build       - Build the application"
    @echo "  just test        - Test validation functionality"
    @echo "  just clean       - Clean build artifacts"
    @echo "  just deps        - Install dependencies"
    @echo ""
    @echo "For detailed help, run: ./bin/template-security --help"