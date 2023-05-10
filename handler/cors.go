// Cross-domain processing

package handler

import (
	"github.com/fine-snow/finesnow/router"
	"net/http"
	"strings"
)

var allowedOrigin = "*"
var allowedMethods = []string{"GET", "POST"}
var allowedHeaders = []string{"Content-Type", "Authorization"}

func SetAllowedOrigin(s string) {
	allowedOrigin = s
}

func SetAllowedMethods(s []string) {
	allowedMethods = s
}

func SetAllowedHeaders(s []string) {
	allowedHeaders = s
}

// allowCORS
func allowCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			// Set CORS headers on response
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ","))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ","))
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			// Handle preflight request
			if r.Method == router.HttpMethodOptions {
				w.Header().Set("Allow", strings.Join(allowedMethods, ","))
				w.Header().Set("Max-Age", "3600")
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		// Call the next middleware/handler in chain
		next.ServeHTTP(w, r)
	})
}
