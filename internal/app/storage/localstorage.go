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

func (s *InMemoryStorage) Lookup(shortURL string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	result, exists := s.urls[shortURL]
	if !exists {
		return "", errors.New("url not found")
	}
	return result.URL, nil
}

func (s *InMemoryStorage) Add(shortURL, originalURL string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	u := entity.InternalStorage{
		Alias: shortURL,
		URL:   originalURL,
	}
	s.urls[shortURL] = u

	return nil
}

func (s *InMemoryStorage) Ping() error {
	return nil
}

func (s *InMemoryStorage) Close() {
}
