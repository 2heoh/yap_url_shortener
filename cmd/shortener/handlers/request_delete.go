package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Response struct {
	result string `json:"result"`
}

func (h *Handler) DeleteBatch(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.ReturnJSONError(w, fmt.Sprintf("can't read body: %v", err))

		return
	}

	var keys []string
	err = json.Unmarshal(body, &keys)
	if err != nil {
		h.ReturnJSONError(w, fmt.Sprintf("bad json: %v", err))
		return
	}
	log.Printf(" keys: %v", keys)

	err = h.urls.DeleteBatch(keys, UserID)
	if err != nil {
		h.ReturnJSONError(w, fmt.Sprintf("can't delete: %s", err))
		return
	}

	h.ReturnJSONResponse(w, "ok", http.StatusAccepted)
}
