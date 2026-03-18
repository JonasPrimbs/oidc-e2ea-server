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
	"log"
	"net/http"
	"time"
)

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s %s %s %s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
