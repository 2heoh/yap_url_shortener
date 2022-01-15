package repositories

import (
	"os"
	"reflect"
	"testing"

	"github.com/2heoh/yap_url_shortener/cmd/shortener/config"
	"github.com/stretchr/testify/require"
)

func TestInitReturnsInMemoryStorage(t *testing.T) {
	cfg := &config.Config{
		FileStoragePath: "",
	}

	repo := Init(cfg)

	require.Equal(t, "*repositories.InMemoryRepository", reflect.TypeOf(repo).String())
}

func TestInitReturnsFileStorage(t *testing.T) {
	defer func() {
		os.Remove("./test")
	}()

	cfg := &config.Config{
		FileStoragePath: "./test",
	}

	repo := Init(cfg)

	require.Equal(t, "*repositories.FileURLRepository", reflect.TypeOf(repo).String())
}
