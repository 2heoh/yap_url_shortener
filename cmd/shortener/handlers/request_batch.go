package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/2heoh/yap_url_shortener/cmd/shortener/entities"
	"io"
	"log"
	"net/http"
)

func (h *Handler) PostBatch(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.ReturnJSONError(w, fmt.Sprintf("can't read body: %v", err))

		return
	}

	var urls []entities.URLItem
	err = json.Unmarshal(body, &urls)
	if err != nil {
		h.ReturnJSONError(w, fmt.Sprintf("bad json: %v", err))
		return
	}

	shortenURLS, err := h.urls.CreateBatch(urls, UserID)
	if err != nil {
		h.ReturnJSONError(w, fmt.Sprintf("can't get add: %s", err))
		return
	}

	var result []entities.ShortenResultURL
	for _, item := range shortenURLS {
		result = append(result, entities.ShortenResultURL{
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
