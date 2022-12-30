package main

import (
	_ "embed"
	"log"

	"github.com/spf13/cobra"
)

//go:generate sh -c "git describe --tags > tag"
//go:embed tag
var tag string

func version() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println(tag)
		},
	}
}
