// Interceptor

package handler

import "net/http"

// Interceptor Abstract Method
// Returning true indicates a release request, while false indicates the opposite.
type Interceptor func(http.ResponseWriter, *http.Request) bool

// defaultInterceptor Default Interceptor Method
// This method does not have any interception processing.
func defaultInterceptor(http.ResponseWriter, *http.Request) bool {
	return true
}
