// Route model collection

package router

import (
	"errors"
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
	errRouteUrlIsNilOrSlash     = errors.New("the route url cannot be null character or '/'")

	fest       = token.NewFileSet()
	astFileMap = make(map[string]*ast.File)

	getRouteModelMap    = make(map[string]RouteModel)
	postRouteModelMap   = make(map[string]RouteModel)
	putRouteModelMap    = make(map[string]RouteModel)
	deleteRouteModelMap = make(map[string]RouteModel)

	allRouteSlice = make([]*traRouteModel, constant.Zero)
)

// checkFun Verify the validity of the fun parameter
func checkFun(t reflect.Type) {
	if t.Kind() != reflect.Func {
		panic(errRouteAddNotFunc)
	}
	if t.NumOut() > constant.One {
		panic(errRouteFuncOutAbnormal)
	}
}

// put Add the routing model to the collection
func put(url string, rm RouteModel, m map[string]RouteModel) {
	if _, ok := m[url]; ok {
		panic(errRouteDuplicateDefinition)
	}
	m[url] = rm
}

// dynamicRoute Dynamic route handling methods
func dynamicRoute(url string) {
	parts := strings.Split(url, constant.Slash)
	prefixRouteTree.insert(parts[constant.One:], constant.Zero)
}

// putSelect Select the corresponding routing model collection to add it
func putSelect(url, method string, rm RouteModel) {
	switch method {
	case http.MethodGet:
		put(url, rm, getRouteModelMap)
		dynamicRoute(url)
	case http.MethodPost:
		put(url, rm, postRouteModelMap)
	case http.MethodPut:
		put(url, rm, putRouteModelMap)
	case http.MethodDelete:
		put(url, rm, deleteRouteModelMap)
	}
}

// Get Based on the URL, method finds the corresponding routing model
func Get(url, method string, r *http.Request) RouteModel {
	switch method {
	case http.MethodGet:
		parts := strings.Split(url, constant.Slash)
		realUrl := prefixRouteTree.search(parts[constant.One:], constant.Zero)
		if realUrl == constant.NullStr {
			return nil
		}
		realUrlParts := strings.Split(realUrl, constant.Slash)
		for i, part := range realUrlParts[constant.One:] {
			if part[constant.Zero] == constant.Colon {
				if r.URL.RawQuery == constant.NullStr {
					r.URL.RawQuery = part[constant.One:] + constant.EqualSign + parts[i+constant.One]
				} else {
					r.URL.RawQuery = r.URL.RawQuery + constant.Ampersand + part[constant.One:] + constant.EqualSign + parts[i+constant.One]
				}
			}
		}
		return getRouteModelMap[realUrl]
	case http.MethodPost:
		return postRouteModelMap[url]
	case http.MethodPut:
		return putRouteModelMap[url]
	case http.MethodDelete:
		return deleteRouteModelMap[url]
	default:
		return nil
	}
}

func dealPrefixSlash(url string) string {
	if strings.HasPrefix(url, constant.Slash) {
		return dealPrefixSlash(url[constant.One:])
	} else {
		return constant.Slash + url
	}
}

func dealSuffixSlash(url string) string {
	if strings.HasSuffix(url, constant.Slash) {
		return dealSuffixSlash(url[:len(url)-constant.One])
	} else {
		return url
	}
}

func checkUrl(url string) string {
	url = strings.ReplaceAll(url, constant.Space, constant.NullStr)
	url = dealPrefixSlash(url)
	url = dealSuffixSlash(url)
	if url == constant.NullStr || url == constant.Slash {
		panic(errRouteUrlIsNilOrSlash)
	}
	return url
}

// dealRoute Processing Transition Routes to Generate Real Route Models
func dealRoute(group, url, method string, fun any) {
	url = checkUrl(url)
	if fun == nil {
		panic(errRouteFuncIsNil)
	}
	t := reflect.TypeOf(fun)
	checkFun(t)
	rm := &routeModel{t: t, hct: textPlain}
	if t.NumOut() > constant.Zero {
		switch t.Out(constant.Zero).Kind() {
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
	if t.NumIn() > constant.Zero {
		pc := rm.v.Pointer()
		funPc := runtime.FuncForPC(pc)
		split := strings.Split(funPc.Name(), constant.Dot)
		funcName := split[len(split)-constant.One]
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
				//doc := fd.Doc
				//if doc != nil && len(doc.List) > 0 {
				//	for _, c := range doc.List {
				//		fmt.Println(c.Text)
				//	}
				//}
				return false
			}
			return true
		})
		var paramNames []string
		for _, param := range funcDecl.Type.Params.List {
			for _, name := range param.Names {
				paramNames = append(paramNames, name.Name)
			}
		}
		rm.paramNames = paramNames
	}
	putSelect(group+url, method, rm)
}

// AddRoute Add Transition Route
func AddRoute(group, url, method string, fun any) {
	allRouteSlice = append(allRouteSlice, &traRouteModel{group: group, url: url, method: method, fun: fun})
}

// DealRoute All transition routing methods for handling external exposure
func DealRoute() {
	for _, v := range allRouteSlice {
		dealRoute(v.group, v.url, v.method, v.fun)
	}
	// Empty and wait for GC to reclaim memory
	allRouteSlice = nil
}
