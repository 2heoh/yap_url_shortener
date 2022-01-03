package services_test

import (
	"errors"
	"github.com/2heoh/yap_url_shortener/cmd/shortener/services"
	"github.com/stretchr/testify/require"
	"testing"
)

type TestableRepo struct{}

func (tr *TestableRepo) Add(url, id string) {}
func (tr *TestableRepo) Get(id string) (string, error) {
	if id == "non-existing" {
		return "", errors.New("id is not found: " + id)
	}

	return "https://example.com/", nil
}
func TestShorterURLCreation(t *testing.T) {
	t.Parallel()

	service := services.NewShorterURL(&TestableRepo{})

	id, err := service.CreateURL("https://example.com")

	require.NoError(t, err)
	require.Equal(t, "96248650", id)
}

func TestShorterURLCreationWhenURLIsEmpty(t *testing.T) {
	t.Parallel()

	service := services.NewShorterURL(&TestableRepo{})

	id, err := service.CreateURL(string([]byte{}))

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
