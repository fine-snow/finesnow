// Framework Launch Method

package snow

import (
	"fmt"
	"github.com/fine-snow/finesnow/handler"
	"github.com/fine-snow/finesnow/logger"
	"net/http"
)

// Run Framework Launch Method
// addr Start address parameter, for example: 127.0.0.1:8088
// intercept Global interceptor parameter, if the interceptor function is not required, this parameter can be passed to nil
func Run(addr string, intercept handler.Interceptor) {
	defer logger.CheckLogChan()
	handle := handler.NewHandle(intercept)
	http.Handle("/", handle)
	fmt.Println("\n    _______           _____                    \n   / ____(_)___  ___ / ___/____  ____ _      __\n  / /_  / / __ \\/ _ \\\\__ \\/ __ \\/ __ \\ | /| / /\n / __/ / / / / /  __/__/ / / / / /_/ / |/ |/ / \n/_/   /_/_/ /_/\\___/____/_/ /_/\\____/|__/|__/  \n                                               ")
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}
