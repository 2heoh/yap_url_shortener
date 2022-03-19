package handlers

import (
	"encoding/hex"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/2heoh/yap_url_shortener/cmd/shortener/services"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

var UserID string

type SignedRequest struct {
	*http.Request
}

func (r *SignedRequest) GetUserID() ([]byte, error) {
	crypto, err := services.NewCrypto()
	if err != nil {
		log.Printf("can't init crypto: %v", err)
		return nil, err
	}

	session, err := r.Cookie("session")
	if err != nil {
		log.Printf("can't find cookie 'session': %v", err)
		return nil, err
	}

	src, err := hex.DecodeString(session.Value)
	if err != nil {
		log.Printf("can't decode: %v", err)
		return nil, err
	}

	log.Printf(" = %v", src)

	return crypto.Decrypt(src)
}

func SignedCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		crypto, err := services.NewCrypto()
		if err != nil {
			log.Printf("can't create machine: %v", err)
			next.ServeHTTP(w, r)
			return
		}
		session, err := r.Cookie("session")
		if err != nil {
			log.Printf("Cant find cookie - set new")

			UserID = RandStringBytes(16)

			cookie := &http.Cookie{
				Name:    "session",
				Path:    "/",
				Expires: time.Now().AddDate(1, 0, 0),
			}

			seal := crypto.Encrypt([]byte(UserID))

			log.Printf("%v >> %v", UserID, seal)

			cookie.Value = hex.EncodeToString(seal) // зашифровываем

			log.Printf("encrypted: %v\n", cookie.Value)
			http.SetCookie(w, cookie)
			next.ServeHTTP(w, r)
			return
		}

		log.Printf("Cookie: %v\n", session.Value)

		src, err := hex.DecodeString(session.Value)
		if err != nil {
			log.Printf("can't decode: %v", err)
			next.ServeHTTP(w, r)
			return
		}

		log.Printf(" = %v", src)

		src2, err := crypto.Decrypt(src)
		if err != nil {
			log.Printf("error: %v\n", err)
			return
		}

		log.Printf("decrypted: %s", string(src2))
		UserID = string(src2)
		next.ServeHTTP(w, r)
	})
}
