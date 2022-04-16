package repositories

import (
	"github.com/2heoh/yap_url_shortener/cmd/shortener/config"
	"log"
)

func Init(cfg *config.Config) Repository {
	if cfg.FileStoragePath != "" {
		log.Printf("used file storage: %s", cfg.FileStoragePath)

		return NewFileURLRepository(cfg.FileStoragePath)
	}

	if cfg.DSN != "" {
		log.Printf("used database: %s", cfg.DSN)

		return NewDatabaseRepository(cfg.DSN)
	}

	log.Println("Use in memory storage")

	return NewInmemoryURLRepository()
}
