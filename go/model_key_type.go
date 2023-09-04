/*
 * OIDC² - Identity Certification Token Endpoint
 *
 * Endpoint for OpenID Connect's Identity Certification Token endpoint.
 *
 * API version: 0.5.0
 */
package ict

type KeyType string

// List of KeyTypes
const (
	EC  KeyType = "EC"
	RSA KeyType = "RSA"
	OKP KeyType = "OKP"
)

func KeyTypeFromString(value string) (KeyType, bool) {
	switch value {
	case "EC":
		return EC, true
	case "RSA":
		return RSA, true
	case "OKP":
		return OKP, true
	default:
		return "", false
	}
}
