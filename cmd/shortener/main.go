package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/2heoh/yap_url_shortener/cmd/shortener/handlers"
	"github.com/2heoh/yap_url_shortener/cmd/shortener/repositories"
	"github.com/2heoh/yap_url_shortener/cmd/shortener/services"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
}

func main() {
	var config Config

	if err := env.Parse(&config); err != nil {
		log.Fatal(err)
	}

	if config.ServerAddress == "" {
		config.ServerAddress = "localhost:8080"
	}
	if config.BaseURL == "" {
		config.BaseURL = fmt.Sprintf("http://%s/", config.ServerAddress)
		log.Printf("env BASE_URL is not set, took server address: %s", config.BaseURL)
	}

	log.Printf("Server started at: http://%s/", config.ServerAddress)
	log.Fatal(
		http.ListenAndServe(
			config.ServerAddress,
			handlers.NewHandler(
				services.NewShorterURL(repositories.NewURLRepository()),
				config.BaseURL,
			),
		),
	)
}
