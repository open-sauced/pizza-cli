package workspaces

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/open-sauced/go-api/client"
	"github.com/open-sauced/pizza-cli/cmd/auth"
	"github.com/open-sauced/pizza-cli/pkg/api"
	"github.com/open-sauced/pizza-cli/pkg/constants"
	"github.com/spf13/cobra"
)

// workspacesOptions are the options for the pizza workspaces command including user
// defined configurations
type workspacesOptions struct {
	// WorkspaceName: name of the workspaces to be created
	WorkspaceName string
	// APIClient: api client to interface with open sauced api
	APIClient *client.APIClient
	// Session: the user session
	Session auth.AccessTokenResponse

	ServerContext context.Context
}

func newDefaultWorkspacesOptions() *workspacesOptions {
	return &workspacesOptions{
		WorkspaceName: fmt.Sprintf("workspace-%s", uuid.NewString()),
		APIClient:     api.NewGoClient(constants.EndpointBeta),
		Session:       auth.AccessTokenResponse{},
		ServerContext: context.TODO(),
	}
}

// NewWorkspacesCommand returns a new cobra command for 'pizza workspaces'
func NewWorkspacesCommand() *cobra.Command {
	opts := newDefaultWorkspacesOptions()
	cmd := &cobra.Command{
		Use:   "workspaces [flags] <command> [flags]",
		Short: "Manage, share, and track open source projects",
		Long:  "Centralized hub for managing, sharing, and tracking open source projects",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
	cmd.Flags().StringVarP(&opts.WorkspaceName, "name", "n", opts.WorkspaceName, "name of the workspace to be created")
	cmd.AddCommand(NewListWorkSpaceCommand(opts))
	cmd.AddCommand(NewAddWorkSpaceCommand(opts))
	return cmd
}
