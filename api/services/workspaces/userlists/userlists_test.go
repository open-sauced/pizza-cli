package userlists

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

func TestGetUserLists(t *testing.T) {
	t.Parallel()
	m := mock.NewMockRoundTripper(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "https://api.example.com/v2/workspaces/abc123/userLists?limit=30&page=1", req.URL.String())
		assert.Equal(t, "GET", req.Method)

		mockResponse := GetUserListsResponse{
			Data: []DbUserList{
				{
					ID:   "abc",
					Name: "userlist1",
				},
				{
					ID:   "xyz",
					Name: "userlist2",
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
	service := NewService(client, "https://api.example.com")

	userlists, resp, err := service.GetUserLists("token", "abc123", 1, 30)

	assert.NoError(t, err)
	assert.NotNil(t, userlists)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, userlists.Data, 2)

	// First workspace
	assert.Equal(t, "abc", userlists.Data[0].ID)
	assert.Equal(t, "userlist1", userlists.Data[0].Name)

	// Second workspace
	assert.Equal(t, "xyz", userlists.Data[1].ID)
	assert.Equal(t, "userlist2", userlists.Data[1].Name)

	// Check the meta information
	assert.Equal(t, 1, userlists.Meta.Page)
	assert.Equal(t, 30, userlists.Meta.Limit)
	assert.Equal(t, 2, userlists.Meta.ItemCount)
	assert.Equal(t, 1, userlists.Meta.PageCount)
	assert.False(t, userlists.Meta.HasPreviousPage)
	assert.False(t, userlists.Meta.HasNextPage)
}

func TestGetUserListForUser(t *testing.T) {
	t.Parallel()
	m := mock.NewMockRoundTripper(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "https://api.example.com/v2/workspaces/abc123/userLists/xyz", req.URL.String())
		assert.Equal(t, "GET", req.Method)

		mockResponse := DbUserList{
			ID:   "abc",
			Name: "userlist1",
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
	service := NewService(client, "https://api.example.com")

	userlists, resp, err := service.GetUserList("token", "abc123", "xyz")

	assert.NoError(t, err)
	assert.NotNil(t, userlists)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "abc", userlists.ID)
	assert.Equal(t, "userlist1", userlists.Name)
}

func TestCreateUserListForUser(t *testing.T) {
	t.Parallel()
	m := mock.NewMockRoundTripper(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "https://api.example.com/v2/workspaces/abc123/userLists", req.URL.String())
		assert.Equal(t, "POST", req.Method)

		mockResponse := CreateUserListResponse{
			ID:         "abc",
			UserListID: "xyz",
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
	service := NewService(client, "https://api.example.com")

	userlists, resp, err := service.CreateUserListForUser("token", "abc123", "userlist1", []string{})

	assert.NoError(t, err)
	assert.NotNil(t, userlists)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "abc", userlists.ID)
	assert.Equal(t, "xyz", userlists.UserListID)
}

func TestPatchUserListForUser(t *testing.T) {
	t.Parallel()
	m := mock.NewMockRoundTripper(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "https://api.example.com/v2/workspaces/abc123/userLists/abc", req.URL.String())
		assert.Equal(t, "PATCH", req.Method)

		mockResponse := DbUserList{
			ID:   "abc",
			Name: "userlist1",
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
	service := NewService(client, "https://api.example.com")

	userlists, resp, err := service.PatchUserListForUser("token", "abc123", "abc", "userlist1", []string{})

	assert.NoError(t, err)
	assert.NotNil(t, userlists)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "abc", userlists.ID)
	assert.Equal(t, "userlist1", userlists.Name)
}
