package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/rwiteshbera/URL-Shortener-Go/routes"
)

func main() {
	config := routes.LoadConfig()

	router := mux.NewRouter()

	server := &http.Server{
		Handler: router,
		Addr:    config.SERVER_BASE_URL,
	}

	// Server Health Check
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}).Methods(http.MethodGet)

	router.HandleFunc("/shorten", config.ShortenURL).Methods(http.MethodPost)
	router.HandleFunc("/{id}", config.ResolveURL).Methods(http.MethodGet)

	// Admin API
	// Delete Expired URL
	router.HandleFunc("/purge", config.PurgeDatabase).Methods(http.MethodPost)

	log.Println("Listening on:", config.SERVER_BASE_URL)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.ListenAndServe(); err != nil {
			log.Println(err.Error())
		}
	}()
	defer wg.Wait()

}
