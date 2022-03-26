package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/2heoh/yap_url_shortener/cmd/shortener/services"
)

type ShortenResultURL struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func (h *Handler) PostBatch(w http.ResponseWriter, r *http.Request) {
	log.Printf("run batch")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.ReturnJSONError(w, fmt.Sprintf("can't read body: %v", err))

		return
	}

	var urls []services.URLItem
	err = json.Unmarshal(body, &urls)
	if err != nil {
		h.ReturnJSONError(w, fmt.Sprintf("bad json: %v", err))
		return
	}

	log.Printf("URLS: %v ", urls)

	shortenURLS, err := h.urls.CreateBatch(urls, UserID)
	if err != nil {
		h.ReturnJSONError(w, fmt.Sprintf("can't get urls: %s", err))
		return
	}

	var result []ShortenResultURL
	for _, item := range shortenURLS {
		result = append(result, ShortenResultURL{
			CorrelationID: item.Key,
			ShortURL:      fmt.Sprintf("%s/%s", h.config.BaseURL, item.Key),
		})
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)

	jsonResponse, err := json.Marshal(result)
	if err != nil {
		log.Printf("Error: %v", err)
	}

	_, err = w.Write(jsonResponse)

	if err != nil {
		log.Printf("Error: %v", err)
	}
}
