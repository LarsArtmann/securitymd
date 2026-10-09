// Command securitymd validates and generates repository SECURITY.md files.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	cmdguard "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
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

// appConfig carries no global settings today; every knob is a per-command
// flag. It exists as cmdguard's typed config slot so fleet-standard CLI
// plumbing (DI scope, validation, lifecycle) has a place to grow into.
type appConfig struct{}

func main() {
	cli, err := cmdguard.NewCLI[appConfig](
		"securitymd",
		"Validate and generate SECURITY.md files",
		appConfig{},
		cmdguard.WithCLILong(`securitymd validates and generates SECURITY.md files.

Detects a missing policy, checks an existing one for the sections GitHub
expects (vulnerability reporting, supported versions, security practices,
contact channel, response commitments), and generates a compliant policy
from the embedded template — never overwriting an existing file.`),
		cmdguard.WithCLIVersion(fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date)),
		cmdguard.WithSignalHandling(),
		cmdguard.WithMiddleware[appConfig](cmdguard.RecoveryMiddleware[appConfig]()),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(exitOperational)
	}

	if err := registerCommands(cli); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(exitOperational)
	}

	// cmdguard prints every execution error exactly once (styled on a TTY,
	// plain when piped); the returned error exists for exit-code mapping only.
	if err := cli.Execute(context.Background()); err != nil {
		os.Exit(exitCodeFor(err))
	}
}

func registerCommands(cli *cmdguard.CLI[appConfig]) error {
	for _, register := range []func(*cmdguard.CLI[appConfig]) error{
		registerSetupCmd,
		registerValidateCmd,
		registerStatusCmd,
	} {
		if err := register(cli); err != nil {
			return err
		}
	}

	return nil
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
