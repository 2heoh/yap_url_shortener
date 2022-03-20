package handlers

import (
	"net/http"
)

func (h *Handler) PingDB(w http.ResponseWriter, r *http.Request) {

	err := h.urls.Ping()
	if err != nil {
		http.Error(w, "can't connect to db", http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)

	return
}
