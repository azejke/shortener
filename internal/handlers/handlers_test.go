package handlers

import (
	"fmt"
	"github.com/azejke/shortener/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSearchURL(t *testing.T) {
	var Store = store.Store{
		"knKvtdNoxw": "https://practicum.yandex.kz/",
	}
	type want struct {
		contentType string
		statusCode  int
		location    string
	}
	tests := []struct {
		name  string
		urlID string
		want  want
	}{
		{
			name:  "Exist id test",
			urlID: "knKvtdNoxw",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  307,
				location:    "https://practicum.yandex.kz/",
			},
		},
		{
			name:  "Doesn't exist id test",
			urlID: "anBvtENHxw",
			want: want{
				statusCode: 400,
			},
		},
		{
			name:  "Empty id test",
			urlID: "",
			want: want{
				statusCode: 400,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:8080/%s", tt.urlID), nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				SearchURL(writer, request, Store)
			})
			h(w, request)
			result := w.Result()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.location, result.Header.Get("Location"))
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			err := result.Body.Close()
			require.NoError(t, err)
		})
	}
}

func TestWriteURL(t *testing.T) {
	var Store = make(store.Store)
	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "Pass correct url",
			url:  "https://practicum.yandex.kz",
			want: want{
				statusCode:  201,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Pass empty url",
			url:  "",
			want: want{
				statusCode: 400,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/", strings.NewReader(tt.url))
			request.Header.Add("Content-Type", "text/plain; charset=utf-8")
			w := httptest.NewRecorder()
			h := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				WriteURL(writer, request, Store)
			})
			h(w, request)
			result := w.Result()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)

			if tt.url != "" {
				assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
				response, err := io.ReadAll(result.Body)
				require.NoError(t, err)
				formattedResponse := strings.Split(string(response), "/")
				id := formattedResponse[len(formattedResponse)-1]
				_, ok := Store[id]
				if !ok {
					t.Errorf("The URL %s are not added", tt.url)
				}
			}

			err := result.Body.Close()
			require.NoError(t, err)
		})
	}
}
