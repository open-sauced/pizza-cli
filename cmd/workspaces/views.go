package workspaces

import (
	"fmt"

	"github.com/charmbracelet/huh"
	client "github.com/open-sauced/go-api/client"
	"github.com/open-sauced/pizza-cli/pkg/utils"
)

func (a *AddCommand) createView() (client.CreateWorkspaceDto, error) {
	formValues := client.CreateWorkspaceDto{}
	var reposInput, membersInput string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Workspace Name").Value(&formValues.Name).Validate(func(s string) error {
				if s == "" {
					return fmt.Errorf("workspace name required")
				}
				return nil
			}),
			huh.NewText().Title("Description").Value(&formValues.Description).Lines(2),
			huh.NewText().Title("Repositories").
				Description("repositories to add to the workspace (yaml file, or comma separated values)").Lines(2).
				Validate(func(input string) error {
					if _, err := utils.ParseFileAndCSV(input); err != nil {
						return err
					}

					return nil
				}).Value(&reposInput),
			huh.NewText().Title("Members").Description("members to add to the workspace (yaml file, or comma separated values)").Lines(2).
				Validate(func(input string) error {
					if _, err := utils.ParseFileAndCSV(input); err != nil {
						return err
					}
					return nil
				}).Value(&membersInput),
		),
	)

	if err := form.Run(); err != nil {
		return formValues, err
	}

	// errors are checked in form
	repos, _ := utils.ParseFileAndCSV(reposInput)
	members, _ := utils.ParseFileAndCSV(membersInput)
	formValues.Repos = repos
	formValues.Members = members

	return formValues, nil
}
