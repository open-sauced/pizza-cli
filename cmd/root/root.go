// Package root initiates and bootstraps the pizza CLI root Cobra command
package root

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/open-sauced/pizza-cli/cmd/auth"
	"github.com/open-sauced/pizza-cli/cmd/bake"
	repoquery "github.com/open-sauced/pizza-cli/cmd/repo-query"
	"github.com/open-sauced/pizza-cli/pkg/api"
)

// NewRootCommand bootstraps a new root cobra command for the pizza CLI
func NewRootCommand() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "pizza <command> <subcommand> [flags]",
		Short: "OpenSauced CLI",
		Long:  `A command line utility for insights, metrics, and all things OpenSauced`,
		RunE:  run,
	}

	cmd.PersistentFlags().StringP("endpoint", "e", api.APIEndpoint, "The API endpoint to send requests to")
	cmd.PersistentFlags().Bool("beta", false, fmt.Sprintf("Shorthand for using the beta OpenSauced API endpoint (\"%s\"). Superceds the '--endpoint' flag", api.BetaAPIEndpoint))
	cmd.PersistentFlags().Bool("disable-telemetry", false, "Disable sending telemetry data to OpenSauced")

	cmd.AddCommand(bake.NewBakeCommand())
	cmd.AddCommand(repoquery.NewRepoQueryCommand())
	cmd.AddCommand(auth.NewLoginCommand())

	return cmd, nil
}

func run(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}
