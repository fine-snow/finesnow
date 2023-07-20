// Log function

package logs

import (
	"log"
	"os"
	"runtime/debug"
)

// LogFunc Log output constraint method
type LogFunc func(...any)

// LogfFunc Log output constraint method
type LogfFunc func(string, ...any)

var (
	outLog   = log.New(os.Stdout, "", 0)
	debugLog = log.New(os.Stdout, "\033[32mDEBUG\033[0m ", log.LstdFlags|log.Lmsgprefix|log.Lshortfile)
	infoLog  = log.New(os.Stdout, "\033[34mINFO\033[0m ", log.LstdFlags|log.Lmsgprefix|log.Lshortfile)
	warnLog  = log.New(os.Stdout, "\033[33mWARN\033[0m ", log.LstdFlags|log.Lmsgprefix|log.Lshortfile)
	errorLog = log.New(os.Stdout, "\033[31mERROR\033[0m ", log.LstdFlags|log.Lmsgprefix)

	OUT    LogFunc
	DEBUG  LogFunc
	DEBUGF LogfFunc
	INFO   LogFunc
	INFOF  LogfFunc
	WARN   LogFunc
	WARNF  LogfFunc
	ERROR  LogFunc
	ERRORF LogfFunc
)

// LogOutput Custom logging middleware constraint interface
type LogOutput interface {
	OUT(...any)
	DEBUG(...any)
	DEBUGF(string, ...any)
	INFO(...any)
	INFOF(string, ...any)
	WARN(...any)
	WARNF(string, ...any)
	ERROR(...any)
	ERRORF(string, ...any)
}

func defaultNewLogFunc(l *log.Logger) LogFunc {
	return func(v ...any) {
		v = append(v, string(debug.Stack()))
		l.Println(v...)
	}
}

func defaultNewLogfFunc(l *log.Logger) LogfFunc {
	return func(format string, v ...any) {
		v = append(v, string(debug.Stack()))
		l.Printf(format+" %s", v...)
	}
}

// init Initialization parameters
func init() {
	if OUT == nil {
		OUT = outLog.Println
	}
	if DEBUG == nil {
		DEBUG = debugLog.Println
	}
	if DEBUGF == nil {
		DEBUGF = debugLog.Printf
	}
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
		ERROR = defaultNewLogFunc(errorLog)
	}
	if ERRORF == nil {
		ERRORF = defaultNewLogfFunc(errorLog)
	}
}

// CustomLogOutput Custom log middleware
func CustomLogOutput(l LogOutput) {
	// Method attribute assignment
	OUT = l.OUT
	DEBUG = l.DEBUG
	DEBUGF = l.DEBUGF
	INFO = l.INFO
	INFOF = l.INFOF
	WARN = l.WARN
	WARNF = l.WARNF
	ERROR = l.ERROR
	ERRORF = l.ERRORF

	// Empty to free memory
	outLog = nil
	debugLog = nil
	infoLog = nil
	warnLog = nil
	errorLog = nil
}
