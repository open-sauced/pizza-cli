package version

import (
	"fmt"

	"github.com/open-sauced/pizza-cli/pkg/utils"
	"github.com/spf13/cobra"
)

func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Displays the build version of the CLI",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %s\nSha: %s\n", utils.Version, utils.Sha)
		},
	}
}
