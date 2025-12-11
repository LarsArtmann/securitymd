package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func newStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show security policy status",
		Long: `Show overview of current security policies,
available templates, and compliance status.`,
		RunE: runStatus,
	}

	return cmd
}

func runStatus(cmd *cobra.Command, args []string) error {
	color.White("📋 Security Policy Status")
	color.White("========================")

	// Check for existing policies
	if _, err := os.Stat("SECURITY.md"); err == nil {
		color.Green("✅ SECURITY.md exists")
	} else {
		color.Red("❌ SECURITY.md not found")
	}

	if _, err := os.Stat("incident-response.md"); err == nil {
		color.Green("✅ Incident response plan exists")
	} else {
		color.Red("❌ Incident response plan not found")
	}

	// Count available templates
	templatesDir := "templates"
	if count := countTemplates(templatesDir); count > 0 {
		color.Cyan("Available templates: %d", count)
	} else {
		color.Yellow("No templates directory found")
	}

	fmt.Println()
	color.Cyan("Next steps:")
	fmt.Println("  1. Run 'template-security setup' to generate policies")
	fmt.Println("  2. Run 'template-security validate' to check policies")
	fmt.Println("  3. Run 'template-security compliance --framework gdpr' for compliance")

	return nil
}

func countTemplates(dir string) int {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}

	count := 0
	for _, entry := range entries {
		if !entry.IsDir() && len(entry.Name()) > 3 && entry.Name()[len(entry.Name())-3:] == ".md" {
			count++
		}
	}
	return count
}
