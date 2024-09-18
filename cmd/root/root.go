// Package root initiates and bootstraps the pizza CLI root Cobra command
package root

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/open-sauced/pizza-cli/v2/cmd/auth"
	"github.com/open-sauced/pizza-cli/v2/cmd/docs"
	"github.com/open-sauced/pizza-cli/v2/cmd/generate"
	"github.com/open-sauced/pizza-cli/v2/cmd/insights"
	"github.com/open-sauced/pizza-cli/v2/cmd/version"
	"github.com/open-sauced/pizza-cli/v2/pkg/constants"
)

// NewRootCommand bootstraps a new root cobra command for the pizza CLI
func NewRootCommand() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "pizza <command> <subcommand> [flags]",
		Short: "OpenSauced CLI",
		Long:  "A command line utility for insights, metrics, and generating CODEOWNERS documentation for your open source projects",
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
	cmd.PersistentFlags().StringP("config", "c", "", "The codeowners config")
	cmd.PersistentFlags().StringP("log-level", "l", "info", "The logging level. Options: error, warn, info, debug")
	cmd.PersistentFlags().Bool("tty-disable", false, "Disable log stylization. Suitable for CI/CD and automation")

	cmd.AddCommand(auth.NewLoginCommand())
	cmd.AddCommand(generate.NewGenerateCommand())
	cmd.AddCommand(insights.NewInsightsCommand())
	cmd.AddCommand(version.NewVersionCommand())

	// The docs command is hidden as it's only used by the pizza-cli maintainers
	docsCmd := docs.NewDocsCommand()
	docsCmd.Hidden = true
	cmd.AddCommand(docsCmd)

	err := cmd.PersistentFlags().MarkHidden(constants.FlagNameEndpoint)
	if err != nil {
		return nil, fmt.Errorf("error marking %s as hidden: %w", constants.FlagNameEndpoint, err)
	}

	err = cmd.PersistentFlags().MarkHidden(constants.FlagNameBeta)
	if err != nil {
		return nil, fmt.Errorf("error marking %s as hidden: %w", constants.FlagNameBeta, err)
	}

	return cmd, nil
}

func run(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}
