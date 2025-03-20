package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/MostovykhGit/urlshortener/datastorage"
)

func main() {
	storageType := flag.String("storage", "memory", "storage type: memory or postgres")
	pgConn := flag.String("pg_conn", "", "postgres connection string")
	port := flag.String("port", "8080", "port to listen to")
	flag.Parse()

	var storage datastorage.Storage
	var err error

	if *storageType == "postgres" {
		if *pgConn == "" {
			log.Fatal("must provide postgres connection")
		}
		storage, err = datastorage.NewPostgresStorage(*pgConn)
		if err != nil {
			log.Fatalf("failed to initialize Postgres storage: %v", err)
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
