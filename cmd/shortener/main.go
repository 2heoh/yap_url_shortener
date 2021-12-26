package main

import (
	"github.com/2heoh/yap_url_shortener/cmd/shortener/handlers"
	"log"
	"net/http"
)

var links = map[string]string{
	"yandex": "https://yandex.ru/",
}

func main() {
	log.Fatal(http.ListenAndServe(":8080", handlers.CreateHandler(links)))
}
