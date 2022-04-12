package handlers

import (
	"log"
	"net/http"

	"github.com/2heoh/yap_url_shortener/cmd/shortener/repositories"
)

func (h *Handler) PingDB(w http.ResponseWriter, r *http.Request) {
	repo := repositories.NewDatabaseRepository(h.config.DSN, nil)
	if err := repo.Ping(); err != nil {
		log.Printf("Ping error: %v", err)
		http.Error(w, "can't connect to db", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
