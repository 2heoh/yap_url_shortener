package handlers_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/2heoh/yap_url_shortener/cmd/shortener/config"
	"github.com/2heoh/yap_url_shortener/cmd/shortener/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (tg *TestableService) DeleteBatch(keys []string, userID string) error {

	if keys[0] == "error" {
		return errors.New("error")
	}

	return nil
}

func TestRequestDeleteHandler(t *testing.T) {
	tests := []struct {
		name     string
		expected expected
		request  request
	}{
		{
			name: "delete request with set of ids",
			request: request{
				method: http.MethodDelete,
				path:   "/api/user/urls",
				body:   strings.NewReader(`["a","b"]`),
			},
			expected: expected{
				code:        202,
				response:    `{"result":"ok"}`,
				contentType: "application/json",
			},
		},
		{
			name: "delete request with set of ids",
			request: request{
				method: http.MethodDelete,
				path:   "/api/user/urls",
				body:   strings.NewReader(`["error"]`),
			},
			expected: expected{
				code:        400,
				response:    "can't delete: error",
				contentType: "application/json",
			},
		},
		{
			name: "delete request with bad json",
			request: request{
				method: http.MethodDelete,
				path:   "/api/user/urls",
				body:   strings.NewReader(`[bad json]`),
			},
			expected: expected{
				code:        400,
				response:    "bad json",
				contentType: "application/json",
			},
		},
		{
			name: "get request with deleted id",
			request: request{
				path:   "/deleted",
				method: http.MethodGet,
				body:   nil,
			},
			expected: expected{
				code:        410,
				response:    "id deleted\n",
				contentType: "text/plain; charset=utf-8",
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
