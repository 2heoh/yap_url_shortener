package main

import (
	"github.com/2heoh/yap_url_shortener/cmd/shortener/handlers"
	"github.com/2heoh/yap_url_shortener/cmd/shortener/repositories"
	"github.com/2heoh/yap_url_shortener/cmd/shortener/services"
	"log"
	"net/http"
)

func main() {
	log.Fatal(
		http.ListenAndServe(":8080",
			handlers.NewHandler(
				repositories.NewUrlRepository(), services.NewIDGenerator(),
			),
		),
	)

}
