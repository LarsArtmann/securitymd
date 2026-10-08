package provider

import (
	"os"
	"os/exec"
	"path/filepath"
)

func runGit(dir string, args ...string) error {
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	cmd.Env = append(os.Environ(), "GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_SYSTEM=/dev/null")

	if out, err := cmd.CombinedOutput(); err != nil {
		return &execError{cmd: cmd.String(), output: string(out), err: err}
	}

	return nil
}

type execError struct {
	cmd    string
	output string
	err    error
}

func (e *execError) Error() string {
	return e.cmd + ": " + e.err.Error() + "\n" + e.output
}

func writeFile(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(content), 0o600)
}

func mustRead(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(content)
}
