/*
 * OIDC² - Identity Certification Token Endpoint
 *
 * Endpoint for OpenID Connect's Identity Certification Token endpoint.
 *
 * API version: 0.5.0
 */
package ict

// Information about ocurred error.
type ErrorStatus struct {
	// Status Code
	Code int `json:"code"`
	// Status Text
	Status string `json:"status"`
	// More detailed description
	Description string `json:"description,omitempty"`
}
