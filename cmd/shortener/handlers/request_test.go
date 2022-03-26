package handlers_test

import (
	"errors"
	"github.com/2heoh/yap_url_shortener/cmd/shortener/config"
	"github.com/2heoh/yap_url_shortener/cmd/shortener/repositories"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/2heoh/yap_url_shortener/cmd/shortener/handlers"
	"github.com/2heoh/yap_url_shortener/cmd/shortener/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestableService struct{}

func (tg *TestableService) CreateBatch(urls []services.URLItem, userID string) ([]services.ShortenURL, error) {

	var items = []services.ShortenURL{{Key: urls[0].CorrelationID}}

	return items, nil
}

func (tg *TestableService) Ping() error {
	//TODO implement me
	panic("implement me")
}

func (tg *TestableService) CreateURL(url string, userID string) (string, error) {
	if url == "" {
		return "", services.ErrEmptyURL
	}

	return "test_url", nil
}

func (tg *TestableService) RetrieveURLsForUser(id string) ([]repositories.LinkItem, error) {
	if id == "non-existing" {
		return nil, errors.New("id is not found: " + id)
	}

	url := repositories.LinkItem{ShortURL: "test", OriginalURL: "https://example.com/"}

	result := []repositories.LinkItem{url}

	return result, nil

}

func (tg *TestableService) RetrieveURL(id string) (string, error) {
	if id == "non-existing" {
		return "", errors.New("id is not found: " + id)
	}

	return "https://example.com/", nil
}

type expected struct {
	code        int
	response    string
	contentType string
}

type request struct {
	method  string
	path    string
	body    io.Reader
	headers map[string]string
}

func TestRequestHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		expected expected
		request  request
	}{
		{
			name: "get request with no id",
			request: request{
				method: http.MethodGet,
				path:   "/",
				body:   nil,
			},
			expected: expected{
				code:        400,
				response:    "empty id\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "get request with existing id",
			request: request{
				path:   "/test",
				method: http.MethodGet,
				body:   nil,
			},
			expected: expected{
				contentType: "text/html; charset=UTF-8",
				code:        200,
			},
		},
		{
			name: "get request with not-existing id",
			request: request{
				path:   "/non-existing",
				method: http.MethodGet,
				body:   nil,
			},
			expected: expected{
				code:        404,
				response:    "id is not found: non-existing\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "post request with no url",
			request: request{
				method: http.MethodPost,
				path:   "/",
				body:   nil,
			},
			expected: expected{
				code:        400,
				response:    "missed url\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "post request with url",
			request: request{
				method: http.MethodPost,
				path:   "/",
				body:   strings.NewReader("https://google.com/"),
			},
			expected: expected{
				code:        201,
				response:    "http://test_host/test_url",
				contentType: "text/html; charset=utf-8",
			},
		},
		{
			name: "post request /api/shorten with broken json body returns 400",
			request: request{
				method: http.MethodPost,
				path:   "/api/shorten",
				body:   strings.NewReader(""),
			},
			expected: expected{
				code:        400,
				response:    `{"error":"bad json:`,
				contentType: "application/json",
			},
		},
		{
			name: "post request /api/shorten with json body returns json with shorten url",
			request: request{
				method: http.MethodPost,
				path:   "/api/shorten",
				body:   strings.NewReader(`{"url": "https://google.com/"}`),
			},
			expected: expected{
				code:        201,
				response:    `{"result":"http://test_host/test_url"}`,
				contentType: "application/json",
			},
		},
		{
			name: "post request /api/shorten with json with empty url returns 400",
			request: request{
				method: http.MethodPost,
				path:   "/api/shorten",
				body:   strings.NewReader(`{"url": ""}`),
			},
			expected: expected{
				code:        400,
				response:    `{"error":"missed url"}`,
				contentType: "application/json",
			},
		},
		{
			name: "unsupported http method",
			request: request{
				method: http.MethodPut,
				path:   "/",
				body:   strings.NewReader("https://google.com/"),
			},
			expected: expected{
				code: 405,
			},
		},
		{
			name: "post request /api/shorten/batch with proper json return batch of shorten urls",
			request: request{
				method: http.MethodPost,
				path:   "/api/shorten/batch",
				body: strings.NewReader(
					`[{"correlation_id": "test_key","original_url": "https://example.com"}]`,
				),
			},
			expected: expected{
				code:        201,
				response:    `[{"correlation_id":"test_key","short_url":"http://test_host/test_key"}]`,
				contentType: "application/json; charset=utf-8",
			},
		},
	}
	testURLService := &TestableService{}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			conf := &config.Config{BaseURL: "http://test_host"}
			r := handlers.NewHandler(testURLService, conf)
			ts := httptest.NewServer(r)
			defer ts.Close()
			req, err := http.NewRequest(tt.request.method, ts.URL+tt.request.path, tt.request.body)

			//for _, key := range tt.request.headers {
			//	req.Header.Set(key, tt.request.headers[key])
			//}

			require.NoError(t, err)
			res, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)
			body, err := ioutil.ReadAll(res.Body)
			defer func() {
				err := res.Body.Close()
				if err != nil {
					t.Fatalf("Error closing body: %v", err)
				}
			}()

			require.NoError(t, err)
			assert.Equal(t, tt.expected.code, res.StatusCode)
			assert.Equal(t, tt.expected.contentType, res.Header.Get("Content-Type"))
			if tt.expected.response != "" {
				assert.Contains(t, string(body), tt.expected.response)
			}
		})
	}
}
