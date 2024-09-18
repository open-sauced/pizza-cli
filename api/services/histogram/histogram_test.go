package histogram

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/open-sauced/pizza-cli/v2/api/mock"
)

func TestPrsHistogram(t *testing.T) {
	t.Parallel()
	m := mock.NewMockRoundTripper(func(req *http.Request) (*http.Response, error) {
		// Check if the URL is correct
		assert.Equal(t, "https://api.example.com/v2/histogram/pull-requests?range=30&repo=testowner%2Ftestrepo", req.URL.String())

		mockResponse := []PrHistogramData{
			{
				PrCount: 1,
			},
			{
				PrCount: 2,
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
	service := NewHistogramService(client, "https://api.example.com")

	prs, resp, err := service.PrsHistogram("testowner/testrepo", 30)

	require.NoError(t, err)
	assert.NotNil(t, prs)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, prs, 2)
	assert.Equal(t, 1, prs[0].PrCount)
	assert.Equal(t, 2, prs[1].PrCount)
}
