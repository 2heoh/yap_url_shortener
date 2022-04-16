package repositories_test

import (
	"os"
	"testing"

	"github.com/2heoh/yap_url_shortener/cmd/shortener/repositories"
	"github.com/stretchr/testify/require"
)

func TestFileStorageNonEmptyStorageWriteAndReadOneURL(t *testing.T) {
	os.WriteFile("./test.db", []byte("a;b;c\n"), 0777)
	defer func() {
		if !t.Failed() {
			os.Remove("./test.db")
		}
	}()

	repository := repositories.NewFileURLRepository("./test.db")

	err := repository.Add("test", "https://example.com", "1")
	require.NoError(t, err)

	url, err := repository.Get("test")

	require.NoError(t, err)
	require.Equal(t, "https://example.com", url.OriginalURL)
}

func TestFileStorageEmptyStorageWriteAndReadTwoURLs(t *testing.T) {

	os.WriteFile("./test_empty.db", []byte(""), 0666)
	defer func() {
		if !t.Failed() {
			os.Remove("./test_empty.db")
		}
	}()

	repository := repositories.NewFileURLRepository("./test.db")

	err := repository.Add("test", "https://example.com", "1")
	require.NoError(t, err)

	url, err := repository.Get("test")

	require.NoError(t, err)
	require.Equal(t, "https://example.com", url.OriginalURL)

	err = repository.Add("test2", "https://example2.com", "2")
	require.NoError(t, err)

	url, err = repository.Get("test2")

	require.NoError(t, err)
	require.Equal(t, "https://example2.com", url.OriginalURL)
}

func TestFileStorageWithExistingStorageAddingSameRecord(t *testing.T) {
	os.WriteFile("./test.db", []byte("user;id;url\n"), 0777)
	defer func() {
		if !t.Failed() {
			os.Remove("./test.db")
		}
	}()

	repository := repositories.NewFileURLRepository("./test.db")

	err := repository.Add("id", "url", "user")
	require.NoError(t, err)

	content, _ := os.ReadFile("./test.db")
	require.Equal(t, "user;id;url\n", string(content))
}

func TestFileStorageReadAllIDs(t *testing.T) {

	defer func() {
		if !t.Failed() {
			os.Remove("./test.db")
		}
	}()

	os.WriteFile("./test.db", []byte("a;b;c\nt;v;u\n1;test;https://example.com\n"), 0666)

	repository := repositories.NewFileURLRepository("./test.db")

	urls := repository.GetAllFor("a")

	require.Equal(t, len(urls), 1)
	require.Equal(t, urls[0].OriginalURL, "c")
	require.Equal(t, urls[0].ShortURL, "b")

	urls = repository.GetAllFor("t")

	require.Equal(t, len(urls), 1)
	require.Equal(t, urls[0].OriginalURL, "u")
	require.Equal(t, urls[0].ShortURL, "v")

	urls = repository.GetAllFor("1")

	require.Equal(t, len(urls), 1)

	require.Equal(t, urls[0].OriginalURL, "https://example.com")
	require.Equal(t, urls[0].ShortURL, "test")

}
