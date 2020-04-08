package handlers

import (
	"github.com/azach/shorturl/lib/pool"
	"github.com/azach/shorturl/lib/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"
)

func TestCreateHandlerNoParams(t *testing.T) {
	urlStorage := storage.NewMemoryStorage()
	urlPool := pool.NewPool(urlStorage, &pool.Options{})
	shortUrlHandler := ShortURLHandler{
		urlStorage: urlStorage,
		urlPool:    urlPool,
	}

	req, err := http.NewRequest("POST", "/create", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(shortUrlHandler.Create)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusUnprocessableEntity)
}

func TestCreateHandlerSuccess(t *testing.T) {
	urlStorage := storage.NewMemoryStorage()
	urlPool := pool.NewPool(urlStorage, &pool.Options{MinPoolSize: 1, MinPoolGenerationSize: 0})
	shortUrlHandler := ShortURLHandler{
		urlStorage: urlStorage,
		urlPool:    urlPool,
	}

	data := url.Values{}
	data.Set("longUrl", "http://www.example.com")

	req, err := http.NewRequest("POST", "/create", strings.NewReader("longUrl=http://www.example.com"))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(shortUrlHandler.Create)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Regexp(t, regexp.MustCompile(`{"shortUrl":".*"}`), rr.Body.String())
}
