// Framework Launch Method

package snow

import (
	"github.com/fine-snow/finesnow/constant"
	"github.com/fine-snow/finesnow/handler"
	"github.com/fine-snow/finesnow/router"
	"net"
	"net/http"
)

// Run Framework Launch Method
// addr Start address parameter, for example: 127.0.0.1:8088
// intercept Global interceptor parameter, if the interceptor function is not required, this parameter can be passed to nil
func Run(addr string, intercept handler.Interceptor) {
	// Capture startup exceptions
	defer handler.CatchRunPanic()
	// Output framework logo version and other information
	outputFrameworkInfo()

	var err error
	// Attempt to listen to the specified address and port
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	handle := handler.NewHandle(intercept)
	server := &http.Server{
		Addr:              addr,
		Handler:           handle,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}
	// Process registered routes
	router.DealRoute()
	defer func(ln net.Listener) {
		err = ln.Close()
		if err != nil {
			panic(err)
		}
	}(ln)
	if certFile != constant.NullStr && keyFile != constant.NullStr {
		err = server.ServeTLS(ln, certFile, keyFile)
	} else {
		err = server.Serve(ln)
	}
	if err != nil {
		panic(err)
	}
}
