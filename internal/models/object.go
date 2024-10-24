package models

import (
	"encoding/json"
	"time"
)

type Item struct {
	Value    string
	Deadline time.Time
}

// GET

type GetItemResponse struct {
	Value string `json:"value"`
}

// ADD

type AddItemRequest struct {
	Key   string        `json:"key"`
	Value string        `json:"value"`
	TTL   time.Duration `json:"ttl"`
}

func (a *AddItemRequest) UnmarshalJSON(data []byte) (err error) {
	type Alias AddItemRequest

	aliasValue := &struct {
		*Alias

		TTL string `json:"ttl"`
	}{
		Alias: (*Alias)(a),
	}

	if err = json.Unmarshal(data, aliasValue); err != nil {
		return
	}

	ttl, err := time.ParseDuration(aliasValue.TTL)
	if err != nil {
		return
	}

	a.TTL = ttl
	return
}

type AddItemResponse struct {
	Message string `json:"message"`
}

// SET

type SetItemRequest struct {
	Value string        `json:"value"`
	TTL   time.Duration `json:"ttl"`
}

func (s *SetItemRequest) UnmarshalJSON(data []byte) (err error) {
	type Alias SetItemRequest

	aliasValue := &struct {
		*Alias

		TTL string `json:"ttl"`
	}{
		Alias: (*Alias)(s),
	}

	if err = json.Unmarshal(data, aliasValue); err != nil {
		return
	}

	ttl, err := time.ParseDuration(aliasValue.TTL)
	if err != nil {
		return
	}

	s.TTL = ttl
	return
}

type SetItemResponse struct {
	Message string `json:"message"`
}

// DELETE

type DeleteItemResponse struct {
	Message string `json:"message"`
}
