package storage

import (
	"log"

	"github.com/msorokin-hash/urlshortener/internal/app/config"
	"github.com/msorokin-hash/urlshortener/internal/app/service"
)

func SelectStorage(cfg *config.Config) (service.Storage, error) {
	if cfg.DatabaseDSN != "" {
		dbStorage, err := NewPostgresStorage(cfg.DatabaseDSN)
		if err == nil {
			log.Println("Используется PostgreSQL")
			return dbStorage, nil
		}
		log.Println("Ошибка подключения к PostgreSQL, пробуем файл:", err)
	}

	if cfg.FileStoragePath != "" {
		fileStorage, err := NewFileStorage(cfg.FileStoragePath)
		if err == nil {
			log.Println("Используется файловое хранилище:", cfg.FileStoragePath)
			return fileStorage, nil
		}
		log.Println("Ошибка работы с файлом, пробуем память:", err)
	}

	log.Println("Используется хранение в памяти (InMemoryStorage)")

	return NewInMemoryStorage(), nil
}
