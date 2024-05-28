package main

import (
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Create your router.
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	router := chi.NewMux()
	server := New(router)
	log.Fatal(server.start())
}
