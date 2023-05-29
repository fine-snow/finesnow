// Route addition method

package snow

import "github.com/fine-snow/finesnow/router"

// AddGetRoute Add a routing method for get requests
func AddGetRoute(url string, fun interface{}) {
	router.AddRoute(url, fun, router.HttpMethodGet)
}

// AddPostRoute Add a routing method for post requests
func AddPostRoute(url string, fun interface{}) {
	router.AddRoute(url, fun, router.HttpMethodPost)
}

func AddPutRoute(url string, fun interface{}) {
	router.AddRoute(url, fun, router.HttpMethodPut)
}
