// Package show contains the bootstrapping and tooling for the pizza show
// cobra command
package show

import (
	"context"
	"errors"

	client "github.com/open-sauced/go-api/client"
	"github.com/open-sauced/pizza-cli/pkg/api"
	"github.com/spf13/cobra"
)

// Options are the options for the pizza show command including user
// defined configurations
type Options struct {
	// RepoName is the git repo name to get statistics
	RepoName string

	// Page is the page to be requested
	Page int

	// Limit is the number of records to be retrieved
	Limit int

	// Range is the number of days to take into account when retrieving statistics
	Range int

	// APIClient is the api client to interface with open sauced api
	APIClient *client.APIClient

	ServerContext context.Context
}

const showLongDesc string = `WARNING: Proof of concept feature.

The show command accepts the name of a git repository and uses OpenSauced api
to retrieve metrics of the repository to be displayed as a TUI.`

// NewShowCommand returns a new cobra command for 'pizza show'
func NewShowCommand() *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "show repository-name [flags]",
		Short: "Get visual metrics of a repository",
		Long:  showLongDesc,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("must specify the URL of a git repository to analyze")
			}
			opts.RepoName = args[0]

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var endpoint string
			customEndpoint, _ := cmd.Flags().GetString("endpoint")
			if customEndpoint != "" {
				endpoint = customEndpoint
			}

			useBeta, _ := cmd.Flags().GetBool("beta")
			if useBeta {
				endpoint = api.BetaAPIEndpoint
			}

			opts.APIClient = api.NewGoClient(endpoint)
			opts.ServerContext = context.TODO()
			return run(opts)
		},
	}

	cmd.Flags().IntVarP(&opts.Range, "range", "r", 30, "The last N number of days to consider")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "l", 10, "The number of records to retrieve")
	cmd.Flags().IntVarP(&opts.Page, "page", "p", 1, "The page number to retrieve")

	return cmd
}

func run(opts *Options) error {
	// Load the pizza TUI
	err := pizzaTUI(opts)
	return err
}
