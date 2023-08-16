package api

import "net/http"

type Client struct {
	// The configured http client for making API requests
	HTTPClient *http.Client

	// The API endpoint to use when making requests
	// Example: https://api.opensauced.pizza or https://beta.api.opensauced.pizza
	Endpoint string
}

// NewClient creates a new OpenSauced API client for making http requests
func NewClient(endpoint string) *Client {
	return &Client{
		HTTPClient: &http.Client{},
		Endpoint:   endpoint,
	}
}
