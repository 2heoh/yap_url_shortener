package repositories

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAddToRepository(t *testing.T) {
	repository := NewURLRepository()
	repository.Add("https://example.com", "test")

	url, err := repository.Get("test")

	require.NoError(t, err)
	require.Equal(t, url, "https://example.com")
}
