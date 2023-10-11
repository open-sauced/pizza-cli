package utils

import (
	"github.com/posthog/posthog-go"
)

var (
	writeOnlyPublicPosthogKey = "dev"
	posthogEndpoint           = "https://app.posthog.com"
)

// PosthogCliClient is a wrapper around the posthog-go client and is used as a
// API entrypoint for sending OpenSauced telemetry data for CLI commands
type PosthogCliClient struct {
	client posthog.Client
}

// NewPosthogCliClient returns a PosthogCliClient which can be used to capture
// telemetry events for CLI users
func NewPosthogCliClient() *PosthogCliClient {
	client, err := posthog.NewWithConfig(
		writeOnlyPublicPosthogKey,
		posthog.Config{
			Endpoint: posthogEndpoint,
		},
	)

	if err != nil {
		// Should never happen since we aren't setting posthog.Config data that
		// would cause its validation to fail
		panic(err)
	}

	return &PosthogCliClient{
		client: client,
	}
}

// Done should always be called in order to flush the Posthog buffers before
// the CLI exits to ensure all events are accurately captured.
//
//nolint:errcheck
func (p *PosthogCliClient) Done() {
	p.client.Close()
}

// CaptureBake gathers telemetry on git repos that are being baked
//
//nolint:errcheck
func (p *PosthogCliClient) CaptureBake(urls []string) {
	p.client.Enqueue(posthog.Capture{
		DistinctId: "pizza-bakers",
		Event:      "cli_user baked repo",
		Properties: posthog.NewProperties().Set("clone_url", urls),
	})
}

// CaptureLogin gathers telemetry on users who log into OpenSauced via the CLI
//
//nolint:errcheck
func (p *PosthogCliClient) CaptureLogin(username string) {
	p.client.Enqueue(posthog.Capture{
		DistinctId: username,
		Event:      "cli_user logged in",
	})
}

// CaptureFailedLogin gathers telemetry on failed logins via the CLI
//
//nolint:errcheck
func (p *PosthogCliClient) CaptureFailedLogin() {
	p.client.Enqueue(posthog.Capture{
		DistinctId: "login-failures",
		Event:      "cli_user failed log in",
	})
}

// CaptureRepoQuery gathers telemetry on users using the repo-query service
//
//nolint:errcheck
func (p *PosthogCliClient) CaptureRepoQuery(url string) {
	p.client.Enqueue(posthog.Capture{
		DistinctId: "repo-queriers",
		Event:      "cli_user used repo-query",
		Properties: posthog.NewProperties().Set("github_url", url),
	})
}
