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

	cfg := config.LoadEnvs()

	log.Printf("Server started at: http://%s/", cfg.ServerAddress)
	log.Fatal(
		http.ListenAndServe(
			cfg.ServerAddress,
			handlers.NewHandler(
				services.NewShorterURL(repositories.NewURLRepository()),
				cfg.BaseURL,
			),
		),
	)
}
