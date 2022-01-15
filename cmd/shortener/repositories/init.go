package repositories

import (
	"log"

	"github.com/2heoh/yap_url_shortener/cmd/shortener/config"
)

func Init(cfg *config.Config) Repository {
	if cfg.FileStoragePath != "" {
		log.Printf("used file storage: %s", cfg.FileStoragePath)
		return NewFileURLRepository(cfg.FileStoragePath)
	}

	log.Println("Use in memory storage")
	return NewInmemoryURLRepository()

}
