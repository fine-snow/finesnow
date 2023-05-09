// Framework Launch Method

package snow

import (
	"github.com/fine-snow/finesnow/handler"
	"net/http"
)

// Run Framework Launch Method
// addr Start address parameter, for example: 127.0.0.1:8088
// intercept Global interceptor parameter, if the interceptor function is not required, this parameter can be passed to nil
func Run(addr string, intercept handler.Interceptor) {
	handle := handler.NewHandle()
	handle.SetIntercept(intercept)
	http.Handle("/", handle)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}
