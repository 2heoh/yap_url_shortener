package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/2heoh/yap_url_shortener/cmd/shortener/services"
)

type JSONRequestBody struct {
	URL string `json:"url"`
}

type JSONResponseBody struct {
	Error  string `json:"error,omitempty"`
	Result string `json:"result,omitempty"`
}

func (h *Handler) PostJSONURL(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.ReturnJSONError(w, fmt.Sprintf("can't read body: %v", err))

		return
	}

	request := JSONRequestBody{}
	err = json.Unmarshal(body, &request)
	if err != nil {
		h.ReturnJSONError(w, fmt.Sprintf("bad json: %v", err))
		return
	}

	srequest := SignedRequest{r}

	userID, err := srequest.GetUserID()
	if err != nil {
		log.Printf("Error: %v", err)
	}

	id, err := h.urls.CreateURLForUser(request.URL, string(userID))

	if errors.Is(err, services.ErrEmptyURL) {
		h.ReturnJSONError(w, "missed url")

		return
	}

	h.ReturnJSONResponse(w, h.baseURL+"/"+id)

	if err != nil {
		log.Printf("Error: %v", err)
	}
}

func (h *Handler) ReturnJSONError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	jsonResponse, err := json.Marshal(JSONResponseBody{Error: message})

	if err != nil {
		log.Printf("Error: %v", err)
	}

	_, err = w.Write(jsonResponse)

	if err != nil {
		log.Printf("Error: %v", err)
	}

}

func (h *Handler) ReturnJSONResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	jsonResponse, err := json.Marshal(JSONResponseBody{Result: message})

	if err != nil {
		log.Printf("Error: %v", err)
	}

	_, err = w.Write(jsonResponse)

	if err != nil {
		log.Printf("Error: %v", err)
	}

}
