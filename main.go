package main

import (
	"log"
	"os"
)

func init() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
}

const escape = "\x1b"

const defaultConfigFile = "kit.yaml"

func main() {
	cmd := upCmd()

	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
