package repositories

import (
	"errors"
	"github.com/2heoh/yap_url_shortener/cmd/shortener/entities"
)

type InMemoryRepository struct {
	links       map[string]string
	linksByUser map[string][]entities.LinkItem
}

func (r *InMemoryRepository) AddBatch(urls []entities.URLItem, userID string) ([]entities.ShortenURL, error) {
	var result []entities.ShortenURL
	for _, item := range urls {
		err := r.Add(item.CorrelationID, item.OriginalURL, userID)
		if err != nil {
			return nil, err
		}
		result = append(result, entities.ShortenURL{Key: item.CorrelationID})
	}
	return result, nil
}

var ErrNotFound = errors.New("id is not found")

func (r *InMemoryRepository) Ping() error {
	return nil
}

func (r *InMemoryRepository) GetAllFor(userID string) []entities.LinkItem {
	if links, found := r.linksByUser[userID]; found {
		return links
	}

	return nil
}

func NewInmemoryURLRepository() Repository {

	return &InMemoryRepository{
		map[string]string{"yandex": "https://yandex.ru/"},
		map[string][]entities.LinkItem{"1": nil},
	}
}

func (r *InMemoryRepository) Add(id string, url string, userID string) error {
	r.links[id] = url
	r.linksByUser[userID] = append(r.linksByUser[userID], entities.LinkItem{ShortURL: id, OriginalURL: url})

	return nil
}

func (r *InMemoryRepository) Get(id string) (string, error) {
	if url, found := r.links[id]; found {
		return url, nil
	}

	return "", ErrNotFound
}
