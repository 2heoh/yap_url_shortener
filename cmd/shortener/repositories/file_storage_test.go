package repositories_test

import (
	"os"
	"testing"

	"github.com/2heoh/yap_url_shortener/cmd/shortener/repositories"
	"github.com/stretchr/testify/require"
)

func TestFileStorageEmptyStorageWriteAndReadOneURL(t *testing.T) {
	os.WriteFile("./test.db", []byte("a;b\n"), 0777)
	defer func() {
		if !t.Failed() {
			os.Remove("./test.db")
		}
	}()

	repository := repositories.NewFileURLRepository("./test.db")

	err := repository.Add("test", "https://example.com")
	require.NoError(t, err)

	url, err := repository.Get("test")

	require.NoError(t, err)
	require.Equal(t, "https://example.com", url)
}

func TestFileStorageWithExistingStorageAddingSameRecord(t *testing.T) {
	os.WriteFile("./test.db", []byte("a;b\n"), 0777)
	defer func() {
		if !t.Failed() {
			os.Remove("./test.db")
		}
	}()

	repository := repositories.NewFileURLRepository("./test.db")

	err := repository.Add("a", "b")
	require.NoError(t, err)

	content, err := os.ReadFile("./test.db")
	require.Equal(t, "a;b\n", string(content))
}

func TestFileStorage(t *testing.T) {

	defer func() {
		if !t.Failed() {
			os.Remove("./test.db")
		}
	}()

	os.WriteFile("./test.db", []byte("a;b\nt;v\ntest;https://example.com"), 0777)

	repository := repositories.NewFileURLRepository("./test.db")

	url, err := repository.Get("a")

	require.NoError(t, err)
	require.Equal(t, url, "b")
}
