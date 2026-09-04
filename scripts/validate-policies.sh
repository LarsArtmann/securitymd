#!/bin/bash
set -euo pipefail

# Security Policy Validation Script
# Validates security policies for completeness and compliance

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() {
	echo -e "${BLUE}ℹ${NC} $1"
}

log_success() {
	echo -e "${GREEN}✓${NC} $1"
}

log_warning() {
	echo -e "${YELLOW}⚠${NC} $1"
}

log_error() {
	echo -e "${RED}✗${NC} $1" >&2
}

validate_security_md() {
	local file="$1"
	local issues=0

	log_info "Validating: $file"

	if [[ ! -f "$file" ]]; then
		log_error "Security policy file not found: $file"
		return 1
	fi

	# Check for required sections
	if ! grep -q "## Reporting a Vulnerability\|## Security" "$file"; then
		log_warning "Missing vulnerability reporting section"
		((issues++))
	else
		log_success "Vulnerability reporting section found"
	fi

	# Check for contact information
	if ! grep -q "@" "$file"; then
		log_warning "No contact email found"
		((issues++))
	else
		log_success "Contact information found"
	fi

	# Check for supported versions
	if ! grep -q "Supported" "$file"; then
		log_warning "No supported versions information"
		((issues++))
	else
		log_success "Supported versions section found"
	fi

	if [[ $issues -eq 0 ]]; then
		log_success "Security policy validation passed"
		return 0
	else
		log_warning "Security policy has $issues potential issues"
		return 1
	fi
}

main() {
	log_info "Validating security policies..."
	echo

	local total_issues=0

	# Check for SECURITY.md
	if [[ -f "SECURITY.md" ]]; then
		validate_security_md "SECURITY.md"
		total_issues=$?
	else
		log_warning "SECURITY.md not found - consider generating one with: just github-security"
		((total_issues++))
	fi

	echo

	if [[ $total_issues -eq 0 ]]; then
		log_success "All security policies validated successfully"
	else
		log_warning "$total_issues validation issues found"
		echo
		log_info "To fix issues:"
		echo "  - Run: just setup (interactive wizard)"
		echo "  - Or: just github-security (quick GitHub policy)"
	fi

	return $total_issues
}

main "$@"
