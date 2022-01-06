package main

import (
	"log"
	"net/http"

	"github.com/2heoh/yap_url_shortener/cmd/shortener/config"
	"github.com/2heoh/yap_url_shortener/cmd/shortener/handlers"
	"github.com/2heoh/yap_url_shortener/cmd/shortener/repositories"
	"github.com/2heoh/yap_url_shortener/cmd/shortener/services"
)

func main() {
	cfg, err := config.LoadArgs()
	if err != nil {
		log.Fatalf("Error parsing args: %v", err)
	}

	cfg, err = config.LoadEnvs(cfg)

	if err != nil {
		log.Fatalf("Error reading envs: %v", err)
	}

	log.Printf("Starting server at: http://%s/", cfg.ServerAddress)

	var repo repositories.Repository = repositories.NewInmemoryURLRepository()

	if cfg.FileStoragePath != "" {
		log.Printf("used file storage: %s", cfg.FileStoragePath)
		repo = repositories.NewFileURLRepository(cfg.FileStoragePath)
	} else {
		log.Println("Use in memory storage")
	}

	log.Fatal(
		http.ListenAndServe(
			cfg.ServerAddress,
			handlers.NewHandler(
				services.NewShorterURL(repo),
				cfg.BaseURL,
			),
		),
	)
}
