package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func newStatusCmd() *cobra.Command {
	cmd := newCommand(
		"status",
		"Show security policy status",
		"Show overview of current security policies,\navailable templates, and compliance status.",
		runStatus,
	)
	return cmd
}

func runStatus(cmd *cobra.Command, args []string) error {
	color.White("📋 Security Policy Status")
	color.White("========================")

	if _, err := os.Stat("SECURITY.md"); err == nil {
		color.Green("✅ SECURITY.md exists")
	} else {
		color.Red("❌ SECURITY.md not found")
		fmt.Println()
		color.Yellow("Run 'template-security setup' to create one")
	}

	templatesDir := "templates"
	if count := countSecurityTemplates(templatesDir); count > 0 {
		printInfo("Available SECURITY.md templates: %d", count)
	} else {
		printWarning("No templates directory found")
	}

	fmt.Println()
	printInfo("Next steps:")
	fmt.Println("  1. Run 'template-security setup --type github' to generate SECURITY.md")
	fmt.Println("  2. Run 'template-security validate' to check SECURITY.md")

	return nil
}

func countSecurityTemplates(dir string) int {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}

	count := 0

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if len(name) <= 3 || name[len(name)-3:] != ".md" {
			continue
		}
		lowerName := strings.ToLower(name)
		if strings.Contains(lowerName, "security") || strings.Contains(lowerName, "github") {
			count++
		}
	}

	return count
}
