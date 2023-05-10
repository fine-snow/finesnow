// Cross-domain configuration

package snow

import "github.com/fine-snow/finesnow/handler"

func SetAllowedOrigin(s string) {
	handler.SetAllowedOrigin(s)
}

func SetAllowedMethods(s []string) {
	handler.SetAllowedMethods(s)
}

func SetAllowedHeaders(s []string) {
	handler.SetAllowedHeaders(s)
}
