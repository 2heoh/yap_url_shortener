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

	log.Fatal(
		http.ListenAndServe(
			cfg.ServerAddress,
			handlers.Zipper(handlers.NewHandler(
				services.NewShorterURL(repositories.Init(cfg)),
				cfg.BaseURL,
			)),
		),
	)
}
