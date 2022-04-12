package repositories

import (
	"github.com/2heoh/yap_url_shortener/cmd/shortener/entities"
	"log"
	"time"

	"github.com/2heoh/yap_url_shortener/cmd/shortener/config"
)

const workersCount = 5

func Init(cfg *config.Config) Repository {
	if cfg.FileStoragePath != "" {
		log.Printf("used file storage: %s", cfg.FileStoragePath)

		return NewFileURLRepository(cfg.FileStoragePath)
	}

	if cfg.DSN != "" {
		log.Printf("used database: %s", cfg.DSN)

		deleteChannel := make(chan entities.DeleteCandidate)

		for i := 0; i < workersCount; i++ {
			go func() {
				defer func() {
					if x := recover(); x != nil {
						log.Printf("run time panic: %v", x)
					}
				}()

				log.Printf("start worker...")
				for job := range deleteChannel {
					repo := NewDatabaseRepository(cfg.DSN, deleteChannel)
					if err := repo.MakeDelete(job); err != nil {
						log.Printf("retry...")
						time.Sleep(time.Second * 2)
						deleteChannel <- job
					}
				}
			}()
		}

		return NewDatabaseRepository(cfg.DSN, deleteChannel)

	}

	log.Println("Use in memory storage")

	return NewInmemoryURLRepository()
}
