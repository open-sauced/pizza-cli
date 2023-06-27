package utils

import "github.com/spf13/cobra"

// SetupRootCommand is a convinence utility for applying templates and nice
// user experience pieces to the root cobra command
func SetupRootCommand(rootCmd *cobra.Command) {
	cobra.AddTemplateFunc("wrappedFlagUsages", wrappedFlagUsages)

	rootCmd.SetUsageTemplate(usageTemplate)
	rootCmd.SetHelpTemplate(helpTemplate)
}
