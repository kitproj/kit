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
	if err := updateExamples(ctx); err != nil {
		log.Fatal(err)
	}
	if err := updateSchema(); err != nil {
		log.Fatal(err)
	}
}
