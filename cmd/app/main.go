package main

import (
	"context"
	"log"

	"github.com/kuromii5/kinescope-test/internal/app"
)

func main() {
	ctx := context.Background()

	app, err := app.New(ctx)
	if err != nil {
		log.Fatalf("failed to initialize app: %s", err.Error())
	}

	if err = app.Run(); err != nil {
		log.Fatalf("failed to run app: %s", err.Error())
	}
}
