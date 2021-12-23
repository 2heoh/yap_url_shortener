package handlers

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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
				path:   "/yandex",
				method: http.MethodGet,
				body:   nil,
			},
			expected: expected{
				code:        307,
				response:    "<a href=\"https://yandex.ru/\">Temporary Redirect</a>.\n\n",
				contentType: "text/html; charset=utf-8",
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
				response:    fmt.Sprintf("http://localhost:8080/%s", GenerateId("https://google.com/")),
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
				code:        400,
				response:    "unknown method\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.request.method, tt.request.path, tt.request.body)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(RequestHandler)
			h.ServeHTTP(w, request)

			res := w.Result()
			defer func() {
				err := res.Body.Close()
				t.Errorf("Error: %v", err)
			}()

			require.Equal(t, tt.expected.code, res.StatusCode)
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			require.Equal(t, tt.expected.response, string(resBody))
			require.Equal(t, tt.expected.contentType, res.Header.Get("Content-Type"))
		})
	}
}
