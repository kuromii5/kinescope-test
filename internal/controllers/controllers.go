package controllers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kuromii5/kinescope-test/internal/models"
	"github.com/kuromii5/kinescope-test/internal/storage"
)

type StorageProvider interface {
	Add(key, value string, ttl time.Duration) error
	Set(key, value string, ttl time.Duration) error
	Get(key string) (string, error)
	Del(key string) error
}

type Controller struct {
	logger   *slog.Logger
	provider StorageProvider
}

func NewController(logger *slog.Logger, provider StorageProvider) *Controller {
	return &Controller{
		logger:   logger,
		provider: provider,
	}
}

func (c *Controller) addItemHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.AddItemRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			c.logger.Error("failed to decode request", slog.String("error", err.Error()))
			http.Error(w, "failed to decode request", http.StatusBadRequest)
			return
		}

		if err := c.provider.Add(req.Key, req.Value, req.TTL); err != nil {
			if errors.Is(err, storage.ErrValueExist) {
				c.logger.Error("key already exists", slog.String("error", err.Error()))
				http.Error(w, "key already exists", http.StatusBadRequest)
				return
			}

			c.logger.Error("internal server error", slog.String("error", err.Error()))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		resp := models.AddItemResponse{
			Message: "Item added successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			c.logger.Error("failed to encode response", slog.String("error", err.Error()))
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}

func (c *Controller) setItemHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.SetItemRequest
		key := chi.URLParam(r, "key")

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			c.logger.Error("failed to decode request", slog.String("error", err.Error()))
			http.Error(w, "failed to decode request", http.StatusBadRequest)
			return
		}

		if err := c.provider.Set(key, req.Value, req.TTL); err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				c.logger.Error("key was not found", slog.String("error", err.Error()))
				http.Error(w, "key was not found", http.StatusNotFound)
				return
			}

			c.logger.Error("failed to set item", slog.String("error", err.Error()))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		resp := models.SetItemResponse{
			Message: "Item set successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			c.logger.Error("failed to encode response", slog.String("error", err.Error()))
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}

func (c *Controller) getItemHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "key")

		value, err := c.provider.Get(key)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				c.logger.Error("key not found", slog.String("error", err.Error()))
				http.Error(w, "key not found", http.StatusNotFound)
				return
			}

			c.logger.Error("failed to get item", slog.String("error", err.Error()))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		resp := models.GetItemResponse{
			Value: value,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			c.logger.Error("failed to encode response", slog.String("error", err.Error()))
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}

func (c *Controller) deleteItemHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "key")

		if err := c.provider.Del(key); err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				c.logger.Error("key not found", slog.String("error", err.Error()))
				http.Error(w, "key not found", http.StatusNotFound)
				return
			}

			c.logger.Error("failed to delete item", slog.String("error", err.Error()))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		resp := models.DeleteItemResponse{
			Message: "Item deleted successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			c.logger.Error("failed to encode response", slog.String("error", err.Error()))
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}
