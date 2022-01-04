package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
}

func LoadEnvs() Config {
	var config Config
	if err := env.Parse(&config); err != nil {
		log.Fatal(err)
	}

	if config.ServerAddress == "" {
		config.ServerAddress = "localhost:8080"
	}

	if config.BaseURL == "" {
		config.BaseURL = fmt.Sprintf("http://%s", config.ServerAddress)
		log.Printf("env BASE_URL is not set, took server address: %s", config.BaseURL)
	}

	return config
}
