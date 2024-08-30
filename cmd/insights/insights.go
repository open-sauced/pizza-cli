package insights

import (
	"github.com/spf13/cobra"

	"github.com/open-sauced/pizza-cli/pkg/constants"
)

// NewInsightsCommand returns a new cobra command for 'pizza insights'
func NewInsightsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "insights <command> [flags]",
		Short: "Gather insights about git contributors, repositories, users and pull requests",
		Long:  "Gather insights about git contributors, repositories, user and pull requests and display the results",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
	cmd.PersistentFlags().StringP(constants.FlagNameOutput, "o", constants.OutputTable, "The formatting for command output. One of: (table, yaml, csv, json)")
	cmd.AddCommand(NewContributorsCommand())
	cmd.AddCommand(NewRepositoriesCommand())
	cmd.AddCommand(NewUserContributionsCommand())
	return cmd
}
