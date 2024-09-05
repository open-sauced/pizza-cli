package config

// The configuration specification
type Spec struct {

	// Attributions are mappings of GitHub usernames to a list of emails. These
	// emails should be the associated addresses used by individual GitHub users.
	// Example: { github_username: [ email1@domain.com, email2@domain.com ]} where
	// "github_username" has 2 emails attributed to them and their work.
	Attributions map[string][]string `yaml:"attribution"`

	// AttributionFallback is the default username to attribute to the filename
	// if no other attributions were found.
	AttributionFallback []string `yaml:"attribution-fallback"`
}
