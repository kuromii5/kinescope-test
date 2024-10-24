package storage

import (
	"errors"
	"sync"
	"time"

	"github.com/kuromii5/kinescope-test/internal/models"
)

var (
	ErrValueExist = errors.New("key already exists")
	ErrNotFound   = errors.New("key not found")
)

type Storage struct {
	items         map[string]models.Item
	mu            sync.RWMutex
	clearInterval time.Duration
}

func NewStorage() *Storage {
	storage := &Storage{
		items:         make(map[string]models.Item),
		clearInterval: time.Minute,
	}

	storage.startClearingHandler()

	return storage
}

// Take value of given key if this key exists and not expired.
// If key was expired - delete from storage
func (s *Storage) Get(key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.items[key]
	if !ok {
		return "", ErrNotFound
	}

	if item.Deadline.Before(time.Now()) {
		delete(s.items, key)
		return "", ErrNotFound
	}

	return item.Value, nil
}

// Replace value for given key if key was set earlier
func (s *Storage) Set(key, value string, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[key]; !ok {
		return ErrNotFound
	}

	item := models.Item{
		Value:    value,
		Deadline: time.Now().Add(ttl),
	}

	s.items[key] = item
	return nil
}

// Append key to storage if this key was not set earlier
func (s *Storage) Add(key, value string, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[key]; ok {
		return ErrValueExist
	}

	s.items[key] = models.Item{
		Value:    value,
		Deadline: time.Now().Add(ttl),
	}

	return nil
}

// Remove key from storage if item with this key was set earlier
func (s *Storage) Del(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[key]; !ok {
		return ErrNotFound
	}

	delete(s.items, key)
	return nil
}

func (s *Storage) startClearingHandler() {
	ticker := time.NewTicker(s.clearInterval)

	go func() {
		for range ticker.C {
			s.mu.Lock()
			for key, item := range s.items {
				if time.Now().After(item.Deadline) {
					delete(s.items, key)
				}
			}
			s.mu.Unlock()
		}
	}()
}
