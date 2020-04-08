package handlers

import (
	"github.com/azach/shorturl/lib/pool"
	"github.com/azach/shorturl/lib/storage"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHandlerNotFound(t *testing.T) {
	urlStorage := storage.NewMemoryStorage()
	urlPool := pool.NewPool(urlStorage, &pool.Options{})
	shortUrlHandler := ShortURLHandler{
		urlStorage: urlStorage,
		urlPool:    urlPool,
	}

	req, err := http.NewRequest("GET", "/foo", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/{shortURL}", shortUrlHandler.Get)
	router.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusNotFound)
}

func TestGetHandlerSuccess(t *testing.T) {
	urlStorage := storage.NewMemoryStorage()
	urlPool := pool.NewPool(urlStorage, &pool.Options{MinPoolSize: 1, MinPoolGenerationSize: 0})
	shortUrlHandler := ShortURLHandler{
		urlStorage: urlStorage,
		urlPool:    urlPool,
	}

	err := urlStorage.Set("foo", "http://www.example.com")
	assert.NoError(t, err)

	req, err := http.NewRequest("GET", "/foo", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/{shortURL}", shortUrlHandler.Get)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusSeeOther, rr.Code)
}

func TestGetHandlerStats(t *testing.T) {
	urlStorage := storage.NewMemoryStorage()
	urlPool := pool.NewPool(urlStorage, &pool.Options{MinPoolSize: 1, MinPoolGenerationSize: 0})
	shortUrlHandler := ShortURLHandler{
		urlStorage: urlStorage,
		urlPool:    urlPool,
	}

	err := urlStorage.Set("foo", "http://www.example.com")
	assert.NoError(t, err)

	req, err := http.NewRequest("GET", "/foo/stats", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/{shortURL}/stats", shortUrlHandler.GetStats)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "{\"all_time\":0,\"daily\":0,\"minute\":0,\"weekly\":0}\n", rr.Body.String())
}
