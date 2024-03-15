package workspaces

import (
	"context"
	"fmt"
	"net/http"

	client "github.com/open-sauced/go-api/client"
	"github.com/open-sauced/pizza-cli/cmd/auth"
	"github.com/open-sauced/pizza-cli/pkg/constants"
	"github.com/open-sauced/pizza-cli/pkg/utils"
	"github.com/spf13/cobra"
)

type AddCommand struct {
	*workspacesOptions
	// Repos is the array of git repository urls
	Repos []string
	// FilePath: the path to the yaml file
	FilePath string
	// TUI: terminal interface mode
	TUI bool
}

func NewAddWorkSpaceCommand(workspaceOpts *workspacesOptions) *cobra.Command {
	addCmd := &AddCommand{FilePath: "", workspacesOptions: workspaceOpts}
	cmd := &cobra.Command{
		Use:   "add url... [flags]",
		Short: "add repositories and contributors to a workspace",
		Long:  "add repositories and contributors to a workspace",
		PreRunE: func(_ *cobra.Command, args []string) error {
			if addCmd.Session.AccessToken == "" {
				session, err := auth.GetUserSession()
				if err != nil {
					return err
				}
				addCmd.Session = session
			}
			return nil
		},
		Args: func(cmd *cobra.Command, args []string) error {
			fileFlag := cmd.Flags().Lookup(constants.FlagNameFile)
			if !fileFlag.Changed && len(args) == 0 && !addCmd.TUI {
				return fmt.Errorf("must specify git repository url argument(s) or provide %s flag", fileFlag.Name)

			}
			addCmd.Repos = append(addCmd.Repos, args...)
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			return addCmd.run()
		},
		TraverseChildren: true,
	}

	cmd.Flags().StringVarP(&addCmd.FilePath, constants.FlagNameFile, "f", "", "Path to yaml file containing an array of git repository urls")
	cmd.Flags().StringVarP(&addCmd.WorkspaceName, "name", "n", addCmd.WorkspaceName, "name of the workspace to be created")
	cmd.Flags().BoolVar(&addCmd.TUI, "tui", addCmd.TUI, "use terminal user interface")
	return cmd
}

func (a *AddCommand) run() error {
	var workspaceData client.CreateWorkspaceDto
	var err error
	if a.TUI {
		workspaceData, err = a.createView()
		if err != nil {
			return err
		}
	} else {
		repos, err := utils.HandleUniqueValues(a.Repos, a.FilePath)
		if err != nil {
			return err
		}
		parsedRepos := make([]interface{}, len(repos))
		i := 0
		for repo := range repos {
			parsedRepos[i] = repo
			i++
		}
		workspaceData = *client.NewCreateWorkspaceDto(a.WorkspaceName, "my workspace", []interface{}{a.Session.User.UserMetadata["user_name"]}, parsedRepos)
	}

	authCtx := context.WithValue(context.Background(), client.ContextAccessToken, a.Session.AccessToken)
	_, r, err := a.APIClient.WorkspacesServiceAPI.CreateWorkspaceForUser(authCtx).CreateWorkspaceDto(workspaceData).Execute()
	if err != nil {
		return err
	}
	if r.StatusCode != http.StatusCreated {
		return fmt.Errorf("HTTP status: %d", r.StatusCode)
	}

	fmt.Printf("Workspace %s, has been created!", workspaceData.Name)
	return nil
}
