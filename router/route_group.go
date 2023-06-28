// Route Group

package router

import (
	"github.com/fine-snow/finesnow/constant"
	"net/http"
	"strings"
)

// RouteGroup Route group abstract interface
type RouteGroup interface {
	Get(url string, fun any) RouteGroup
	Post(url string, fun any) RouteGroup
	Put(url string, fun any) RouteGroup
	Delete(url string, fun any) RouteGroup
}

// routeGroup Routing group functional structure
type routeGroup struct {
	url string
}

// NewGroup Create a routing group struct object
func NewGroup(url string) RouteGroup {
	url = strings.ReplaceAll(url, constant.Space, constant.NullStr)
	url = dealPrefixSlash(url)
	url = dealSuffixSlash(url)
	if url == constant.NullStr || url == constant.Slash {
		panic(errRouteUrlIsNilOrSlash)
	}
	return &routeGroup{url: url}
}

// Get The routing group adds an HTTP GET request
func (rg *routeGroup) Get(url string, fun any) RouteGroup {
	AddRoute(rg.url, url, http.MethodGet, fun)
	return rg
}

// Post The routing group adds an HTTP POST request
func (rg *routeGroup) Post(url string, fun any) RouteGroup {
	AddRoute(rg.url, url, http.MethodPost, fun)
	return rg
}

// Put The routing group adds an HTTP PUT request
func (rg *routeGroup) Put(url string, fun any) RouteGroup {
	AddRoute(rg.url, url, http.MethodPut, fun)
	return rg
}

// Delete The routing group adds an HTTP DELETE request
func (rg *routeGroup) Delete(url string, fun any) RouteGroup {
	AddRoute(rg.url, url, http.MethodDelete, fun)
	return rg
}
