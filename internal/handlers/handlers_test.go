package handlers

import (
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
	ts := httptest.NewServer(RoutesBuilder(Store))
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
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
			urlID: "/knKvtdNoxw",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  307,
				location:    "https://practicum.yandex.kz/",
			},
		},
		{
			name:  "Doesn't exist id test",
			urlID: "/anBvtENHxw",
			want: want{
				statusCode: 400,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodGet, ts.URL+tt.urlID, nil)
			require.NoError(t, err)
			response, err := ts.Client().Do(request)
			defer response.Body.Close()
			require.NoError(t, err)
			assert.Equal(t, tt.want.statusCode, response.StatusCode)
			assert.Equal(t, tt.want.location, response.Header.Get("Location"))
			assert.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"))
		})
	}
}

func TestWriteURL(t *testing.T) {
	Store := make(store.Store)
	ts := httptest.NewServer(RoutesBuilder(Store))
	defer ts.Close()

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
			request, err := http.NewRequest(http.MethodPost, ts.URL, strings.NewReader(tt.url))
			require.NoError(t, err)
			request.Header.Add("Content-Type", "text/plain; charset=utf-8")
			response, err := ts.Client().Do(request)
			require.NoError(t, err)
			assert.Equal(t, tt.want.statusCode, response.StatusCode)

			if tt.url != "" {
				assert.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"))
				responseBody, err := io.ReadAll(response.Body)
				require.NoError(t, err)
				formattedResponse := strings.Split(string(responseBody), "/")
				id := formattedResponse[len(formattedResponse)-1]
				_, ok := Store[id]
				if !ok {
					t.Errorf("The URL %s are not added", tt.url)
				}
			}

			err = response.Body.Close()
			require.NoError(t, err)
		})
	}
}
