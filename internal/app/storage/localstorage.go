package storage

import (
	"errors"
	"sync"

	"github.com/msorokin-hash/urlshortener/internal/app/entity"
)

type InMemoryStorage struct {
	mu   sync.Mutex
	urls map[string]entity.InternalStorage
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		urls: make(map[string]entity.InternalStorage),
	}
}

func (s *InMemoryStorage) Lookup(hash string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	result, exists := s.urls[hash]
	if !exists {
		return "", errors.New("url not found")
	}
	return result.URL, nil
}

func (s *InMemoryStorage) Add(hash, url string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	u := entity.InternalStorage{
		Alias: hash,
		URL:   url,
	}
	s.urls[hash] = u

	return nil
}
