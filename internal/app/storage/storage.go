package storage

import (
	"errors"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type Data struct {
	HashURL    string
	RequestURL string
}

type Storage struct {
	mu   sync.Mutex
	urls map[string]Data
}

func NewStorage() *Storage {
	return &Storage{
		urls: make(map[string]Data),
	}
}

func (s *Storage) GetURLByHash(hash string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	result, exists := s.urls[hash]
	if !exists {
		return "", errors.New("url not found")
	}
	return result.RequestURL, nil
}

func (s *Storage) CreateURL(hash string, url string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	u := Data{
		HashURL:    hash,
		RequestURL: url,
	}
	s.urls[hash] = u

	return nil
}
