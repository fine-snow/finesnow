// Route addition method

package snow

import "github.com/fine-snow/finesnow/router"

// AddGetRoute Add a routing method for GET request
func AddGetRoute(url string, fun interface{}) {
	router.AddRoute(url, fun, router.HttpMethodGet)
}

// AddPostRoute Add a routing method for POST request
func AddPostRoute(url string, fun interface{}) {
	router.AddRoute(url, fun, router.HttpMethodPost)
}

// AddPutRoute Add a routing method for PUT request
func AddPutRoute(url string, fun interface{}) {
	router.AddRoute(url, fun, router.HttpMethodPut)
}

// AddDeleteRoute Add a routing method for DELETE request
func AddDeleteRoute(url string, fun interface{}) {
	router.AddRoute(url, fun, router.HttpMethodDelete)
}
