package workspaces

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/open-sauced/pizza-cli/api/services/workspaces/userlists"
)

// Service is used to access the "v2/workspaces" endpoints and services.
// It has a child service UserListService used for accessing workspace contributor insights
type Service struct {
	UserListService *userlists.Service

	httpClient *http.Client
	endpoint   string
}

// NewWorkspacesService returns a new workspace Service
func NewWorkspacesService(httpClient *http.Client, endpoint string) *Service {
	userListService := userlists.NewService(httpClient, endpoint)

	return &Service{
		UserListService: userListService,
		httpClient:      httpClient,
		endpoint:        endpoint,
	}
}

// GetWorkspaces calls the "GET v2/workspaces" endpoint for the authenticated user
func (s *Service) GetWorkspaces(token string, page, limit int) (*DbWorkspacesResponse, *http.Response, error) {
	baseURL := s.endpoint + "/v2/workspaces"

	// Create URL with query parameters
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing URL: %v", err)
	}

	q := u.Query()
	q.Set("page", strconv.Itoa(page))
	q.Set("limit", strconv.Itoa(limit))
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, resp, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var workspacesResp DbWorkspacesResponse
	if err := json.NewDecoder(resp.Body).Decode(&workspacesResp); err != nil {
		return nil, resp, fmt.Errorf("error decoding response: %w", err)
	}

	return &workspacesResp, resp, nil
}

// CreateWorkspaceForUser calls the "POST v2/workspaces" endpoint for the authenticated user
func (s *Service) CreateWorkspaceForUser(token string, name string, description string, repos []string) (*DbWorkspace, *http.Response, error) {
	url := s.endpoint + "/v2/workspaces"

	repoReqs := []CreateWorkspaceRequestRepoInfo{}
	for _, repo := range repos {
		repoReqs = append(repoReqs, CreateWorkspaceRequestRepoInfo{FullName: repo})
	}

	req := CreateWorkspaceRequest{
		Name:         name,
		Description:  description,
		Repos:        repoReqs,
		Members:      []string{},
		Contributors: []string{},
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, nil, fmt.Errorf("error marshaling request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, nil, fmt.Errorf("error creating request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, resp, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, resp, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var createdWorkspace DbWorkspace
	if err := json.NewDecoder(resp.Body).Decode(&createdWorkspace); err != nil {
		return nil, resp, fmt.Errorf("error decoding response: %w", err)
	}

	return &createdWorkspace, resp, nil
}
