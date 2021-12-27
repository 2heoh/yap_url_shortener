package handlers

import (
	"fmt"
	"github.com/2heoh/yap_url_shortener/cmd/shortener/repositories"
	"github.com/2heoh/yap_url_shortener/cmd/shortener/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"log"
	"net/http"
)

type Handler struct {
	*chi.Mux
	urlRepository repositories.Repository
	idGenerator   services.Generator
}

func NewHandler(repo repositories.Repository, generator services.Generator) *Handler {
	var h = &Handler{
		Mux:           chi.NewMux(),
		urlRepository: repo,
		idGenerator:   generator,
	}

	h.Use(middleware.Logger)
	h.Post("/", h.PostURL)
	h.Get("/{id}", h.GetURL)
	h.Get("/", func(w http.ResponseWriter, request *http.Request) {
		http.Error(w, "empty id", http.StatusBadRequest)
	})

	return h
}

func (h *Handler) GetURL(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	url, err := h.urlRepository.Get(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)

		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) PostURL(w http.ResponseWriter, r *http.Request) {
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

	id := h.idGenerator.Generate(url)
	h.urlRepository.Add(url, id)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(fmt.Sprintf("http://localhost:8080/%s", id)))

	if err != nil {
		log.Printf("Error: %v", err)
	}
}
