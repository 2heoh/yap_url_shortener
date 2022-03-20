package repositories

import (
	"errors"
)

type Repository interface {
	Get(id string) (string, error)
	Add(id string, url string, userID string) error
	GetAllFor(userID string) []LinkItem
}

type LinkItem struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type InMemoryRepository struct {
	links       map[string]string
	linksByUser map[string][]LinkItem
}

func (r *InMemoryRepository) GetAllFor(userID string) []LinkItem {

	if links, found := r.linksByUser[userID]; found {
		return links
	}

	return nil
}

var ErrNotFound = errors.New("id is not found")

func NewInmemoryURLRepository() Repository {

	return &InMemoryRepository{
		map[string]string{"yandex": "https://yandex.ru/"},
		map[string][]LinkItem{"1": nil},
	}
}

func (r *InMemoryRepository) Add(id string, url string, userID string) error {
	r.links[id] = url
	r.linksByUser[userID] = append(r.linksByUser[userID], LinkItem{id, url})

	return nil
}

func (r *InMemoryRepository) Get(id string) (string, error) {
	if url, found := r.links[id]; found {
		return url, nil
	}

	return "", ErrNotFound
}
