package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/azejke/shortener/internal/config"
	"github.com/azejke/shortener/internal/models"
	"github.com/azejke/shortener/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var cfg *config.Config
var storage *store.Store

func TestMain(m *testing.M) {
	cfg = config.InitConfig()
	storage = store.InitStore()
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestURLHandler_SearchURL(t *testing.T) {
	storage.Insert("knKvtdNoxw", "https://practicum.yandex.kz/")
	ts := httptest.NewServer(RoutesBuilder(cfg, storage))
	defer ts.Close()
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
			require.NoError(t, err)
			defer response.Body.Close()
			assert.Equal(t, tt.want.statusCode, response.StatusCode)
			assert.Equal(t, tt.want.location, response.Header.Get("Location"))
			assert.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"))
		})
	}
}

func TestURLHandler_WriteURL(t *testing.T) {
	ts := httptest.NewServer(RoutesBuilder(cfg, storage))
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
			assert.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"))
			err = response.Body.Close()
			require.NoError(t, err)
		})
	}
}

func TestURLHandler_Shorten(t *testing.T) {
	ts := httptest.NewServer(RoutesBuilder(cfg, storage))
	defer ts.Close()

	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name string
		body models.ShortenRequest
		want want
	}{
		{
			name: "Pass correct url",
			body: models.ShortenRequest{
				URL: "https://practicum.yandex.kz",
			},
			want: want{
				statusCode:  201,
				contentType: "application/json",
			},
		},
		{
			name: "Pass empty url",
			body: models.ShortenRequest{
				URL: "",
			},
			want: want{
				statusCode: 500,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.body)
			require.NoError(t, err)
			request, err := http.NewRequest(http.MethodPost, ts.URL+"/api/shorten", bytes.NewBuffer(jsonData))
			require.NoError(t, err)
			request.Header.Add("Content-Type", "application/json")
			response, err := ts.Client().Do(request)
			require.NoError(t, err)
			assert.Equal(t, tt.want.statusCode, response.StatusCode)
			if response.StatusCode == http.StatusCreated {
				assert.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"))
			}
			err = response.Body.Close()
			require.NoError(t, err)
		})
	}
}
