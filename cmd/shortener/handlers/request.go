package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"hash/crc32"
	"io"
	"log"
	"net/http"
)

type Handler struct {
	*chi.Mux
	links map[string]string
}

func CreateHandler(links map[string]string) *Handler {
	h := &Handler{
		Mux:   chi.NewMux(),
		links: links,
	}

	h.Use(middleware.Logger)

	h.Post("/", h.PostUrlHandler)
	h.Get("/{id}", h.GetUrlHandler)
	h.Get("/", func(w http.ResponseWriter, request *http.Request) {
		http.Error(w, "empty id", http.StatusBadRequest)
	})

	return h
}

func (h *Handler) GetUrlHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if url, found := h.links[id]; found {
		http.Redirect(w, r, url, http.StatusFound)
		return
	}
	http.Error(w, "id is not found: "+id, http.StatusNotFound)
	return
}

func (h *Handler) PostUrlHandler(w http.ResponseWriter, r *http.Request) {
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
	h.links[id] = url
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(fmt.Sprintf("http://localhost:8080/%s", id)))
	if err != nil {
		log.Printf("Error: %v", err)
	}
	return
}

func GenerateId(url string) string {
	crc32q := crc32.MakeTable(0xD5828281)
	return fmt.Sprintf("%08x", crc32.Checksum([]byte(url), crc32q))
}
