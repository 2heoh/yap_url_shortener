package main

import (
	"github.com/2heoh/yap_url_shortener/cmd/shortener/config"
	"github.com/2heoh/yap_url_shortener/cmd/shortener/entities"
	"github.com/2heoh/yap_url_shortener/cmd/shortener/handlers"
	"github.com/2heoh/yap_url_shortener/cmd/shortener/repositories"
	"github.com/2heoh/yap_url_shortener/cmd/shortener/services"
	"log"
	"net/http"
)

const (
	workersCount = 5
	maxRetries   = 3
)

func main() {
	cfg, err := config.LoadArgs()
	if err != nil {
		log.Fatalf("Error parsing args: %v", err)
	}

	cfg, err = config.LoadEnvs(cfg)
	if err != nil {
		log.Fatalf("Error reading envs: %v", err)
	}

	log.Printf("Starting server at: http://%s/", cfg.ServerAddress)

	repository := repositories.Init(cfg)

	deleteChannel := make(chan entities.DeleteCandidate)

	for i := 0; i < workersCount; i++ {
		go func(workerID int) {
			defer func() {
				if x := recover(); x != nil {
					log.Printf("run time panic: %v", x)
				}
			}()

			log.Printf("start worker #%d...", workerID)
			for job := range deleteChannel {
				if err := repository.MakeDelete(job); err != nil {
					log.Printf("Deletion error: %v\n", err)

					if job.RetryCount < maxRetries {
						log.Printf("retry (left: %d)", maxRetries-job.RetryCount)
						deleteChannel <- entities.DeleteCandidate{
							Key:        job.Key,
							UserID:     job.UserID,
							RetryCount: job.RetryCount + 1,
						}
					}
				}
			}
		}(i)
	}

	log.Fatal(
		http.ListenAndServe(
			cfg.ServerAddress,
			handlers.NewHandler(
				services.NewShorterURL(repository, deleteChannel),
				cfg,
			),
		),
	)
}
