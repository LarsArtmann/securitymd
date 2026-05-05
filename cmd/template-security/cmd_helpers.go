// Package main provides CLI command helpers for template-security.
package main

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type RunEFunc func(*cobra.Command, []string) error

func newCommand(use, short, long string, runE RunEFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		RunE:  runE,
	}

	return cmd
}

func printInfof(format string, args ...any) {
	color.Cyan(format, args...)
}

func printWarningf(format string, args ...any) {
	color.Yellow(format, args...)
}
