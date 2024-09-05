package userlists

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Service is used to access the "v2/workspaces/:workspaceId/userLists"
// endpoints and services
type Service struct {
	httpClient *http.Client
	endpoint   string
}

// NewService returns a new UserListsService
func NewService(httpClient *http.Client, endpoint string) *Service {
	return &Service{
		httpClient: httpClient,
		endpoint:   endpoint,
	}
}

// GetUserLists calls the "GET v2/workspaces/:workspaceId/userLists" endpoint
// for the authenticated user
func (uss *Service) GetUserLists(token string, workspaceID string, page, limit int) (*GetUserListsResponse, *http.Response, error) {
	baseURL := fmt.Sprintf("%s/v2/workspaces/%s/userLists", uss.endpoint, workspaceID)

	// Create URL with query parameters
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing URL: %v", err)
	}

	q := u.Query()
	q.Set("page", fmt.Sprintf("%d", page))
	q.Set("limit", fmt.Sprintf("%d", limit))
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := uss.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, resp, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var userListsResp GetUserListsResponse
	if err := json.NewDecoder(resp.Body).Decode(&userListsResp); err != nil {
		return nil, resp, fmt.Errorf("error decoding response: %w", err)
	}

	return &userListsResp, resp, nil
}

// GetUserList calls the "GET v2/workspaces/:workspaceId/userLists" endpoint
// for the authenticated user
func (uss *Service) GetUserList(token string, workspaceID string, userlistID string) (*DbUserList, *http.Response, error) {
	url := fmt.Sprintf("%s/v2/workspaces/%s/userLists/%s", uss.endpoint, workspaceID, userlistID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := uss.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, resp, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var userList DbUserList
	if err := json.NewDecoder(resp.Body).Decode(&userList); err != nil {
		return nil, resp, fmt.Errorf("error decoding response: %w", err)
	}

	return &userList, resp, nil
}

// CreateUserListForUser calls the "POST v2/workspaces/:workspaceId/userLists" endpoint
// for the authenticated user
func (uss *Service) CreateUserListForUser(token string, workspaceID string, name string, logins []string) (*CreateUserListResponse, *http.Response, error) {
	url := fmt.Sprintf("%s/v2/workspaces/%s/userLists", uss.endpoint, workspaceID)

	loginReqs := []CreateUserListRequestContributor{}
	for _, login := range logins {
		loginReqs = append(loginReqs, CreateUserListRequestContributor{Login: login})
	}

	req := CreatePatchUserListRequest{
		Name:         name,
		IsPublic:     false,
		Contributors: loginReqs,
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, nil, fmt.Errorf("error marshaling request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, nil, fmt.Errorf("error creating request: %w", err)
	}

	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := uss.httpClient.Do(httpReq)
	if err != nil {
		return nil, resp, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, resp, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var createdUserList CreateUserListResponse
	if err := json.NewDecoder(resp.Body).Decode(&createdUserList); err != nil {
		return nil, resp, fmt.Errorf("error decoding response: %w", err)
	}

	return &createdUserList, resp, nil
}

// CreateUserListForUser calls the "PATCH v2/lists/:listId" endpoint
// for the authenticated user
func (uss *Service) PatchUserListForUser(token string, workspaceID string, listID string, name string, logins []string) (*DbUserList, *http.Response, error) {
	url := fmt.Sprintf("%s/v2/workspaces/%s/userLists/%s", uss.endpoint, workspaceID, listID)

	loginReqs := []CreateUserListRequestContributor{}
	for _, login := range logins {
		loginReqs = append(loginReqs, CreateUserListRequestContributor{Login: login})
	}

	req := CreatePatchUserListRequest{
		Name:         name,
		IsPublic:     false,
		Contributors: loginReqs,
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, nil, fmt.Errorf("error marshaling request: %w", err)
	}

	httpReq, err := http.NewRequest("PATCH", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, nil, fmt.Errorf("error creating request: %w", err)
	}

	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := uss.httpClient.Do(httpReq)
	if err != nil {
		return nil, resp, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, resp, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var createdUserList DbUserList
	if err := json.NewDecoder(resp.Body).Decode(&createdUserList); err != nil {
		return nil, resp, fmt.Errorf("error decoding response: %w", err)
	}

	return &createdUserList, resp, nil
}
