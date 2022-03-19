package repositories

import (
	"errors"
	"fmt"
	"log"
)

type Repository interface {
	Get(id string) (string, error)
	AddBy(id string, url string, userID string) error
	GetAllBy(userID string) []LinkItem
}

type LinkItem struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type InMemoryRepository struct {
	links       map[string]string
	linksByUser map[string][]LinkItem
}

func (r *InMemoryRepository) GetAllBy(userID string) []LinkItem {

	fmt.Printf(" ==%v \n", r.linksByUser[userID])

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

func (r *InMemoryRepository) AddBy(id string, url string, userID string) error {
	r.links[id] = url
	r.linksByUser[userID] = append(r.linksByUser[userID], LinkItem{id, url})
	log.Printf(" [%s => %s] ", userID, id)
	return nil
}

//func (r *InMemoryRepository) Add(id string, url string) error {
//	r.links[id] = url
//	log.Printf(" // %s => %s ", url, id)
//	return nil
//}

func (r *InMemoryRepository) Get(id string) (string, error) {
	if url, found := r.links[id]; found {
		return url, nil
	}

	return "", ErrNotFound
}
