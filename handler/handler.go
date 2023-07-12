// Http request reading and parsing

package handler

import (
	"encoding/json"
	"github.com/fine-snow/finesnow/logs"
	"github.com/fine-snow/finesnow/router"
	"io"
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

// Partial Constant
// contentType HTTP request header content type constant string.
// applicationJson The json property constant of the http request header Content-Type.
// maxMemory The maximum number of bytes allowed to directly store part of the file content in memory
// when Golang's underlying processing of multipart/form-data data type requests and carrying file parameters.
const (
	contentType           = "Content-Type"
	applicationJson       = "application/json"
	maxMemory       int64 = 1048576
)

// globalHandle HTTP Request Receiving And Processing Abstract Interface Implementation Body
// intercept Global interceptor properties.
type globalHandle struct {
	intercept Interceptor
}

func NewHandle(intercept Interceptor) http.Handler {
	handle := allowCORS(&globalHandle{intercept})
	return handle
}

// ServeHTTP Implement the golang underlying Handler interface
func (gh *globalHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	method := r.Method
	route := router.Get(path, method, r)
	if route == nil {
		text := http.StatusText(http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(text))
		return
	}
	rt := route.GetType()
	rv := route.GetValue()
	numIn := rt.NumIn()
	names := route.GetParamNames()
	w.Header().Set(contentType, string(*route.GetHttpContentType()))
	if gh.intercept != nil && !gh.intercept(w, r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	defer catchHttpPanic(w)
	if numIn == 0 {
		outParam := rv.Call(nil)
		if outParam == nil {
			return
		}
		_, _ = w.Write(convertToByteArray(outParam[0]))
		return
	}
	if method == http.MethodGet {
		outParam := rv.Call(dealInParam(names, rt, r.URL.Query(), nil, w, r))
		if outParam == nil {
			return
		}
		_, _ = w.Write(convertToByteArray(outParam[0]))
		return
	}
	ct := r.Header.Get(contentType)
	if strings.Contains(ct, applicationJson) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		in := make([]reflect.Value, numIn, numIn)
		for i := 0; i < numIn; i++ {
			t := rt.In(i)
			if t.String() == pRequest {
				in[i] = reflect.ValueOf(r)
				continue
			}
			if t.String() == response {
				in[i] = reflect.ValueOf(w)
				continue
			}
			pointer := reflect.New(t)
			err = json.Unmarshal(body, pointer.Interface())
			if err != nil {
				panic(err)
			}
			in[i] = pointer.Elem()
		}
		outParam := rv.Call(in)
		if outParam == nil {
			return
		}
		_, _ = w.Write(convertToByteArray(outParam[0]))
		return
	}
	_ = r.ParseMultipartForm(maxMemory)
	multipartForm := r.MultipartForm
	if multipartForm != nil {
		outParam := rv.Call(dealInParam(names, rt, multipartForm.Value, multipartForm.File, w, r))
		if outParam == nil {
			return
		}
		_, _ = w.Write(convertToByteArray(outParam[0]))
		return
	}
	postForm := r.PostForm
	if postForm != nil {
		outParam := rv.Call(dealInParam(names, rt, postForm, nil, w, r))
		if outParam == nil {
			return
		}
		_, _ = w.Write(convertToByteArray(outParam[0]))
		return
	}
}

// ErrHandleFunc Abstract Method
type ErrHandleFunc func(err any) any

// globalErrHandleFunc Global Exception Handling Function Variables
var globalErrHandleFunc ErrHandleFunc

// SetGlobalErrHandleFunc Set global exception handling functions
func SetGlobalErrHandleFunc(fun ErrHandleFunc) {
	globalErrHandleFunc = fun
}

// catchHttpPanic Capture exceptions thrown during http request processing
func catchHttpPanic(w http.ResponseWriter) {
	err := recover()
	if err != nil {
		logs.ERROR(err)
		w.WriteHeader(http.StatusInternalServerError)
		switch err.(type) {
		case runtime.Error:
			_, _ = w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		default:
			if globalErrHandleFunc != nil {
				err = globalErrHandleFunc(err)
			}
			errBytes := convertToByteArray(reflect.ValueOf(err))
			_, _ = w.Write(errBytes)
		}
		return
	}
}

// CatchRunPanic Capture exceptions generated during framework startup process
func CatchRunPanic() {
	err := recover()
	if err != nil {
		logs.ERROR(err)
		return
	}
}
