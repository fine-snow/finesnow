// Log bus

package logger

import (
	"io"
	"log"
	"os"
	"time"
)

// Log bus parameter declaration
var (
	logChan   = make(chan *snowLog, 300)
	infoFile  *os.File
	warnFile  *os.File
	errorFile *os.File
	infoLog   *log.Logger
	warnLog   *log.Logger
	errorLog  *log.Logger
)

// Log structure model
type snowLog struct {
	t int
	v string
}

// Log Service Bus
func logServer() {
	for {
		l := <-logChan
		switch l.t {
		case INFO_:
			infoLog.Println(l.v)
		case WARN_:
			warnLog.Println(l.v)
		case ERROR_:
			errorLog.Println(l.v)
		default:
			infoLog.Println(l.v)
		}
	}
}

func init() {
	go logServer()
	_, err := os.Stat("./log")
	if err != nil {
		os.Mkdir("log", 0777)
	}
	// Initialize Log Service parameters
	infoFile, _ = os.OpenFile("./log/snow-info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	warnFile, _ = os.OpenFile("./log/snow-warn.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	errorFile, _ = os.OpenFile("./log/snow-error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	infoLog = log.New(io.MultiWriter(infoFile, os.Stdout), "[INFO] ", 3)
	warnLog = log.New(io.MultiWriter(warnFile, os.Stdout), "[WARN] ", 3)
	errorLog = log.New(io.MultiWriter(errorFile, os.Stdout), "[ERROR] ", 3)
}

// How to use exposed logs

func INFO(val string) {
	logChan <- &snowLog{t: INFO_, v: val}
}

func WARN(val string) {
	logChan <- &snowLog{t: WARN_, v: val}
}

func ERROR(val string) {
	logChan <- &snowLog{t: ERROR_, v: val}
}

// CheckLogChan Check whether all messages in the log pipeline have been consumed
func CheckLogChan() {
	// Sleep for 1 second to ensure that all logs are in the pipeline
	time.Sleep(time.Second)
	for {
		// All messages in the log pipeline are consumed and the infinite loop ends
		if len(logChan) == 0 {
			break
		}
	}
}
