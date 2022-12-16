/*
 * OIDC IAT Userinfo Endpoint
 *
 * Endpoint for OpenID Connect's ID Assertion Token endpoint for userinfo.
 *
 * API version: 0.2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package iat

type IatResponse struct {
	IdAssertionToken string `json:"id_assertion_token"`
	// Number of seconds until the ID Assertion Token expires.
	ExpiresIn int32 `json:"expires_in,omitempty"`
	// Space delimited claims provided in the ID Assertion Token.
	Claims string `json:"claims,omitempty"`
}
