package main

import (
	"context"
	"log"
	"os"
	"os/signal"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	switch os.Args[1] {
	case "update-examples":
		if err := updateExamples(ctx); err != nil {
			log.Fatal(err)
		}
	case "update-schema":
		if err := updateSchema(); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("unknown command %q", os.Args[1])
	}
}
