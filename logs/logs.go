// Log function

package logs

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
)

type LogFunc func(...any)

type LogfFunc func(string, ...any)

var (
	infoLog  = log.New(os.Stdout, "\033[34mINFO\033[0m ", log.LstdFlags|log.Lmsgprefix|log.Lshortfile)
	warnLog  = log.New(os.Stdout, "\033[33mWARN\033[0m ", log.LstdFlags|log.Lmsgprefix|log.Lshortfile)
	errorLog = log.New(os.Stdout, "\033[31mERROR\033[0m ", log.LstdFlags|log.Lmsgprefix)

	INFO   LogFunc
	INFOF  LogfFunc
	WARN   LogFunc
	WARNF  LogfFunc
	ERROR  LogFunc
	ERRORF LogfFunc
)

func newLogFunc(l *log.Logger) LogFunc {
	return func(v ...any) {
		v = append(v, string(debug.Stack()))
		_ = l.Output(3, fmt.Sprintln(v...))
	}
}

func newLogfFunc(l *log.Logger) LogfFunc {
	return func(format string, v ...any) {
		v = append(v, string(debug.Stack()))
		_ = l.Output(3, fmt.Sprintln(v...))
	}
}

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
		ERROR = newLogFunc(errorLog)
	}
	if ERRORF == nil {
		ERRORF = newLogfFunc(errorLog)
	}
}
