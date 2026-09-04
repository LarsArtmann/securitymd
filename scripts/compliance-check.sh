#!/bin/bash
set -euo pipefail

# Compliance Framework Validation Script
# Validates policies against various compliance frameworks

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Compliance frameworks
FRAMEWORKS=("gdpr" "soc2" "iso27001" "nist" "ccpa")

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

validate_gdpr() {
	local file="$1"
	local score=0
	local max_score=8

	log_info "Validating GDPR compliance..."

	# Check for GDPR requirements
	if grep -q -i "data protection" "$file"; then
		((score++))
		log_success "Data protection principles found"
	else
		log_warning "Missing data protection principles"
	fi

	if grep -q -i "right to access\|subject rights" "$file"; then
		((score++))
		log_success "Data subject rights mentioned"
	else
		log_warning "Missing data subject rights"
	fi

	if grep -q -i "consent\|legal basis" "$file"; then
		((score++))
		log_success "Legal basis for processing found"
	else
		log_warning "Missing legal basis documentation"
	fi

	if grep -q -i "data breach\|incident notification" "$file"; then
		((score++))
		log_success "Breach notification procedures found"
	else
		log_warning "Missing breach notification procedures"
	fi

	if grep -q -i "data protection officer\|DPO" "$file"; then
		((score++))
		log_success "DPO designation mentioned"
	else
		log_warning "Missing DPO information"
	fi

	if grep -q -i "international.*transfer\|cross.border" "$file"; then
		((score++))
		log_success "International data transfer provisions found"
	else
		log_warning "Missing international transfer provisions"
	fi

	if grep -q -i "data retention\|retention period" "$file"; then
		((score++))
		log_success "Data retention policy found"
	else
		log_warning "Missing data retention policy"
	fi

	if grep -q -i "privacy.*impact.*assessment\|PIA" "$file"; then
		((score++))
		log_success "PIA requirements mentioned"
	else
		log_warning "Missing privacy impact assessment"
	fi

	echo "GDPR Compliance Score: $score/$max_score"
	return $((max_score - score))
}

validate_soc2() {
	local file="$1"
	local score=0
	local max_score=5

	log_info "Validating SOC 2 compliance..."

	if grep -q -i "security.*controls\|control.*activities" "$file"; then
		((score++))
		log_success "Security controls documentation found"
	else
		log_warning "Missing security controls"
	fi

	if grep -q -i "availability\|uptime\|service.*availability" "$file"; then
		((score++))
		log_success "Availability commitments found"
	else
		log_warning "Missing availability commitments"
	fi

	if grep -q -i "confidentiality\|data.*confidential" "$file"; then
		((score++))
		log_success "Confidentiality measures found"
	else
		log_warning "Missing confidentiality measures"
	fi

	if grep -q -i "audit\|audit.*trail\|independent.*audit" "$file"; then
		((score++))
		log_success "Audit provisions found"
	else
		log_warning "Missing audit provisions"
	fi

	if grep -q -i "incident.*response\|security.*incident" "$file"; then
		((score++))
		log_success "Incident response procedures found"
	else
		log_warning "Missing incident response procedures"
	fi

	echo "SOC 2 Compliance Score: $score/$max_score"
	return $((max_score - score))
}

validate_iso27001() {
	local file="$1"
	local score=0
	local max_score=6

	log_info "Validating ISO 27001 compliance..."

	if grep -q -i "information.*security.*policy\|ISMS" "$file"; then
		((score++))
		log_success "Information security policy found"
	else
		log_warning "Missing information security policy"
	fi

	if grep -q -i "risk.*assessment\|risk.*management" "$file"; then
		((score++))
		log_success "Risk management framework found"
	else
		log_warning "Missing risk management framework"
	fi

	if grep -q -i "access.*control\|access.*management" "$file"; then
		((score++))
		log_success "Access control measures found"
	else
		log_warning "Missing access control measures"
	fi

	if grep -q -i "security.*awareness\|training" "$file"; then
		((score++))
		log_success "Security awareness program found"
	else
		log_warning "Missing security awareness program"
	fi

	if grep -q -i "business.*continuity\|disaster.*recovery" "$file"; then
		((score++))
		log_success "Business continuity provisions found"
	else
		log_warning "Missing business continuity provisions"
	fi

	if grep -q -i "supplier.*relationship\|vendor.*management" "$file"; then
		((score++))
		log_success "Supplier management found"
	else
		log_warning "Missing supplier management"
	fi

	echo "ISO 27001 Compliance Score: $score/$max_score"
	return $((max_score - score))
}

validate_framework() {
	local framework="$1"
	local file="$2"

	case "$framework" in
	"gdpr")
		validate_gdpr "$file"
		;;
	"soc2")
		validate_soc2 "$file"
		;;
	"iso27001")
		validate_iso27001 "$file"
		;;
	*)
		log_error "Unsupported framework: $framework"
		return 1
		;;
	esac
}

generate_compliance_report() {
	local output_file="${PROJECT_DIR}/compliance-report-$(date +%Y%m%d).json"

	log_info "Generating comprehensive compliance report..."

	local total_issues=0
	local report="{"
	report+="\"generated\":\"$(date -Iseconds)\","
	report+="\"frameworks\":{"

	for framework in "${FRAMEWORKS[@]}"; do
		if [[ "$framework" == "nist" || "$framework" == "ccpa" ]]; then
			continue # Skip unsupported for now
		fi

		local framework_issues=0
		local framework_report="{"

		for file in "$PROJECT_DIR"/*.md; do
			if [[ -f "$file" ]]; then
				validate_framework "$framework" "$file" >/dev/null 2>&1 || ((framework_issues++))
			fi
		done

		framework_report+="\"issues\":$framework_issues,"
		framework_report+="\"status\":\"$([ $framework_issues -eq 0 ] && echo "compliant" || echo "non-compliant")\""
		framework_report+="}"

		report+="\"$framework\":$framework_report,"
		total_issues=$((total_issues + framework_issues))
	done

	report="${report%,}"
	report+="},"
	report+="\"total_issues\":$total_issues,"
	report+="\"overall_status\":\"$([ $total_issues -eq 0 ] && echo "compliant" || echo "non-compliant")\""
	report+="}"

	echo "$report" >"$output_file"
	log_success "Compliance report generated: $output_file"
}

main() {
	local framework=""
	local report_mode=false

	while [[ $# -gt 0 ]]; do
		case $1 in
		--framework)
			framework="$2"
			shift 2
			;;
		--report)
			report_mode=true
			shift
			;;
		-h | --help)
			echo "Usage: $0 [--framework FRAMEWORK] [--report]"
			echo "Frameworks: ${FRAMEWORKS[*]}"
			exit 0
			;;
		*)
			log_error "Unknown option: $1"
			exit 1
			;;
		esac
	done

	if [[ "$report_mode" == "true" ]]; then
		generate_compliance_report
		exit 0
	fi

	if [[ -z "$framework" ]]; then
		log_error "Please specify a framework with --framework"
		echo "Available frameworks: ${FRAMEWORKS[*]}"
		exit 1
	fi

	log_info "Validating $framework compliance..."
	echo

	local total_issues=0

	for file in "$PROJECT_DIR"/*.md; do
		if [[ -f "$file" ]]; then
			echo "Validating: $(basename "$file")"
			validate_framework "$framework" "$file"
			total_issues=$((total_issues + $?))
			echo
		fi
	done

	if [[ $total_issues -eq 0 ]]; then
		log_success "All policies are $framework compliant!"
	else
		log_warning "Found $total_issues compliance issues"
	fi

	return $total_issues
}

main "$@"
