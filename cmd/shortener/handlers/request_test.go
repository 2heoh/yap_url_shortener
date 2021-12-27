package handlers

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestableRepo struct{}

func (tr *TestableRepo) Add(url, id string) {

}

func (tr *TestableRepo) Get(id string) (string, error) {
	if id == "non-existing" {
		return "", errors.New("id is not found: " + id)
	}

	return "https://example.com/", nil
}

type TestableGenerator struct{}

func (tg *TestableGenerator) Generate(url string) string {
	return "test_url"
}

func TestRequestHandler(t *testing.T) {
	type expected struct {
		code        int
		response    string
		contentType string
	}

	type request struct {
		method string
		path   string
		body   io.Reader
	}

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
				// когда text/html charset почему-то капсом
				contentType: "text/html; charset=UTF-8",
				code:        200,
			},
			// chi каким-то магическим способом сразу умадряется отдавать страницу
			// что тут проверять не очень понятно
			//expected: expected{
			//	code:        307,
			//	response:    "<a href=\"https://example.com/\">Temporary Redirect</a>.",
			//
			//},
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
				response:    "http://localhost:8080/test_url",
				contentType: "text/html; charset=utf-8",
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
	}
	testRepo := &TestableRepo{}
	testGenerator := &TestableGenerator{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewHandler(testRepo, testGenerator)
			ts := httptest.NewServer(r)
			defer ts.Close()
			req, err := http.NewRequest(tt.request.method, ts.URL+tt.request.path, tt.request.body)
			require.NoError(t, err)

			res, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
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
				assert.Equal(t, tt.expected.response, string(body))
			}
		})
	}
}
