package main

import (
	"log"

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
		log.Fatal(err)
	}
}
