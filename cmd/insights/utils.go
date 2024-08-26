package insights

import (
	"context"
	"fmt"
	"net/http"

	"github.com/open-sauced/go-api/client"

	"github.com/open-sauced/pizza-cli/pkg/utils"
)

// findRepositoryByOwnerAndRepoName returns an API client Db Repo
// based on the given repository URL
func findRepositoryByOwnerAndRepoName(ctx context.Context, apiClient *client.APIClient, repoURL string) (*client.DbRepo, error) {
	owner, repoName, err := utils.GetOwnerAndRepoFromURL(repoURL)
	if err != nil {
		return nil, fmt.Errorf("could not extract owner and repo from url: %w", err)
	}

	repo, response, err := apiClient.RepositoryServiceAPI.FindOneByOwnerAndRepo(ctx, owner, repoName).Execute()
	if err != nil {
		if response != nil && response.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("repository %s is either non-existent, private, or has not been indexed yet", repoURL)
		}
		return nil, fmt.Errorf("error while calling 'RepositoryServiceAPI.FindOneByOwnerAndRepo' with owner %q and repo %q: %w", owner, repoName, err)
	}

	return repo, nil
}
