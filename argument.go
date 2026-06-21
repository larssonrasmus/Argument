package main

import (
	"log"
	"net/http"
	"time"

	"argument/internal/app"
	"argument/internal/config"
)

var conf config.Config

func main() {
	srv := &http.Server{
		Handler: app.SetupHandler(),
		Addr:    "127.0.0.1:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
