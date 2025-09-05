#!/bin/bash
set -euo pipefail

# Template Security Setup Script
# Interactive wizard for generating security policies
# Usage: ./security-setup.sh [options]

VERSION="1.0.0"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &> /dev/null && pwd)"
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
    cat << EOF
Usage: $0 [OPTIONS]

Interactive security policy generation wizard.

OPTIONS:
    -t, --type TYPE         Policy type (github, enterprise, bug-bounty)
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
    github      - GitHub SECURITY.md for open source projects
    enterprise  - Comprehensive organizational security policy
    bug-bounty  - Bug bounty and responsible disclosure program

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
    echo "$content" > "$output_file"
    
    log_success "Created ${WHITE}$output_file${NC}"
    log_info "Security contact: ${WHITE}security@${ORGANIZATION}.com${NC}"
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
    echo
    
    local selection
    selection=$(prompt_with_default "Select policy type (1-3)" "1")
    
    case "$selection" in
        1) POLICY_TYPE="github" ;;
        2) POLICY_TYPE="enterprise" ;;
        3) POLICY_TYPE="bug-bounty" ;;
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
        -t|--type)
            POLICY_TYPE="$2"
            shift 2
            ;;
        -o|--organization)
            ORGANIZATION="$2"
            shift 2
            ;;
        -e|--email)
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
        -h|--help)
            print_usage
            exit 0
            ;;
        -v|--version)
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
        log_warning "Enterprise policy generation not yet implemented"
        log_info "Using GitHub SECURITY.md as fallback"
        generate_github_security
        ;;
    "bug-bounty")
        log_warning "Bug bounty program generation not yet implemented" 
        log_info "Using GitHub SECURITY.md as fallback"
        generate_github_security
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