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
	"context"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// DefaultAPIService is a service that implements the logic for the DefaultAPIServicer
// This service should implement the business logic for every endpoint for the DefaultAPI API.
// Include any external packages or services that will be required by this service.
type DefaultAPIService struct {
}

// NewDefaultAPIService creates a default api service
func NewDefaultAPIService() *DefaultAPIService {
	return &DefaultAPIService{}
}

// TokenRequest - Request a token.
func (s *DefaultAPIService) TokenRequest(ctx context.Context, grantType string, clientId string, code string, redirectUri string, clientSecret string, refreshToken string, scope string, resource string, audience string, subjectToken string, subjectTokenType string, requestedTokenType string, actorToken string, actorTokenType string, dpop string, authorizationHeader string, requestURL string) (ImplResponse, error) {
	cfg := Configuration

	if !isICTRequest(grantType, requestedTokenType) {
		log.Printf("Received non-ICT token request with grant_type=%q requested_token_type=%q, forwarding to token endpoint", grantType, requestedTokenType)
		return forwardTokenRequest(ctx, cfg, map[string]string{
			"grant_type":         grantType,
			"client_id":          clientId,
			"code":               code,
			"redirect_uri":       redirectUri,
			"client_secret":      clientSecret,
			"refresh_token":      refreshToken,
			"scope":              scope,
			"resource":           resource,
			"audience":           audience,
			"subject_token":      subjectToken,
			"subject_token_type": subjectTokenType,
			"requested_token_type": requestedTokenType,
			"actor_token":        actorToken,
			"actor_token_type":   actorTokenType,
		}, dpop, authorizationHeader)
	}

	if subjectToken == "" {
		return oauthError(http.StatusBadRequest, "invalid_request", "subject_token is required for ICT requests"), nil
	}
	if subjectTokenType != "" && subjectTokenType != "urn:ietf:params:oauth:token-type:access_token" {
		return oauthError(http.StatusBadRequest, "invalid_request", "subject_token_type must identify an access token"), nil
	}

	introspection, errResp := introspectToken(ctx, cfg, subjectToken)
	if errResp != nil {
		return *errResp, nil
	}
	if !introspection.Active {
		return oauthError(http.StatusBadRequest, "invalid_grant", "subject_token is inactive"), nil
	}
	if scope != "" && !hasRequiredScopes(introspection.Scope, scope) {
		return oauthError(http.StatusForbidden, "insufficient_scope", "subject_token does not contain the required scope"), nil
	}

	subjectJkt := resolveSubjectTokenJkt(introspection, subjectToken)
	if subjectJkt == "" {
		return oauthError(
			http.StatusBadRequest,
			"invalid_token",
			"subject_token must be a DPoP-bound access token (cnf.jkt required for ICT exchange)",
		), nil
	}
	if strings.TrimSpace(dpop) == "" {
		return oauthError(
			http.StatusBadRequest,
			"invalid_dpop_proof",
			"DPoP proof header is required for ICT token exchange",
		), nil
	}
	if _, err := validateDPoPProof(dpop, subjectJkt, requestURL); err != nil {
		return oauthError(http.StatusBadRequest, "invalid_dpop_proof", err.Error()), nil
	}
	introspection.Cnf.Jkt = subjectJkt
	userInfo := userInfoFromIntrospection(introspection)

	ctxClaims, errResp := resolveCtxClaims(cfg, introspection)
	if errResp != nil {
		return *errResp, nil
	}

	issuedAt := time.Now().UTC()
	expiresAt := issuedAt.Add(time.Duration(cfg.TokenPeriod) * time.Second)
	ict, err := issueICT(cfg, clientId, introspection, userInfo, ctxClaims, issuedAt, expiresAt)
	if err != nil {
		return Response(http.StatusInternalServerError, ErrorStatus{
			Error:            "server_error",
			ErrorDescription: err.Error(),
		}), nil
	}

	respScope := strings.TrimSpace(introspection.Scope)
	if scope != "" {
		respScope = scope
	}

	return Response(http.StatusOK, TokenRequest200Response{
		AccessToken:        ict,
		TokenType:          "DPoP",
		ExpiresIn:          int64(cfg.TokenPeriod),
		Scope:              respScope,
		DpopJkt:            introspection.Cnf.Jkt,
		IssuedTokenType:    "urn:ietf:params:oauth:token-type:ic_token",
		RequestedTokenType: "urn:ietf:params:oauth:token-type:ic_token",
		ClientId:           clientId,
	}), nil
}

// TokenRequestPreflight - CORS preflight request
func (s *DefaultAPIService) TokenRequestPreflight(ctx context.Context) (ImplResponse, error) {
	// TODO - update TokenRequestPreflight with the required logic for this service method.
	// Add api_default_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(204, {}) or use other options such as http.Ok ...
	// return Response(204, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("TokenRequestPreflight method not implemented")
}

type introspectionResponse struct {
	Active    bool                   `json:"active"`
	Scope     string                 `json:"scope"`
	ClientID  string                 `json:"client_id"`
	Azp       string                 `json:"azp"`
	Username  string                 `json:"username"`
	TokenType string                 `json:"token_type"`
	Sub       string                 `json:"sub"`
	Aud       any                    `json:"aud"`
	Iss       string                 `json:"iss"`
	Exp       int64                  `json:"exp"`
	Iat       int64                  `json:"iat"`
	Nbf       int64                  `json:"nbf"`
	Cnf       introspectionCNF       `json:"cnf"`
	Extra     map[string]interface{} `json:"-"`
}

type introspectionCNF struct {
	Jkt string `json:"jkt"`
}

func (r *introspectionResponse) UnmarshalJSON(data []byte) error {
	type alias introspectionResponse
	aux := alias{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	extra := map[string]interface{}{}
	if err := json.Unmarshal(data, &extra); err != nil {
		return err
	}
	delete(extra, "active")
	delete(extra, "scope")
	delete(extra, "client_id")
	delete(extra, "azp")
	delete(extra, "username")
	delete(extra, "token_type")
	delete(extra, "sub")
	delete(extra, "aud")
	delete(extra, "iss")
	delete(extra, "exp")
	delete(extra, "iat")
	delete(extra, "nbf")
	delete(extra, "cnf")
	*r = introspectionResponse(aux)
	r.Extra = extra
	return nil
}

func (r introspectionResponse) IsDPoPToken() bool {
	return strings.EqualFold(r.TokenType, "DPoP") || r.Cnf.Jkt != ""
}

// Returns the confirmation thumbprint for subjectToken
func resolveSubjectTokenJkt(introspection *introspectionResponse, subjectToken string) string {
	if introspection.Cnf.Jkt != "" {
		return introspection.Cnf.Jkt
	}
	return jktFromAccessTokenJWT(subjectToken)
}

func jktFromAccessTokenJWT(accessToken string) string {
	parts := strings.Split(accessToken, ".")
	if len(parts) != 3 {
		return ""
	}
	payload, err := decodeBase64URLField(parts[1])
	if err != nil {
		return ""
	}
	var claims map[string]interface{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return ""
	}
	cnf, ok := claims["cnf"].(map[string]interface{})
	if !ok {
		return ""
	}
	jkt, _ := cnf["jkt"].(string)
	return strings.TrimSpace(jkt)
}

func jktLogPrefix(jkt string) string {
	if len(jkt) <= 8 {
		return jkt
	}
	return jkt[:8]
}

type dpopProofClaims struct {
	HTM string `json:"htm"`
	HTU string `json:"htu"`
	JTI string `json:"jti"`
	jwt.RegisteredClaims
}

func isICTRequest(grantType, requestedTokenType string) bool {
	if grantType != "urn:ietf:params:oauth:grant-type:token-exchange" {
		return false
	}
	// When requested_token_type is absent, default to ICT for backward compatibility
	// with clients that omit the field. When it is present it MUST be ic_token.
	if requestedTokenType == "" {
		return true
	}
	return requestedTokenType == "urn:ietf:params:oauth:token-type:ic_token"
}

func forwardTokenRequest(ctx context.Context, cfg *AppConfiguration, fields map[string]string, dpopProof string, authorizationHeader string) (ImplResponse, error) {
	form := url.Values{}
	for key, value := range fields {
		if value != "" {
			form.Set(key, value)
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.TokenEndpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return Response(http.StatusInternalServerError, ErrorStatus{
			Error:            "server_error",
			ErrorDescription: err.Error(),
		}), nil
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// Forward the DPoP proof as an HTTP header so Keycloak can
	// bind the issued access token to the client's public key when DPoP is used
	if dpopProof != "" {
		req.Header.Set("DPoP", dpopProof)
	}
	if authorizationHeader != "" {
		req.Header.Set("Authorization", authorizationHeader)
	}

	return doJSONRequest(req)
}

func introspectToken(ctx context.Context, cfg *AppConfiguration, token string) (*introspectionResponse, *ImplResponse) {
	form := url.Values{}
	form.Set("token", token)
	form.Set("token_type_hint", "access_token")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.IntrospectionEndpoint, strings.NewReader(form.Encode()))
	if err != nil {
		resp := Response(http.StatusInternalServerError, ErrorStatus{
			Error:            "server_error",
			ErrorDescription: err.Error(),
		})
		return nil, &resp
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if credentials := os.Getenv("INTROSPECTION_CREDENTIALS"); credentials != "" {
		req.Header.Set("Authorization", credentials)
	}

	httpResp, err := http.DefaultClient.Do(req)
	if err != nil {
		resp := Response(http.StatusInternalServerError, ErrorStatus{
			Error:            "server_error",
			ErrorDescription: err.Error(),
		})
		return nil, &resp
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode == http.StatusUnauthorized {
		resp := oauthError(http.StatusUnauthorized, "invalid_client", "introspection request was rejected")
		return nil, &resp
	}
	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		resp := Response(http.StatusInternalServerError, ErrorStatus{
			Error:            "server_error",
			ErrorDescription: fmt.Sprintf("introspection endpoint returned status %d", httpResp.StatusCode),
		})
		return nil, &resp
	}

	var introspection introspectionResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&introspection); err != nil {
		resp := Response(http.StatusInternalServerError, ErrorStatus{
			Error:            "server_error",
			ErrorDescription: err.Error(),
		})
		return nil, &resp
	}

	return &introspection, nil
}

func hasRequiredScopes(granted, requested string) bool {
	grantedSet := map[string]struct{}{}
	for _, part := range strings.Fields(granted) {
		grantedSet[part] = struct{}{}
	}
	for _, part := range strings.Fields(requested) {
		if _, ok := grantedSet[part]; !ok {
			return false
		}
	}
	return true
}

// dpopProofMaxAge is the maximum acceptable age of a DPoP proof
var dpopProofMaxAge = getDPoPProofMaxAge()

func getDPoPProofMaxAge() time.Duration {
	const defaultMaxAge = 5 * time.Minute

	raw := strings.TrimSpace(os.Getenv("DPOP_MAX_AGE"))
	if raw == "" {
		return defaultMaxAge
	}

	if seconds, err := strconv.Atoi(raw); err == nil && seconds > 0 {
		return time.Duration(seconds) * time.Second
	}

	return defaultMaxAge
}

// jtiEntry holds the expiry time for a DPoP proof jti stored in the replay cache
type jtiEntry struct {
	expiresAt time.Time
}

// dpopJTICache stores recently seen DPoP proof JTIs to prevent replay attacks
var dpopJTICache sync.Map

func init() {
	// GC to maintain clear cache
	go func() {
		for range time.Tick(time.Minute) {
			now := time.Now()
			dpopJTICache.Range(func(k, v interface{}) bool {
				if e, ok := v.(jtiEntry); ok && now.After(e.expiresAt) {
					dpopJTICache.Delete(k)
				}
				return true
			})
		}
	}()
}

func validateDPoPProof(raw string, expectedJKT string, expectedHTU string) (*dpopProofClaims, error) {
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{"ES256", "ES384", "ES512", "RS256", "RS384", "RS512", "PS256", "PS384", "PS512", "EdDSA"}),
		jwt.WithIssuedAt(),
		jwt.WithLeeway(10*time.Second),
	)

	claims := &dpopProofClaims{}
	token, err := parser.ParseWithClaims(raw, claims, func(token *jwt.Token) (interface{}, error) {
		// The typ JOSE header MUST be "dpop+jwt"
		typHeader, _ := token.Header["typ"].(string)
		if !strings.EqualFold(typHeader, "dpop+jwt") {
			return nil, fmt.Errorf("dpop proof typ must be dpop+jwt, got %q", typHeader)
		}
		jwkHeader, ok := token.Header["jwk"]
		if !ok {
			return nil, errors.New("missing jwk header")
		}
		jwk, ok := jwkHeader.(map[string]interface{})
		if !ok {
			return nil, errors.New("invalid jwk header")
		}
		pub, err := publicKeyFromJWK(jwk)
		if err != nil {
			return nil, err
		}
		if expectedJKT != "" {
			thumbprint, err := jwkThumbprint(jwk)
			if err != nil {
				return nil, err
			}
			if thumbprint != expectedJKT {
				return nil, errors.New("dpop proof does not match subject token binding")
			}
		}
		return pub, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid dpop proof")
	}
	if !strings.EqualFold(claims.HTM, http.MethodPost) {
		return nil, errors.New("invalid htm claim in dpop proof")
	}
	// htu MUST match the target URI of the HTTP request
	if expectedHTU != "" {
		if !strings.EqualFold(
			strings.TrimRight(claims.HTU, "/"),
			strings.TrimRight(expectedHTU, "/"),
		) {
			return nil, fmt.Errorf("invalid htu claim: got %q, expected %q", claims.HTU, expectedHTU)
		}
	}
	if claims.JTI == "" {
		return nil, errors.New("missing jti claim in dpop proof")
	}
	if claims.IssuedAt == nil || time.Since(claims.IssuedAt.Time) > dpopProofMaxAge {
		return nil, errors.New("dpop proof is stale")
	}
	// Reject replayed DPoP proofs by checking the jti uniqueness
	expiresAt := claims.IssuedAt.Time.Add(dpopProofMaxAge)
	if _, alreadySeen := dpopJTICache.LoadOrStore(claims.JTI, jtiEntry{expiresAt: expiresAt}); alreadySeen {
		return nil, errors.New("dpop proof jti has already been used (replay attack)")
	}
	return claims, nil
}

// Maps introspection fields into ICT claim inputs
func (r *introspectionResponse) authorizedParty() string {
	if strings.TrimSpace(r.Azp) != "" {
		return strings.TrimSpace(r.Azp)
	}
	return strings.TrimSpace(r.ClientID)
}

func resolveCtxClaims(cfg *AppConfiguration, introspection *introspectionResponse) ([]string, *ImplResponse) {
	if len(cfg.AthCtx) == 0 {
		return nil, nil
	}
	azp := introspection.authorizedParty()
	if azp == "" {
		resp := oauthError(http.StatusBadRequest, "invalid_token", "subject_token is missing authorized party (azp)")
		return nil, &resp
	}
	ctxClaims, ok := cfg.AthCtx[azp]
	if !ok {
		resp := oauthError(http.StatusForbidden, "invalid_client", fmt.Sprintf("no ctx mapping configured for authorized party %q", azp))
		return nil, &resp
	}
	return ctxClaims, nil
}

func userInfoFromIntrospection(introspection *introspectionResponse) map[string]interface{} {
	userInfo := map[string]interface{}{}
	if introspection.Sub != "" {
		userInfo["sub"] = introspection.Sub
	}
	if introspection.Username != "" {
		userInfo["preferred_username"] = introspection.Username
	}
	for key, value := range introspection.Extra {
		userInfo[key] = value
	}
	return userInfo
}

func fetchUserInfo(ctx context.Context, cfg *AppConfiguration, subjectToken string) (map[string]interface{}, *ImplResponse) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cfg.UserInfoEndpoint, nil)
	if err != nil {
		resp := Response(http.StatusInternalServerError, ErrorStatus{
			Error:            "server_error",
			ErrorDescription: err.Error(),
		})
		return nil, &resp
	}
	req.Header.Set("Authorization", "Bearer "+subjectToken)

	httpResp, err := http.DefaultClient.Do(req)
	if err != nil {
		resp := Response(http.StatusInternalServerError, ErrorStatus{
			Error:            "server_error",
			ErrorDescription: err.Error(),
		})
		return nil, &resp
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode == http.StatusUnauthorized {
		resp := oauthError(http.StatusUnauthorized, "invalid_token", "userinfo request was rejected")
		return nil, &resp
	}
	if httpResp.StatusCode == http.StatusForbidden {
		resp := oauthError(http.StatusForbidden, "insufficient_scope", "userinfo request is not permitted for the subject token")
		return nil, &resp
	}
	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		resp := Response(http.StatusInternalServerError, ErrorStatus{
			Error:            "server_error",
			ErrorDescription: fmt.Sprintf("userinfo endpoint returned status %d", httpResp.StatusCode),
		})
		return nil, &resp
	}

	payload := map[string]interface{}{}
	if err := json.NewDecoder(httpResp.Body).Decode(&payload); err != nil {
		resp := Response(http.StatusInternalServerError, ErrorStatus{
			Error:            "server_error",
			ErrorDescription: err.Error(),
		})
		return nil, &resp
	}
	return payload, nil
}

func issueICT(cfg *AppConfiguration, clientID string, introspection *introspectionResponse, userInfo map[string]interface{}, ctxClaims []string, issuedAt time.Time, expiresAt time.Time) (string, error) {
	signingKey, err := loadSigningKey(cfg)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"iss": cfg.Issuer,
		"sub": introspection.Sub,
		"iat": issuedAt.Unix(),
		"nbf": issuedAt.Unix(),
		"exp": expiresAt.Unix(),
		"cnf": map[string]interface{}{
			"jkt": introspection.Cnf.Jkt,
		},
	}
	if clientID != "" {
		claims["azp"] = clientID
	}
	if introspection.ClientID != "" {
		claims["client_id"] = introspection.ClientID
	}
	if len(ctxClaims) > 0 {
		claims["ctx"] = ctxClaims
	}
	for key, value := range userInfo {
		if _, exists := claims[key]; !exists {
			claims[key] = value
		}
	}

	token := jwt.NewWithClaims(signingMethod(cfg.Algorithm), claims)
	token.Header["kid"] = cfg.KeyID
	token.Header["typ"] = "jwt+ict"
	return token.SignedString(signingKey)
}

func loadSigningKey(cfg *AppConfiguration) (interface{}, error) {
	pemData, err := os.ReadFile(cfg.KeyFile)
	if err != nil {
		return nil, err
	}

	switch cfg.Algorithm {
	case "RS256", "RS384", "RS512", "PS256", "PS384", "PS512":
		return jwt.ParseRSAPrivateKeyFromPEM(pemData)
	case "ES256", "ES384", "ES512":
		return jwt.ParseECPrivateKeyFromPEM(pemData)
	case "EdDSA":
		return jwt.ParseEdPrivateKeyFromPEM(pemData)
	default:
		return nil, fmt.Errorf("unsupported signing algorithm %q", cfg.Algorithm)
	}
}

func signingMethod(alg string) jwt.SigningMethod {
	switch alg {
	case "RS256":
		return jwt.SigningMethodRS256
	case "RS384":
		return jwt.SigningMethodRS384
	case "RS512":
		return jwt.SigningMethodRS512
	case "PS256":
		return jwt.SigningMethodPS256
	case "PS384":
		return jwt.SigningMethodPS384
	case "PS512":
		return jwt.SigningMethodPS512
	case "ES256":
		return jwt.SigningMethodES256
	case "ES384":
		return jwt.SigningMethodES384
	case "ES512":
		return jwt.SigningMethodES512
	case "EdDSA":
		return jwt.SigningMethodEdDSA
	default:
		return jwt.SigningMethodRS256
	}
}

func oauthError(code int, errCode string, description string) ImplResponse {
	return Response(code, TokenResponseError{
		Error:            errCode,
		ErrorDescription: description,
	})
}

func doJSONRequest(req *http.Request) (ImplResponse, error) {
	httpResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Response(http.StatusInternalServerError, ErrorStatus{
			Error:            "server_error",
			ErrorDescription: err.Error(),
		}), nil
	}
	defer httpResp.Body.Close()

	var body interface{}
	if err := json.NewDecoder(httpResp.Body).Decode(&body); err != nil {
		body = map[string]interface{}{
			"error":             "server_error",
			"error_description": err.Error(),
		}
	}
	return Response(httpResp.StatusCode, body), nil
}

func publicKeyFromJWK(jwk map[string]interface{}) (interface{}, error) {
	kty, _ := jwk["kty"].(string)
	switch kty {
	case "EC":
		crv, _ := jwk["crv"].(string)
		xStr, _ := jwk["x"].(string)
		yStr, _ := jwk["y"].(string)
		xBytes, err := decodeBase64URLField(xStr)
		if err != nil {
			return nil, err
		}
		yBytes, err := decodeBase64URLField(yStr)
		if err != nil {
			return nil, err
		}
		curve, err := ellipticCurveFromName(crv)
		if err != nil {
			return nil, err
		}
		return &ecdsa.PublicKey{
			Curve: curve,
			X:     new(big.Int).SetBytes(xBytes),
			Y:     new(big.Int).SetBytes(yBytes),
		}, nil
	case "RSA":
		nStr, _ := jwk["n"].(string)
		eStr, _ := jwk["e"].(string)
		nBytes, err := decodeBase64URLField(nStr)
		if err != nil {
			return nil, err
		}
		eBytes, err := decodeBase64URLField(eStr)
		if err != nil {
			return nil, err
		}
		e := 0
		for _, b := range eBytes {
			e = e<<8 + int(b)
		}
		return &rsa.PublicKey{
			N: new(big.Int).SetBytes(nBytes),
			E: e,
		}, nil
	case "OKP":
		crv, _ := jwk["crv"].(string)
		if crv != "Ed25519" {
			return nil, fmt.Errorf("unsupported OKP curve %q", crv)
		}
		xStr, _ := jwk["x"].(string)
		xBytes, err := decodeBase64URLField(xStr)
		if err != nil {
			return nil, err
		}
		return ed25519.PublicKey(xBytes), nil
	default:
		return nil, fmt.Errorf("unsupported jwk key type %q", kty)
	}
}

func jwkThumbprint(jwk map[string]interface{}) (string, error) {
	canonical, err := canonicalJWK(jwk)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256([]byte(canonical))
	return base64.RawURLEncoding.EncodeToString(sum[:]), nil
}

func trimJWKEncoding(s string) string {
	return strings.TrimRight(strings.TrimSpace(s), "=")
}

func canonicalJWK(jwk map[string]interface{}) (string, error) {
	kty, _ := jwk["kty"].(string)
	switch kty {
	case "EC":
		crv, _ := jwk["crv"].(string)
		x := trimJWKEncoding(jwk["x"].(string))
		y := trimJWKEncoding(jwk["y"].(string))
		return fmt.Sprintf("{\"crv\":\"%s\",\"kty\":\"EC\",\"x\":\"%s\",\"y\":\"%s\"}", crv, x, y), nil
	case "RSA":
		e := trimJWKEncoding(jwk["e"].(string))
		n := trimJWKEncoding(jwk["n"].(string))
		return fmt.Sprintf("{\"e\":\"%s\",\"kty\":\"RSA\",\"n\":\"%s\"}", e, n), nil
	case "OKP":
		crv, _ := jwk["crv"].(string)
		x := trimJWKEncoding(jwk["x"].(string))
		return fmt.Sprintf("{\"crv\":\"%s\",\"kty\":\"OKP\",\"x\":\"%s\"}", crv, x), nil
	default:
		return "", fmt.Errorf("unsupported jwk key type %q", kty)
	}
}

func ellipticCurveFromName(name string) (elliptic.Curve, error) {
	switch name {
	case "P-256":
		return elliptic.P256(), nil
	case "P-384":
		return elliptic.P384(), nil
	case "P-521":
		return elliptic.P521(), nil
	default:
		return nil, fmt.Errorf("unsupported elliptic curve %q", name)
	}
}
