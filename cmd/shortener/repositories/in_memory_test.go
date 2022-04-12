package repositories_test

import (
	"testing"

	"github.com/2heoh/yap_url_shortener/cmd/shortener/repositories"
	"github.com/stretchr/testify/require"
)

func TestAddToRepository(t *testing.T) {
	t.Parallel()

	repository := repositories.NewInmemoryURLRepository()

	repository.Add("test", "https://example.com", "id")

	url, err := repository.Get("test")

	require.NoError(t, err)
	require.Equal(t, "https://example.com", url.OriginalURL)
}

//func TestMarkItemDeleted(t *testing.T) {
//	repository := &repositories.InMemoryRepository{
//		links: map[string]string{"yandex": "https://yandex.ru/"},
//		map[string][]entities.LinkItem{"1": nil},
//		make(chan entities.DeleteCandidate),
//		&sync.WaitGroup{},
//	}
//
//	repository.Add("a", "https://test1.com/", "1")
//	repository.Add("b", "https://test2.com/", "1")
//	repository.Add("c", "https://test3.com/", "1")
//
//	keys := []string{"a", "b"}
//	err := repository.DeleteBatch(keys, "1")
//
//	link, _ := repository.Get("a")
//	require.NoError(t, err)
//	require.False(t, link.IsDeleted)
//
//	repository.ProcessDelete()
//	repository.ProcessDelete()
//
//	link1 := repository.GetShortenURL("a", "1")
//	link2 := repository.GetShortenURL("b", "1")
//	link3 := repository.GetShortenURL("c", "1")
//	require.True(t, link1.IsDeleted)
//	require.True(t, link2.IsDeleted)
//	require.False(t, link3.IsDeleted)
//}
