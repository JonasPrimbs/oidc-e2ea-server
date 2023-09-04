/*
 * OIDC² - Identity Certification Token Endpoint
 *
 * Endpoint for OpenID Connect's Identity Certification Token endpoint.
 *
 * API version: 0.5.0
 */
package ict

type RsaSigningAlgorithm string

// List of RsaSigningAlgorithms
const (
	RS256 RsaSigningAlgorithm = "RS256"
	RS384 RsaSigningAlgorithm = "RS384"
	RS512 RsaSigningAlgorithm = "RS512"
)
