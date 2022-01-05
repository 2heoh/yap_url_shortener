package config_test

import (
	"testing"

	"github.com/2heoh/yap_url_shortener/cmd/shortener/config"
	"github.com/stretchr/testify/require"
)

func TestDefaultValueForBaseURL(t *testing.T) {
	t.Parallel()

	cfg, err := config.LoadEnvs()

	require.NoError(t, err)
	require.Equal(t, "http://localhost:8080", cfg.BaseURL)
}

func TestDefaultValueForServerAddress(t *testing.T) {
	t.Parallel()

	cfg, err := config.LoadEnvs()

	require.NoError(t, err)
	require.Equal(t, "localhost:8080", cfg.ServerAddress)
}

func TestSetValueForBaseURL(t *testing.T) {
	t.Setenv("BASE_URL", "test://base/url")

	cfg, err := config.LoadEnvs()

	require.NoError(t, err)
	require.Equal(t, "test://base/url", cfg.BaseURL)
}

func TestSetValueForServerAddress(t *testing.T) {
	t.Setenv("SERVER_ADDRESS", "HOST:PORT")

	cfg, err := config.LoadEnvs()

	require.NoError(t, err)
	require.Equal(t, "HOST:PORT", cfg.ServerAddress)
}

func TestSetValueForFileStoragePath(t *testing.T) {
	t.Setenv("FILE_STORAGE_PATH", "/path/to/file")

	cfg, err := config.LoadEnvs()

	require.NoError(t, err)
	require.Equal(t, "/path/to/file", cfg.FileStoragePath)
}
