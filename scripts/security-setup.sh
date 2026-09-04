#!/bin/bash
set -euo pipefail

# Template Security Setup Script
# Interactive wizard for generating security policies
# Usage: ./security-setup.sh [options]

VERSION="1.0.0"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"
TEMPLATES_DIR="${SCRIPT_DIR}/templates"
OUTPUT_DIR="."
POLICY_TYPE=""
ORGANIZATION=""
CONTACT_EMAIL=""
INTERACTIVE=false

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
WHITE='\033[1;37m'
NC='\033[0m' # No Color

print_header() {
	echo -e "${PURPLE}╔══════════════════════════════════════════════════════════════════╗${NC}"
	echo -e "${PURPLE}║                    🔒 Security Policy Setup v${VERSION}                 ║${NC}"
	echo -e "${PURPLE}║              Let's secure your project with proper policies!     ║${NC}"
	echo -e "${PURPLE}╚══════════════════════════════════════════════════════════════════╝${NC}"
	echo
}

print_usage() {
	cat <<EOF
Usage: $0 [OPTIONS]

Interactive security policy generation wizard.

OPTIONS:
    -t, --type TYPE         Policy type (github, enterprise, bug-bounty, incident-response, privacy-policy)
    -o, --organization ORG  Organization name
    -e, --email EMAIL       Security contact email
    --output DIR           Output directory (default: current directory)
    --interactive          Force interactive mode
    --quick               Quick setup with minimal prompts
    --update              Update existing policies
    -h, --help             Show this help message
    -v, --version          Show version information

EXAMPLES:
    $0 --interactive                        # Full interactive wizard
    $0 --type github --quick               # Quick GitHub SECURITY.md
    $0 --type enterprise --organization "Acme Corp"  # Enterprise policy
    $0 --update --policies SECURITY.md    # Update existing policy

SUPPORTED POLICY TYPES:
    github           - GitHub SECURITY.md for open source projects
    enterprise       - Comprehensive organizational security policy
    bug-bounty       - Bug bounty and responsible disclosure program
    incident-response- Incident response and management plan
    privacy-policy   - GDPR/CCPA privacy policy

EOF
}

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

prompt_with_default() {
	local prompt="$1"
	local default="$2"
	local response

	if [[ -n "$default" ]]; then
		echo -n -e "${CYAN}?${NC} $prompt ${YELLOW}(default: $default)${NC} › "
	else
		echo -n -e "${CYAN}?${NC} $prompt › "
	fi

	read -r response
	echo "${response:-$default}"
}

generate_github_security() {
	log_info "Generating GitHub SECURITY.md..."

	local template_file="${TEMPLATES_DIR}/github-security.md"
	local output_file="${OUTPUT_DIR}/SECURITY.md"

	if [[ ! -f "$template_file" ]]; then
		log_error "Template not found: $template_file"
		return 1
	fi

	# Read template and substitute variables
	local content
	content=$(<"$template_file")

	# Replace placeholders
	content="${content//\{\{ORGANIZATION\}\}/$ORGANIZATION}"
	content="${content//\{\{CONTACT_EMAIL\}\}/$CONTACT_EMAIL}"

	# Write to output file
	echo "$content" >"$output_file"

	log_success "Created ${WHITE}$output_file${NC}"
	log_info "Security contact: ${WHITE}$CONTACT_EMAIL${NC}"
}

generate_enterprise_policy() {
	log_info "Generating Enterprise Security Policy..."

	local template_file="${TEMPLATES_DIR}/enterprise-policy.md"
	local output_file="${OUTPUT_DIR}/security-policy.md"

	if [[ ! -f "$template_file" ]]; then
		log_warning "Enterprise template not found, creating comprehensive policy..."
		create_enterprise_template "$template_file"
	fi

	# Read template and substitute variables
	local content
	content=$(<"$template_file")

	# Replace placeholders
	content="${content//\{\{ORGANIZATION\}\}/$ORGANIZATION}"
	content="${content//\{\{CONTACT_EMAIL\}\}/$CONTACT_EMAIL}"

	# Write to output file
	echo "$content" >"$output_file"

	log_success "Created ${WHITE}$output_file${NC}"
	log_info "Security contact: ${WHITE}$CONTACT_EMAIL${NC}"
}

generate_bug_bounty_program() {
	log_info "Generating Bug Bounty Program..."

	local template_file="${TEMPLATES_DIR}/bug-bounty-program.md"
	local output_file="${OUTPUT_DIR}/bug-bounty-policy.md"

	if [[ ! -f "$template_file" ]]; then
		log_warning "Bug bounty template not found, creating comprehensive program..."
		create_bug_bounty_template "$template_file"
	fi

	# Read template and substitute variables
	local content
	content=$(<"$template_file")

	# Replace placeholders
	content="${content//\{\{ORGANIZATION\}\}/$ORGANIZATION}"
	content="${content//\{\{CONTACT_EMAIL\}\}/$CONTACT_EMAIL}"

	# Write to output file
	echo "$content" >"$output_file"

	log_success "Created ${WHITE}$output_file${NC}"
	log_info "Security contact: ${WHITE}$CONTACT_EMAIL${NC}"
}

generate_incident_response() {
	log_info "Generating Incident Response Plan..."

	local template_file="${TEMPLATES_DIR}/incident-response.md"
	local output_file="${OUTPUT_DIR}/incident-response.md"

	if [[ ! -f "$template_file" ]]; then
		log_error "Incident response template not found: $template_file"
		return 1
	fi

	# Read template and substitute variables
	local content
	content=$(<"$template_file")

	# Replace placeholders
	content="${content//\{\{ORGANIZATION\}\}/$ORGANIZATION}"
	content="${content//\{\{CONTACT_EMAIL\}\}/$CONTACT_EMAIL}"

	# Write to output file
	echo "$content" >"$output_file"

	log_success "Created ${WHITE}$output_file${NC}"
	log_info "Security contact: ${WHITE}$CONTACT_EMAIL${NC}"
}

generate_privacy_policy() {
	log_info "Generating Privacy Policy..."

	local template_file="${TEMPLATES_DIR}/privacy-policy.md"
	local output_file="${OUTPUT_DIR}/privacy-policy.md"

	if [[ ! -f "$template_file" ]]; then
		log_error "Privacy policy template not found: $template_file"
		return 1
	fi

	# Read template and substitute variables
	local content
	content=$(<"$template_file")

	# Replace placeholders
	content="${content//\{\{ORGANIZATION\}\}/$ORGANIZATION}"
	content="${content//\{\{CONTACT_EMAIL\}\}/$CONTACT_EMAIL}"

	# Write to output file
	echo "$content" >"$output_file"

	log_success "Created ${WHITE}$output_file${NC}"
	log_info "Privacy contact: ${WHITE}$CONTACT_EMAIL${NC}"
}

create_enterprise_template() {
	local template_file="$1"

	cat >"$template_file" <<'EOF'
# Enterprise Security Policy

## Purpose

This document outlines the comprehensive security policies and procedures for {{ORGANIZATION}}.

## Information Security Policy

### Policy Statement

{{ORGANIZATION}} is committed to protecting the confidentiality, integrity, and availability of information assets.

### Scope

This policy applies to all employees, contractors, and third parties with access to {{ORGANIZATION}} information systems.

### Security Roles and Responsibilities

#### Management
- Ensure adequate resources for security program
- Approve security policies and procedures
- Review security program effectiveness

#### IT Department
- Implement technical security controls
- Monitor systems for security incidents
- Maintain security infrastructure

#### Employees
- Follow security policies and procedures
- Report security incidents promptly
- Attend security awareness training

## Risk Management

### Risk Assessment Process

1. Identify assets and risks
2. Analyze and evaluate risks
3. Implement risk treatment measures
4. Monitor and review risks

### Risk Categories

- Strategic Risk
- Operational Risk
- Financial Risk
- Compliance Risk

## Access Control

### Access Control Policy

Access to information systems shall be granted based on:

- Principle of least privilege
- Job role requirements
- Business need-to-know

### User Access Management

- User accounts created for authorized personnel only
- Access rights reviewed quarterly
- Immediate deactivation for terminated employees

## Incident Response

### Incident Classification

- **Low**: Minimal impact, localized effect
- **Medium**: Significant impact, department-wide effect
- **High**: Critical impact, organization-wide effect
- **Critical**: Severe impact, regulatory breach

### Response Procedures

1. **Detection and Analysis**
   - Monitor security alerts
   - Analyze potential incidents
   - Classify incident severity

2. **Containment**
   - Isolate affected systems
   - Prevent further damage
   - Preserve evidence

3. **Eradication**
   - Remove root cause
   - Patch vulnerabilities
   - Clean affected systems

4. **Recovery**
   - Restore normal operations
   - Monitor for recurrence
   - Document lessons learned

## Security Awareness Training

### Training Program

- Initial security orientation for all new hires
- Annual refresher training
- Role-specific training for sensitive positions

### Training Topics

- Password security
- Phishing awareness
- Social engineering
- Physical security
- Data classification

## Compliance

### Regulatory Requirements

- GDPR (General Data Protection Regulation)
- SOC 2 Type II
- ISO 27001
- Industry-specific regulations

### Compliance Monitoring

- Quarterly compliance reviews
- Annual security assessments
- Third-party audits

## Contact Information

Security Team: {{CONTACT_EMAIL}}
Security Hotline: [Add phone number]

Last Updated: $(date +%Y-%m-%d)
EOF
}

create_bug_bounty_template() {
	local template_file="$1"

	cat >"$template_file" <<'EOF'
# Bug Bounty Program

## Welcome to {{ORGANIZATION}}'s Bug Bounty Program!

Thank you for helping keep {{ORGANIZATION}} and our users safe! We appreciate the security community's efforts in identifying and reporting vulnerabilities.

## Program Rules

### Scope

The following assets are in scope:

- Web applications: *.{{ORGANIZATION}}.com
- Mobile applications: iOS and Android apps
- API endpoints: api.{{ORGANIZATION}}.com
- Infrastructure: Public-facing servers and services

### Out of Scope

- Third-party services and dependencies
- Denial of service attacks
- Physical attacks on facilities or employees
- Social engineering attacks on employees
- Vulnerabilities in outdated browsers

### Reward Ranges

- **Critical**: $1,000 - $5,000
- **High**: $500 - $1,000
- **Medium**: $100 - $500
- **Low**: $50 - $100
- **Informational**: $0 - $50

### Submission Guidelines

Please include the following in your report:

- Clear, concise description of the vulnerability
- Steps to reproduce the issue
- Proof of concept or exploit
- Potential impact assessment
- Suggested remediation (if applicable)

### Safe Harbor

We commit to not taking legal action against researchers who:

- Follow this policy
- Report vulnerabilities in good faith
- Do not access data beyond what's necessary
- Do not degrade service performance
- Disclose findings responsibly

### Reporting Process

**Please report vulnerabilities to:** {{CONTACT_EMAIL}}

You should receive an acknowledgment within 24 hours. If you don't hear back within 48 hours, please follow up.

### Review Process

1. **Triage**: We review and prioritize submissions
2. **Validation**: We reproduce and validate the vulnerability
3. **Coordination**: We work with you on disclosure timing
4. **Reward**: We determine and pay rewards
5. **Recognition**: We acknowledge your contribution (with permission)

### Recognition

Participants who contribute valid vulnerabilities may be recognized in:

- Our Security Researcher Hall of Fame
- Annual security reports
- Conference presentations (with permission)

## Questions?

If you have questions about this program, please contact us at {{CONTACT_EMAIL}}.

Thank you for helping make {{ORGANIZATION}} more secure!

Program Last Updated: $(date +%Y-%m-%d)
EOF
}

interactive_setup() {
	print_header

	# Detect project context
	local project_name=""
	if [[ -f "package.json" ]]; then
		project_name=$(grep -o '"name": *"[^"]*"' package.json 2>/dev/null | cut -d'"' -f4 || echo "")
	elif [[ -f "go.mod" ]]; then
		project_name=$(head -1 go.mod 2>/dev/null | awk '{print $2}' | xargs basename || echo "")
	elif [[ -f "Cargo.toml" ]]; then
		project_name=$(grep -o '^name = *"[^"]*"' Cargo.toml 2>/dev/null | cut -d'"' -f2 || echo "")
	fi

	if [[ -n "$project_name" ]]; then
		log_info "Detected project: ${WHITE}$project_name${NC}"
		echo
	fi

	# Policy type selection
	echo -e "${WHITE}Security Policy Type:${NC}"
	echo -e "  ${CYAN}1)${NC} GitHub SECURITY.md (Open source projects)"
	echo -e "  ${CYAN}2)${NC} Enterprise Security Policy (Organizations)"
	echo -e "  ${CYAN}3)${NC} Bug Bounty Program (Public-facing applications)"
	echo -e "  ${CYAN}4)${NC} Incident Response Plan (Security incident management)"
	echo -e "  ${CYAN}5)${NC} Privacy Policy (GDPR/CCPA compliance)"
	echo

	local selection
	selection=$(prompt_with_default "Select policy type (1-5)" "1")

	case "$selection" in
	1) POLICY_TYPE="github" ;;
	2) POLICY_TYPE="enterprise" ;;
	3) POLICY_TYPE="bug-bounty" ;;
	4) POLICY_TYPE="incident-response" ;;
	5) POLICY_TYPE="privacy-policy" ;;
	*) POLICY_TYPE="github" ;;
	esac

	echo
	log_info "Selected: ${WHITE}$POLICY_TYPE${NC}"

	# Organization information
	ORGANIZATION=$(prompt_with_default "Organization/Project name" "${project_name:-MyProject}")
	CONTACT_EMAIL=$(prompt_with_default "Security contact email" "security@${ORGANIZATION}.com")

	echo
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
	case $1 in
	-t | --type)
		POLICY_TYPE="$2"
		shift 2
		;;
	-o | --organization)
		ORGANIZATION="$2"
		shift 2
		;;
	-e | --email)
		CONTACT_EMAIL="$2"
		shift 2
		;;
	--output)
		OUTPUT_DIR="$2"
		shift 2
		;;
	--interactive)
		INTERACTIVE=true
		shift
		;;
	--quick)
		INTERACTIVE=false
		shift
		;;
	-h | --help)
		print_usage
		exit 0
		;;
	-v | --version)
		echo "Security Setup v$VERSION"
		exit 0
		;;
	*)
		log_error "Unknown option: $1"
		print_usage
		exit 1
		;;
	esac
done

# Create templates directory if it doesn't exist
mkdir -p "$TEMPLATES_DIR"

# Main execution
if [[ "$INTERACTIVE" == "true" ]] || [[ -z "$POLICY_TYPE" ]]; then
	interactive_setup
fi

# Set defaults if not provided
if [[ -z "$ORGANIZATION" ]]; then
	ORGANIZATION="YourOrganization"
fi

if [[ -z "$CONTACT_EMAIL" ]]; then
	CONTACT_EMAIL="security@${ORGANIZATION}.com"
fi

if [[ -z "$POLICY_TYPE" ]]; then
	POLICY_TYPE="github"
fi

# Generate the appropriate policy
case "$POLICY_TYPE" in
"github")
	generate_github_security
	;;
"enterprise")
	generate_enterprise_policy
	;;
"bug-bounty")
	generate_bug_bounty_program
	;;
"incident-response")
	generate_incident_response
	;;
"privacy-policy")
	generate_privacy_policy
	;;
*)
	log_error "Unknown policy type: $POLICY_TYPE"
	exit 1
	;;
esac

echo
log_success "Security policy setup complete!"
log_info "Next steps:"
echo "  1. Review and customize the generated policy"
echo "  2. Set up ${CONTACT_EMAIL} email address"
echo "  3. Test the security reporting process"
echo "  4. Consider additional security measures"
