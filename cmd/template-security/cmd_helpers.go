package main

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type RunEFunc func(*cobra.Command, []string) error

func newCommand(use, short, long string, runE RunEFunc) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		RunE:  runE,
	}
}

func printInfo(format string, args ...any) {
	color.Cyan(format, args...)
}

func printSuccess(format string, args ...any) {
	color.Green(format, args...)
}

func printWarning(format string, args ...any) {
	color.Yellow(format, args...)
}

func printError(format string, args ...any) {
	color.Red(format, args...)
}
