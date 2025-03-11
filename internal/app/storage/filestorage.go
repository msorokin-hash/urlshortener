package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"

	"github.com/google/uuid"
	"github.com/msorokin-hash/urlshortener/internal/app/entity"
)

type FileStorage struct {
	file    *os.File
	encoder *json.Encoder
}

func NewFileStorage(filePath string) (*FileStorage, error) {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return &FileStorage{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (fs *FileStorage) Add(hash, url string) error {
	record := entity.FileStorage{
		UUID:        uuid.New().String(),
		ShortURL:    hash,
		OriginalURL: url,
	}

	file, err := os.OpenFile(fs.file.Name(), os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(record)
}

func (fs *FileStorage) Lookup(hash string) (string, error) {
	file, err := os.Open(fs.file.Name())
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var record entity.FileStorage
		if err := json.Unmarshal(scanner.Bytes(), &record); err == nil {
			if record.ShortURL == hash {
				return record.OriginalURL, nil
			}
		}
	}
	return "", errors.New("url not found")
}
