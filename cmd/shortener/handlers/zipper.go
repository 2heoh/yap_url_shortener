package handlers

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	log.Printf("GZIP BODY: %s", string(b))

	return w.Writer.Write(b)
}

func Zipper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Printf("Accept-Encoding: %v", r.Header.Get("Accept-Encoding"))
		log.Printf("Content-Encoding: %v", r.Header.Get("Content-Encoding"))

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// не забывайте потом закрыть *gzip.Reader
			defer gz.Close()

			// при чтении вернётся распакованный слайс байт
			body, err := io.ReadAll(gz)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			log.Printf("* Decoded body: %v", body)
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		w.Header().Set("Content-Encoding", "gzip")
		if err != nil {
			log.Printf("Error: %v", err)
		}

		defer gz.Close()

		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

func DebugRequest(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf, bodyErr := ioutil.ReadAll(r.Body)
		if bodyErr != nil {
			log.Print("bodyErr ", bodyErr.Error())
			http.Error(w, bodyErr.Error(), http.StatusInternalServerError)
			return
		}

		rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
		rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))
		log.Printf("\n= DEBUG ============\n%s %s\nAccept-Encoding: %s\nContent-Encoding: %s\n= BODY: ==============\n%v\n= END BODY: ==========", r.Method, r.URL, r.Header.Get("Accept-Encoding"), r.Header.Get("Content-Encoding"), rdr1)
		r.Body = rdr2
		next.ServeHTTP(w, r)
	})
}
