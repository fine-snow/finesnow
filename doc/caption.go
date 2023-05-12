// Api Doc

package doc

import (
	"net/http"
	"os"
	"runtime"
	"strings"
)

var enableApiDoc = false

func GetEnableApiDoc() bool {
	return enableApiDoc
}

func EnableApiDoc() {
	enableApiDoc = true
}

const name = "index.html"

func ApiDoc(w http.ResponseWriter, r *http.Request) {
	_, f, _, _ := runtime.Caller(0)
	file, err := os.Open(strings.Replace(f, "doc/caption.go", name, -1))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	stat, err := file.Stat()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	http.ServeContent(w, r, name, stat.ModTime(), file)
}
