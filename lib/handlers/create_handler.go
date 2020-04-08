package handlers

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
)

func (h *ShortURLHandler) Create(w http.ResponseWriter, r *http.Request) {
	longURL := r.FormValue("longUrl")

	_, err := url.ParseRequestURI(longURL)
	if err != nil {
		logrus.Infof("invalid long url: %s", longURL)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	shortURL := h.urlPool.Get()
	err = h.urlStorage.Set(shortURL, longURL)
	if err != nil {
		logrus.Errorf("error setting short url: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logrus.Infof("creating short url for long url: %s", longURL)

	err = json.NewEncoder(w).Encode(map[string]string{"shortUrl": shortURL})
	if err != nil {
		logrus.Errorf("error encoding JSON response: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
