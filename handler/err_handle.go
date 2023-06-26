// Global Exception Handling Function

package handler

import (
	"log"
	"net/http"
	"reflect"
	"runtime"
)

// ErrHandleFunc Abstract Method
type ErrHandleFunc func(err any) any

// globalErrHandleFunc Global Exception Handling Function Variables
var globalErrHandleFunc ErrHandleFunc = nil

// SetGlobalErrHandleFunc Set global exception handling functions
func SetGlobalErrHandleFunc(fun ErrHandleFunc) {
	globalErrHandleFunc = fun
}

// catchPanic Capture exceptions thrown during http request processing
func catchPanic(w http.ResponseWriter, path, method string) {
	err := recover()
	if err != nil {
		log.Println(err)
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
	}
}
