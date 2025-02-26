package config

import (
	"flag"
	"os"
)

type Config struct {
	Address      string
	BaseShortURL string
}

func InitConfig() *Config {
	addressEnv := os.Getenv("SERVER_ADDRESS")
	baseURLEnv := os.Getenv("BASE_URL")

	addressFlag := flag.String("a", "localhost:8080", "Адрес HTTP-сервера")
	baseURLFlag := flag.String("b", "http://localhost:8080", "Базовый адрес коротких ссылок")

	flag.Parse()

	address := addressEnv
	if address == "" {
		if *addressFlag != "" {
			address = *addressFlag
		} else {
			address = ":8080"
		}
	}

	baseURL := baseURLEnv
	if baseURL == "" {
		if *baseURLFlag != "" {
			baseURL = *baseURLFlag
		} else {
			baseURL = "http://localhost:8080"
		}
	}

	return &Config{
		Address:      address,
		BaseShortURL: baseURL,
	}
}
