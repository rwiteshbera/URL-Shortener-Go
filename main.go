package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/rwiteshbera/URL-Shortener-Go/routes"
)

func main() {
	config, err := routes.LoadConfig()
	if err != nil {
		log.Println(err.Error())
	}

	router := mux.NewRouter()

	server := &http.Server{
		Handler: router,
		Addr:    config.SERVER_BASE_URL,
	}

	router.HandleFunc("/shorten", config.ShortenURL).Methods(http.MethodPost)
	router.HandleFunc("/{id}", config.ResolveURL).Methods(http.MethodGet)

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
