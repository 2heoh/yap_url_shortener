package repositories

import (
	"errors"
)

type Repository interface {
	Add(id string, url string) error
	Get(id string) (string, error)
}

type URLRepository struct {
	links map[string]string
}

var ErrNotFound = errors.New("id is not found")

func NewURLRepository() *URLRepository {
	return &URLRepository{map[string]string{"yandex": "https://yandex.ru/"}}
}

func (r *URLRepository) Add(id string, url string) {
	r.links[id] = url
}

func (r *URLRepository) Get(id string) (string, error) {
	if url, found := r.links[id]; found {
		return url, nil
	}

	return "", ErrNotFound
}
