package services

import (
	"errors"
	"log"

	"github.com/2heoh/yap_url_shortener/cmd/shortener/repositories"
)

type Shorter interface {
	CreateURL(url string, userID string) (string, error)
	RetrieveURL(id string) (string, error)
	RetrieveURLsForUser(id string) ([]repositories.LinkItem, error)
	Ping() error
}

var (
	ErrEmptyURL     = errors.New("url is empty")
	ErrEmptyID      = errors.New("id is empty")
	ErrIDIsNotFound = errors.New("id is not found")
)

type ShorterURL struct {
	repository repositories.Repository
}

func (s *ShorterURL) Ping() error {
	return s.repository.Ping()
}

func NewShorterURL(repo repositories.Repository) *ShorterURL {
	return &ShorterURL{repo}
}

func (s *ShorterURL) CreateURL(url string, userID string) (string, error) {
	if url == "" {
		return "", ErrEmptyURL
	}

	id := GenerateID(url)

	log.Printf("userID: %v", userID)

	s.repository.Add(id, url, userID)

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

func (s *ShorterURL) RetrieveURLsForUser(id string) ([]repositories.LinkItem, error) {
	result := s.repository.GetAllFor(id)

	return result, nil
}
