package agent

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/kuromii5/kinescope-test/internal/logger"
)

type Agent struct {
	client http.Client
	logger *slog.Logger
}

func NewAgent(ctx context.Context) (*Agent, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	logger := logger.New("local", "info")

	return &Agent{
		client: client,
		logger: logger,
	}, nil
}

func (a *Agent) Run(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(1)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		defer wg.Done()
		a.performRequests()
	}()

	<-stopChan
	a.logger.Info("Shutting down the agent...")
	wg.Wait()
}

func (a *Agent) performRequests() {
	baseURL := "http://localhost:8080"

	a.addItem(baseURL, "key1", "value1", 10*time.Second)

	a.getItem(baseURL, "key1")

	a.setItem(baseURL, "key1", "value2", 20*time.Second)

	a.deleteItem(baseURL, "key1")
}
