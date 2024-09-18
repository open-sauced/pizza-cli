package api

import (
	"net/http"
	"time"

	"github.com/open-sauced/pizza-cli/v2/api/services/contributors"
	"github.com/open-sauced/pizza-cli/v2/api/services/histogram"
	"github.com/open-sauced/pizza-cli/v2/api/services/repository"
	"github.com/open-sauced/pizza-cli/v2/api/services/workspaces"
)

// Client is the API client for OpenSauced API
type Client struct {
	// API services
	RepositoryService  *repository.Service
	ContributorService *contributors.Service
	HistogramService   *histogram.Service
	WorkspacesService  *workspaces.Service

	// The configured http client for making API requests
	httpClient *http.Client

	// The API endpoint to use when making requests
	// Example: https://api.opensauced.pizza
	endpoint string
}

// NewClient returns a new API Client based on provided inputs
func NewClient(endpoint string) *Client {
	httpClient := &http.Client{
		// TODO (jpmcb): in the future, we can allow users to configure the API
		// timeout via some global flag
		Timeout: time.Second * 10,
	}

	client := Client{
		httpClient: httpClient,
		endpoint:   endpoint,
	}

	client.ContributorService = contributors.NewContributorsService(client.httpClient, client.endpoint)
	client.RepositoryService = repository.NewRepositoryService(client.httpClient, client.endpoint)
	client.HistogramService = histogram.NewHistogramService(client.httpClient, client.endpoint)
	client.WorkspacesService = workspaces.NewWorkspacesService(client.httpClient, client.endpoint)

	return &client
}
