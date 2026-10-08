// Command securitymd validates and generates repository SECURITY.md files.
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
		Use:   "securitymd",
		Short: "Validate and generate SECURITY.md files",
		Long: `securitymd validates and generates SECURITY.md files.

Detects a missing policy, checks an existing one for the sections GitHub
expects (vulnerability reporting, supported versions, security practices,
contact channel, response commitments), and generates a compliant policy
from the embedded template — never overwriting an existing file.`,
		Version:           fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
		DisableAutoGenTag: true,
		SilenceUsage:      true,
	}

	rootCmd.AddCommand(newSetupCmd())
	rootCmd.AddCommand(newValidateCmd())
	rootCmd.AddCommand(newStatusCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
