package services_test

import (
	"errors"
	"testing"

	"github.com/2heoh/yap_url_shortener/cmd/shortener/entities"
	"github.com/2heoh/yap_url_shortener/cmd/shortener/services"
	"github.com/stretchr/testify/require"
)

type TestableRepo struct{}

func (tr *TestableRepo) DeleteBatch(keys []string, userID string) error {
	//TODO implement me
	panic("implement me")
}

func (tr *TestableRepo) AddBatch(urls []entities.URLItem, userID string) ([]entities.ShortenURL, error) {
	//TODO implement me
	panic("implement me")
}

func (tr *TestableRepo) Ping() error {
	//TODO implement me
	panic("implement me")
}

func (tr *TestableRepo) Add(id string, url string, userID string) error {
	return nil
}

func (tr *TestableRepo) GetAllFor(userID string) []entities.LinkItem {
	//TODO implement me
	panic("implement me")
}

func (tr *TestableRepo) GetAll(userID string) []entities.LinkItem {
	//TODO implement me
	panic("implement me")
}

func (tr *TestableRepo) Get(id string) (*entities.LinkItem, error) {
	if id == "non-existing" {
		return nil, errors.New("id is not found: " + id)
	}

	return &entities.LinkItem{id, "https://example.com/", false}, nil
}
func TestShorterURLCreation(t *testing.T) {
	t.Parallel()

	service := services.NewShorterURL(&TestableRepo{})

	id, err := service.CreateURL("https://example.com", "1")

	require.NoError(t, err)
	require.Equal(t, "96248650", id)
}

func TestShorterURLCreationWhenURLIsEmpty(t *testing.T) {
	t.Parallel()

	service := services.NewShorterURL(&TestableRepo{})

	id, err := service.CreateURL(string([]byte{}), "1")

	t.Logf(id, err)
	require.Equal(t, services.ErrEmptyURL, err)
	require.Equal(t, "", id)
}

func TestShorterURLRetrieving(t *testing.T) {
	t.Parallel()

	service := services.NewShorterURL(&TestableRepo{})

	url, err := service.RetrieveURL("test")

	require.NoError(t, err)
	require.Equal(t, "https://example.com/", url)
}

func TestShorterURLRetrievingEmptyID(t *testing.T) {
	t.Parallel()

	service := services.NewShorterURL(&TestableRepo{})

	url, err := service.RetrieveURL("")

	require.Equal(t, err, services.ErrEmptyID)
	require.Equal(t, "", url)
}
