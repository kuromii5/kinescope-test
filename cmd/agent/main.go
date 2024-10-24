package main

import (
	"context"
	"log"

	"github.com/kuromii5/kinescope-test/internal/agent"
)

func main() {
	ctx := context.Background()

	agent, err := agent.NewAgent(ctx)
	if err != nil {
		log.Fatalf("failed to initialize agent: %s", err.Error())
	}

	agent.Run(ctx)
}
