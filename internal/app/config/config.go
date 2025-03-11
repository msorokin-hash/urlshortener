package config

import (
	"flag"
	"os"
)

type Config struct {
	BaseShortURL    string
	Address         string
	LogLevel        string
	FileStoragePath string
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Parse() {

	c.parseFlags()

	serverAdd := os.Getenv("SERVER_ADDRESS")
	if serverAdd != "" {
		c.Address = serverAdd
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL != "" {
		c.BaseShortURL = baseURL
	}

	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		c.LogLevel = envLogLevel
	}

	if envPathDB := os.Getenv("FILE_STORAGE_PATH"); envPathDB != "" {
		c.FileStoragePath = envPathDB
	}
}

func (c *Config) parseFlags() {
	flag.StringVar(&c.Address, "a", "localhost:8080", "Адрес HTTP-сервера")
	flag.StringVar(&c.BaseShortURL, "b", "http://localhost:8080", "Базовый адрес коротких ссылок")
	flag.StringVar(&c.FileStoragePath, "f", "./urls.json", "Уровень логирования")
	flag.StringVar(&c.LogLevel, "l", "info", "Путь к файлу хранилища")
	flag.Parse()
}
