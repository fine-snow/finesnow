// Http request processing tool method

package handler

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"strconv"
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
		panic("argument out of range")
	}
}

// catchPanic Capture exceptions thrown during http request processing
func catchPanic(w http.ResponseWriter) {
	err := recover()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		switch err.(type) {
		case runtime.Error:
			_, _ = w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		default:
			_, _ = w.Write(convertToByteArray(reflect.ValueOf(err)))
		}
	}
}

type Files map[string][]*multipart.FileHeader

// dealInParam Http request input processing method
func dealInParam(paramNames []string, rt reflect.Type, values url.Values, files Files) []reflect.Value {
	var in []reflect.Value
	for i, k := range paramNames {
		t := rt.In(i)
		if t.String() == "*multipart.FileHeader" {
			headers := files[k]
			if headers != nil {
				in = append(in, reflect.ValueOf(headers[0]))
			} else {
				in = append(in, reflect.ValueOf(&multipart.FileHeader{}))
			}
			continue
		}
		if t.String() == "multipart.FileHeader" {
			headers := files[k]
			if headers != nil {
				in = append(in, reflect.ValueOf(*(headers[0])))
			} else {
				in = append(in, reflect.ValueOf(multipart.FileHeader{}))
			}
			continue
		}
		if t.String() == "[]*multipart.FileHeader" {
			in = append(in, reflect.ValueOf(files[k]))
			continue
		}
		if t.String() == "[]multipart.FileHeader" {
			var fs []multipart.FileHeader
			for _, f := range files[k] {
				fs = append(fs, *f)
			}
			in = append(in, reflect.ValueOf(fs))
			continue
		}
		switch t.Kind() {
		case reflect.Bool:
			v, _ := strconv.ParseBool(values.Get(k))
			in = append(in, reflect.ValueOf(v))
		case reflect.String:
			in = append(in, reflect.ValueOf(values.Get(k)))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v, _ := strconv.ParseInt(values.Get(k), 0, 64)
			elem := reflect.New(t).Elem()
			elem.SetInt(v)
			in = append(in, elem)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v, _ := strconv.ParseUint(values.Get(k), 0, 64)
			elem := reflect.New(t).Elem()
			elem.SetUint(v)
			in = append(in, elem)
		case reflect.Float32, reflect.Float64:
			v, _ := strconv.ParseFloat(values.Get(k), 64)
			elem := reflect.New(t).Elem()
			elem.SetFloat(v)
			in = append(in, elem)
		default:
			panic("argument out of range")
		}
	}
	return in
}
