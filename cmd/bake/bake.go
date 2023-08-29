// Package bake contains the bootstrapping and tooling for the pizza bake
// cobra command
package bake

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/open-sauced/pizza-cli/pkg/api"
	"github.com/open-sauced/pizza-cli/pkg/utils"
	"github.com/spf13/cobra"

	"gopkg.in/yaml.v3"
)

// Options are the options for the pizza bake command including user
// defined configurations
type Options struct {
	// The API Client for the calls to bake git repos
	APIClient *api.Client

	// URLs are the git repo URLs that will be sourced via 'pizza bake'
	URLs []string

	// Wait defines the client choice to wait for /bake to finish processing
	Wait bool

	// FilePath is the location of the file containing a batch of repos to be baked
	FilePath string

	// telemetry for capturing CLI events
	telemetry *utils.PosthogCliClient
}

const bakeLongDesc string = `WARNING: Proof of concept feature.

The bake command accepts one or multiple URLs to a git repository and uses a pizza-oven service
to source those commits. These commits will then be used for insights on OpenSauced.`

// NewBakeCommand returns a new cobra command for 'pizza bake'
func NewBakeCommand() *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "bake url [flags]",
		Short: "Use a pizza-oven to source git commits into OpenSauced",
		Long:  bakeLongDesc,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 && opts.FilePath == "" {
				return errors.New("must specify the URL(s) of a git repository or provide a batch file")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			disableTelem, _ := cmd.Flags().GetBool("disable-telemetry")
			endpoint, _ := cmd.Flags().GetString("endpoint")
			useBeta, _ := cmd.Flags().GetBool("beta")

			if useBeta {
				fmt.Printf("Using beta API endpoint - %s\n", api.BetaAPIEndpoint)
				endpoint = api.BetaAPIEndpoint
			}
			opts.APIClient = api.NewClient(endpoint)

			opts.URLs = append(opts.URLs, args...)

			if !disableTelem {
				opts.telemetry = utils.NewPosthogCliClient()
				defer opts.telemetry.Done()

				opts.telemetry.CaptureBake(opts.URLs)
			}

			return run(opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.Wait, "wait", "w", false, "Wait for bake processing to finish")
	cmd.Flags().StringVarP(&opts.FilePath, "file", "f", "", "The yaml file containing a series of repos to batch to /bake")

	return cmd
}

type bakePostRequest struct {
	URL  string `json:"url"`
	Wait bool   `json:"wait"`
}

type repos struct {
	URLs []string `yaml:"repos"`
}

func run(opts *Options) error {
	var repos repos
	uniqueURLs := make(map[string]bool)

	if opts.FilePath != "" {
		configFile, err := os.ReadFile(opts.FilePath)
		if err != nil {
			return err
		}

		err = yaml.Unmarshal(configFile, &repos)
		if err != nil {
			return err
		}

		for _, url := range repos.URLs {
			uniqueURLs[url] = true
		}
	}

	// make sure there are no duplicated queries to the same URL
	for _, url := range opts.URLs {
		if _, ok := uniqueURLs[url]; !ok {
			uniqueURLs[url] = true
			continue
		}
		fmt.Printf("Warning: duplicated URL (%s) would not be processed again\n", url)
	}

	for url := range uniqueURLs {
		bodyPostReq := bakePostRequest{
			URL:  url,
			Wait: opts.Wait,
		}

		err := bakeRepo(bodyPostReq, opts.APIClient)
		if err != nil {
			fmt.Printf("Error: failed fetch of %s repository (%s)\n", url, err.Error())
		}
	}

	return nil
}

func bakeRepo(bodyPostReq bakePostRequest, apiClient *api.Client) error {
	bodyPostJSON, err := json.Marshal(bodyPostReq)
	if err != nil {
		return err
	}

	responseBody := bytes.NewBuffer(bodyPostJSON)
	resp, err := apiClient.HTTPClient.Post(fmt.Sprintf("%s/bake", apiClient.Endpoint), "application/json", responseBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		fmt.Printf("Resp body: %v\n", string(body))
	}

	return nil
}
