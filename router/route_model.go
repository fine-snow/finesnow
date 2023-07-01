// Http request routing model

package router

import "reflect"

// RouteModel Route features abstract interface model
type RouteModel interface {
	GetHttpContentType() *httpContentType
	GetType() reflect.Type
	GetValue() reflect.Value
	GetParamNames() []string
}

// routeModel The route model implements the structure
type routeModel struct {
	hct        *httpContentType
	t          reflect.Type
	v          reflect.Value
	paramNames []string
}

func (rm *routeModel) GetHttpContentType() *httpContentType {
	return rm.hct
}

func (rm *routeModel) GetType() reflect.Type {
	return rm.t
}

func (rm *routeModel) GetValue() reflect.Value {
	return rm.v
}

func (rm *routeModel) GetParamNames() []string {
	return rm.paramNames
}
