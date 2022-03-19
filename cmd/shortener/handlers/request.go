package handlers

import (
	"encoding/json"
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
	urls    services.Shorter
	baseURL string
}

func NewHandler(service services.Shorter, baseURL string) *Handler {
	h := &Handler{
		Mux:     chi.NewMux(),
		urls:    service,
		baseURL: baseURL,
	}

	h.Use(middleware.Logger)
	h.Use(SignedCookie)
	h.Use(Zipper)

	h.Post("/", h.PostURL)
	h.Post("/api/shorten", h.PostJSONURL)
	h.Get("/{id}", h.GetURL)
	h.Get("/api/user/urls", h.GetURLSForUser)
	h.Get("/", func(w http.ResponseWriter, request *http.Request) {
		http.Error(w, "empty id", http.StatusBadRequest)
	})

	return h
}

func (h *Handler) GetURLSForUser(w http.ResponseWriter, r *http.Request) {
	request := SignedRequest{r}

	id, err := request.GetUserID()
	if err != nil {
		log.Printf("Error: %v", err)
	}

	log.Printf(" UserID: %s ", id)
	urls, err := h.urls.RetrieveURLsForUser(string(id))
	if err != nil {
		log.Printf("can't get urls: %v", err)
	}

	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	for i, url := range urls {
		fmt.Printf(" * %v \n", url.ShortURL)
		urls[i].ShortURL = h.baseURL + "/" + url.ShortURL
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	body, err := json.Marshal(urls)
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
	srequest := SignedRequest{r}
	userID, err := srequest.GetUserID()
	if err != nil {
		log.Printf("Error: %v", err)
	}
	id, err := h.urls.CreateURLForUser(string(b), string(userID))

	if errors.Is(err, services.ErrEmptyURL) {
		http.Error(w, "missed url", http.StatusBadRequest)

		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(fmt.Sprintf("%s/%s", h.baseURL, id)))
	if err != nil {
		log.Printf("Error: %v", err)
	}
}
