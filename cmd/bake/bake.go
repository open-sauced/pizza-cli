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

	"github.com/spf13/cobra"

	"gopkg.in/yaml.v3"
)

// Options are the options for the pizza bake command including user
// defined configurations
type Options struct {
	// Endpoint is the service endpoint to reach out to
	Endpoint string

	// URLs are the git repo URLs that will be sourced via 'pizza bake'
	URLs []string

	// Wait defines the client choice to wait for /bake to finish processing
	Wait bool

	// FilePath is the location of the file containing a batch of repos to be baked
	FilePath string
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
			opts.URLs = append(opts.URLs, args...)
			return run(opts)
		},
	}

	// TODO - this will need to be a live service URL by default.
	// For now, localhost is fine.
	cmd.Flags().StringVarP(&opts.Endpoint, "endpoint", "e", "http://localhost:8080", "The endpoint to send requests to")
	cmd.Flags().BoolVarP(&opts.Wait, "wait", "w", false, "Wait for bake processing to finish")
	cmd.Flags().StringVarP(&opts.FilePath, "file", "f", "", "The yaml file containing a series of repos to batch to /bake")

	return cmd
}

type bakePostRequest struct {
	URL  string `json:"url"`
	Wait bool   `json:"wait,omitempty"`
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

		err := bakeRepo(bodyPostReq, opts.Endpoint)
		if err != nil {
			fmt.Printf("Error: failed fetch of %s repository (%s)\n", url, err.Error())
		}
	}

	return nil
}

func bakeRepo(bodyPostReq bakePostRequest, endpoint string) error {
	bodyPostJSON, err := json.Marshal(bodyPostReq)
	if err != nil {
		return err
	}

	requestBody := bytes.NewBuffer(bodyPostJSON)
	resp, err := http.Post(fmt.Sprintf("%s/bake", endpoint), "application/json", requestBody)

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
