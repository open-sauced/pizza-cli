package workspaces

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/open-sauced/pizza-cli/api/mock"
	"github.com/open-sauced/pizza-cli/api/services"
)

func TestGetWorkspaces(t *testing.T) {
	t.Parallel()
	m := mock.NewMockRoundTripper(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "https://api.example.com/v2/workspaces?limit=30&page=1", req.URL.String())
		assert.Equal(t, "GET", req.Method)

		mockResponse := DbWorkspacesResponse{
			Data: []DbWorkspace{
				{
					ID:   "abc123",
					Name: "workspace1",
				},
				{
					ID:   "xyz987",
					Name: "workspace2",
				},
			},
			Meta: services.MetaData{
				Page:            1,
				Limit:           30,
				ItemCount:       2,
				PageCount:       1,
				HasPreviousPage: false,
				HasNextPage:     false,
			},
		}

		// Convert the mock response to JSON
		responseBody, _ := json.Marshal(mockResponse)

		// Return the mock response
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBuffer(responseBody)),
		}, nil
	})

	client := &http.Client{Transport: m}
	service := NewWorkspacesService(client, "https://api.example.com")

	workspaces, resp, err := service.GetWorkspaces("token", 1, 30)

	assert.NoError(t, err)
	assert.NotNil(t, workspaces)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, workspaces.Data, 2)

	// First workspace
	assert.Equal(t, "abc123", workspaces.Data[0].ID)
	assert.Equal(t, "workspace1", workspaces.Data[0].Name)

	// Second workspace
	assert.Equal(t, "xyz987", workspaces.Data[1].ID)
	assert.Equal(t, "workspace2", workspaces.Data[1].Name)

	// Check the meta information
	assert.Equal(t, 1, workspaces.Meta.Page)
	assert.Equal(t, 30, workspaces.Meta.Limit)
	assert.Equal(t, 2, workspaces.Meta.ItemCount)
	assert.Equal(t, 1, workspaces.Meta.PageCount)
	assert.False(t, workspaces.Meta.HasPreviousPage)
	assert.False(t, workspaces.Meta.HasNextPage)
}

func TestCreateWorkspaceForUser(t *testing.T) {
	t.Parallel()
	m := mock.NewMockRoundTripper(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "https://api.example.com/v2/workspaces", req.URL.String())
		assert.Equal(t, "POST", req.Method)

		mockResponse := DbWorkspace{
			ID:   "abc123",
			Name: "workspace1",
		}

		// Convert the mock response to JSON
		responseBody, _ := json.Marshal(mockResponse)

		// Return the mock response
		return &http.Response{
			StatusCode: http.StatusCreated,
			Body:       io.NopCloser(bytes.NewBuffer(responseBody)),
		}, nil
	})

	client := &http.Client{Transport: m}
	service := NewWorkspacesService(client, "https://api.example.com")

	workspace, resp, err := service.CreateWorkspaceForUser("token", "test workspace", "a workspace for testing", []string{})

	assert.NoError(t, err)
	assert.NotNil(t, workspace)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "abc123", workspace.ID)
	assert.Equal(t, "workspace1", workspace.Name)
}
