// Cross-domain processing

package handler

import (
	"github.com/fine-snow/finesnow/constant"
	"net/http"
	"strings"
)

// Constant declarations
const (
	allow         = "Allow"
	origin        = "Origin"
	maxAge        = "Max-Age"
	maxAgeValue   = "3600"
	allowMethods  = "Access-Control-Allow-Methods"
	allowHeaders  = "Access-Control-Allow-Headers"
	allowOrigin   = "Access-Control-Allow-Origin"
	authorization = "Authorization"
)

var (
	allowedOrigin  = "*"
	allowedMethods = []string{http.MethodGet, http.MethodPost}
	allowedHeaders = []string{contentType, authorization}
)

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
		ori := r.Header.Get(origin)
		if ori != constant.NullStr {
			// Set CORS headers on response
			w.Header().Set(allowMethods, strings.Join(allowedMethods, constant.Comma))
			w.Header().Set(allowHeaders, strings.Join(allowedHeaders, constant.Comma))
			w.Header().Set(allowOrigin, allowedOrigin)
			// Handle preflight request
			if r.Method == http.MethodOptions {
				w.Header().Set(allow, strings.Join(allowedMethods, constant.Comma))
				w.Header().Set(maxAge, maxAgeValue)
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		// Call the next middleware/handler in chain
		next.ServeHTTP(w, r)
	})
}
