package repositories_test

import (
	"testing"

	"github.com/2heoh/yap_url_shortener/cmd/shortener/repositories"
	"github.com/stretchr/testify/require"
)

func TestAddToRepository(t *testing.T) {
	t.Parallel()

	repository := repositories.NewInmemoryURLRepository()

	repository.AddBy("test", "https://example.com", "id")

	url, err := repository.Get("test")

	require.NoError(t, err)
	require.Equal(t, url, "https://example.com")
}
