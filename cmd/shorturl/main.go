package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"

	"github.com/azach/shorturl/lib/cache"
	"github.com/azach/shorturl/lib/pool"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var (
	port string
)

func init() {
	flag.StringVar(&port, "port", "8080", "web service port")
	flag.Parse()
}

type shortURLHandler struct {
	urlCache cache.Cache
	urlPool  *pool.Pool
}

func main() {
	r := mux.NewRouter()

	urlCache := cache.NewMemoryCache()
	urlPool := pool.NewPool(urlCache)

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

	svcPort := fmt.Sprintf(":%s", port)
	logrus.Infof("running service on port %s", port)
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
		logrus.Infof("redirecting short url: %s", shortURL)
		http.Redirect(w, r, longURL, http.StatusSeeOther)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}
