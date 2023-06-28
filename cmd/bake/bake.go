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

	"github.com/spf13/cobra"
)

// Options are the options for the pizza bake command including user
// defined configurations
type Options struct {
	// Endpoint is the service endpoint to reach out to
	Endpoint string

	// URL is the git repo URL that will be sourced via 'pizza bake'
	URL string
}

const bakeLongDesc string = `WARNING: Proof of concept feature.

The bake command takes a URL to a git repository and uses a pizza-oven service
to source those commits. These commits will then be used for insights on OpenSauced.`

// NewBakeCommand returns a new cobra command for 'pizza bake'
func NewBakeCommand() *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "bake url [flags]",
		Short: "Use a pizza-oven to source git commits into OpenSauced",
		Long:  bakeLongDesc,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return errors.New("only a single url can be ingested at a time")
			}
			if len(args) == 0 {
				return errors.New("must specify the URL of a git repository")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.URL = args[0]
			return run(opts)
		},
	}

	// TODO - this will need to be a live service URL by default.
	// For now, localhost is fine.
	cmd.Flags().StringVarP(&opts.Endpoint, "endpoint", "e", "http://localhost:8080", "The endpoint to send requests to")

	return cmd
}

type bakePostRequest struct {
	URL string `json:"url"`
}

func run(opts *Options) error {
	bodyPostReq := &bakePostRequest{
		URL: opts.URL,
	}

	bodyPostJSON, err := json.Marshal(bodyPostReq)
	if err != nil {
		return err
	}

	responseBody := bytes.NewBuffer(bodyPostJSON)
	resp, err := http.Post(fmt.Sprintf("%s/bake", opts.Endpoint), "application/json", responseBody)
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
