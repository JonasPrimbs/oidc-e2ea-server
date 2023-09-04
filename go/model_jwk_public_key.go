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

type JwkPublicKey struct {
	JwkEcPublicKey
	JwkRsaPublicKey
	JwkEdPublicKey
}

func PublicJwkFromJson(json map[string]interface{}, alg jwt.SigningMethod) (JwkPublicKey, error) {
	switch alg {
	// Elliptic Curve:
	case jwt.SigningMethodES256:
		fallthrough
	case jwt.SigningMethodES384:
		fallthrough
	case jwt.SigningMethodES512:
		ecJwk, err := EcJwkFromJson(json, alg)
		if err != nil {
			return JwkPublicKey{}, errors.New("failed to read EC public key: " + err.Error())
		}
		return JwkPublicKey{JwkEcPublicKey: ecJwk}, nil
	// RSA:
	case jwt.SigningMethodRS256:
		fallthrough
	case jwt.SigningMethodRS384:
		fallthrough
	case jwt.SigningMethodRS512:
		rsaJwk, err := RsaJwkFromJson(json)
		if err != nil {
			return JwkPublicKey{}, errors.New("failed to read RSA public key: " + err.Error())
		}
		return JwkPublicKey{JwkRsaPublicKey: rsaJwk}, nil
	case jwt.SigningMethodEdDSA:
		edDsaJwk, err := EdJwkFromJson(json)
		if err != nil {
			return JwkPublicKey{}, errors.New("failed to read Eduard curve public key: " + err.Error())
		}
		return JwkPublicKey{JwkEdPublicKey: edDsaJwk}, nil
	// Not supported:
	default:
		return JwkPublicKey{}, errors.New("signing algorithm '" + alg.Alg() + "' not supported")
	}
}

func KeyTypeFromJson(json map[string]interface{}, attributeName string) (KeyType, error) {
	// Read key type from json object
	ktyString, err := StringFromJson(json, attributeName)
	if err != nil {
		return "", err
	}

	// Convert key type to KeyTypes
	kty, ok := KeyTypeFromString(ktyString)
	if !ok {
		return "", errors.New("key type '" + ktyString + "' not supported")
	}

	return kty, nil
}
