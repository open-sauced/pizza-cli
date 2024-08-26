// Package bake contains the bootstrapping and tooling for the pizza bake
// cobra command
package bake

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/open-sauced/go-api/client"
	"github.com/spf13/cobra"

	"github.com/open-sauced/pizza-cli/pkg/api"
	"github.com/open-sauced/pizza-cli/pkg/constants"
	"github.com/open-sauced/pizza-cli/pkg/utils"
)

// Options are the options for the pizza bake command including user
// defined configurations
type Options struct {
	// APIClient is the http client for making calls to the open-sauced api
	APIClient *client.APIClient

	// Repos is the array of git repository urls
	Repos []string

	// Wait defines the client choice to wait for /bake to finish processing
	Wait bool

	// FilePath is the path to yaml file containing an array of git repository urls
	FilePath string

	// telemetry for capturing CLI events
	telemetry *utils.PosthogCliClient
}

const bakeLongDesc string = `WARNING: Proof of concept feature.

The bake command accepts one or multiple Repos to a git repository and uses a pizza-oven service
to source those commits. These commits will then be used for insights on OpenSauced.`

// NewBakeCommand returns a new cobra command for 'pizza bake'
func NewBakeCommand() *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "bake url... [flags]",
		Short: "Use a pizza-oven to source git commits into OpenSauced",
		Long:  bakeLongDesc,
		Args: func(cmd *cobra.Command, args []string) error {
			fileFlag := cmd.Flags().Lookup(constants.FlagNameFile)
			if !fileFlag.Changed && len(args) == 0 {
				return fmt.Errorf("must specify git repository url argument(s) or provide %s flag", fileFlag.Name)
			}
			opts.Repos = append(opts.Repos, args...)
			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			endpointURL, _ := cmd.Flags().GetString(constants.FlagNameEndpoint)
			opts.APIClient = api.NewGoClient(endpointURL)
			disableTelem, _ := cmd.Flags().GetBool(constants.FlagNameTelemetry)

			if !disableTelem {
				opts.telemetry = utils.NewPosthogCliClient()
				defer opts.telemetry.Done()

				opts.telemetry.CaptureBake(opts.Repos)
			}

			return run(opts)
		},
	}
	cmd.Flags().StringVarP(&opts.FilePath, constants.FlagNameFile, "f", "", "Path to yaml file containing an array of git repository urls")
	cmd.Flags().BoolVarP(&opts.Wait, constants.FlagNameWait, "w", false, "Wait for bake processing to finish")
	return cmd
}

func run(opts *Options) error {
	repositories, err := utils.HandleRepositoryValues(opts.Repos, opts.FilePath)
	if err != nil {
		return err
	}
	var (
		waitGroup = new(sync.WaitGroup)
		errorChan = make(chan error, len(repositories))
	)
	for url := range repositories {
		waitGroup.Add(1)
		go func(repoURL string) {
			defer waitGroup.Done()
			err = bakeRepository(context.TODO(), opts.APIClient, repoURL, opts.Wait)
			if err != nil {
				errorChan <- err
				return
			}
			fmt.Println("successfully baked repository", repoURL)
		}(url)
	}
	waitGroup.Wait()
	close(errorChan)
	var allErrors error
	for err = range errorChan {
		allErrors = errors.Join(allErrors, err)
	}
	return allErrors
}

func bakeRepository(ctx context.Context, apiClient *client.APIClient, repoURL string, wait bool) error {
	body := client.BakeRepoDto{
		Url:  repoURL,
		Wait: wait,
	}
	_, err := apiClient.PizzaOvenServiceAPI.
		BakeARepositoryWithThePizzaOvenMicroservice(ctx).
		BakeRepoDto(body).
		Execute()
	if err != nil {
		return fmt.Errorf("error while calling 'PizzaOvenServiceAPI.BakeARepositoryWithThePizzaOvenMicroservice' with repository %q: %w", repoURL, err)
	}
	return nil
}
