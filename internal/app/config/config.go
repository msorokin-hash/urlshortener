package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env"
)

type Config struct {
	BaseShortURL    string `env:"BASE_URL"`
	Address         string `env:"SERVER_ADDRESS"`
	LogLevel        string `env:"LOG_LEVEL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
}

func NewConfig() *Config {
	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		log.Fatal("Ошибка парсинга переменных окружения:", err)
	}

	addressFlag := flag.String("a", "localhost:8080", "Адрес HTTP-сервера")
	baseURLFlag := flag.String("b", "http://localhost:8080", "Базовый адрес коротких ссылок")
	logLevel := flag.String("l", "info", "Уровень логирования")
	fileStorageFlag := flag.String("f", "default.json", "Путь к хранилищу")
	dbDsn := flag.String("d", "", "Строка подключения к БД")

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

	if cfg.FileStoragePath == "" {
		cfg.FileStoragePath = *fileStorageFlag
	}

	if cfg.DatabaseDSN == "" {
		cfg.DatabaseDSN = *dbDsn
	}

	log.Println("Конфигурация сервера:")
	log.Println("  Адрес сервера:", cfg.Address)
	log.Println("  Базовый URL:", cfg.BaseShortURL)
	log.Println("  Уровень логирования:", cfg.LogLevel)

	return &cfg
}
