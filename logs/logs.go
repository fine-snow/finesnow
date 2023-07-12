// Log function

package logs

import (
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

type LogOutput interface {
	INFO(...any)
	INFOF(string, ...any)
	WARN(...any)
	WARNF(string, ...any)
	ERROR(...any)
	ERRORF(string, ...any)
}

func newLogFunc(l *log.Logger) LogFunc {
	return func(v ...any) {
		v = append(v, string(debug.Stack()))
		l.Println(v)
	}
}

func newLogfFunc(l *log.Logger) LogfFunc {
	return func(format string, v ...any) {
		v = append(v, string(debug.Stack()))
		l.Printf(format, v)
	}
}

func newErrorLogFunc(l LogOutput) LogFunc {
	return func(v ...any) {
		v = append(v, string(debug.Stack()))
		l.ERROR(v)
	}
}

func newErrorLogfFunc(l LogOutput) LogfFunc {
	return func(format string, v ...any) {
		v = append(v, string(debug.Stack()))
		l.ERRORF(format, v)
	}
}

func init() {
	INFO = infoLog.Println
	INFOF = infoLog.Printf
	WARN = warnLog.Println
	WARNF = warnLog.Printf
	ERRORF = newLogfFunc(errorLog)
	ERROR = newLogFunc(errorLog)
}

func SetLogOutput(l LogOutput) {
	// Method attribute assignment
	INFO = l.INFO
	INFOF = l.INFOF
	WARN = l.WARN
	WARNF = l.WARNF
	ERROR = newErrorLogFunc(l)
	ERRORF = newErrorLogfFunc(l)

	// Empty to free memory
	infoLog = nil
	warnLog = nil
	errorLog = nil
}
