// Route Group

package router

import (
	"github.com/fine-snow/finesnow/constant"
	"net/http"
	"strings"
)

type RouteGroup interface {
	Get(url string, fun any) RouteGroup
	Post(url string, fun any) RouteGroup
	Put(url string, fun any) RouteGroup
	Delete(url string, fun any) RouteGroup
}

type routeGroup struct {
	url string
}

func NewGroup(url string) RouteGroup {
	url = strings.ReplaceAll(url, constant.Space, constant.NullStr)
	url = dealPrefixSlash(url)
	url = dealSuffixSlash(url)
	if url == constant.NullStr || url == constant.Slash {
		panic(errRouteUrlIsNilOrSlash)
	}
	return &routeGroup{url: url}
}

func (rg *routeGroup) Get(url string, fun any) RouteGroup {
	AddRoute(rg.url, url, http.MethodGet, fun)
	return rg
}

func (rg *routeGroup) Post(url string, fun any) RouteGroup {
	AddRoute(rg.url, url, http.MethodPost, fun)
	return rg
}

func (rg *routeGroup) Put(url string, fun any) RouteGroup {
	AddRoute(rg.url, url, http.MethodPut, fun)
	return rg
}

func (rg *routeGroup) Delete(url string, fun any) RouteGroup {
	AddRoute(rg.url, url, http.MethodDelete, fun)
	return rg
}
