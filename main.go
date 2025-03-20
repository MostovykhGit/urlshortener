package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"time"

	"./datastorage"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	storageType := flag.String("storage", "memory", "Storage type: memory or postgres")
	pgConn := flag.String("pg_conn", "", "Postgres connection string")
	port := flag.String("port", "8080", "Port to listen on")
	flag.Parse()

	var storage datastorage.Storage
	var err error

	if *storageType == "postgres" {
		if *pgConn == "" {
			log.Fatal("Postgres connection string must be provided")
		}
		storage, err = datastorage.NewPostgresStorage(*pgConn)
		if err != nil {
			log.Fatalf("Failed to initialize Postgres storage: %v", err)
		}
	} else {
		storage = datastorage.NewInMemoryStorage()
	}

	app := &App{storage: storage}
	http.HandleFunc("/shorten", app.shortenHandler)
	http.HandleFunc("/", app.redirectHandler)

	log.Printf("Listening on port %s...", *port)
	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
