package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func (a *Agent) addItem(baseURL, key, value string, ttl time.Duration) {
	url := fmt.Sprintf("%s/items", baseURL)

	reqBody, _ := json.Marshal(map[string]interface{}{
		"key":   key,
		"value": value,
		"ttl":   ttl.String(),
	})

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		a.logger.Error("failed to create request", slog.String("error", err.Error()))
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		a.logger.Error("failed to send request", slog.String("error", err.Error()))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		a.logger.Error("failed to add item", slog.String("status", resp.Status), slog.String("url", url))
		return
	}

	a.logger.Info("added item", slog.String("status", resp.Status))
}

func (a *Agent) getItem(baseURL, key string) {
	url := fmt.Sprintf("%s/items/%s", baseURL, key)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		a.logger.Error("failed to create request", slog.String("error", err.Error()))
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		a.logger.Error("failed to send request", slog.String("error", err.Error()))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		a.logger.Error("failed to get item", slog.String("status", resp.Status), slog.String("url", url))
		return
	}

	a.logger.Info("got item", slog.String("status", resp.Status))
}

func (a *Agent) setItem(baseURL, key, value string, ttl time.Duration) {
	url := fmt.Sprintf("%s/items/%s", baseURL, key)

	reqBody, _ := json.Marshal(map[string]interface{}{
		"value": value,
		"ttl":   ttl.String(),
	})

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(reqBody))
	if err != nil {
		a.logger.Error("failed to create request", slog.String("error", err.Error()))
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		a.logger.Error("failed to send request", slog.String("error", err.Error()))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		a.logger.Error("failed to set item", slog.String("status", resp.Status), slog.String("url", url))
		return
	}

	a.logger.Info("updated item", slog.String("status", resp.Status))
}

func (a *Agent) deleteItem(baseURL, key string) {
	url := fmt.Sprintf("%s/items/%s", baseURL, key)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		a.logger.Error("failed to create request", slog.String("error", err.Error()))
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		a.logger.Error("failed to send request", slog.String("error", err.Error()))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		a.logger.Error("failed to delete item", slog.String("status", resp.Status), slog.String("url", url))
		return
	}

	a.logger.Info("deleted item", slog.String("status", resp.Status))
}
