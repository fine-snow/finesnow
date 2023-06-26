// Global Exception Handling Function Configuration

package snow

import "github.com/fine-snow/finesnow/handler"

// SetGlobalErrHandle Set global exception handling functions
func SetGlobalErrHandle(fun handler.ErrHandleFunc) {
	handler.SetGlobalErrHandleFunc(fun)
}
