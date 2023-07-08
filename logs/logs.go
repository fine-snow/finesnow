// Log function

package logs

import (
	"log"
	"os"
)

type LogFunc func(...any)
type LogfFunc func(string, ...any)

var (
	infoLog  = log.New(os.Stdout, "\033[34mINFO\033[0m ", log.LstdFlags|log.Lmsgprefix)
	warnLog  = log.New(os.Stdout, "\033[33mWARN\033[0m ", log.LstdFlags|log.Lmsgprefix)
	errorLog = log.New(os.Stdout, "\033[31mERROR\033[0m ", log.LstdFlags|log.Lmsgprefix)

	INFO   LogFunc
	INFOF  LogfFunc
	WARN   LogFunc
	WARNF  LogfFunc
	ERROR  LogFunc
	ERRORF LogfFunc
)

func init() {
	if INFO == nil {
		INFO = infoLog.Println
	}
	if INFOF == nil {
		INFOF = infoLog.Printf
	}
	if WARN == nil {
		WARN = warnLog.Println
	}
	if WARNF == nil {
		WARNF = warnLog.Printf
	}
	if ERROR == nil {
		ERROR = errorLog.Println
	}
	if ERRORF == nil {
		ERRORF = errorLog.Printf
	}
}
