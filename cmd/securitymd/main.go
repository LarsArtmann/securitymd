// Command securitymd validates and generates repository SECURITY.md files.
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Exit codes per README's contract: findings are the policy failing, not the
// tool crashing, and CI keys its decisions on the distinction.
const (
	exitOK          = 0
	exitFindings    = 1
	exitOperational = 2
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
		os.Exit(exitCodeFor(err))
	}
}

// exitCodeFor implements the documented contract: error-severity findings
// exit 1 (the tool worked; the policy failed validation), every other error
// exits 2 (the tool could not run: bad flags, IO failures, cancellations).
func exitCodeFor(err error) int {
	if errors.Is(err, errPolicyFindings) {
		return exitFindings
	}

	return exitOperational
}
