package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/msorokin-hash/urlshortener/internal/app/entity"
)

type FileStorage struct {
	mu       sync.Mutex
	filePath string
	urls     map[string]entity.FileStorage
}

func NewFileStorage(filePath string) (*FileStorage, error) {
	fs := &FileStorage{
		filePath: filePath,
		urls:     make(map[string]entity.FileStorage),
	}

	if err := fs.load(); err != nil {
		return nil, err
	}
	return fs, nil
}

func (fs *FileStorage) load() error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	file, err := os.Open(fs.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var record entity.FileStorage
		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			return err
		}
		fs.urls[record.ShortURL] = record
	}
	return scanner.Err()
}

func (fs *FileStorage) saveLine(record entity.FileStorage) error {
	file, err := os.OpenFile(fs.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(record)
}

func (fs *FileStorage) Add(hash, url string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	record := entity.FileStorage{
		UUID:        uuid.New().String(),
		ShortURL:    hash,
		OriginalURL: url,
	}

	fs.urls[hash] = record
	return fs.saveLine(record)
}

func (fs *FileStorage) Lookup(hash string) (string, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	result, exists := fs.urls[hash]
	if !exists {
		return "", errors.New("url not found")
	}
	return result.OriginalURL, nil
}
