package config

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

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

type NetAddress struct {
	Host string
	Port int
}

func (na *NetAddress) String() string {
	return fmt.Sprintf("%s:%d", na.Host, na.Port)
}

func (na *NetAddress) Set(flagValue string) error {
	parts := strings.Split(flagValue, ":")
	na.Host = parts[0]
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}
	na.Port = port
	return nil
}

func LoadEnvs(config *Config) (*Config, error) {

	if err := env.Parse(config); err != nil {
		return config, ErrConfigIsNotLoaded
	}

	if config.ServerAddress == "" {
		config.ServerAddress = serverAddress
	}

	if config.BaseURL == "" {
		config.BaseURL = fmt.Sprintf("http://%s", config.ServerAddress)
		log.Printf("env BASE_URL is not set, took server address: %s", config.BaseURL)
	}

	return config, nil
}
