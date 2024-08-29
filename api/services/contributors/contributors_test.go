package contributors

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

func TestNewPullrequestContributors(t *testing.T) {
	t.Parallel()
	m := mock.NewMockRoundTripper(func(req *http.Request) (*http.Response, error) {
		// Check if the URL is correct
		assert.Equal(t, "https://api.example.com/v2/contributors/insights/new?range=30&repos=testowner%2Ftestrepo", req.URL.String())

		mockResponse := ContribResponse{
			Data: []DbContributor{
				{
					AuthorLogin: "contributor1",
				},
				{
					AuthorLogin: "contributor2",
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
	service := NewContributorsService(client, "https://api.example.com")

	newContribs, resp, err := service.NewPullRequestContributors([]string{"testowner/testrepo"}, 30)

	assert.NoError(t, err)
	assert.NotNil(t, newContribs)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	assert.Equal(t, newContribs.Data[0].AuthorLogin, "contributor1")
	assert.Equal(t, newContribs.Data[1].AuthorLogin, "contributor2")

	// Check the meta information
	assert.Equal(t, 1, newContribs.Meta.Page)
	assert.Equal(t, 30, newContribs.Meta.Limit)
	assert.Equal(t, 2, newContribs.Meta.ItemCount)
	assert.Equal(t, 1, newContribs.Meta.PageCount)
	assert.False(t, newContribs.Meta.HasPreviousPage)
	assert.False(t, newContribs.Meta.HasNextPage)
}

func TestRecentPullRequestContributors(t *testing.T) {
	t.Parallel()
	m := mock.NewMockRoundTripper(func(req *http.Request) (*http.Response, error) {
		// Check if the URL is correct
		assert.Equal(t, "https://api.example.com/v2/contributors/insights/recent?range=30&repos=testowner%2Ftestrepo", req.URL.String())

		mockResponse := ContribResponse{
			Data: []DbContributor{
				{
					AuthorLogin: "contributor1",
				},
				{
					AuthorLogin: "contributor2",
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
	service := NewContributorsService(client, "https://api.example.com")

	recentContribs, resp, err := service.RecentPullRequestContributors([]string{"testowner/testrepo"}, 30)

	assert.NoError(t, err)
	assert.NotNil(t, recentContribs)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	assert.Equal(t, recentContribs.Data[0].AuthorLogin, "contributor1")
	assert.Equal(t, recentContribs.Data[1].AuthorLogin, "contributor2")

	// Check the meta information
	assert.Equal(t, 1, recentContribs.Meta.Page)
	assert.Equal(t, 30, recentContribs.Meta.Limit)
	assert.Equal(t, 2, recentContribs.Meta.ItemCount)
	assert.Equal(t, 1, recentContribs.Meta.PageCount)
	assert.False(t, recentContribs.Meta.HasPreviousPage)
	assert.False(t, recentContribs.Meta.HasNextPage)
}

func TestAlumniPullRequestContributors(t *testing.T) {
	t.Parallel()
	m := mock.NewMockRoundTripper(func(req *http.Request) (*http.Response, error) {
		// Check if the URL is correct
		assert.Equal(t, "https://api.example.com/v2/contributors/insights/alumni?range=30&repos=testowner%2Ftestrepo", req.URL.String())

		mockResponse := ContribResponse{
			Data: []DbContributor{
				{
					AuthorLogin: "contributor1",
				},
				{
					AuthorLogin: "contributor2",
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
	service := NewContributorsService(client, "https://api.example.com")

	alumniContribs, resp, err := service.AlumniPullRequestContributors([]string{"testowner/testrepo"}, 30)

	assert.NoError(t, err)
	assert.NotNil(t, alumniContribs)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	assert.Equal(t, alumniContribs.Data[0].AuthorLogin, "contributor1")
	assert.Equal(t, alumniContribs.Data[1].AuthorLogin, "contributor2")

	// Check the meta information
	assert.Equal(t, 1, alumniContribs.Meta.Page)
	assert.Equal(t, 30, alumniContribs.Meta.Limit)
	assert.Equal(t, 2, alumniContribs.Meta.ItemCount)
	assert.Equal(t, 1, alumniContribs.Meta.PageCount)
	assert.False(t, alumniContribs.Meta.HasPreviousPage)
	assert.False(t, alumniContribs.Meta.HasNextPage)
}

func TestRepeatPullRequestContributors(t *testing.T) {
	t.Parallel()
	m := mock.NewMockRoundTripper(func(req *http.Request) (*http.Response, error) {
		// Check if the URL is correct
		assert.Equal(t, "https://api.example.com/v2/contributors/insights/repeat?range=30&repos=testowner%2Ftestrepo", req.URL.String())

		mockResponse := ContribResponse{
			Data: []DbContributor{
				{
					AuthorLogin: "contributor1",
				},
				{
					AuthorLogin: "contributor2",
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
	service := NewContributorsService(client, "https://api.example.com")

	repeatContribs, resp, err := service.RepeatPullRequestContributors([]string{"testowner/testrepo"}, 30)

	assert.NoError(t, err)
	assert.NotNil(t, repeatContribs)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	assert.Equal(t, repeatContribs.Data[0].AuthorLogin, "contributor1")
	assert.Equal(t, repeatContribs.Data[1].AuthorLogin, "contributor2")

	// Check the meta information
	assert.Equal(t, 1, repeatContribs.Meta.Page)
	assert.Equal(t, 30, repeatContribs.Meta.Limit)
	assert.Equal(t, 2, repeatContribs.Meta.ItemCount)
	assert.Equal(t, 1, repeatContribs.Meta.PageCount)
	assert.False(t, repeatContribs.Meta.HasPreviousPage)
	assert.False(t, repeatContribs.Meta.HasNextPage)
}

func TestSearchPullRequestContributors(t *testing.T) {
	t.Parallel()
	m := mock.NewMockRoundTripper(func(req *http.Request) (*http.Response, error) {
		// Check if the URL is correct
		assert.Equal(t, "https://api.example.com/v2/contributors/search?range=30&repos=testowner%2Ftestrepo", req.URL.String())

		mockResponse := ContribResponse{
			Data: []DbContributor{
				{
					AuthorLogin: "contributor1",
				},
				{
					AuthorLogin: "contributor2",
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
	service := NewContributorsService(client, "https://api.example.com")

	repeatContribs, resp, err := service.SearchPullRequestContributors([]string{"testowner/testrepo"}, 30)

	assert.NoError(t, err)
	assert.NotNil(t, repeatContribs)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	assert.Equal(t, repeatContribs.Data[0].AuthorLogin, "contributor1")
	assert.Equal(t, repeatContribs.Data[1].AuthorLogin, "contributor2")

	// Check the meta information
	assert.Equal(t, 1, repeatContribs.Meta.Page)
	assert.Equal(t, 30, repeatContribs.Meta.Limit)
	assert.Equal(t, 2, repeatContribs.Meta.ItemCount)
	assert.Equal(t, 1, repeatContribs.Meta.PageCount)
	assert.False(t, repeatContribs.Meta.HasPreviousPage)
	assert.False(t, repeatContribs.Meta.HasNextPage)
}
