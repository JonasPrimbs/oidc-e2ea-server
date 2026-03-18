/*
 * Open Identity Certification with OpenID Connect (OIDC²)
 *
 * Authorization Server middleware for requesting Identity Certification Tokens (ICT).
 *
 * API version: 0.2.0
 * Contact: mail@jonasprimbs.de
 */
package oidc2middleware

// ImplResponse defines an implementation response with error code and the associated body
type ImplResponse struct {
	Code int
	Body interface{}
}
