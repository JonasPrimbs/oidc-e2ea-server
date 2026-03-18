/*
 * Open Identity Certification with OpenID Connect (OIDC²)
 *
 * Authorization Server middleware for requesting Identity Certification Tokens (ICT).
 *
 * API version: 0.2.0
 * Contact: mail@jonasprimbs.de
 */

package main

import (
	"log"
	"net/http"
	"strconv"

	oidc2middleware "github.com/JonasPrimbs/oidc-e2ea-server/go"
)

func main() {
	log.Printf("Starting server...")

	// Load configuration from environment variables
	log.Printf("Loading configuration...")
	appConfig, err := oidc2middleware.LoadConfigurationFromEnv()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Printf("Configuration loaded: %+v", appConfig)

	log.Printf("Initializing API services...")
	DefaultAPIService := oidc2middleware.NewDefaultAPIService()
	DefaultAPIController := oidc2middleware.NewDefaultAPIController(DefaultAPIService)
	router := oidc2middleware.NewRouter(DefaultAPIController)

	for _, listener := range appConfig.Listeners {
		addr := listener + ":" + strconv.Itoa(appConfig.Port)
		log.Printf("Server started on http://%s", addr)
		go func(addr string) {
			if err := http.ListenAndServe(addr, router); err != nil {
				log.Fatalf("Server failed on http://%s: %v", addr, err)
			}
		}(addr)
	}
	select {}
}
