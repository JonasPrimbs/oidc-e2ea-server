/*
 * OIDC² - Identity Certification Token Endpoint
 *
 * Endpoint for OpenID Connect's Identity Certification Token endpoint.
 *
 * API version: 0.5.0
 */
package main

import (
	"log"
	"net/http"
	"os"

	ict "ict/go"
)

func main() {
	log.Printf("Starting Server...")

	// Load configuration
	log.Printf("Loading configuration...")
	ict.Initialize()

	// Load router
	router := ict.NewRouter()

	// Get port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Configuration loaded")

	log.Printf("Running on port " + port)

	log.Fatal(http.ListenAndServe(":"+port, router))
}
