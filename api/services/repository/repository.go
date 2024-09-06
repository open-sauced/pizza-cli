package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// Service is used to access the "v2/repos" endpoints and services
type Service struct {
	httpClient *http.Client
	endpoint   string
}

// NewRepositoryService returns a new repository Service
func NewRepositoryService(httpClient *http.Client, endpoint string) *Service {
	return &Service{
		httpClient: httpClient,
		endpoint:   endpoint,
	}
}

// FindOneByOwnerAndRepo calls the "v2/repos/:owner/:name" endpoint
func (rs *Service) FindOneByOwnerAndRepo(owner string, repo string) (*DbRepository, *http.Response, error) {
	url := fmt.Sprintf("%s/v2/repos/%s/%s", rs.endpoint, owner, repo)

	resp, err := rs.httpClient.Get(url)
	if err != nil {
		return nil, resp, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, resp, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var repository DbRepository
	if err := json.NewDecoder(resp.Body).Decode(&repository); err != nil {
		return nil, resp, fmt.Errorf("error decoding response: %v", err)
	}

	return &repository, resp, nil
}

// FindContributorsByOwnerAndRepo calls the "v2/repos/:owner/:name/contributors" endpoint
func (rs *Service) FindContributorsByOwnerAndRepo(owner string, repo string, rangeVal int) (*ContributorsResponse, *http.Response, error) {
	baseURL := fmt.Sprintf("%s/v2/repos/%s/%s/contributors", rs.endpoint, owner, repo)

	// Create URL with query parameters
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing URL: %v", err)
	}

	q := u.Query()
	q.Set("range", strconv.Itoa(rangeVal))
	u.RawQuery = q.Encode()

	resp, err := rs.httpClient.Get(u.String())
	if err != nil {
		return nil, resp, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, err
	}

	var contributorsResp ContributorsResponse
	err = json.Unmarshal(body, &contributorsResp)
	if err != nil {
		return nil, resp, err
	}

	return &contributorsResp, resp, nil
}
