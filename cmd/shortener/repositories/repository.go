package repositories

import (
	"github.com/2heoh/yap_url_shortener/cmd/shortener/entities"
)

type Repository interface {
	Get(id string) (string, error)
	Add(id string, url string, userID string) error
	AddBatch(urls []entities.URLItem, userID string) ([]entities.ShortenURL, error)
	GetAllFor(userID string) []entities.LinkItem
	Ping() error
}
