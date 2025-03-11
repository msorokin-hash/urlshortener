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

	// Создаём in-memory хранилище
	store, err := storage.NewFileStorage("urls.json")
	if err != nil {
		log.Fatal("Ошибка инициализации файлового хранилища:", err)
	}

	// Создаём приложение
	handler := handler.NewHandler(cfg, store)

	// Создаём сервер
	server := server.NewServer(cfg, *handler)

	// Запускаем сервер
	server.StartServer()
}
