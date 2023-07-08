// Framework Launch Method

package snow

import (
	"github.com/fine-snow/finesnow/constant"
	"github.com/fine-snow/finesnow/handler"
	"github.com/fine-snow/finesnow/logs"
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
	logs.InitLogFunc()
	var err error
	if certFile != constant.NullStr && keyFile != constant.NullStr {
		err = server.ListenAndServeTLS(certFile, keyFile)
	} else {
		err = server.ListenAndServe()
	}
	if err != nil {
		panic(err)
	}
}
