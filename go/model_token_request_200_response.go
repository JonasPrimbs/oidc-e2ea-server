/*
 * Open Identity Certification with OpenID Connect (OIDC²)
 *
 * Authorization Server middleware for requesting Identity Certification Tokens (ICT).
 *
 * API version: 0.2.0
 * Contact: mail@jonasprimbs.de
 */
package oidc2middleware

type TokenRequest200Response struct {

	// The access token issued by the authorization server.
	AccessToken string `json:"access_token" validate:"regexp=^[a-zA-Z0-9._~\\\\-\\/+=]+$"`

	// Token type (typically \"Bearer\" or \"DPoP\").
	TokenType string `json:"token_type"`

	// Lifetime in seconds of the access token.
	ExpiresIn int64 `json:"expires_in,omitempty"`

	// Refresh token for obtaining new access tokens.
	RefreshToken string `json:"refresh_token,omitempty" validate:"regexp=^[a-zA-Z0-9._~\\\\-\\/+=]+$"`

	// Scope of the access token.
	Scope string `json:"scope,omitempty" validate:"regexp=^[!#-\\\\[\\\\]-~]+$"`

	// OPTIONAL (OpenID Connect). ID Token (JWT).
	IdToken string `json:"id_token,omitempty" validate:"regexp=^[a-zA-Z0-9\\\\-_]+\\\\.[a-zA-Z0-9\\\\-_]+\\\\.[a-zA-Z0-9\\\\-_]+$"`

	// OPTIONAL (RFC 9449). DPoP proof key thumbprint.
	DpopJkt string `json:"dpop_jkt,omitempty" validate:"regexp=^[a-zA-Z0-9\\\\-_]{43}$"`

	// OPTIONAL (OpenID Connect). Code hash for authorization code flow.
	CHash string `json:"c_hash,omitempty" validate:"regexp=^([a-zA-Z0-9\\\\-_]{22}|[a-zA-Z0-9\\\\-_]{32}|[a-zA-Z0-9\\\\-_]{43})$"`

	// OPTIONAL (RFC 8693). The type of the issued token.
	IssuedTokenType string `json:"issued_token_type,omitempty" validate:"regexp=^urn:ietf:params:oauth:token-type:[a-zA-Z0-9._\\\\-]+$"`

	// OPTIONAL (RFC 8693). The type of the requested token for confirmation.
	RequestedTokenType string `json:"requested_token_type,omitempty" validate:"regexp=^urn:ietf:params:oauth:token-type:[a-zA-Z0-9._\\\\-]+$"`

	// OPTIONAL. State value from the authorization request.
	State string `json:"state,omitempty" validate:"regexp=^[a-zA-Z0-9._~\\\\-]{8,512}$"`

	// OPTIONAL. Client identifier for which the token was issued for confirmation.
	ClientId string `json:"client_id,omitempty" validate:"regexp=^[!-~]+$"`

	// OPTIONAL (proprietary for Keycloak). The time before which the token MUST NOT be accepted for processing.
	NotBeforePolicy int64 `json:"not-before-policy,omitempty"`
}

// AssertTokenRequest200ResponseRequired checks if the required fields are not zero-ed
func AssertTokenRequest200ResponseRequired(obj TokenRequest200Response) error {
	elements := map[string]interface{}{
		"access_token": obj.AccessToken,
		"token_type":   obj.TokenType,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertTokenRequest200ResponseConstraints checks if the values respects the defined constraints
func AssertTokenRequest200ResponseConstraints(obj TokenRequest200Response) error {
	return nil
}
