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
	handle := handler.NewHandle(intercept)
	server := &http.Server{
		Addr:              addr,
		Handler:           handle,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}
	outputFrameworkInfo()
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
