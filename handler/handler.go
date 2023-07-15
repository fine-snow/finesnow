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

// Interceptor Abstract Method
// Returning true indicates a release request, while false indicates the opposite.
type Interceptor func(http.ResponseWriter, *http.Request) bool

var intercept Interceptor

func SetIntercept(i Interceptor) {
	intercept = i
}

// PostProcessor Abstract Method
// The routing function is called when the return value is finished processing data.
type PostProcessor func(v any) any

var postProcess PostProcessor

func SetPostProcess(p PostProcessor) {
	postProcess = p
}

// snowHandler HTTP Request Receiving And Processing Abstract Interface Implementation Body
type snowHandler struct {
}

func NewHandle() http.Handler {
	handle := allowCORS(&snowHandler{})
	return handle
}

// ServeHTTP Implement the golang underlying Handler interface
func (sh *snowHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	if intercept != nil && !intercept(w, r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	defer catchHttpPanic(w)
	if numIn == 0 {
		outParam := rv.Call(nil)
		if outParam == nil {
			return
		}
		if postProcess != nil {
			outParam[0] = reflect.ValueOf(postProcess(outParam[0].Interface()))
		}
		_, _ = w.Write(convertToByteArray(outParam[0]))
		return
	}
	if method == http.MethodGet {
		outParam := rv.Call(dealInParam(names, rt, r.URL.Query(), nil, w, r))
		if outParam == nil {
			return
		}
		if postProcess != nil {
			outParam[0] = reflect.ValueOf(postProcess(outParam[0].Interface()))
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
		if postProcess != nil {
			outParam[0] = reflect.ValueOf(postProcess(outParam[0].Interface()))
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
		if postProcess != nil {
			outParam[0] = reflect.ValueOf(postProcess(outParam[0].Interface()))
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
		if postProcess != nil {
			outParam[0] = reflect.ValueOf(postProcess(outParam[0].Interface()))
		}
		_, _ = w.Write(convertToByteArray(outParam[0]))
		return
	}
}

// ErrHandleFunc Abstract Method
type ErrHandleFunc func(err any) any

// errHandleFunc Global Exception Handling Function Variables
var errHandleFunc ErrHandleFunc

// SetErrHandleFunc Set global exception handling functions
func SetErrHandleFunc(fun ErrHandleFunc) {
	errHandleFunc = fun
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
			if errHandleFunc != nil {
				err = errHandleFunc(err)
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
