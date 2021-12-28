package repositories_test

import (
	"testing"

	"github.com/2heoh/yap_url_shortener/cmd/shortener/repositories"
	"github.com/stretchr/testify/require"
)

func TestAddToRepository(t *testing.T) {
	t.Parallel()

	repository := repositories.NewURLRepository()

	repository.Add("test", "https://example.com")

	url, err := repository.Get("test")

	require.NoError(t, err)
	require.Equal(t, url, "https://example.com")
}
