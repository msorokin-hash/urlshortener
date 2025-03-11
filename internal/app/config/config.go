package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/caarlos0/env"
)

type Config struct {
	BaseShortURL    string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	Address         string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	LogLevel        string `env:"LOG_LEVEL" envDefault:"info"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"default.json"`
}

func NewConfig() *Config {
	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		log.Fatal("Ошибка парсинга переменных окружения:", err)
	}

	addressFlag := flag.String("a", "localhost:8080", "Адрес HTTP-сервера")
	baseURLFlag := flag.String("b", "http://localhost:8080", "Базовый адрес коротких ссылок")
	logLevel := flag.String("l", "info", "Уровень логирования")
	fileStorageFlag := flag.String("f", "./urls.json", "Путь к хранилищу")

	flag.Parse()

	if *addressFlag != "" {
		cfg.Address = *addressFlag
	}

	if *baseURLFlag != "" {
		cfg.BaseShortURL = *baseURLFlag
	}

	if *logLevel != "" {
		cfg.LogLevel = *logLevel
	}

	fileStorageEnv := os.Getenv("FILE_STORAGE_PATH")
	if fileStorageEnv == "" {
		if *fileStorageFlag != "" {
			cfg.FileStoragePath = *fileStorageFlag
		}
	}

	fmt.Println("Конфигурация сервера:")
	fmt.Println("  Адрес сервера:", cfg.Address)
	fmt.Println("  Базовый URL:", cfg.BaseShortURL)
	fmt.Println("  Файл хранения данных:", cfg.FileStoragePath)

	return &cfg
}
