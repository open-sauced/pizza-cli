// Package root initiates and bootstraps the pizza CLI root Cobra command
package root

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/open-sauced/pizza-cli/cmd/auth"
	"github.com/open-sauced/pizza-cli/cmd/generate"
	"github.com/open-sauced/pizza-cli/cmd/insights"
	"github.com/open-sauced/pizza-cli/cmd/version"
	"github.com/open-sauced/pizza-cli/pkg/constants"
)

// NewRootCommand bootstraps a new root cobra command for the pizza CLI
func NewRootCommand() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "pizza <command> <subcommand> [flags]",
		Short: "OpenSauced CLI",
		Long:  "A command line utility for insights, metrics, and all things OpenSauced",
		RunE:  run,
		Args: func(cmd *cobra.Command, _ []string) error {
			betaFlag := cmd.Flags().Lookup(constants.FlagNameBeta)
			if betaFlag.Changed {
				err := cmd.Flags().Lookup(constants.FlagNameEndpoint).Value.Set(constants.EndpointBeta)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}

	cmd.PersistentFlags().StringP(constants.FlagNameEndpoint, "e", constants.EndpointProd, "The API endpoint to send requests to")
	cmd.PersistentFlags().Bool(constants.FlagNameBeta, false, fmt.Sprintf("Shorthand for using the beta OpenSauced API endpoint (\"%s\"). Supersedes the '--%s' flag", constants.EndpointBeta, constants.FlagNameEndpoint))
	cmd.PersistentFlags().Bool(constants.FlagNameTelemetry, false, "Disable sending telemetry data to OpenSauced")
	cmd.PersistentFlags().StringP("config", "c", "~/.sauced.yaml", "The saucectl config")
	cmd.PersistentFlags().StringP("log-level", "l", "info", "The logging level. Options: error, warn, info, debug")
	cmd.PersistentFlags().Bool("tty-disable", false, "Disable log stylization. Suitable for CI/CD and automation")

	cmd.AddCommand(auth.NewLoginCommand())
	cmd.AddCommand(generate.NewGenerateCommand())
	cmd.AddCommand(insights.NewInsightsCommand())
	cmd.AddCommand(version.NewVersionCommand())

	return cmd, nil
}

func run(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}
