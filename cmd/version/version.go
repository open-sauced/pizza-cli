package version

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/open-sauced/pizza-cli/pkg/utils"
)

func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Displays the build version of the CLI",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf("Version: %s\nSha: %s\nBuilt at: %s\n", utils.Version, utils.Sha, utils.Datetime)
		},
	}
}
