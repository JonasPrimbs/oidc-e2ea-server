/*
 * Open Identity Certification with OpenID Connect (OIDC²)
 *
 * Authorization Server middleware for requesting Identity Certification Tokens (ICT).
 *
 * API version: 0.2.0
 * Contact: mail@jonasprimbs.de
 */
package oidc2middleware

// TokenResponseError - Error status response compliant with RFC 6749 and OpenID Connect Core 1.0
type TokenResponseError struct {

	// Error code as defined in RFC 6749 and OpenID Connect Core 1.0
	Error string `json:"error"`

	// Human-readable ASCII encoded text providing additional information about the error
	ErrorDescription string `json:"error_description,omitempty"`

	// URI identifying a human-readable web page with information about the error
	ErrorUri string `json:"error_uri,omitempty"`
}

// AssertTokenResponseErrorRequired checks if the required fields are not zero-ed
func AssertTokenResponseErrorRequired(obj TokenResponseError) error {
	elements := map[string]interface{}{
		"error": obj.Error,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertTokenResponseErrorConstraints checks if the values respects the defined constraints
func AssertTokenResponseErrorConstraints(obj TokenResponseError) error {
	return nil
}
