package services

import (
	"errors"

	"github.com/2heoh/yap_url_shortener/cmd/shortener/repositories"
)

type Shorter interface {
	CreateURL(url string) (string, error)
	RetrieveURL(id string) (string, error)
}

var (
	ErrEmptyURL     = errors.New("url is empty")
	ErrEmptyID      = errors.New("id is empty")
	ErrIDIsNotFound = errors.New("id is not found")
)

type ShorterURL struct {
	repository repositories.Repository
}

func NewShorterURL(repo repositories.Repository) *ShorterURL {
	return &ShorterURL{repo}
}

func (s *ShorterURL) CreateURL(url string) (string, error) {
	if url == "" {
		return "", ErrEmptyURL
	}

	id := GenerateID(url)

	s.repository.Add(id, url)

	return id, nil
}

func (s *ShorterURL) RetrieveURL(id string) (string, error) {
	if id == "" {
		return "", ErrEmptyID
	}

	url, err := s.repository.Get(id)
	if err != nil {
		return "", ErrIDIsNotFound
	}

	return url, nil
}
