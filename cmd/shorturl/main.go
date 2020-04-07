package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/azach/shorturl/lib/cache"
	"github.com/azach/shorturl/lib/pool"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type shortURLHandler struct {
	urlCache cache.Cache
	urlPool  *pool.Pool
}

func main() {
	r := mux.NewRouter()

	urlCache := cache.NewRedisCache()
	urlPool := pool.NewPool(urlCache, &pool.Options{})

	shortURLHandler := &shortURLHandler{
		urlCache: urlCache,
		urlPool:  urlPool,
	}

	// Continuously generate candidate words in the background
	go func() {
		for {
			shortURLHandler.urlPool.Generate()
		}
	}()

	actionsRouter := r.Methods("POST").Subrouter()
	actionsRouter.HandleFunc("/create", shortURLHandler.Create)

	r.HandleFunc("/{cache:\\w+}", shortURLHandler.Get)

	svcPort := fmt.Sprintf(":%s", os.Getenv("PORT"))
	logrus.Infof("running service on port %s", svcPort)
	logrus.Fatal(http.ListenAndServe(svcPort, r))
}

func (h *shortURLHandler) Create(w http.ResponseWriter, r *http.Request) {
	longURL := r.FormValue("longurl")

	logrus.Infof("longurl: %s", longURL)

	_, err := url.ParseRequestURI(longURL)
	if err != nil {
		logrus.Infof("invalid long url: %s", longURL)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	shortURL := h.urlPool.Get()
	err = h.urlCache.Set(shortURL, longURL)
	if err != nil {
		logrus.Errorf("error setting short url: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]string{"shorturl": shortURL})
	if err != nil {
		logrus.Errorf("error encoding JSON response: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *shortURLHandler) Get(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Path[1:]
	longURL, found := h.urlCache.Get(shortURL)

	if found {
		now := time.Now()
		h.urlCache.Hit(shortURL, now)
		logrus.Infof("redirecting short url: %s", shortURL)
		hits, err := h.urlCache.GetHits(shortURL, now, cache.Minute)
		if err != nil {
			logrus.Errorf("error getting hits: %s", err)
		} else {
			logrus.Infof("total hits: %v", hits)
		}

		http.Redirect(w, r, longURL, http.StatusSeeOther)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}
