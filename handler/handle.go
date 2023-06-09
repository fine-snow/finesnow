// Http request reading and parsing

package handler

import (
	"encoding/json"
	"github.com/fine-snow/finesnow/router"
	"io"
	"net/http"
	"reflect"
	"strings"
)

// Partial Constant
// contentType HTTP request header content type constant string.
// maxMemory The maximum number of bytes allowed to directly store part of the file content in memory
// when Golang's underlying processing of multipart/form-data data type requests and carrying file parameters.
const (
	contentType       = "Content-Type"
	maxMemory   int64 = 1048576
)

// Handle Http Request Receive Processing Abstract Interface
// ServeHTTP The method of processing HTTP requests at the bottom of Golang.
type Handle interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

// globalHandle HTTP Request Receiving And Processing Abstract Interface Implementation Body
// intercept Global interceptor properties.
type globalHandle struct {
	intercept Interceptor
}

func NewHandle(intercept Interceptor) http.Handler {
	// If the interceptor parameter passed in by the startup method is empty, supplement the default interceptor method
	if intercept == nil {
		intercept = defaultInterceptor
	}
	handle := allowCORS(&globalHandle{intercept})
	return handle
}

func (gh *globalHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	method := r.Method
	route := router.Get(path)
	if route == nil {
		text := http.StatusText(http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(text))
		return
	}
	if method != string(*(route.GetHttpMethod())) {
		text := http.StatusText(http.StatusMethodNotAllowed)
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte(text))
		return
	}
	w.Header().Set(contentType, string(*route.GetHttpContentType()))
	if !gh.intercept(w, r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	defer catchPanic(w, path, method)
	if route.GetType().NumIn() == 0 {
		outParam := route.GetValue().Call(nil)
		if outParam == nil {
			return
		}
		_, _ = w.Write(convertToByteArray(outParam[0]))
		return
	}

	if method == string(*router.HttpMethodGet) {
		outParam := route.GetValue().Call(dealInParam(route.GetParamNames(), route.GetType(), r.URL.Query(), nil, w, r))
		if outParam == nil {
			return
		}
		_, _ = w.Write(convertToByteArray(outParam[0]))
		return
	}

	ct := r.Header.Get(contentType)
	if strings.Contains(ct, "application/json") {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		rt := route.GetType()
		var in []reflect.Value
		for i := 0; i < rt.NumIn(); i++ {
			t := rt.In(i)
			if t.String() == "*http.Request" {
				in = append(in, reflect.ValueOf(r))
				continue
			}
			if t.String() == "http.ResponseWriter" {
				in = append(in, reflect.ValueOf(w))
				continue
			}
			pointer := reflect.New(t)
			err = json.Unmarshal(body, pointer.Interface())
			if err != nil {
				panic(err)
			}
			in = append(in, pointer.Elem())
		}
		outParam := route.GetValue().Call(in)
		if outParam == nil {
			return
		}
		_, _ = w.Write(convertToByteArray(outParam[0]))
		return
	}

	_ = r.ParseMultipartForm(maxMemory)

	multipartForm := r.MultipartForm
	if multipartForm != nil {
		outParam := route.GetValue().Call(dealInParam(route.GetParamNames(), route.GetType(), multipartForm.Value, multipartForm.File, w, r))
		if outParam == nil {
			return
		}
		_, _ = w.Write(convertToByteArray(outParam[0]))
		return
	}

	postForm := r.PostForm
	if postForm != nil {
		outParam := route.GetValue().Call(dealInParam(route.GetParamNames(), route.GetType(), postForm, nil, w, r))
		if outParam == nil {
			return
		}
		_, _ = w.Write(convertToByteArray(outParam[0]))
		return
	}
}
