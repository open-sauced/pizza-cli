package mock

import "net/http"

// RoundTripper is a custom, mock http.RoundTripper used for testing and mocking
// purposes ONLY.
type RoundTripper struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

// NewMockRoundTripper returns a new RoundTripper which will execut the given
// roundTripFunc provided by the caller
func NewMockRoundTripper(roundTripFunc func(req *http.Request) (*http.Response, error)) *RoundTripper {
	return &RoundTripper{
		RoundTripFunc: roundTripFunc,
	}
}

// RoundTrip fufills the http.Client interface and executes the provided RoundTripFunc
// given by the caller in the NewMockRoundTripper
func (m *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}
