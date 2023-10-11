package main

import (
	"log"
	"os"

	"github.com/open-sauced/pizza-cli/cmd/root"
	"github.com/open-sauced/pizza-cli/pkg/utils"
)

func main() {
	rootCmd, err := root.NewRootCommand()
	if err != nil {
		log.Fatal(err)
	}
	utils.SetupRootCommand(rootCmd)
	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
