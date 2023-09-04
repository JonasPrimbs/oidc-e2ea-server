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
	"os"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
)

type AppConfiguration struct {
	KeyFilePath        string            `json:"keyFilePath"`
	KeyId              string            `json:"keyId"`
	SigningAlgorithm   jwt.SigningMethod `json:"alg"`
	UserinfoEndpoint   string            `json:"UserinfoEndpoint"`
	Issuer             string            `json:"Issuer"`
	DefaultTokenPeriod uint64            `json:"DefaultTokenPeriod"`
	MaxTokenPeriod     uint32            `json:"MaxTokenPeriod"`
}

func LoadAppConfigurationFromEnv() (AppConfiguration, error) {
	// Parse key file path
	keyFilePath := os.Getenv("KEY_FILE")
	if keyFilePath == "" {
		return AppConfiguration{}, errors.New("failed to load key file path: environment variable 'KEY_FILE' not found")
	}

	// Parse key ID
	keyId := os.Getenv("KID")
	if keyId == "" {
		return AppConfiguration{}, errors.New("failed to load key id: environment variable 'KID' not found")
	}

	// Parse signing algorithm
	signingAlgorithmString := os.Getenv("ALG")
	if signingAlgorithmString == "" {
		signingAlgorithmString = "ES256"
	}
	var signingAlgorithm jwt.SigningMethod
	switch signingAlgorithmString {
	case "ES256":
		signingAlgorithm = jwt.SigningMethodES256
	case "ES384":
		signingAlgorithm = jwt.SigningMethodES384
	case "ES512":
		signingAlgorithm = jwt.SigningMethodES512
	case "RS256":
		signingAlgorithm = jwt.SigningMethodRS256
	case "RS384":
		signingAlgorithm = jwt.SigningMethodRS384
	case "RS512":
		signingAlgorithm = jwt.SigningMethodRS512
	case "EdDSA":
		signingAlgorithm = jwt.SigningMethodEdDSA
	default:
		return AppConfiguration{}, errors.New("failed to load signing algorithm: signing algorithm '" + signingAlgorithmString + "' is not supported")
	}

	// Parse userinfo endpoint
	userinfoEndpoint := os.Getenv("USERINFO")
	if userinfoEndpoint == "" {
		return AppConfiguration{}, errors.New("failed to read userinfo endpoint: environment variable 'USERINFO' not found")
	}

	// Parse issuer
	issuer := os.Getenv("ISSUER")
	if userinfoEndpoint == "" {
		return AppConfiguration{}, errors.New("failed to load issuer: environment variable 'ISSUER' not found")
	}

	// Parse default token period
	defaultTokenPeriodString := os.Getenv("DEFAULT_TOKEN_PERIOD")
	if defaultTokenPeriodString == "" {
		defaultTokenPeriodString = "3600"
	}
	defaultTokenPeriodInt, err := strconv.Atoi(defaultTokenPeriodString)
	if err != nil {
		return AppConfiguration{}, errors.New("Failed to load default token period: value '" + defaultTokenPeriodString + "' is not an integer")
	}
	defaultTokenPeriod := uint64(defaultTokenPeriodInt)

	// Parse maximum token period
	maxTokenPeriodString := os.Getenv("MAX_TOKEN_PERIOD")
	if maxTokenPeriodString == "" {
		maxTokenPeriodString = "2592000"
	}
	maxTokenPeriodInt, err := strconv.Atoi(maxTokenPeriodString)
	if err != nil {
		return AppConfiguration{}, errors.New("failed load maximum token period: value '" + maxTokenPeriodString + "' is not an integer")
	}
	maxTokenPeriod := uint32(maxTokenPeriodInt)

	// Return result
	return AppConfiguration{
		KeyFilePath:        keyFilePath,
		KeyId:              keyId,
		SigningAlgorithm:   signingAlgorithm,
		UserinfoEndpoint:   userinfoEndpoint,
		Issuer:             issuer,
		DefaultTokenPeriod: defaultTokenPeriod,
		MaxTokenPeriod:     maxTokenPeriod,
	}, nil
}
