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
	client    posthog.Client
	activated bool
}

// NewPosthogCliClient returns a PosthogCliClient which can be used to capture
// telemetry events for CLI users
func NewPosthogCliClient(activated bool) *PosthogCliClient {
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
		client:    client,
		activated: activated,
	}
}

// Done should always be called in order to flush the Posthog buffers before
// the CLI exits to ensure all events are accurately captured.
//
//nolint:errcheck
func (p *PosthogCliClient) Done() {
	p.client.Close()
}

// CaptureLogin gathers telemetry on users who log into OpenSauced via the CLI
//
//nolint:errcheck
func (p *PosthogCliClient) CaptureLogin(username string) {
	if p.activated {
		p.client.Enqueue(posthog.Capture{
			DistinctId: username,
			Event:      "cli_user logged in",
		})
	}
}

// CaptureFailedLogin gathers telemetry on failed logins via the CLI
//
//nolint:errcheck
func (p *PosthogCliClient) CaptureFailedLogin() {
	if p.activated {
		p.client.Enqueue(posthog.Capture{
			DistinctId: "login-failures",
			Event:      "cli_user failed log in",
		})
	}
}

//nolint:errcheck
func (p *PosthogCliClient) CaptureCodeownersGenerate() {
	if p.activated {
		p.client.Enqueue(posthog.Capture{
			DistinctId: "codeowners-generated",
			Event:      "cli generated codeowners",
		})
	}
}

//nolint:errcheck
func (p *PosthogCliClient) CaptureFailedCodeownersGenerate() {
	if p.activated {
		p.client.Enqueue(posthog.Capture{
			DistinctId: "failed-codeowners-generated",
			Event:      "cli failed to generate codeowners",
		})
	}
}

//nolint:errcheck
func (p *PosthogCliClient) CaptureCodeownersGenerateAuth(username string) {
	if p.activated {
		p.client.Enqueue(posthog.Capture{
			DistinctId: username,
			Event:      "user authenticated during generate codeowners flow",
		})
	}
}

//nolint:errcheck
func (p *PosthogCliClient) CaptureFailedCodeownersGenerateAuth() {
	if p.activated {
		p.client.Enqueue(posthog.Capture{
			DistinctId: "codeowners-generate-auth-failed",
			Event:      "user failed to authenticate during generate codeowners flow",
		})
	}
}

//nolint:errcheck
func (p *PosthogCliClient) CaptureCodeownersGenerateContributorInsight() {
	if p.activated {
		p.client.Enqueue(posthog.Capture{
			DistinctId: "codeowners-generate-contributor-insight",
			Event:      "cli created/updated contributor list for user",
		})
	}
}

//nolint:errcheck
func (p *PosthogCliClient) CaptureFailedCodeownersGenerateContributorInsight() {
	if p.activated {
		p.client.Enqueue(posthog.Capture{
			DistinctId: "failed-codeowners-generation-contributor-insight",
			Event:      "cli failed to create/update contributor insight for user",
		})
	}
}

//nolint:errcheck
func (p *PosthogCliClient) CaptureInsights() {
	if p.activated {
		p.client.Enqueue(posthog.Capture{
			DistinctId: "insights",
			Event:      "cli called insights command",
		})
	}
}

//nolint:errcheck
func (p *PosthogCliClient) CaptureFailedInsights() {
	if p.activated {
		p.client.Enqueue(posthog.Capture{
			DistinctId: "failed-insight",
			Event:      "cli failed to call insights command",
		})
	}
}
