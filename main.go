package main

import (
	"log"
	"os"
)

func init() {
	_ = os.Mkdir("logs", 0o777)
	f, err := os.OpenFile("logs/kit.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	log.SetOutput(f)
	log.Println(tag)
}

const escape = "\x1b"

const defaultConfigFile = "kit.yaml"

func main() {
	cmd := upCmd()

	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
