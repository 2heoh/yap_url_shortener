package services

import (
	"errors"
	"log"

	"github.com/2heoh/yap_url_shortener/cmd/shortener/entities"
	"github.com/2heoh/yap_url_shortener/cmd/shortener/repositories"
)

type Shorter interface {
	CreateURL(url string, userID string) (string, error)
	CreateBatch(urls []entities.URLItem, userID string) ([]entities.ShortenURL, error)
	RetrieveURL(id string) (string, error)
	RetrieveURLsForUser(id string) ([]entities.LinkItem, error)
	Ping() error
	DeleteBatch(keys []string, userID string) error
}

var (
	ErrEmptyURL     = errors.New("url is empty")
	ErrEmptyID      = errors.New("id is empty")
	ErrIDIsNotFound = errors.New("id is not found")
	ErrDeletedID    = errors.New("id deleted")
)

type ShorterURL struct {
	repository repositories.Repository
}

func (s *ShorterURL) DeleteBatch(keys []string, userID string) error {
	return s.repository.DeleteBatch(keys, userID)
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

	var id = GenerateID(url)

	log.Printf("userID: %v", userID)

	err := s.repository.Add(id, url, userID)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (s *ShorterURL) RetrieveURL(id string) (string, error) {
	if id == "" {
		return "", ErrEmptyID
	}

	url, err := s.repository.Get(id)

	if url.IsDeleted {
		return "", ErrDeletedID
	}

	if err != nil {
		return "", ErrIDIsNotFound
	}

	return url.OriginalURL, nil
}

func (s *ShorterURL) RetrieveURLsForUser(id string) ([]entities.LinkItem, error) {
	result := s.repository.GetAllFor(id)

	return result, nil
}

func (s *ShorterURL) CreateBatch(urls []entities.URLItem, userID string) ([]entities.ShortenURL, error) {
	return s.repository.AddBatch(urls, userID)
}
