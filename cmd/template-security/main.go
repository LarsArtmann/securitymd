package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

//nolint:gochecknoglobals // Version info from build flags
func main() {
	rootCmd := &cobra.Command{
		Use:   "template-security",
		Short: "Validate and generate SECURITY.md files",
		Long: `Template-Security validates and generates SECURITY.md files for completeness 
and compliance. Ensures your security policy contains all essential sections 
and follows industry best practices for vulnerability disclosure.`,
		Example: `  # Validate existing SECURITY.md
  template-security validate

  # Validate with JSON output (CI/CD friendly)
  template-security validate --format json

  # Generate a new SECURITY.md
  template-security setup --type github --organization MyOrg --email security@example.com

  # Check current status
  template-security status`,
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
	}

	rootCmd.AddCommand(newSetupCmd())
	rootCmd.AddCommand(newValidateCmd())
	rootCmd.AddCommand(newStatusCmd())

	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
