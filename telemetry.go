//go:build telemetry
// +build telemetry

package main

import "github.com/open-sauced/pizza-cli/pkg/utils"

// This alternate main is used as a one-shot for bootstrapping Posthog events:
// the various events called herein do not exist in Posthog's datalake until the
// event has landed.
//
// Therefore, this is useful for when there are new events for Posthog that need
// a dashboard bootstrapped for them.

func main() {
	println("Started bootstrapping Posthog events")
	client := utils.NewPosthogCliClient(true)

	err := client.CaptureLogin("test-user")
	if err != nil {
		panic(err)
	}

	err = client.CaptureFailedLogin()
	if err != nil {
		panic(err)
	}

	err = client.CaptureCodeownersGenerate()
	if err != nil {
		panic(err)
	}

	err = client.CaptureFailedCodeownersGenerate()
	if err != nil {
		panic(err)
	}

	err = client.CaptureConfigGenerate()
	if err != nil {
		panic(err)
	}

	err = client.CaptureFailedConfigGenerate()
	if err != nil {
		panic(err)
	}

	err = client.CaptureConfigGenerateMode("interactive")
	if err != nil {
		panic(err)
	}

	err = client.CaptureCodeownersGenerateAuth("test-user")
	if err != nil {
		panic(err)
	}

	err = client.CaptureFailedCodeownersGenerateAuth()
	if err != nil {
		panic(err)
	}

	err = client.CaptureCodeownersGenerateContributorInsight()
	if err != nil {
		panic(err)
	}

	err = client.CaptureFailedCodeownersGenerateContributorInsight()
	if err != nil {
		panic(err)
	}

	err = client.CaptureInsights()
	if err != nil {
		panic(err)
	}

	err = client.CaptureFailedInsights()
	if err != nil {
		panic(err)
	}

	err = client.Done()
	if err != nil {
		panic(err)
	}

	println("Done bootstrapping Posthog events")
}
