package handlers_test

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/2heoh/yap_url_shortener/cmd/shortener/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestZipperMiddleware(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		expected expected
		request  request
	}{
		{
			name: "Content-Encoding: gzip",
			request: request{
				method: http.MethodPost,
				path:   "/",
				body:   strings.NewReader("https://google.com/"),
				headers: map[string]string{
					"Accept-Encoding": "gzip",
				},
			},
			expected: expected{
				code:        201,
				response:    "test_url",
				contentType: "text/html; charset=utf-8",
			},
		},
	}
	testURLService := &TestableService{}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			r := handlers.Zipper(handlers.NewHandler(testURLService, "http://test"))
			ts := httptest.NewServer(r)
			defer ts.Close()
			req, err := http.NewRequest(tt.request.method, ts.URL+tt.request.path, tt.request.body)

			for _, key := range tt.request.headers {
				req.Header.Set(key, tt.request.headers[key])
			}

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
			gz, _ := gzip.NewReader(bytes.NewBuffer(body))
			defer gz.Close()
			decripted, _ := io.ReadAll(gz)

			require.NoError(t, err)
			assert.Equal(t, tt.expected.code, res.StatusCode)
			assert.Equal(t, tt.expected.contentType, res.Header.Get("Content-Type"))
			if tt.expected.response != "" {
				assert.Contains(t, string(decripted), tt.expected.response)
			}
		})
	}
}
