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
)

type EdCurve string

// List of EdCurves
const (
	ED25519 EdCurve = "Ed25519"
)

func EdCurveFromName(name string) (EdCurve, bool) {
	switch name {
	case "Ed25519":
		return ED25519, true
	default:
		return "", false
	}
}

func EdCurveFromJson(json map[string]interface{}, attributeName string) (EdCurve, error) {
	// Read curve name from json object
	crvString, err := StringFromJson(json, attributeName)
	if err != nil {
		return "", err
	}

	// Convert attribute value to EcCurve
	crv, ok := EdCurveFromName(crvString)
	if !ok {
		return crv, errors.New("elliptic curve name '" + string(crv) + "' not supported")
	}

	// Return result
	return crv, nil
}
