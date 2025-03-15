package main

import (
	"log"

	"github.com/msorokin-hash/urlshortener/internal/app/config"
	"github.com/msorokin-hash/urlshortener/internal/app/handler"
	"github.com/msorokin-hash/urlshortener/internal/app/server"
	"github.com/msorokin-hash/urlshortener/internal/app/storage"
)

func main() {
	// Инициализируем конфигурацию
	cfg := config.NewConfig()

	// Создаем storage
	store, err := storage.SelectStorage(cfg)
	if err != nil {
		log.Fatal("Ошибка при выборе хранилища:", err)
	}

	// Создаём приложение
	handler := handler.NewHandler(cfg, store)

	// Создаём сервер
	server := server.NewServer(cfg, *handler)

	// Запускаем сервер
	server.StartServer()
}
