package repositories

import (
	"errors"
	"github.com/2heoh/yap_url_shortener/cmd/shortener/entities"
	"log"
)

type InMemoryRepository struct {
	links       map[string]string
	linksByUser map[string][]entities.LinkItem
}

func NewInmemoryURLRepository() Repository {
	return &InMemoryRepository{
		map[string]string{},
		map[string][]entities.LinkItem{},
	}
}

func (r *InMemoryRepository) MakeDelete(candidate entities.DeleteCandidate) error {
	log.Printf("  -> %v \n", candidate)

	if urls, found := r.linksByUser[candidate.UserID]; found {
		for i, item := range urls {
			if candidate.Key == item.ShortURL {
				log.Printf("  \\_/ %v \n", candidate)
				r.linksByUser[candidate.UserID][i].IsDeleted = true
			}
		}
	}

	return nil
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

func (r *InMemoryRepository) Add(id string, url string, userID string) error {
	r.links[id] = url
	r.linksByUser[userID] = append(r.linksByUser[userID], entities.LinkItem{ShortURL: id, OriginalURL: url})

	return nil
}

func (r *InMemoryRepository) Get(id string) (*entities.LinkItem, error) {
	if url, found := r.links[id]; found {
		return &entities.LinkItem{ShortURL: id, OriginalURL: url}, nil
	}

	return nil, ErrNotFound
}

func (r *InMemoryRepository) GetShortenURL(key string, userID string) *entities.LinkItem {
	if urls, found := r.linksByUser[userID]; found {

		for _, link := range urls {

			if link.ShortURL == key {
				return &link
			}
		}
	}
	return nil
}
