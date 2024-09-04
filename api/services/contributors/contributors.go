package contributors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Service is the Contributors service used for accessing the "v2/contributors"
// endpoint and API services
type Service struct {
	httpClient *http.Client
	endpoint   string
}

// NewContributorsService returns a new contributors Service
func NewContributorsService(httpClient *http.Client, endpoint string) *Service {
	return &Service{
		httpClient: httpClient,
		endpoint:   endpoint,
	}
}

// NewPullRequestContributors calls the "v2/contributors/insights/new" API endpoint
func (s *Service) NewPullRequestContributors(repos []string, rangeVal int) (*ContribResponse, *http.Response, error) {
	baseURL := fmt.Sprintf("%s/v2/contributors/insights/new", s.endpoint)

	// Create URL with query parameters
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing URL: %v", err)
	}

	q := u.Query()
	q.Set("range", fmt.Sprintf("%d", rangeVal))
	q.Set("repos", strings.Join(repos, ","))
	u.RawQuery = q.Encode()

	resp, err := s.httpClient.Get(u.String())
	if err != nil {
		return nil, resp, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, resp, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var newContributorsResponse ContribResponse
	if err := json.NewDecoder(resp.Body).Decode(&newContributorsResponse); err != nil {
		return nil, resp, fmt.Errorf("error decoding response: %v", err)
	}

	return &newContributorsResponse, resp, nil
}

// RecentPullRequestContributors calls the "v2/contributors/insights/recent" API endpoint
func (s *Service) RecentPullRequestContributors(repos []string, rangeVal int) (*ContribResponse, *http.Response, error) {
	baseURL := fmt.Sprintf("%s/v2/contributors/insights/recent", s.endpoint)

	// Create URL with query parameters
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing URL: %v", err)
	}

	q := u.Query()
	q.Set("range", fmt.Sprintf("%d", rangeVal))
	q.Set("repos", strings.Join(repos, ","))
	u.RawQuery = q.Encode()

	resp, err := s.httpClient.Get(u.String())
	if err != nil {
		return nil, resp, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, resp, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var recentContributorsResponse ContribResponse
	if err := json.NewDecoder(resp.Body).Decode(&recentContributorsResponse); err != nil {
		return nil, resp, fmt.Errorf("error decoding response: %v", err)
	}

	return &recentContributorsResponse, resp, nil
}

// AlumniPullRequestContributors calls the "v2/contributors/insights/alumni" API endpoint
func (s *Service) AlumniPullRequestContributors(repos []string, rangeVal int) (*ContribResponse, *http.Response, error) {
	baseURL := fmt.Sprintf("%s/v2/contributors/insights/alumni", s.endpoint)

	// Create URL with query parameters
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing URL: %v", err)
	}

	q := u.Query()
	q.Set("range", fmt.Sprintf("%d", rangeVal))
	q.Set("repos", strings.Join(repos, ","))
	u.RawQuery = q.Encode()

	resp, err := s.httpClient.Get(u.String())
	if err != nil {
		return nil, resp, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, resp, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var alumniContributorsResponse ContribResponse
	if err := json.NewDecoder(resp.Body).Decode(&alumniContributorsResponse); err != nil {
		return nil, resp, fmt.Errorf("error decoding response: %v", err)
	}

	return &alumniContributorsResponse, resp, nil
}

// RepeatPullRequestContributors calls the "v2/contributors/insights/repeat" API endpoint
func (s *Service) RepeatPullRequestContributors(repos []string, rangeVal int) (*ContribResponse, *http.Response, error) {
	baseURL := fmt.Sprintf("%s/v2/contributors/insights/repeat", s.endpoint)

	// Create URL with query parameters
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing URL: %v", err)
	}

	q := u.Query()
	q.Set("range", fmt.Sprintf("%d", rangeVal))
	q.Set("repos", strings.Join(repos, ","))
	u.RawQuery = q.Encode()

	resp, err := s.httpClient.Get(u.String())
	if err != nil {
		return nil, resp, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, resp, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var repeatContributorsResponse ContribResponse
	if err := json.NewDecoder(resp.Body).Decode(&repeatContributorsResponse); err != nil {
		return nil, resp, fmt.Errorf("error decoding response: %v", err)
	}

	return &repeatContributorsResponse, resp, nil
}

// SearchPullRequestContributors calls the "v2/contributors/search"
func (s *Service) SearchPullRequestContributors(repos []string, rangeVal int) (*ContribResponse, *http.Response, error) {
	baseURL := fmt.Sprintf("%s/v2/contributors/search", s.endpoint)

	// Create URL with query parameters
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing URL: %v", err)
	}

	q := u.Query()
	q.Set("range", fmt.Sprintf("%d", rangeVal))
	q.Set("repos", strings.Join(repos, ","))
	u.RawQuery = q.Encode()

	resp, err := s.httpClient.Get(u.String())
	if err != nil {
		return nil, resp, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, resp, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var searchContributorsResponse ContribResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchContributorsResponse); err != nil {
		return nil, resp, fmt.Errorf("error decoding response: %v", err)
	}

	return &searchContributorsResponse, resp, nil
}
