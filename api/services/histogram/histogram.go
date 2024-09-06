package histogram

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// Service is used to access the API "v2/histogram" endpoints and services
type Service struct {
	httpClient *http.Client
	endpoint   string
}

// NewHistogramService returns a new histogram Service
func NewHistogramService(httpClient *http.Client, endpoint string) *Service {
	return &Service{
		httpClient: httpClient,
		endpoint:   endpoint,
	}
}

// PrsHistogram calls the "v2/histogram/pull-requests" endpoints
func (s *Service) PrsHistogram(repo string, rangeVal int) ([]PrHistogramData, *http.Response, error) {
	baseURL := s.endpoint + "/v2/histogram/pull-requests"

	// Create URL with query parameters
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing URL: %v", err)
	}

	q := u.Query()
	q.Set("range", strconv.Itoa(rangeVal))
	q.Set("repo", repo)
	u.RawQuery = q.Encode()

	resp, err := s.httpClient.Get(u.String())
	if err != nil {
		return nil, resp, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, resp, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var prHistogramData []PrHistogramData
	if err := json.NewDecoder(resp.Body).Decode(&prHistogramData); err != nil {
		return nil, resp, fmt.Errorf("error decoding response: %v", err)
	}

	return prHistogramData, resp, nil
}
