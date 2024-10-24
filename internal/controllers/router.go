package controllers

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kuromii5/kinescope-test/internal/storage"
)

func NewRouter(logger *slog.Logger, storage *storage.Storage) http.Handler {
	r := chi.NewRouter()

	c := NewController(logger, storage)

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/items", func(r chi.Router) {
		r.Post("/", c.addItemHandler())
		r.Get("/{key}", c.getItemHandler())
		r.Put("/{key}", c.setItemHandler())
		r.Delete("/{key}", c.deleteItemHandler())
	})

	return r
}
