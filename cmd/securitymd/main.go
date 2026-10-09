// Command securitymd validates and generates repository SECURITY.md files.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	cmdguard "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
	"github.com/mattn/go-isatty"
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
	disableForcedColorWhenPiped()

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

// disableForcedColorWhenPiped keeps piped output ANSI-free even in ambient
// environments that force color: fang renders through colorprofile, which
// honors CLICOLOR_FORCE/TTY_FORCE for pipes and ignores NO_COLOR there, so a
// shell that exports CLICOLOR_FORCE=1 would otherwise leak escape sequences
// into the machine-readable output CI and the BuildFlow provider consume
// (pinned by TestCLI_piped_output_is_ansi_free). TTY detection must go
// through go-isatty; ModeCharDevice misclassifies /dev/null as a terminal.
func disableForcedColorWhenPiped() {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		return
	}

	for _, variable := range []string{"CLICOLOR_FORCE", "TTY_FORCE"} {
		_ = os.Unsetenv(variable)
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
