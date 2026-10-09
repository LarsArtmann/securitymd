package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestZZZEnvProbe(t *testing.T) {
	out := filepath.Join(os.TempDir(), "securitymd-buildflow-env.txt")
	content := "IS_TTY_STDOUT=false\n"
	if fi, _ := os.Stdout.Stat(); fi != nil && fi.Mode()&os.ModeCharDevice != 0 {
		content = "IS_TTY_STDOUT=true\n"
	}
	content += strings.Join(os.Environ(), "\n")
	if err := os.WriteFile(out, []byte(content), 0o600); err != nil {
		t.Fatalf("write env probe: %v", err)
	}
}
