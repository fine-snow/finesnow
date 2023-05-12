// Routing model collection

package router

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"runtime"
	"strings"
)

var (
	errRouteFuncIsNil           = errors.New("the route mapping add function 'fun' parameter value cannot be nil")
	errRouteAddNotFunc          = errors.New("the route mapping add function 'fun' parameter value is not a function")
	errRouteFuncOutAbnormal     = errors.New("the route mapping add function 'fun' parameter value function return value amount is greater than 1")
	errRouteDuplicateDefinition = errors.New("duplicate definition of routing path")
)

var fest = token.NewFileSet()

var astFileMap = make(map[string]*ast.File)

var routeModelMap = make(map[string]RouteModel)

func checkFun(t reflect.Type) {
	if t.Kind() != reflect.Func {
		panic(errRouteAddNotFunc)
	}
	if t.NumOut() > 1 {
		panic(errRouteFuncOutAbnormal)
	}
}

func put(url string, rm RouteModel) {
	if _, ok := routeModelMap[url]; ok {
		panic(errRouteDuplicateDefinition)
	}
	routeModelMap[url] = rm
}

func Get(url string) RouteModel {
	return routeModelMap[url]
}

func AddRoute(url string, fun interface{}, hms ...*httpMethod) {
	if fun == nil {
		panic(errRouteFuncIsNil)
	}
	t := reflect.TypeOf(fun)
	checkFun(t)
	rm := &routeModel{hm: HttpMethodPost, t: t, hct: textPlain}
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
		split := strings.Split(funPc.Name(), ".")
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
				// TODO 解析接口注释生成文档模型
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
			paramNames = append(paramNames, param.Names[0].Name)
		}
		rm.paramNames = paramNames
	}
	if len(hms) > 0 {
		rm.hm = hms[0]
	}
	put(url, rm)
}
