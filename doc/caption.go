// Api Doc

package doc

import (
	"net/http"
	"runtime"
	"strings"
)

// enableApiDoc Enable the interface document feature identification
var enableApiDoc = false

var IsApiJson *ApiJson

func GetEnableApiDoc() bool {
	return enableApiDoc
}

func EnableApiDoc() {
	enableApiDoc = true
}

func GetApiJson() *ApiJson {
	return IsApiJson
}

func HandleDoc() {
	_, f, _, _ := runtime.Caller(0)
	url := strings.Replace(f, "doc/caption.go", "", -1)
	http.Handle("/public/", http.FileServer(http.Dir(url)))
}
