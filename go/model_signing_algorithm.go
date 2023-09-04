/*
 * OIDC² - Identity Certification Token Endpoint
 *
 * Endpoint for OpenID Connect's Identity Certification Token endpoint.
 *
 * API version: 0.5.0
 */
package ict

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

type SigningAlgorithm struct {
	EcSigningAlgorithm
	RsaSigningAlgorithm
}

func SigningAlgorithmFromJwa(jwa string) (jwt.SigningMethod, bool) {
	switch jwa {
	case "ES256":
		return jwt.SigningMethodES256, true
	case "ES384":
		return jwt.SigningMethodES384, true
	case "ES512":
		return jwt.SigningMethodES512, true
	case "RS256":
		return jwt.SigningMethodRS256, true
	case "RS384":
		return jwt.SigningMethodRS384, true
	case "RS512":
		return jwt.SigningMethodRS512, true
	case "EdDSA":
		return jwt.SigningMethodEdDSA, true
	default:
		return jwt.SigningMethodNone, false
	}
}

func SigningAlgorithmFromJson(json map[string]interface{}, attributeName string) (jwt.SigningMethod, error) {
	// Read algorithm from json object
	algString, err := StringFromJson(json, attributeName)
	if err != nil {
		return jwt.SigningMethodNone, err
	}

	// Convert attribute value to jwt.SigningMethod
	alg, ok := SigningAlgorithmFromJwa(algString)
	if !ok {
		return jwt.SigningMethodNone, errors.New("signing algorithm '" + algString + "' not supported")
	}

	// Return result
	return alg, nil
}
