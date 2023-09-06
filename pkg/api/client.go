package api

import (
	"net/http"

	"github.com/open-sauced/go-api/client"
)

type Client struct {
	// The configured http client for making API requests
	HTTPClient *http.Client

	// The API endpoint to use when making requests
	// Example: https://api.opensauced.pizza
	Endpoint string
}

// NewClient creates a new OpenSauced API client for making http requests
func NewClient(endpoint string) *Client {
	return &Client{
		HTTPClient: &http.Client{},
		Endpoint:   endpoint,
	}
}

func NewGoClient(endpoint string) *client.APIClient {
	configuration := client.NewConfiguration()
	configuration.Servers = client.ServerConfigurations{
		{
			URL: endpoint,
		},
	}
	return client.NewAPIClient(configuration)
}
