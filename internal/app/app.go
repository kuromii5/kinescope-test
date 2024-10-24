package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/kuromii5/kinescope-test/internal/controllers"
	"github.com/kuromii5/kinescope-test/internal/logger"
	"github.com/kuromii5/kinescope-test/internal/storage"
)

type App struct {
	logger *slog.Logger
	server *http.Server
}

func New(ctx context.Context) (*App, error) {
	app := &App{}

	app.logger = logger.New("local", "info")

	if err := app.initDeps(); err != nil {
		return nil, fmt.Errorf("failed to init dependencies: %w", err)
	}

	return app, nil
}

func (a *App) initDeps() error {
	storage := storage.NewStorage()

	if err := a.initServer(storage); err != nil {
		return fmt.Errorf("failed to init server: %w", err)
	}

	return nil
}

func (a *App) initServer(storage *storage.Storage) error {
	r := controllers.NewRouter(a.logger, storage)

	a.server = &http.Server{
		Addr:         "localhost:8080",
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		IdleTimeout:  time.Second * 10,
		Handler:      r,
	}

	return nil
}

func (a *App) Run() error {
	a.logger.Info("Starting server on localhost:8080")

	if err := a.server.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to run app: %w", err)
	}

	return nil
}
