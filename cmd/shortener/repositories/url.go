package repositories

import (
	"fmt"
)

type Repository interface {
	Add(url string, id string)
	Get(id string) (string, error)
}

type URLRepository struct{}

var links = map[string]string{
	"yandex": "https://yandex.ru/",
}

func NewURLRepository() *URLRepository {
	return &URLRepository{}
}

func (r *URLRepository) Add(url string, id string) {
	links[id] = url
}

func (r *URLRepository) Get(id string) (string, error) {
	if url, found := links[id]; found {
		return url, nil
	}

	return "", fmt.Errorf("id is not found: %v", id)
}
