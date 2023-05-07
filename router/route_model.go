// Http request routing model

package router

import "reflect"

type RouteModel interface {
	GetHttpMethod() *httpMethod
	GetHttpContentType() *httpContentType
	GetType() reflect.Type
	GetValue() reflect.Value
	GetParamNames() []string
}

type routeModel struct {
	hm         *httpMethod
	hct        *httpContentType
	t          reflect.Type
	v          reflect.Value
	paramNames []string
}

func (rm *routeModel) GetHttpMethod() *httpMethod {
	return rm.hm
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
