package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_, err = w.Write([]byte(`{"result":"ok"}`))

	if err != nil {
		log.Printf("Error: %v", err)
	}
	return
}
