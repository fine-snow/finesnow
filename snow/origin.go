// Cross-domain configuration

package snow

import "github.com/fine-snow/finesnow/handler"

// SetAllowedOrigin Set domains that allow cross domain access
func SetAllowedOrigin(s string) {
	handler.SetAllowedOrigin(s)
}

// SetAllowedMethods Set request methods that allow cross domain requests
func SetAllowedMethods(s []string) {
	handler.SetAllowedMethods(s)
}

// SetAllowedHeaders Set request headers that allow cross domain connections
func SetAllowedHeaders(s []string) {
	handler.SetAllowedHeaders(s)
}
