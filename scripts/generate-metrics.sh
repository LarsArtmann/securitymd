#!/bin/bash
set -euo pipefail

# Security Metrics Generation Script
# Generates security metrics and compliance dashboards

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

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

generate_security_metrics() {
	local output_format="${1:-text}"
	local timestamp=$(date -Iseconds)

	log_info "Generating security metrics..."

	# Count security policies
	local policy_count=$(find "$PROJECT_DIR" -name "*.md" -type f | wc -l | tr -d ' ')

	# Check for key security documents
	local has_security_md=0
	local has_incident_response=0
	local has_privacy_policy=0
	local has_enterprise_policy=0

	[[ -f "$PROJECT_DIR/SECURITY.md" ]] && has_security_md=1
	[[ -f "$PROJECT_DIR/incident-response.md" ]] && has_incident_response=1
	[[ -f "$PROJECT_DIR/privacy-policy.md" ]] && has_privacy_policy=1
	[[ -f "$PROJECT_DIR/security-policy.md" ]] && has_enterprise_policy=1

	# Compliance scores (simplified)
	local gdpr_score=$((has_privacy_policy * 80 + has_security_md * 20))
	local soc2_score=$((has_security_md * 30 + has_incident_response * 40 + has_enterprise_policy * 30))
	local iso27001_score=$((has_enterprise_policy * 35 + has_incident_response * 35 + has_security_md * 30))

	# Security maturity score
	local maturity_score=$(((has_security_md + has_incident_response + has_privacy_policy + has_enterprise_policy) * 25))

	case "$output_format" in
	"json")
		cat <<EOF
{
  "generated": "$timestamp",
  "security_policies": {
    "total_count": $policy_count,
    "security_md": $has_security_md,
    "incident_response": $has_incident_response,
    "privacy_policy": $has_privacy_policy,
    "enterprise_policy": $has_enterprise_policy
  },
  "compliance_scores": {
    "gdpr": $gdpr_score,
    "soc2": $soc2_score,
    "iso27001": $iso27001_score
  },
  "security_maturity": {
    "score": $maturity_score,
    "level": "$(get_maturity_level $maturity_score)"
  },
  "recommendations": [
    $(get_recommendations $has_security_md $has_incident_response $has_privacy_policy $has_enterprise_policy)
  ]
}
EOF
		;;
	"prometheus")
		cat <<EOF
# HELP security_policies_total Total number of security policies
# TYPE security_policies_total gauge
security_policies_total $policy_count

# HELP security_policy_exists Whether specific security policies exist (1 = yes, 0 = no)
# TYPE security_policy_exists gauge
security_policy_exists{policy="security_md"} $has_security_md
security_policy_exists{policy="incident_response"} $has_incident_response
security_policy_exists{policy="privacy_policy"} $has_privacy_policy
security_policy_exists{policy="enterprise_policy"} $has_enterprise_policy

# HELP compliance_score Compliance score by framework (0-100)
# TYPE compliance_score gauge
compliance_score{framework="gdpr"} $gdpr_score
compliance_score{framework="soc2"} $soc2_score
compliance_score{framework="iso27001"} $iso27001_score

# HELP security_maturity_score Security maturity score (0-100)
# TYPE security_maturity_score gauge
security_maturity_score $maturity_score

# HELP security_metrics_last_update Last update timestamp
# TYPE security_metrics_last_update gauge
security_metrics_last_update $(date +%s)
EOF
		;;
	*)
		cat <<EOF
🔒 Security Metrics Report
Generated: $timestamp

📋 Security Policies Overview:
   Total Policies: $policy_count
   ✓ SECURITY.md: $([ $has_security_md -eq 1 ] && echo "Present" || echo "Missing")
   ✓ Incident Response: $([ $has_incident_response -eq 1 ] && echo "Present" || echo "Missing")
   ✓ Privacy Policy: $([ $has_privacy_policy -eq 1 ] && echo "Present" || echo "Missing")
   ✓ Enterprise Policy: $([ $has_enterprise_policy -eq 1 ] && echo "Present" || echo "Missing")

📊 Compliance Scores:
   GDPR: $gdpr_score/100
   SOC 2: $soc2_score/100
   ISO 27001: $iso27001_score/100

🎯 Security Maturity: $maturity_score/100 ($(get_maturity_level $maturity_score))

💡 Recommendations:
$(get_recommendations $has_security_md $has_incident_response $has_privacy_policy $has_enterprise_policy | sed 's/,/\n   /g' | sed 's/"//g')
EOF
		;;
	esac
}

get_maturity_level() {
	local score=$1

	if [[ $score -ge 75 ]]; then
		echo "Advanced"
	elif [[ $score -ge 50 ]]; then
		echo "Intermediate"
	elif [[ $score -ge 25 ]]; then
		echo "Basic"
	else
		echo "Initial"
	fi
}

get_recommendations() {
	local has_security_md=$1
	local has_incident_response=$2
	local has_privacy_policy=$3
	local has_enterprise_policy=$4
	local recommendations=()

	[[ $has_security_md -eq 0 ]] && recommendations+=('"Generate SECURITY.md with: just github-security"')
	[[ $has_incident_response -eq 0 ]] && recommendations+=('"Create incident response plan"')
	[[ $has_privacy_policy -eq 0 ]] && recommendations+=('"Develop privacy policy for GDPR compliance"')
	[[ $has_enterprise_policy -eq 0 ]] && recommendations+='"Generate enterprise security policy"'

	if [[ ${#recommendations[@]} -eq 0 ]]; then
		recommendations+=('"All critical policies are in place - consider advanced security measures"')
	fi

	local IFS=','
	echo "${recommendations[*]}"
}

generate_executive_report() {
	local output_file="${PROJECT_DIR}/executive-security-report-$(date +%Y%m%d).pdf"

	log_info "Generating executive security report..."

	# This would typically generate a PDF using a tool like wkhtmltopdf or pandoc
	# For now, we'll create a markdown version
	local markdown_file="${output_file%.pdf}.md"

	cat >"$markdown_file" <<EOF
# Executive Security Report

**Generated:** $(date +%B %d, %Y)  
**Prepared for:** {{ORGANIZATION}} Leadership

## Executive Summary

This report provides an overview of our current security posture and compliance status.

## Security Maturity Assessment

$(generate_security_metrics "text")

## Risk Assessment

### High Priority Risks
1. **Policy Gaps**: Missing security documentation creates compliance risk
2. **Incident Response**: Lack of formal incident response procedures
3. **Data Protection**: Insufficient privacy controls for customer data

### Mitigation Strategies
- Implement comprehensive security policies
- Establish incident response procedures
- Enhance data protection measures

## Compliance Status

Our compliance programs are progressing towards full implementation:
- **GDPR**: In progress with privacy policy development
- **SOC 2**: Preparing for Type II audit
- **ISO 27001**: Implementing security management system

## Budget Recommendations

Based on current assessment, recommend investment in:
1. Security policy development tools and training
2. Incident response platform and team
3. Compliance automation solutions
4. Security awareness training program

## Next Steps

1. **Immediate** (0-30 days):
   - Complete missing security policies
   - Establish incident response procedures
   
2. **Short-term** (30-90 days):
   - Implement compliance monitoring
   - Conduct security awareness training
   
3. **Long-term** (90+ days):
   - Pursue formal compliance certifications
   - Establish continuous security improvement program

---

*This report should be reviewed quarterly and updated as our security program matures.*
EOF

	log_success "Executive report generated: $markdown_file"
	log_info "To convert to PDF: pandoc \"$markdown_file\" -o \"$output_file\""
}

main() {
	local format="text"
	local executive_mode=false

	while [[ $# -gt 0 ]]; do
		case $1 in
		--format)
			format="$2"
			shift 2
			;;
		--executive)
			executive_mode=true
			shift
			;;
		-h | --help)
			echo "Usage: $0 [--format FORMAT] [--executive]"
			echo "Formats: text, json, prometheus"
			exit 0
			;;
		*)
			log_error "Unknown option: $1"
			exit 1
			;;
		esac
	done

	if [[ "$executive_mode" == "true" ]]; then
		generate_executive_report
	else
		generate_security_metrics "$format"
	fi
}

main "$@"
