/*
 * OIDC² - Identity Certification Token Endpoint
 *
 * Endpoint for OpenID Connect's Identity Certification Token endpoint.
 *
 * API version: 0.5.0
 */
package ict

type IctResponse struct {
	IdentityCertificationToken string `json:"identity_certification_token"`
	// Number of seconds until the Identity Certification Token expires.
	ExpiresIn int32 `json:"expires_in"`
	// Space delimited claims provided in the Identity Certification Token.
	Claims []string `json:"claims"`
	// Array of authorized end-to-end authentication contexts.
	E2eAuthContexts []string `json:"e2e_auth_contexts"`
}
