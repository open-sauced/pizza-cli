package utils

import (
	"os"

	"github.com/open-sauced/pizza-cli/pkg/constants"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/term"
)

// SetupRootCommand is a convenience utility for applying templates and nice
// user experience pieces to the root cobra command
func SetupRootCommand(rootCmd *cobra.Command) {
	cobra.AddTemplateFunc("wrappedFlagUsages", wrappedFlagUsages)
	rootCmd.SetUsageTemplate(constants.UsageTemplate)
	rootCmd.SetHelpTemplate(constants.HelpTemplate)
}

// Uses the users terminal size or width of 80 if cannot determine users width
func wrappedFlagUsages(cmd *pflag.FlagSet) string {
	fd := int(os.Stdout.Fd())
	width := 80

	// Get the terminal width and dynamically set
	termWidth, _, err := term.GetSize(fd)
	if err == nil {
		width = termWidth
	}

	return cmd.FlagUsagesWrapped(width - 1)
}
