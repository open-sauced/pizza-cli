package bake

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/open-sauced/go-api/client"
)

func TestSendsPost(t *testing.T) {
	tests := []struct {
		name string
		opts *Options
	}{
		{
			name: "Sends post request",
			opts: &Options{
				Repos: []string{"https://test.com"},
			},
		},
		{
			name: "Sends post request with multiple Repos",
			opts: &Options{
				Repos: []string{"https://test.com", "https://github.com/open-sauced/pizza", "https://github.com/open-sauced/insights"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Fail()
				}
				// Always return an ok status with a dummy body from the mock server
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte("body"))
				if err != nil {
					// Writing back to the client shouldn't fail
					t.Fail()
				}
			}))
			defer testServer.Close()
			configuration := client.NewConfiguration()
			configuration.Servers = client.ServerConfigurations{{URL: testServer.URL}}
			tt.opts.APIClient = client.NewAPIClient(configuration)
			err := run(tt.opts)
			if err != nil {
				t.Fail()
			}
		})
	}
}
