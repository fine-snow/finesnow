// Http request processing tool method

package handler

import (
	"encoding/json"
	"errors"
	"github.com/fine-snow/finesnow/constant"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

var outRange = errors.New("argument out of range")

const (
	multipartFileHeader   = "multipart.FileHeader"
	pMultipartFileHeader  = "*multipart.FileHeader"
	multipartFileHeaders  = "[]multipart.FileHeader"
	pMultipartFileHeaders = "[]*multipart.FileHeader"

	response = "http.ResponseWriter"
	pRequest = "*http.Request"
)

// convertToByteArray Convert the reflection.Value of a value into a byte array
func convertToByteArray(value reflect.Value) []byte {
	switch value.Kind() {
	case reflect.Bool,
		reflect.Struct,
		reflect.Float32, reflect.Float64,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		bytes, _ := json.Marshal(value.Interface())
		return bytes
	case reflect.String:
		return []byte(value.String())
	case reflect.Pointer, reflect.Interface:
		return convertToByteArray(value.Elem())
	default:
		panic(outRange)
	}
}

// catchPanic Capture exceptions thrown during http request processing
func catchPanic(w http.ResponseWriter, path, method string) {
	err := recover()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		switch err.(type) {
		case runtime.Error:
			_, _ = w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		default:
			errBytes := convertToByteArray(reflect.ValueOf(err))
			_, _ = w.Write(errBytes)
		}
	}
}

type Files map[string][]*multipart.FileHeader

// dealInParam Http request input processing method
func dealInParam(paramNames []string, rt reflect.Type, values url.Values, files Files, w http.ResponseWriter, r *http.Request) []reflect.Value {
	var in []reflect.Value
	for i, k := range paramNames {
		t := rt.In(i)
		if strings.Contains(t.String(), multipartFileHeader) {
			if t.String() == pMultipartFileHeader {
				headers := files[k]
				if headers != nil {
					in = append(in, reflect.ValueOf(headers[constant.Zero]))
				} else {
					in = append(in, reflect.ValueOf(&multipart.FileHeader{}))
				}
				continue
			}
			if t.String() == pMultipartFileHeaders {
				in = append(in, reflect.ValueOf(files[k]))
				continue
			}
			if t.String() == multipartFileHeader {
				headers := files[k]
				if headers != nil {
					in = append(in, reflect.ValueOf(*(headers[constant.Zero])))
				} else {
					in = append(in, reflect.ValueOf(multipart.FileHeader{}))
				}
				continue
			}
			if t.String() == multipartFileHeaders {
				var fs []multipart.FileHeader
				for _, f := range files[k] {
					fs = append(fs, *f)
				}
				in = append(in, reflect.ValueOf(fs))
				continue
			}
		}
		if t.String() == response {
			in = append(in, reflect.ValueOf(w))
			continue
		}
		if t.String() == pRequest {
			in = append(in, reflect.ValueOf(r))
			continue
		}
		switch t.Kind() {
		case reflect.Bool:
			v, _ := strconv.ParseBool(values.Get(k))
			in = append(in, reflect.ValueOf(v))
		case reflect.String:
			in = append(in, reflect.ValueOf(values.Get(k)))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v, _ := strconv.ParseInt(values.Get(k), constant.Zero, constant.SixtyFour)
			elem := reflect.New(t).Elem()
			elem.SetInt(v)
			in = append(in, elem)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v, _ := strconv.ParseUint(values.Get(k), constant.Zero, constant.SixtyFour)
			elem := reflect.New(t).Elem()
			elem.SetUint(v)
			in = append(in, elem)
		case reflect.Float32, reflect.Float64:
			v, _ := strconv.ParseFloat(values.Get(k), constant.SixtyFour)
			elem := reflect.New(t).Elem()
			elem.SetFloat(v)
			in = append(in, elem)
		default:
			panic(outRange)
		}
	}
	return in
}
