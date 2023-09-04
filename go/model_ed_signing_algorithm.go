/*
 * OIDC² - Identity Certification Token Endpoint
 *
 * Endpoint for OpenID Connect's Identity Certification Token endpoint.
 *
 * API version: 0.5.0
 */
package ict

type EdSigningAlgorithm string

// List of EdSigningAlgorithm
const (
	ED_DSA EdSigningAlgorithm = "EdDSA"
)
