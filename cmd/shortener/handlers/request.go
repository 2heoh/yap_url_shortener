package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/2heoh/yap_url_shortener/cmd/shortener/config"
	"github.com/2heoh/yap_url_shortener/cmd/shortener/entities"
	"github.com/2heoh/yap_url_shortener/cmd/shortener/repositories"
	"io"
	"log"
	"net/http"

	"github.com/2heoh/yap_url_shortener/cmd/shortener/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handler struct {
	*chi.Mux
	urls   services.Shorter
	config *config.Config
}

func NewHandler(service services.Shorter, config *config.Config) *Handler {
	h := &Handler{
		Mux:    chi.NewMux(),
		urls:   service,
		config: config,
	}

	h.Use(middleware.Logger)
	h.Use(HandleSignedCookie)
	h.Use(Zipper)

	h.Post("/", h.PostURL)
	h.Get("/ping", h.PingDB)
	h.Post("/api/shorten", h.PostJSONURL)
	h.Get("/{id}", h.GetURL)
	h.Get("/api/user/urls", h.GetURLSForUser)
	h.Post("/api/shorten/batch", h.PostBatch)
	h.Get("/", func(w http.ResponseWriter, request *http.Request) {
		http.Error(w, "empty id", http.StatusBadRequest)
	})

	return h
}

func (h *Handler) GetURLSForUser(w http.ResponseWriter, r *http.Request) {
	log.Printf(" UserID: %s ", UserID)
	urls, err := h.urls.RetrieveURLsForUser(UserID)

	if err != nil {
		log.Printf("can't get urls: %v", err)
	}

	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var result []entities.LinkItem
	for _, url := range urls {
		result = append(result, entities.LinkItem{
			OriginalURL: url.OriginalURL,
			ShortURL:    h.config.BaseURL + "/" + url.ShortURL,
		})
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	body, err := json.Marshal(result)
	if err != nil {
		log.Printf("json serialization error: %v", err)
	}

	_, err = w.Write(body)
	if err != nil {
		log.Printf("Error: %v", err)
	}
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
	log.Printf("Body: '%s'", string(b))
	id, err := h.urls.CreateURL(string(b), UserID)
	if errors.Is(err, services.ErrEmptyURL) {
		http.Error(w, "missed url", http.StatusBadRequest)

		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if errors.Is(err, repositories.ErrKeyExists) {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	_, err = w.Write([]byte(fmt.Sprintf("%s/%s", h.config.BaseURL, id)))

	if err != nil {
		log.Printf("Error: %v", err)
	}
}
