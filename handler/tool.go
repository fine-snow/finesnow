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

type Files map[string][]*multipart.FileHeader

// dealInParam Http request input processing method
func dealInParam(paramNames []string, rt reflect.Type, values url.Values, files Files, w http.ResponseWriter, r *http.Request) []reflect.Value {
	l := len(paramNames)
	in := make([]reflect.Value, l, l)
	for i, k := range paramNames {
		t := rt.In(i)
		if strings.Contains(t.String(), multipartFileHeader) {
			if t.String() == pMultipartFileHeader {
				headers := files[k]
				if headers != nil {
					in[i] = reflect.ValueOf(headers[constant.Zero])
				} else {
					var header *multipart.FileHeader
					in[i] = reflect.ValueOf(header)
				}
				continue
			}
			if t.String() == pMultipartFileHeaders {
				in[i] = reflect.ValueOf(files[k])
				continue
			}
			if t.String() == multipartFileHeader {
				headers := files[k]
				if headers != nil {
					in[i] = reflect.ValueOf(*(headers[constant.Zero]))
				} else {
					in[i] = reflect.ValueOf(multipart.FileHeader{})
				}
				continue
			}
			if t.String() == multipartFileHeaders {
				var fs []multipart.FileHeader
				headers := files[k]
				for _, f := range headers {
					fs = append(fs, *f)
				}
				in[i] = reflect.ValueOf(fs)
				continue
			}
		}
		if t.String() == response {
			in[i] = reflect.ValueOf(w)
			continue
		}
		if t.String() == pRequest {
			in[i] = reflect.ValueOf(r)
			continue
		}
		switch t.Kind() {
		case reflect.Bool:
			v, _ := strconv.ParseBool(values.Get(k))
			in[i] = reflect.ValueOf(v)
		case reflect.String:
			in[i] = reflect.ValueOf(values.Get(k))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v, _ := strconv.ParseInt(values.Get(k), constant.Zero, constant.SixtyFour)
			elem := reflect.New(t).Elem()
			elem.SetInt(v)
			in[i] = elem
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v, _ := strconv.ParseUint(values.Get(k), constant.Zero, constant.SixtyFour)
			elem := reflect.New(t).Elem()
			elem.SetUint(v)
			in[i] = elem
		case reflect.Float32, reflect.Float64:
			v, _ := strconv.ParseFloat(values.Get(k), constant.SixtyFour)
			elem := reflect.New(t).Elem()
			elem.SetFloat(v)
			in[i] = elem
		default:
			panic(outRange)
		}
	}
	return in
}
