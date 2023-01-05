package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
}

const escape = "\x1b"

const defaultConfigFile = "kit.yaml"

func main() {
	cmd := &cobra.Command{Use: "kit"}
	cmd.AddCommand(initCmd())
	cmd.AddCommand(lintCmd())
	cmd.AddCommand(upCmd())
	cmd.AddCommand(versionCmd())

	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
