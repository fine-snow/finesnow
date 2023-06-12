// Routing model collection

package router

import (
	"errors"
	"fmt"
	"github.com/fine-snow/finesnow/constant"
	"go/ast"
	"go/parser"
	"go/token"
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

var (
	errRouteFuncIsNil           = errors.New("the route mapping add function 'fun' parameter value cannot be nil")
	errRouteAddNotFunc          = errors.New("the route mapping add function 'fun' parameter value is not a function")
	errRouteFuncOutAbnormal     = errors.New("the route mapping add function 'fun' parameter value function return value amount is greater than 1")
	errRouteDuplicateDefinition = errors.New("duplicate definition of routing path")
	errRouteNoHttpMethod        = errors.New("there is no maintenance routing HTTP request type")
	errRouteUrlIsNilOrSlash     = errors.New("the route url cannot be null character or '/'")

	fest       = token.NewFileSet()
	astFileMap = make(map[string]*ast.File)

	getRouteModelMap    = make(map[string]RouteModel)
	postRouteModelMap   = make(map[string]RouteModel)
	putRouteModelMap    = make(map[string]RouteModel)
	deleteRouteModelMap = make(map[string]RouteModel)
)

func checkFun(t reflect.Type) {
	if t.Kind() != reflect.Func {
		panic(errRouteAddNotFunc)
	}
	if t.NumOut() > 1 {
		panic(errRouteFuncOutAbnormal)
	}
}

func put(url string, rm RouteModel, m map[string]RouteModel) {
	if _, ok := m[url]; ok {
		panic(errRouteDuplicateDefinition)
	}
	m[url] = rm
}

func dynamicRoute(url string) {
	parts := strings.Split(url, constant.Slash)
	trieRouteTree.insert(parts[1:], 0)
}

func putSelect(url string, rm RouteModel) {
	switch rm.GetHttpMethod() {
	case HttpMethodGet:
		put(url, rm, getRouteModelMap)
		dynamicRoute(url)
	case HttpMethodPost:
		put(url, rm, postRouteModelMap)
	case HttpMethodPut:
		put(url, rm, putRouteModelMap)
	case HttpMethodDelete:
		put(url, rm, deleteRouteModelMap)
	}
}

func Get(url, method string, r *http.Request) RouteModel {
	switch method {
	case string(*HttpMethodGet):
		parts := strings.Split(url, constant.Slash)
		realUrl := trieRouteTree.search(parts[1:], 0, r)
		if realUrl == "" {
			return nil
		}
		return getRouteModelMap[realUrl]
	case string(*HttpMethodPost):
		return postRouteModelMap[url]
	case string(*HttpMethodPut):
		return putRouteModelMap[url]
	case string(*HttpMethodDelete):
		return deleteRouteModelMap[url]
	default:
		return nil
	}
}

func dealPrefixSlash(url string) string {
	if strings.HasPrefix(url, constant.Slash) {
		return dealPrefixSlash(url[1:])
	} else {
		return constant.Slash + url
	}
}

func dealSuffixSlash(url string) string {
	if strings.HasSuffix(url, constant.Slash) {
		return dealSuffixSlash(url[:len(url)-1])
	} else {
		return url
	}
}

func AddRoute(url string, fun interface{}, hms *httpMethod) {
	url = strings.ReplaceAll(url, constant.Space, constant.NullStr)
	url = dealPrefixSlash(url)
	url = dealSuffixSlash(url)
	if url == constant.NullStr || url == constant.Slash {
		panic(errRouteUrlIsNilOrSlash)
	}
	if fun == nil {
		panic(errRouteFuncIsNil)
	}
	if hms == nil {
		panic(errRouteNoHttpMethod)
	}
	t := reflect.TypeOf(fun)
	checkFun(t)
	rm := &routeModel{hm: hms, t: t, hct: textPlain}
	if t.NumOut() > 0 {
		switch t.Out(0).Kind() {
		case reflect.Bool,
			reflect.String,
			reflect.Float32, reflect.Float64,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			rm.hct = textPlain
		case reflect.Struct, reflect.Pointer, reflect.Interface:
			rm.hct = applicationJson
		default:
			rm.hct = textPlain
		}
	}
	rm.v = reflect.ValueOf(fun)
	if t.NumIn() > 0 {
		pc := rm.v.Pointer()
		funPc := runtime.FuncForPC(pc)
		split := strings.Split(funPc.Name(), constant.Dot)
		funcName := split[len(split)-1]
		fileName, _ := funPc.FileLine(pc)
		var af *ast.File
		if f, ok := astFileMap[fileName]; ok {
			af = f
		} else {
			af, _ = parser.ParseFile(fest, fileName, nil, parser.ParseComments)
			astFileMap[fileName] = af
		}
		var funcDecl *ast.FuncDecl
		ast.Inspect(af, func(n ast.Node) bool {
			if fd, ok := n.(*ast.FuncDecl); ok && fd.Name.Name == funcName {
				funcDecl = fd
				// TODO Parsing interface annotations generates document models
				doc := fd.Doc
				if doc != nil && len(doc.List) > 0 {
					for _, c := range doc.List {
						fmt.Println(c.Text)
					}
				}
				return false
			}
			return true
		})
		var paramNames []string
		for _, param := range funcDecl.Type.Params.List {
			paramNames = append(paramNames, param.Names[0].Name)
		}
		rm.paramNames = paramNames
	}
	putSelect(url, rm)
}
