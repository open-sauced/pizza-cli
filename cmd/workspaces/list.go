package workspaces

import (
	"context"
	"fmt"
	"net/http"

	sw "github.com/open-sauced/go-api/client"
	"github.com/open-sauced/pizza-cli/cmd/auth"
	"github.com/spf13/cobra"
)

type ListCommandOpts struct {
	*workspacesOptions
}

func NewListWorkSpaceCommand(workspaceOpts *workspacesOptions) *cobra.Command {
	opts := &ListCommandOpts{workspaceOpts}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list all workspaces",
		Long:  "retrieve all the workspaces",
		PreRunE: func(_ *cobra.Command, args []string) error {
			if opts.Session.AccessToken == "" {
				session, err := auth.GetUserSession()
				if err != nil {
					return err
				}
				opts.Session = session
			}
			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			return opts.run()
		},
	}
	return cmd
}

func (opts *ListCommandOpts) run() error {
	authCtx := context.WithValue(context.Background(), sw.ContextAccessToken, opts.Session.AccessToken)
	// here should return an array of DbWorkspace
	_, r, err := opts.APIClient.WorkspacesServiceAPI.GetWorkspaceForUser(authCtx).Execute()
	if err != nil {
		return err
	}
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP status: %d", r.StatusCode)
	}
	return nil
}
