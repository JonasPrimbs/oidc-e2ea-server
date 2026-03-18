/*
 * Open Identity Certification with OpenID Connect (OIDC²)
 *
 * Authorization Server middleware for requesting Identity Certification Tokens (ICT).
 *
 * API version: 0.2.0
 * Contact: mail@jonasprimbs.de
 */
package oidc2middleware

import (
	"encoding/json"
	"net/http"
)

type OidcDiscoveryDocument struct {
	IntrospectionEndpoint string `json:"introspection_endpoint"`
	UserInfoEndpoint      string `json:"userinfo_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
}

// discoverEndpoints retrieves the OpenID Provider's discovery document from the well-known configuration endpoint.
// It makes an HTTP GET request to the issuer's /.well-known/openid-configuration endpoint,
// parses the JSON response, and returns a OidcDiscoveryDocument containing the introspection and user info endpoints.
//
// Parameters:
//   - issuer: The base URL of the OpenID Provider (e.g., "https://example.com")
//
// Returns:
//   - *OidcDiscoveryDocument: A pointer to a OidcDiscoveryDocument containing the IntrospectionEndpoint and UserInfoEndpoint
//   - error: An error if the HTTP request fails or if JSON decoding fails
//
// Example:
//
//	doc, err := discoverEndpoints("https://accounts.google.com")
//	if err != nil {
//	    log.Fatal(err)
//	}
func discoverEndpoints(issuer string) (*OidcDiscoveryDocument, error) {
	// Prepare request to the OpenID Provider's discovery endpoint
	resp, err := http.Get(issuer + "/.well-known/openid-configuration")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse discovery document
	var config OidcDiscoveryDocument
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
