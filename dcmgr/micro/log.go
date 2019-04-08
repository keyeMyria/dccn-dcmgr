package micro2

import (
"fmt"
"log"
"path"
"runtime"
"time"
)

type logWriter struct {
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(time.Now().UTC().Format("2006-01-02 15:04:05") + string(bytes))
}

func WriteLog(msg string) {

	pathInfo := ""
	if _, file, line, ok := runtime.Caller(1); ok {
		pathInfo = fmt.Sprintf("    [%s:%v]", path.Base(file), line)
	}

	log.SetFlags(0)
	log.SetOutput(new(logWriter))
	log.Println("    " + msg + pathInfo)
}

