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

type JwkEcPublicKey struct {
	KeyType   KeyType `json:"kty"`
	CurveName EcCurve `json:"crv"`
	X         string  `json:"x"`
	Y         string  `json:"y"`
}

func EcJwkFromJson(json map[string]interface{}, alg jwt.SigningMethod) (JwkEcPublicKey, error) {
	// Parse key type
	keyType, err := KeyTypeFromJson(json, "kty")
	if err != nil {
		return JwkEcPublicKey{}, errors.New("failed to parse key type: " + err.Error())
	}
	if keyType != EC {
		return JwkEcPublicKey{}, errors.New("failed to parse key type: expected attribute 'kty' to be 'EC' but found '" + string(keyType) + "'")
	}

	// Parse curve name
	curveName, err := EcCurveFromJson(json, "crv")
	if err != nil {
		return JwkEcPublicKey{}, errors.New("failed to parse curve name: " + err.Error())
	}
	switch alg {
	case jwt.SigningMethodES256:
		if curveName != "P-256" {
			return JwkEcPublicKey{}, errors.New("failed to parse curve name: expected curve name 'P-256' for signing algorithm '" + alg.Alg() + "'")
		}
	case jwt.SigningMethodES384:
		if curveName != "P-384" {
			return JwkEcPublicKey{}, errors.New("failed to parse curve name: expected curve name 'P-384' for signing algorithm '" + alg.Alg() + "'")
		}
	case jwt.SigningMethodES512:
		if curveName != "P-521" {
			return JwkEcPublicKey{}, errors.New("failed to parse curve name: expected curve name 'P-521' for signing algorithm '" + alg.Alg() + "'")
		}
	}

	// Parse x coordinate
	x, err := StringFromJson(json, "x")
	if err != nil {
		return JwkEcPublicKey{}, errors.New("failed to parse x coordinate: " + err.Error())
	}

	// Parse y coordinate
	y, err := StringFromJson(json, "y")
	if err != nil {
		return JwkEcPublicKey{}, errors.New("failed to parse y coordinate: " + err.Error())
	}

	// Return as struct
	return JwkEcPublicKey{
		KeyType:   keyType,
		CurveName: curveName,
		X:         x,
		Y:         y,
	}, nil
}
