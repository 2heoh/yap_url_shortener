package handlers

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/2heoh/yap_url_shortener/cmd/shortener/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handler struct {
	*chi.Mux
	urls services.Shorter
}

func NewHandler(service services.Shorter) *Handler {
	h := &Handler{
		Mux:  chi.NewMux(),
		urls: service,
	}

	h.Use(middleware.Logger)
	h.Post("/", h.PostURL)
	h.Post("/api/shorten", h.PostJSONURL)
	h.Get("/{id}", h.GetURL)
	h.Get("/", func(w http.ResponseWriter, request *http.Request) {
		http.Error(w, "empty id", http.StatusBadRequest)
	})

	return h
}

func (h *Handler) GetURL(w http.ResponseWriter, r *http.Request) {
	url, err := h.urls.RetrieveURL(chi.URLParam(r, "id"))
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

	id, err := h.urls.CreateURL(string(b))

	if errors.Is(err, services.ErrEmptyURL) {
		http.Error(w, "missed url", http.StatusBadRequest)

		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(fmt.Sprintf("http://localhost:8080/%s", id)))

	if err != nil {
		log.Printf("Error: %v", err)
	}
}
