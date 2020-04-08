package main

import (
	"fmt"
	"github.com/azach/shorturl/lib/handlers"
	"github.com/azach/shorturl/lib/pool"
	"github.com/azach/shorturl/lib/storage"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func main() {
	r := mux.NewRouter()

	urlStorage := storage.NewRedisStorage()
	urlPool := pool.NewPool(urlStorage, &pool.Options{})
	shortURLHandler := handlers.NewShortUrlHandler(urlStorage, urlPool)
	// Continuously generate candidate words in the background
	go shortURLHandler.GeneratePool()

	actionsRouter := r.Methods("POST").Subrouter()
	actionsRouter.HandleFunc("/create", shortURLHandler.Create)

	r.HandleFunc("/{shortURL}", shortURLHandler.Get)
	r.HandleFunc("/{shortURL}/stats", shortURLHandler.GetStats)

	svcPort := fmt.Sprintf(":%s", os.Getenv("PORT"))
	logrus.Infof("running service on port %s", svcPort)
	logrus.Fatal(http.ListenAndServe(svcPort, r))
}
