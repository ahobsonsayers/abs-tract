package main

import (
	"log"
	"net/http"

	"github.com/ahobsonsayers/abs-goodreads/server"
)

//go:generate go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen -config .oapigen.yaml schema/openapi.yaml

const serverAddress = "0.0.0.0:5555"

func main() {
	router, err := server.NewRouter()
	if err != nil {
		log.Fatalf("Failed to create router: %s", err)
	}

	// Start the Server
	log.Printf("Server listening on %s\n", serverAddress)
	err = http.ListenAndServe(serverAddress, router)
	if err != nil {
		log.Fatalf("Server exited with error: %s", err)
	}
}
