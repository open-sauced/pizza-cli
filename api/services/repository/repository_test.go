package repository

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/open-sauced/pizza-cli/api/mock"
	"github.com/open-sauced/pizza-cli/api/services"
)

func TestFindOneByOwnerAndRepo(t *testing.T) {
	t.Parallel()
	m := mock.NewMockRoundTripper(func(req *http.Request) (*http.Response, error) {
		// Check if the URL is correct
		assert.Equal(t, "https://api.example.com/v2/repos/testowner/testrepo", req.URL.String())

		mockResponse := DbRepository{
			ID:       1,
			FullName: "testowner/testrepo",
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
	service := NewRepositoryService(client, "https://api.example.com")

	repo, resp, err := service.FindOneByOwnerAndRepo("testowner", "testrepo")

	require.NoError(t, err)
	assert.NotNil(t, repo)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 1, repo.ID)
	assert.Equal(t, "testowner/testrepo", repo.FullName)
}

func TestFindContributorsByOwnerAndRepo(t *testing.T) {
	t.Parallel()
	m := mock.NewMockRoundTripper(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "https://api.example.com/v2/repos/testowner/testrepo/contributors?range=30", req.URL.String())

		mockResponse := ContributorsResponse{
			Data: []DbContributorInfo{
				{
					ID:    1,
					Login: "contributor1",
				},
				{
					ID:    2,
					Login: "contributor2",
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
	service := NewRepositoryService(client, "https://api.example.com")

	contributors, resp, err := service.FindContributorsByOwnerAndRepo("testowner", "testrepo", 30)

	require.NoError(t, err)
	assert.NotNil(t, contributors)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, contributors.Data, 2)

	// Check the first contributor
	assert.Equal(t, 1, contributors.Data[0].ID)
	assert.Equal(t, "contributor1", contributors.Data[0].Login)

	// Check the second contributor
	assert.Equal(t, 2, contributors.Data[1].ID)
	assert.Equal(t, "contributor2", contributors.Data[1].Login)

	// Check the meta information
	assert.Equal(t, 1, contributors.Meta.Page)
	assert.Equal(t, 30, contributors.Meta.Limit)
	assert.Equal(t, 2, contributors.Meta.ItemCount)
	assert.Equal(t, 1, contributors.Meta.PageCount)
	assert.False(t, contributors.Meta.HasPreviousPage)
	assert.False(t, contributors.Meta.HasNextPage)
}
