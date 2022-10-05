/*
 * OIDC RIDT Userinfo Endpoint
 *
 * Endpoint for OpenID Connect's Remote ID Token endpoint for userinfo.
 *
 * API version: 0.2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package ridt

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func Base64ToBigInt(s string) (*big.Int, error) {
	// Parse base64url encoded string to bytes
	data, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return nil, errors.New("failed to parse base64url encoded big integer '" + s + "': " + err.Error())
	}

	// Convert bytes to big integer
	i := new(big.Int)
	i.SetBytes(data)
	return i, nil
}

func Base64ToInt(s string) (int, error) {
	// Parse base64url encoded string to bytes.
	data, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return 0, errors.New("failed to parse base64url encoded big integer '" + s + "': " + err.Error())
	}

	// Convert bytes to integer
	i := binary.LittleEndian.Uint32(append(data, 0))
	return int(i), nil
}

func ReadEcPrivateKey(fileName string) (*ecdsa.PrivateKey, error) {
	privateData, err := os.ReadFile(fileName)
	if err != nil {
		return nil, errors.New("Failed to read private key file: " + err.Error())
	}
	privateKey, err := jwt.ParseECPrivateKeyFromPEM(privateData)
	if err != nil {
		return nil, errors.New("Failed to parse private key: " + err.Error())
	}
	return privateKey, nil
}

func ReadEcPublicKey(fileName string) (*ecdsa.PublicKey, error) {
	publicData, err := os.ReadFile(fileName)
	if err != nil {
		return nil, errors.New("Failed to read public key file: " + err.Error())
	}
	publicKey, err := jwt.ParseECPublicKeyFromPEM(publicData)
	if err != nil {
		return nil, errors.New("Failed to parse public key: " + err.Error())
	}
	return publicKey, nil
}

func ReadRsaPrivateKey(fileName string) (*rsa.PrivateKey, error) {
	privateData, err := os.ReadFile(fileName)
	if err != nil {
		return nil, errors.New("Failed to read private key file: " + err.Error())
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateData)
	if err != nil {
		return nil, errors.New("Failed to parse private key: " + err.Error())
	}
	return privateKey, nil
}

func ReadRsaPublicKey(fileName string) (*rsa.PublicKey, error) {
	publicData, err := os.ReadFile(fileName)
	if err != nil {
		return nil, errors.New("Failed to read public key file: " + err.Error())
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicData)
	if err != nil {
		return nil, errors.New("Failed to parse public key: " + err.Error())
	}
	return publicKey, nil
}

func Header(r *http.Request, headerName string) (string, error) {
	// Read header and ensure that it is exists
	headerValue := r.Header.Get(headerName)
	if headerValue == "" {
		return "", errors.New("header '" + headerName + "' not found")
	}

	return headerValue, nil
}

func AuthorizationHeader(r *http.Request) (string, error) {
	// Read authorization header
	authorizationHeader, err := Header(r, "Authorization")
	if err != nil {
		return "", errors.New("authorization header not found")
	}

	return authorizationHeader, nil
}

func delimiterFn(c rune) func(rune) bool {
	fn := func(char rune) bool {
		return c == char
	}
	return fn
}

func BearerTokenFromAuthorizationHeader(r *http.Request) (string, error) {
	// Get authorization header
	authorizationHeader, err := AuthorizationHeader(r)
	if err != nil {
		return "", err
	}
	authorizationHeaderParts := strings.FieldsFunc(authorizationHeader, delimiterFn(' '))
	if len(authorizationHeaderParts) < 2 {
		return "", errors.New("authorization header value missing")
	}

	// Ensure that authorization is a bearer authorization
	if authorizationType := authorizationHeaderParts[0]; strings.ToLower(authorizationType) != "bearer" {
		return "", errors.New("authorization type is '" + authorizationType + "' but expected 'bearer'")
	}

	// Return bearer token
	return authorizationHeaderParts[1], nil
}

func LogAndSendError(w http.ResponseWriter, statusCode int, status string, description string, details string) {
	// Log error
	log.Print("[ERROR] " + details)

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(ErrorStatus{
		Code:        statusCode,
		Status:      status,
		Description: description,
	})
}

func RequestUserinfo(bearerToken string, uri string, issuer string) (map[string]interface{}, error, bool) {
	// Create new http client
	client := &http.Client{}

	// Create new http request
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, errors.New("failed to create userinfo request to '" + uri + "': " + err.Error()), false
	}
	req.Header.Set("Authorization", "Bearer "+bearerToken)
	// Set Host header
	issuerParts := strings.Split(issuer, "/")
	hostname := issuerParts[2]
	req.Host = hostname

	// Send http request and validate response
	res, err := client.Do(req)
	if err != nil {
		return nil, errors.New("failed to send userinfo request to '" + uri + "': " + err.Error()), false
	}
	if res.StatusCode != 200 {
		return nil, errors.New("failed to get userinfo response from '" + uri + "'. status code: " + fmt.Sprint(res.StatusCode) + ", status: '" + res.Status + "'"), res.StatusCode == 401
	}

	// Parse response
	var claims map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&claims)
	if err != nil {
		return nil, errors.New("failed to parse userinfo response: " + err.Error()), false
	}

	// Return parsed claims
	return claims, nil, false
}

func ReadRequestBody(r *http.Request) (string, error) {
	// Read request body
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", errors.New("failed to read request body: " + err.Error())
	}

	// Convert to string and return
	return string(bodyBytes), nil
}

func EllipticCurveFromString(crv string) (elliptic.Curve, error) {
	switch crv {
	case "P-256":
		return elliptic.P256(), nil
	case "P-384":
		return elliptic.P384(), nil
	case "P-521":
		return elliptic.P521(), nil
	default:
		return nil, errors.New("curve '" + crv + "' is not supported")
	}
}

func EllipticCurveFromJson(json map[string]interface{}, attributeName string) (elliptic.Curve, error) {
	// Parse curve from json
	crv, err := StringFromJson(json, attributeName)
	if err != nil {
		return nil, err
	}

	// Parse curve name
	curve, err := EllipticCurveFromString(crv)
	if err != nil {
		return nil, err
	}

	// Return curve
	return curve, nil
}

func EcdsaPublicKeyFromJson(jwk map[string]interface{}) (*ecdsa.PublicKey, map[string]interface{}, error) {
	publicKey := new(ecdsa.PublicKey)
	publicKeyJwk := make(map[string]interface{})
	publicKeyJwk["kty"] = "EC"

	// Parse curve name
	crv, err := EllipticCurveFromJson(jwk, "crv")
	if err != nil {
		return nil, nil, errors.New("failed to read curve name: " + err.Error())
	}
	publicKey.Curve = crv
	crvString, err := StringFromJson(jwk, "crv")
	if err != nil {
		return nil, nil, errors.New("curve name not found: " + err.Error())
	}
	publicKeyJwk["crv"] = crvString

	// Parse x value
	x, err := BigIntFromJsonBase64(jwk, "x")
	if err != nil {
		return nil, nil, errors.New("failed to read x value: " + err.Error())
	}
	publicKey.X = x
	xString, err := StringFromJson(jwk, "x")
	if err != nil {
		return nil, nil, errors.New("x value not found: " + err.Error())
	}
	publicKeyJwk["x"] = xString

	// Parse y value
	y, err := BigIntFromJsonBase64(jwk, "y")
	if err != nil {
		return nil, nil, errors.New("failed to read y value: " + err.Error())
	}
	publicKey.Y = y
	yString, err := StringFromJson(jwk, "y")
	if err != nil {
		return nil, nil, errors.New("y value not found: " + err.Error())
	}
	publicKeyJwk["y"] = yString

	// Return public key
	return publicKey, publicKeyJwk, nil
}
func RsaPublicKeyFromJson(jwk map[string]interface{}) (*rsa.PublicKey, map[string]interface{}, error) {
	publicKey := new(rsa.PublicKey)
	publicKeyJwk := make(map[string]interface{})
	publicKeyJwk["kty"] = "RSA"

	// Parse x value
	e, err := IntFromJsonBase64(jwk, "e")
	if err != nil {
		return nil, nil, errors.New("failed to read exponent: " + err.Error())
	}
	publicKey.E = e
	eString, err := StringFromJson(jwk, "e")
	if err != nil {
		return nil, nil, errors.New("exponent not found: " + err.Error())
	}
	publicKeyJwk["e"] = eString

	// Parse modulus value
	n, err := BigIntFromJsonBase64(jwk, "n")
	if err != nil {
		return nil, nil, errors.New("failed to read modulus: " + err.Error())
	}
	publicKey.N = n
	nString, err := StringFromJson(jwk, "n")
	if err != nil {
		return nil, nil, errors.New("modulus not found: " + err.Error())
	}
	publicKeyJwk["n"] = nString

	// Return public key
	return publicKey, publicKeyJwk, nil
}

func PublicKeyFromJwt(token *jwt.Token) (interface{}, map[string]interface{}, error) {
	// Get public key from proof of possession header
	jwk, err := JsonFromJson(token.Header, "jwk")
	if err != nil {
		return nil, nil, errors.New("failed to get public key from header: " + err.Error())
	}

	// Extract public key
	switch token.Method {
	// Elliptic Curve:
	case jwt.SigningMethodES256:
		fallthrough
	case jwt.SigningMethodES384:
		fallthrough
	case jwt.SigningMethodES512:
		// Parse elliptic curve key from "jwk" claim in JWT header
		publicKey, publicKeyJwk, err := EcdsaPublicKeyFromJson(jwk)
		if err != nil {
			return nil, nil, errors.New("failed to parse ECDSA public key: " + err.Error())
		}
		return publicKey, publicKeyJwk, nil
	// RSA:
	case jwt.SigningMethodRS256:
		fallthrough
	case jwt.SigningMethodRS384:
		fallthrough
	case jwt.SigningMethodRS512:
		publicKey, publicKeyJwk, err := RsaPublicKeyFromJson(jwk)
		if err != nil {
			return nil, nil, errors.New("failed to parse RSA public key: " + err.Error())
		}
		return publicKey, publicKeyJwk, nil
	// Not supported:
	default:
		return nil, nil, errors.New("signing algorithm '" + token.Method.Alg() + "' not supported")
	}
}

func ParseProofOfPossessionFromRequestBody(r *http.Request) (*jwt.Token, jwt.MapClaims, map[string]interface{}, error) {
	// Read request body
	requestBody, err := ReadRequestBody(r)
	if err != nil {
		return nil, jwt.MapClaims{}, nil, err
	}

	// Parse token
	claims := jwt.MapClaims{}
	var publicKeyJwk map[string]interface{}
	token, err := jwt.ParseWithClaims(requestBody, &claims, func(token *jwt.Token) (interface{}, error) {
		key, keyJwk, err := PublicKeyFromJwt(token)
		if err != nil {
			return nil, err
		}
		publicKeyJwk = keyJwk
		return key, nil
	})
	if err != nil {
		return nil, jwt.MapClaims{}, nil, errors.New("failed to parse proof of possession token: " + err.Error())
	}

	return token, claims, publicKeyJwk, nil
}

func ValidateProofOfPossession(popToken *jwt.Token, popClaims jwt.MapClaims, userinfoClaims map[string]interface{}, config AppConfiguration, now int64) error {
	// Validate proof of possession token
	if !popToken.Valid {
		return errors.New("proof of possession token is not valid")
	}

	// Validate and compare subject
	userinfoSub, err := StringFromJson(userinfoClaims, "sub")
	if err != nil {
		return errors.New("subject claim in userinfo response not found")
	}
	popSub, err := StringFromJson(popClaims, "sub")
	if err != nil {
		return errors.New("subject claim in proof of possession token not found")
	}
	if userinfoSub != popSub {
		return errors.New("invalid subject claim in proof of possession token")
	}

	// Validate audience
	ok := popClaims.VerifyAudience(config.Issuer, true)
	if !ok {
		return errors.New("invalid audience claim in proof of possession token")
	}

	// Verify expiration
	if !popClaims.VerifyExpiresAt(now, true) ||
		!popClaims.VerifyNotBefore(now, false) ||
		!popClaims.VerifyIssuedAt(now, true) {
		return errors.New("token expired or is not yet valid")
	}

	return nil
}

func GenerateRidt(privateKey interface{}, algorithm jwt.SigningMethod, tokenClaims jwt.MapClaims, publicKeyJwk map[string]interface{}, userinfoClaims map[string]interface{}, config AppConfiguration) (string, []string, int64, error) {
	// Compute token validity
	expiresIn := config.DefaultTokenPeriod
	if tokenLifetime, ok := tokenClaims["token_lifetime"]; ok {
		popTokenLifetime := uint32(tokenLifetime.(float64))
		if popTokenLifetime > config.MaxTokenPeriod {
			expiresIn = config.MaxTokenPeriod
		} else if popTokenLifetime > 0 {
			expiresIn = popTokenLifetime
		}
	}

	var nonce string
	if tokenNonce, err := StringFromJson(tokenClaims, "token_nonce"); err == nil {
		nonce = tokenNonce
	} else {
		randBytes := make([]byte, 24)
		rand.Read(randBytes)
		nonce = base64.URLEncoding.EncodeToString(randBytes)
	}

	// Compose claims for remote id token
	requestedClaims := jwt.MapClaims{}
	if tokenClaims, err := StringFromJson(tokenClaims, "token_claims"); err == nil {
		// Include only selected claims
		reqClaimKeys := strings.FieldsFunc(tokenClaims, delimiterFn(' '))
		for i := 0; i < len(reqClaimKeys); i++ {
			claimName := reqClaimKeys[i]
			claimValue, ok := userinfoClaims[claimName]
			if ok {
				requestedClaims[claimName] = claimValue
			}
		}
	} else {
		// Include all claims
		userinfoClaimKeys := reflect.ValueOf(userinfoClaims).MapKeys()
		for i := 0; i < len(userinfoClaimKeys); i++ {
			claimName := userinfoClaimKeys[i].String()
			claimValue, ok := userinfoClaims[claimName]
			if ok {
				requestedClaims[claimName] = claimValue
			}
		}
	}
	identityClaimNames := reflect.ValueOf(requestedClaims).MapKeys()
	claimNames := make([]string, len(identityClaimNames))
	for i := 0; i < len(claimNames); i++ {
		claimNames[i] = identityClaimNames[i].String()
	}

	_, err := StringFromJson(userinfoClaims, "sub")
	if err != nil {
		return "", nil, 0, errors.New("subject not found")
	}
	// requestedClaims["sub"] = subject
	requestedClaims["iss"] = config.Issuer
	requestedClaims["nonce"] = nonce

	// Set time constraints
	now := time.Now().Unix()
	expiresAt := now + int64(expiresIn)
	requestedClaims["iat"] = now
	requestedClaims["nbf"] = now
	requestedClaims["exp"] = expiresAt

	// Add confirmation header
	confirmation := make(map[string]interface{})
	confirmation["jwk"] = publicKeyJwk
	requestedClaims["cnf"] = confirmation

	// Generate JWT
	ridt := jwt.NewWithClaims(algorithm, requestedClaims)
	ridt.Header["kid"] = config.KeyId
	ridtString, err := ridt.SignedString(privateKey)
	if err != nil {
		return "", nil, 0, errors.New("failed to sign remote id token: " + err.Error())
	}

	return ridtString, claimNames, expiresAt, nil
}

func GenRidt(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	log.Print("origin: " + origin)
	if origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

	// Load configuration
	appConfig, err := LoadAppConfigurationFromEnv()
	if err != nil {
		LogAndSendError(w, http.StatusInternalServerError, "internal server error", "unknown internal server error", "failed to load configuration: "+err.Error())
		return
	}
	privateKey, err := ReadRsaPrivateKey(appConfig.KeyFilePath)
	if err != nil {
		LogAndSendError(w, http.StatusInternalServerError, "internal server error", "unknown internal server error", "failed to load private key file: "+err.Error())
		return
	}

	// Get bearer token from authorization header
	bearerToken, err := BearerTokenFromAuthorizationHeader(r)
	if err != nil {
		LogAndSendError(w, http.StatusUnauthorized, "unauthorized", "bearer authentication required", "failed to read bearer token: "+err.Error())
		return
	}

	// Get identity claims from userinfo endpoint
	userinfoClaims, err, authFailed := RequestUserinfo(bearerToken, appConfig.UserinfoEndpoint, appConfig.Issuer)
	if err != nil {
		if authFailed {
			LogAndSendError(w, http.StatusUnauthorized, "unauthorized", "invalid bearer token", err.Error())
		} else {
			LogAndSendError(w, http.StatusInternalServerError, "internal server error", "unknown internal server error", err.Error())
		}
		return
	}

	// Read proof of possession from request body
	popToken, popClaims, publicKeyJwk, err := ParseProofOfPossessionFromRequestBody(r)
	if err != nil {
		LogAndSendError(w, http.StatusInternalServerError, "internal server error", "unknown internal server error", "failed to validate proof of possession: "+err.Error())
		return
	}

	// Validate proof of possession
	err = ValidateProofOfPossession(popToken, popClaims, userinfoClaims, appConfig, time.Now().Unix())
	if err != nil {
		LogAndSendError(w, http.StatusForbidden, "forbidden", "invalid proof of possession", "failed to validate proof of possession: "+err.Error())
		return
	}

	// Generate remote id token
	ridt, identityClaims, expiresAt, err := GenerateRidt(privateKey, appConfig.SigningAlgorithm, popClaims, publicKeyJwk, userinfoClaims, appConfig)
	if err != nil {
		LogAndSendError(w, http.StatusInternalServerError, "internal server error", "unknown internal server error", "failed to generate remote id token: "+err.Error())
		return
	}

	// Encode response
	expiresIn := expiresAt - time.Now().Unix()
	response := IatResponse{
		RemoteIdToken: ridt,
		ExpiresIn:     int32(expiresIn),
		Claims:        strings.Join(identityClaims, " "),
	}

	// Write response
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	json.NewEncoder(w).Encode(response)
}

func RidtOptions(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	log.Print("origin: " + origin)
	if origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
	w.WriteHeader(http.StatusNoContent)
}