// Package root initiates and bootstraps the pizza CLI root Cobra command
package root

import (
	"github.com/spf13/cobra"

	"github.com/open-sauced/pizza-cli/cmd/bake"
	"github.com/open-sauced/pizza-cli/cmd/repo-query"
)

// NewRootCommand bootstraps a new root cobra command for the pizza CLI
func NewRootCommand() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "pizza <command> <subcommand> [flags]",
		Short: "OpenSauced CLI",
		Long:  `A command line utility for insights, metrics, and all things OpenSauced`,
		RunE:  run,
	}

	cmd.AddCommand(bake.NewBakeCommand())
	cmd.AddCommand(repoQuery.NewRepoQueryCommand())

	return cmd, nil
}

func run(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}
