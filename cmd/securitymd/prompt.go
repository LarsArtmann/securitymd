package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// errEmptyPromptAnswer marks an unusable interactive answer; the message
// points at the non-interactive flags.
var errEmptyPromptAnswer = errors.New("no answer given")

// promptIdentity asks for the GitHub coordinates when no git remote exists.
// Both answers are required: an incomplete identity renders a broken policy.
func promptIdentity(reader io.Reader, writer io.Writer) (org string, repo string, err error) {
	in := bufio.NewReader(reader)

	org, err = promptValue(in, writer, "Organization (e.g. AcmeCorp)")
	if err != nil {
		return "", "", err
	}

	repo, err = promptValue(in, writer, "Repository (e.g. widget)")
	if err != nil {
		return "", "", err
	}

	return org, repo, nil
}

func promptValue(in *bufio.Reader, writer io.Writer, label string) (string, error) {
	fmt.Fprintf(writer, "%s: ", label)

	line, err := in.ReadString('\n')
	if err != nil && line == "" {
		return "", fmt.Errorf("%s: %w", label, err)
	}

	answer := strings.TrimSpace(line)
	if answer == "" {
		return "", fmt.Errorf("%s: %w (pass --organization and --repository instead)", label, errEmptyPromptAnswer)
	}

	return answer, nil
}

// stdinIsInteractive reports whether stdin is a terminal, so interactive
// prompts only run for humans — never in CI or piped contexts.
func stdinIsInteractive() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	return stat.Mode()&os.ModeCharDevice != 0
}
