package handlers

import (
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"net/http"
	"strings"
)

var links = map[string]string{
	"yandex": "https://yandex.ru/",
}

func RequestHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		url := string(b)
		if url == "" {
			http.Error(w, "missed url", http.StatusBadRequest)
			return
		}
		id := GenerateId(url)
		log.Printf("url: %s, id: %s", url, id)
		links[id] = url
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		_, err = fmt.Fprintf(w, fmt.Sprintf("http://localhost:8080/%s", id))
		if err != nil {
			log.Printf("Error: %v", err)
		}
		return
	case "GET":
		parts := strings.Split(r.URL.String(), "/")
		if parts[1] == "" {
			http.Error(w, "empty id", http.StatusBadRequest)
			return
		}

		if url, found := links[parts[1]]; found {
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		} else {
			http.Error(w, "id is not found: "+parts[1], http.StatusNotFound)
		}
		return
	default:
		http.Error(w, "unknown method", http.StatusBadRequest)
	}
}

func GenerateId(url string) string {
	crc32q := crc32.MakeTable(0xD5828281)
	return fmt.Sprintf("%08x", crc32.Checksum([]byte(url), crc32q))
}
