/*
 * Open Identity Certification with OpenID Connect (OIDC²)
 *
 * Authorization Server middleware for requesting Identity Certification Tokens (ICT).
 *
 * API version: 0.2.0
 * Contact: mail@jonasprimbs.de
 */
package oidc2middleware

// ErrorStatus - Generic error status response
type ErrorStatus struct {

	// Error code identifying the error condition
	Error string `json:"error"`

	// Human-readable ASCII encoded text providing additional information about the error
	ErrorDescription string `json:"error_description,omitempty"`
}

// AssertErrorStatusRequired checks if the required fields are not zero-ed
func AssertErrorStatusRequired(obj ErrorStatus) error {
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

// AssertErrorStatusConstraints checks if the values respects the defined constraints
func AssertErrorStatusConstraints(obj ErrorStatus) error {
	return nil
}
