package main

import (
	"github.com/2heoh/yap_url_shortener/cmd/shortener/handlers"
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("/favicon.ico", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", handlers.RequestHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
