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

func main() {
	rootCmd := &cobra.Command{
		Use:   "template-security",
		Short: "Enterprise-grade security policy generator",
		Long: `Template-Security is a comprehensive security policy template collection 
that provides enterprise-grade security documentation, vulnerability reporting 
procedures, incident response plans, and compliance frameworks.`,
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
	}

	// Add subcommands
	rootCmd.AddCommand(newSetupCmd())
	rootCmd.AddCommand(newValidateCmd())
	rootCmd.AddCommand(newComplianceCmd())
	rootCmd.AddCommand(newMetricsCmd())
	rootCmd.AddCommand(newStatusCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
