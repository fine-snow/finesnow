// Route addition method

package snow

import (
	"github.com/fine-snow/finesnow/constant"
	"github.com/fine-snow/finesnow/router"
	"net/http"
)

// Get Add a routing method for GET request
func Get(url string, fun any) {
	router.AddRoute(constant.NullStr, url, http.MethodGet, fun)
}

// Post Add a routing method for POST request
func Post(url string, fun any) {
	router.AddRoute(constant.NullStr, url, http.MethodPost, fun)
}

// Put Add a routing method for PUT request
func Put(url string, fun any) {
	router.AddRoute(constant.NullStr, url, http.MethodPut, fun)
}

// Delete Add a routing method for DELETE request
func Delete(url string, fun any) {
	router.AddRoute(constant.NullStr, url, http.MethodDelete, fun)
}

// Group Create a route group
func Group(url string) router.RouteGroup {
	return router.NewGroup(url)
}
