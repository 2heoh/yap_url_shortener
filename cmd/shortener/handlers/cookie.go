package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/2heoh/yap_url_shortener/cmd/shortener/services"
)

var UserID string

type SignedRequest struct {
	*http.Request
}

func HandleSignedCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		crypto, err := services.NewCrypto()
		if err != nil {
			log.Printf("can't create crypto service: %v", err)
			next.ServeHTTP(w, r)

			return
		}
		session, err := r.Cookie("session")
		if err != nil {
			log.Printf("\nCant find cookie - set new")
			UserID = crypto.GenerateUserID()
			http.SetCookie(w, &http.Cookie{
				Name:    "session",
				Path:    "/",
				Expires: time.Now().AddDate(1, 0, 0),
				Value:   crypto.GetEncodedSessionValue(UserID),
			})
			next.ServeHTTP(w, r)

			return
		}

		UserID, err = crypto.GetDecodedUserID(session.Value)
		if err != nil {
			log.Printf("can't get UserID: %v", err)
			next.ServeHTTP(w, r)

			return
		}

		next.ServeHTTP(w, r)
	})
}
