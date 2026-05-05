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
	"net/http"
	"strings"
)

// DefaultAPIController binds http requests to an api service and writes the service results to the http response
type DefaultAPIController struct {
	service      DefaultAPIServicer
	errorHandler ErrorHandler
}

// DefaultAPIOption for how the controller is set up.
type DefaultAPIOption func(*DefaultAPIController)

// WithDefaultAPIErrorHandler inject ErrorHandler into controller
func WithDefaultAPIErrorHandler(h ErrorHandler) DefaultAPIOption {
	return func(c *DefaultAPIController) {
		c.errorHandler = h
	}
}

// NewDefaultAPIController creates a default api controller
func NewDefaultAPIController(s DefaultAPIServicer, opts ...DefaultAPIOption) *DefaultAPIController {
	controller := &DefaultAPIController{
		service:      s,
		errorHandler: DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// Routes returns all the api routes for the DefaultAPIController
func (c *DefaultAPIController) Routes() Routes {
	return Routes{
		"TokenRequest": Route{
			"TokenRequest",
			strings.ToUpper("Post"),
			"/",
			c.TokenRequest,
		},
		"TokenRequestPreflight": Route{
			"TokenRequestPreflight",
			strings.ToUpper("Options"),
			"/",
			c.TokenRequestPreflight,
		},
	}
}

// OrderedRoutes returns all the api routes in a deterministic order for the DefaultAPIController
func (c *DefaultAPIController) OrderedRoutes() []Route {
	return []Route{
		Route{
			"TokenRequest",
			strings.ToUpper("Post"),
			"/",
			c.TokenRequest,
		},
		Route{
			"TokenRequestPreflight",
			strings.ToUpper("Options"),
			"/",
			c.TokenRequestPreflight,
		},
	}
}

// TokenRequest - Request a token.
func (c *DefaultAPIController) TokenRequest(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}

	grantTypeParam := r.FormValue("grant_type")

	clientIdParam := r.FormValue("client_id")

	codeParam := r.FormValue("code")

	redirectUriParam := r.FormValue("redirect_uri")

	clientSecretParam := r.FormValue("client_secret")

	refreshTokenParam := r.FormValue("refresh_token")

	scopeParam := r.FormValue("scope")

	resourceParam := r.FormValue("resource")

	audienceParam := r.FormValue("audience")

	subjectTokenParam := r.FormValue("subject_token")

	subjectTokenTypeParam := r.FormValue("subject_token_type")
	
	requestedTokenTypeParam := r.FormValue("requested_token_type")

	actorTokenParam := r.FormValue("actor_token")

	actorTokenTypeParam := r.FormValue("actor_token_type")

	// Accept DPoP proof from the standard DPoP HTTP header
	// Fallback: The dpop form field for backward compatibility
	dpopParam := r.Header.Get("DPoP")
	if dpopParam == "" {
		dpopParam = r.FormValue("dpop")
	}

	authorizationParam := r.Header.Get("Authorization")

	// Reconstruct the full request URL so the service can validate the DPoP htu claim
	scheme := "https"
	if fwdProto := r.Header.Get("X-Forwarded-Proto"); fwdProto != "" {
		scheme = fwdProto
	} else if r.TLS == nil {
		scheme = "http"
	}
	host := r.Host
	if fwdHost := r.Header.Get("X-Forwarded-Host"); fwdHost != "" {
		host = fwdHost
	}
	requestURLParam := scheme + "://" + host + r.URL.Path

	result, err := c.service.TokenRequest(r.Context(), grantTypeParam, clientIdParam, codeParam, redirectUriParam, clientSecretParam, refreshTokenParam, scopeParam, resourceParam, audienceParam, subjectTokenParam, subjectTokenTypeParam, requestedTokenTypeParam, actorTokenParam, actorTokenTypeParam, dpopParam, authorizationParam, requestURLParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	_ = EncodeJSONResponse(result.Body, &result.Code, w)
}

// TokenRequestPreflight - CORS preflight request
func (c *DefaultAPIController) TokenRequestPreflight(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, DPoP")
	w.WriteHeader(http.StatusNoContent)
}
