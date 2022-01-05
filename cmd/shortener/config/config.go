package config

import (
	"errors"
	"fmt"
	"log"

	"github.com/caarlos0/env/v6"
)

var (
	ErrConfigIsNotLoaded = errors.New("can't parse env")
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

func LoadEnvs() (Config, error) {
	var config Config
	if err := env.Parse(&config); err != nil {
		return config, ErrConfigIsNotLoaded
	}

	if config.ServerAddress == "" {
		config.ServerAddress = "localhost:8080"
	}

	if config.BaseURL == "" {
		config.BaseURL = fmt.Sprintf("http://%s", config.ServerAddress)
		log.Printf("env BASE_URL is not set, took server address: %s", config.BaseURL)
	}

	return config, nil
}
