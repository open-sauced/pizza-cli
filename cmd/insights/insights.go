package insights

import (
	"github.com/spf13/cobra"
)

// NewInsightsCommand returns a new cobra command for 'pizza insights'
func NewInsightsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "insights <command> [flags]",
		Short: "Gather insights about git contributors, repositories, users and pull requests",
		Long:  "Gather insights about git contributors, repositories, user and pull requests and display the results",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(NewContributorsCommand())
	return cmd
}
