package main

import (
	"log"
	"net/http"

	"github.com/2heoh/yap_url_shortener/cmd/shortener/handlers"
	"github.com/2heoh/yap_url_shortener/cmd/shortener/repositories"
	"github.com/2heoh/yap_url_shortener/cmd/shortener/services"
)

func main() {
	log.Fatal(
		http.ListenAndServe(":8080",
			handlers.NewHandler(
				repositories.NewURLRepository(), services.NewIDGenerator(),
			),
		),
	)
}
