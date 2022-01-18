package repositories

import (
	"errors"
)

type Repository interface {
	Add(id string, url string) error
	Get(id string) (string, error)
}

type InMemoryRepository struct {
	links map[string]string
}

var ErrNotFound = errors.New("id is not found")

func NewInmemoryURLRepository() Repository {
	return &InMemoryRepository{map[string]string{"yandex": "https://yandex.ru/"}}
}

func (r *InMemoryRepository) Add(id string, url string) error {
	r.links[id] = url
	return nil
}

func (r *InMemoryRepository) Get(id string) (string, error) {
	if url, found := r.links[id]; found {
		return url, nil
	}

	return "", ErrNotFound
}
