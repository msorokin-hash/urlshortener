package config

import (
	"flag"
	"fmt"
)

type Config struct {
	Address      string
	BaseShortURL string
}

func InitConfig() *Config {
	address := flag.String("a", "localhost:8080", "Адрес HTTP-сервера")
	baseURL := flag.String("b", "http://localhost:8080", "Базовый адрес коротких ссылок")

	flag.Parse()

	fmt.Println("Конфигурация сервиса:")
	fmt.Println("  Адрес сервера:", *address)
	fmt.Println("  Базовый URL:", *baseURL)

	return &Config{
		Address:      *address,
		BaseShortURL: *baseURL,
	}
}
