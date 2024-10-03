package utils

import (
	"fmt"

	"github.com/posthog/posthog-go"
)

var (
	writeOnlyPublicPosthogKey = "dev"
	posthogEndpoint           = "https://us.i.posthog.com"
)

// PosthogCliClient is a wrapper around the posthog-go client and is used as a
// API entrypoint for sending OpenSauced telemetry data for CLI commands
type PosthogCliClient struct {
	// client is the Posthog Go client
	client posthog.Client

	// activated denotes if the user has enabled or disabled telemetry
	activated bool

	// uniqueID is the user's unique, anonymous identifier
	uniqueID string
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

	uniqueID, err := getOrCreateUniqueID()
	if err != nil {
		fmt.Printf("could not build anonymous telemetry client: %s\n", err)
	}

	return &PosthogCliClient{
		client:    client,
		activated: activated,
		uniqueID:  uniqueID,
	}
}

// Done should always be called in order to flush the Posthog buffers before
// the CLI exits to ensure all events are accurately captured.
func (p *PosthogCliClient) Done() error {
	return p.client.Close()
}

// CaptureLogin gathers telemetry on users who log into OpenSauced via the CLI
func (p *PosthogCliClient) CaptureLogin(username string) error {
	if p.activated {
		return p.client.Enqueue(posthog.Capture{
			DistinctId: username,
			Event:      "pizza_cli_user_logged_in",
		})
	}

	return nil
}

// CaptureFailedLogin gathers telemetry on failed logins via the CLI
func (p *PosthogCliClient) CaptureFailedLogin() error {
	if p.activated {
		return p.client.Enqueue(posthog.Capture{
			DistinctId: p.uniqueID,
			Event:      "pizza_cli_user_failed_log_in",
		})
	}

	return nil
}

// CaptureCodeownersGenerate gathers telemetry on successful codeowners generation
func (p *PosthogCliClient) CaptureCodeownersGenerate() error {
	if p.activated {
		return p.client.Enqueue(posthog.Capture{
			DistinctId: p.uniqueID,
			Event:      "pizza_cli_generated_codeowners",
		})
	}

	return nil
}

// CaptureFailedCodeownersGenerate gathers telemetry on failed codeowners generation
func (p *PosthogCliClient) CaptureFailedCodeownersGenerate() error {
	if p.activated {
		return p.client.Enqueue(posthog.Capture{
			DistinctId: p.uniqueID,
			Event:      "pizza_cli_failed_to_generate_codeowners",
		})
	}

	return nil
}

// CaptureCodeownersGenerateAuth gathers telemetry on successful auth flows during codeowners generation
func (p *PosthogCliClient) CaptureCodeownersGenerateAuth(username string) error {
	if p.activated {
		return p.client.Enqueue(posthog.Capture{
			DistinctId: username,
			Event:      "pizza_cli_user_authenticated_during_generate_codeowners_flow",
		})
	}

	return nil
}

// CaptureFailedCodeownersGenerateAuth gathers telemetry on failed auth flows during codeowners generations
func (p *PosthogCliClient) CaptureFailedCodeownersGenerateAuth() error {
	if p.activated {
		return p.client.Enqueue(posthog.Capture{
			DistinctId: p.uniqueID,
			Event:      "pizza_cli_user_failed_to_authenticate_during_generate_codeowners_flow",
		})
	}

	return nil
}

// CaptureCodeownersGenerateContributorInsight gathers telemetry on successful
// Contributor Insights creation/update during codeowners generation
func (p *PosthogCliClient) CaptureCodeownersGenerateContributorInsight() error {
	if p.activated {
		return p.client.Enqueue(posthog.Capture{
			DistinctId: p.uniqueID,
			Event:      "pizza_cli_created_updated_contributor_list",
		})
	}

	return nil
}

// CaptureFailedCodeownersGenerateContributorInsight gathers telemetry on failed
// Contributor Insights during codeowners generation
func (p *PosthogCliClient) CaptureFailedCodeownersGenerateContributorInsight() error {
	if p.activated {
		return p.client.Enqueue(posthog.Capture{
			DistinctId: p.uniqueID,
			Event:      "pizza_cli_failed_to_create_update_contributor_insight_for_user",
		})
	}

	return nil
}

// CaptureConfigGenerate gathers telemetry on success
func (p *PosthogCliClient) CaptureConfigGenerate() error {
	if p.activated {
		return p.client.Enqueue(posthog.Capture{
			DistinctId: p.uniqueID,
			Event:      "pizza_cli_generated_config",
		})
	}

	return nil
}

// CaptureConfigGenerateMode gathers what mode a user is in when generating
// either 'Automatic' (default) or 'Interactive'
func (p *PosthogCliClient) CaptureConfigGenerateMode(mode string) error {
	properties := make(map[string]interface{})

	properties["mode"] = mode

	if p.activated {
		return p.client.Enqueue(posthog.Capture{
			DistinctId: p.uniqueID,
			Event:      "pizza_cli_generated_config_mode",
			Properties: properties,
		})
	}

	return nil
}

// CaptureFailedConfigGenerate gathers telemetry on failed
func (p *PosthogCliClient) CaptureFailedConfigGenerate() error {
	if p.activated {
		return p.client.Enqueue(posthog.Capture{
			DistinctId: p.uniqueID,
			Event:      "pizza_cli_failed_to_generate_config",
		})
	}

	return nil
}

func (p *PosthogCliClient) CaptureOffboard() error {
	if p.activated {
		return p.client.Enqueue(posthog.Capture{
			DistinctId: p.uniqueID,
			Event:      "pizza_cli_offboard",
		})
	}

	return nil
}

func (p *PosthogCliClient) CaptureFailedOffboard() error {
	if p.activated {
		return p.client.Enqueue(posthog.Capture{
			DistinctId: p.uniqueID,
			Event:      "pizza_cli_failed_to_offboard",
		})
	}

	return nil
}

// CaptureInsights gathers telemetry on successful Insights command runs
func (p *PosthogCliClient) CaptureInsights() error {
	if p.activated {
		return p.client.Enqueue(posthog.Capture{
			DistinctId: p.uniqueID,
			Event:      "pizza_cli_called_insights_command",
		})
	}

	return nil
}

// CaptureFailedInsights gathers telemetry on failed Insights command runs
func (p *PosthogCliClient) CaptureFailedInsights() error {
	if p.activated {
		return p.client.Enqueue(posthog.Capture{
			DistinctId: p.uniqueID,
			Event:      "pizza_cli_failed_to_call_insights_command",
		})
	}

	return nil
}
