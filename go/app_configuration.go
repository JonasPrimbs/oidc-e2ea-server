/*
 * Open Identity Certification with OpenID Connect (OIDC²)
 *
 * Authorization Server middleware for requesting Identity Certification Tokens (ICT).
 *
 * API version: 0.2.0
 * Contact: mail@jonasprimbs.de
 */
package oidc2middleware

import (
	"errors"
	"net"
	"os"
	"slices"
	"strconv"
	"strings"
)

type AppConfiguration struct {
	KeyFile               string
	KeyID                 string
	Issuer                string
	IntrospectionEndpoint string
	UserInfoEndpoint      string
	TokenEndpoint         string
	Algorithm             string
	TokenPeriod           int
	Port                  int
	Listeners             []string
}

var Configuration *AppConfiguration

func LoadConfigurationFromEnv() (*AppConfiguration, error) {
	// Ensure that required parameters are set
	keyFile := os.Getenv("KEY_FILE")
	if keyFile == "" {
		return nil, errors.New("KEY_FILE environment variable is required")
	}
	keyID := os.Getenv("KEY_ID")
	if keyID == "" {
		return nil, errors.New("KEY_ID environment variable is required")
	}

	// Validate algorithm
	algorithm := getEnvOrDefault("ALG", "RS256")
	if !isValidAlgorithm(algorithm) {
		return nil, errors.New("invalid algorithm: " + algorithm + ". must be one of: ES256, ES384, ES512, RS256, RS384, RS512, PS256, PS384, PS512, EdDSA")
	}

	// Validate token period
	tokenPeriod := getEnvAsInt("TOKEN_PERIOD", 300)
	if !isValidTokenPeriod(tokenPeriod) {
		return nil, errors.New("invalid token period: " + strconv.Itoa(tokenPeriod) + ". must be greater than 0")
	}

	// Validate port
	port := getEnvAsInt("PORT", 8080)
	if !isValidPort(port) {
		return nil, errors.New("invalid port number: " + strconv.Itoa(port) + ". must be between 1 and 65535")
	}

	// Read issuer and endpoints from environment variables
	issuer := os.Getenv("ISSUER")
	introspectionEndpoint := os.Getenv("INTROSPECTION_ENDPOINT")
	userInfoEndpoint := os.Getenv("USERINFO_ENDPOINT")
	tokenEndpoint := os.Getenv("TOKEN_ENDPOINT")

	// Discover endpoints if not configured
	if introspectionEndpoint == "" || userInfoEndpoint == "" || tokenEndpoint == "" {
		discoveryDocument, err := discoverEndpoints(issuer)
		if err != nil {
			return nil, err
		}

		// Override endpoints with discovered values if they were not set in the environment
		if introspectionEndpoint == "" {
			introspectionEndpoint = discoveryDocument.IntrospectionEndpoint
		}
		if userInfoEndpoint == "" {
			userInfoEndpoint = discoveryDocument.UserInfoEndpoint
		}
		if tokenEndpoint == "" {
			tokenEndpoint = discoveryDocument.TokenEndpoint
		}
	}

	return &AppConfiguration{
		KeyFile:               keyFile,
		KeyID:                 keyID,
		Issuer:                issuer,
		IntrospectionEndpoint: introspectionEndpoint,
		UserInfoEndpoint:      userInfoEndpoint,
		TokenEndpoint:         tokenEndpoint,
		Algorithm:             algorithm,
		TokenPeriod:           tokenPeriod,
		Port:                  port,
		Listeners:             parseListeners(getEnvOrDefault("HOSTS", "0.0.0.0")),
	}, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func isValidAlgorithm(alg string) bool {
	validAlgorithms := []string{"ES256", "ES384", "ES512", "RS256", "RS384", "RS512", "PS256", "PS384", "PS512", "EdDSA"}
	return slices.Contains(validAlgorithms, alg)
}

func isValidPort(port int) bool {
	return port > 0 && port <= 65535
}

func isValidTokenPeriod(period int) bool {
	return period > 0
}

func parseListeners(hostsStr string) []string {
	if hostsStr == "" {
		return []string{
			"0.0.0.0",
		}
	}

	listeners := strings.Split(hostsStr, ",")
	validListeners := []string{}

	for _, listener := range listeners {
		listener = strings.TrimSpace(listener)
		if isValidCIDR(listener) || net.ParseIP(listener) != nil {
			validListeners = append(validListeners, listener)
		}
	}

	return validListeners
}

func isValidCIDR(cidr string) bool {
	_, _, err := net.ParseCIDR(cidr)
	return err == nil
}
