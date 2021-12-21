package main

import (
	"fmt"
	"hash/crc32"
	"log"
	"net/http"
	"strings"
)

var links = map[string]string{
	"yandex": "https://yandex.ru/",
}

func RequestHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("%v", links)

	switch r.Method {
	case "POST":
		url := r.FormValue("url")
		w.WriteHeader(http.StatusCreated)
		id := GenerateId(url)
		log.Printf("url: %s, id: %s", url, id)
		links[id] = url
		fmt.Fprintln(w, id)
		return
	case "GET":
		parts := strings.Split(r.URL.String(), "/")
		if parts[1] == "" {
			http.Error(w, "empty id", http.StatusBadRequest)
			return
		}

		log.Printf(" get id: %s", parts[1])
		if url, found := links[parts[1]]; found {
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		} else {
			http.Error(w, "id is not found: "+parts[1], http.StatusNotFound)
		}
	default:
		http.Error(w, "unknown method", http.StatusBadRequest)
	}
}

func GenerateId(url string) string {
	crc32q := crc32.MakeTable(0xD5828281)
	return fmt.Sprintf("%08x", crc32.Checksum([]byte(url), crc32q))
}

func main() {
	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("/favicon.ico", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", RequestHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
