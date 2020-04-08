package handlers

import (
	"encoding/json"
	"github.com/azach/shorturl/lib/storage"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func (h *ShortURLHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortURL := vars["shortURL"]
	now := time.Now()

	minuteHits, err := h.urlStorage.GetHits(shortURL, now, storage.Minute)
	dailyHits, err := h.urlStorage.GetHits(shortURL, now, storage.Daily)
	weeklyHits, err := h.urlStorage.GetHits(shortURL, now, storage.Weekly)
	allTimeHits, err := h.urlStorage.GetHits(shortURL, now, storage.AllTime)

	err = json.NewEncoder(w).Encode(map[string]int64{
		"minute":   minuteHits,
		"daily":    dailyHits,
		"weekly":   weeklyHits,
		"all_time": allTimeHits,
	})

	if err != nil {
		logrus.Errorf("error generating stats response: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *ShortURLHandler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortURL := vars["shortURL"]
	longURL, found := h.urlStorage.Get(shortURL)

	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	now := time.Now()
	h.urlStorage.Hit(shortURL, now)
	logrus.Infof("redirecting short url: %s", shortURL)
	http.Redirect(w, r, longURL, http.StatusSeeOther)
	return
}
