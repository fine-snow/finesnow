// Provide methods for user operation configuration framework

package snow

import (
	"github.com/fine-snow/finesnow/handler"
	"github.com/fine-snow/finesnow/logs"
	"github.com/fine-snow/finesnow/router"
	"net"
	"net/http"
	"time"
)

/**
 * Log middleware
 * Start
 */

// SetLogOutput Custom log middleware
func SetLogOutput(l logs.LogOutput) {
	logs.CustomLogOutput(l)
}

/**
 * Log middleware
 * End
 */

/**
 * Register Route
 * Start
 */

// Get Add a routing method for GET request
func Get(url string, fun any) {
	router.AddRoute("", url, http.MethodGet, fun)
}

// Post Add a routing method for POST request
func Post(url string, fun any) {
	router.AddRoute("", url, http.MethodPost, fun)
}

// Put Add a routing method for PUT request
func Put(url string, fun any) {
	router.AddRoute("", url, http.MethodPut, fun)
}

// Delete Add a routing method for DELETE request
func Delete(url string, fun any) {
	router.AddRoute("", url, http.MethodDelete, fun)
}

// Group Create a route group
func Group(url string) router.RouteGroup {
	return router.NewGroup(url)
}

/**
 * Register Route
 * End
 */

/**
 * Timeout configuration
 * Start
 */

var (
	readTimeout       = 3 * time.Second
	readHeaderTimeout = 3 * time.Second
	writeTimeout      = 3 * time.Second
	idleTimeout       = time.Minute
)

// SetReadTimeout Configure read timeout time parameter
func SetReadTimeout(t time.Duration) {
	readTimeout = t
}

// SetReadHeaderTimeout Configure read header timeout time parameter
func SetReadHeaderTimeout(t time.Duration) {
	readHeaderTimeout = t
}

// SetWriteTimeout Configure write timeout time parameter
func SetWriteTimeout(t time.Duration) {
	writeTimeout = t
}

// SetIdleTimeout Configure idle timeout time parameter
func SetIdleTimeout(t time.Duration) {
	idleTimeout = t
}

/**
 * Timeout configuration
 * End
 */

/**
 * Server TLS
 * Start
 */

// certFile Certificate file path
// keyFile SecretKeyFile path
var (
	certFile string
	keyFile  string
)

// SetCertFile Configure certificate file path
func SetCertFile(url string) {
	certFile = url
}

// SetKeyFile Configure secretKeyFile path
func SetKeyFile(url string) {
	keyFile = url
}

/**
 * Server TLS
 * End
 */

/**
 * Cross-domain configuration
 * Start
 */

// SetAllowedOrigin Set domains that allow cross domain access
func SetAllowedOrigin(s string) {
	handler.SetAllowedOrigin(s)
}

// SetAllowedMethods Set request methods that allow cross domain requests
func SetAllowedMethods(s []string) {
	handler.SetAllowedMethods(s)
}

// SetAllowedHeaders Set request headers that allow cross domain connections
func SetAllowedHeaders(s []string) {
	handler.SetAllowedHeaders(s)
}

/**
 * Cross-domain configuration
 * End
 */

/**
 * Global exception handling function
 * Start
 */

// SetGlobalErrHandle Set the global exception handling function
func SetGlobalErrHandle(fun handler.ErrHandleFunc) {
	handler.SetErrHandleFunc(fun)
}

/**
 * Global Exception Handling Function
 * End
 */

/**
 * Interceptor function
 * Start
 */

// SetInterceptor Set the global interceptor function
func SetInterceptor(intercept handler.Interceptor) {
	handler.SetIntercept(intercept)
}

/**
 * Interceptor function
 * End
 */

/**
 * Postprocessor function
 * Start
 */

// SetPostProcessor Set the global post processor function
func SetPostProcessor(postProcess handler.PostProcessor) {
	handler.SetPostProcess(postProcess)
}

/**
 * Postprocessor function
 * End
 */

// Run Framework Launch Method
// addr Start address parameter, for example: 127.0.0.1:8088.
// intercept Global interceptor parameter, if the interceptor function is not required, this parameter can be passed to nil.
func Run(addr string) {

	// Capture startup exceptions
	defer handler.CatchRunPanic()

	// Output framework logo version and other information
	outputFrameworkInfo()

	// addr check
	if addr == "" {
		addr = ":9801"
		logs.WARN("You are using the default port 9801 to start the framework")
	}

	// Receive abnormal parameter declaration
	var err error

	// Attempt to listen to the specified address and port
	ln, err := net.Listen("tcp", addr)

	// Abnormal information judgment
	if err != nil {
		panic(err)
	}

	// Generate HTTP processing function parameters
	handle := handler.NewHandle()

	// Configure service startup parameters
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

	// Startup method judgment
	if certFile != "" && keyFile != "" {
		err = server.ServeTLS(ln, certFile, keyFile)
	} else {
		err = server.Serve(ln)
	}

	// Abnormal information judgment
	if err != nil {
		panic(err)
	}
}
