package main

import (
	"log"
	"os"

	"github.com/heppu/contact-storage"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("$DATABASE_URL must be set")
	}

	port := os.Getenv("PORT")
	if dbURL == "" {
		log.Fatal("$PORT must be set")
	}

	log.Fatal(collector.Run(":"+port, dbURL))
}
