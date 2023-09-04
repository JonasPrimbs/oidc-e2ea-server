/*
 * OIDC² - Identity Certification Token Endpoint
 *
 * Endpoint for OpenID Connect's Identity Certification Token endpoint.
 *
 * API version: 0.5.0
 */
package ict

type EcSigningAlgorithm string

// List of EcSigningAlgorithms
const (
	ES256 EcSigningAlgorithm = "ES256"
	ES384 EcSigningAlgorithm = "ES384"
	ES512 EcSigningAlgorithm = "ES512"
)
