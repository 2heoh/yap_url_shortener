package repositories

import (
	"errors"
)

type Repository interface {
	Add(url string, id string)
	Get(id string) (string, error)
}

type UrlRepository struct{}

var links = map[string]string{
	"yandex": "https://yandex.ru/",
}

func NewURLRepository() *UrlRepository {
	return &UrlRepository{}
}

func (r *UrlRepository) Add(url string, id string) {
	links[id] = url
}

func (r *UrlRepository) Get(id string) (string, error) {
	if url, found := links[id]; found {
		return url, nil
	}

	return "", errors.New("id is not found: " + id)
}
