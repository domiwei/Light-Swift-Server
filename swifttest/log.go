package swifttest

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type ServerLogger struct {
	filename string
	out *os.File
	logger *log.Logger
}

// Create logger
func OpenLogger(filename string) *ServerLogger {
	logFile, err := os.Create(filename)
	if err != nil {
		return nil
	}

	return &ServerLogger{
		filename: filename,
		out:      logFile,
		logger:   log.New(logFile,"[Debug]",log.LstdFlags),
	}
}

func (serverlogger *ServerLogger) Print(prefix string, format string, v ...interface{}) {
	serverlogger.logger.SetPrefix(prefix)
	serverlogger.logger.Printf(format, v)
}

// Close logger
func (serverlogger *ServerLogger) CloseLogger() {
	serverlogger.out.Close()
}

func fatalf(code int, codeStr string, errf string, a ...interface{}) {
	if DEBUG {
		log.Printf("statusCode %q Code %s Message %s ", code, codeStr, fmt.Sprintf(errf, a...))
	}
	panic(&swiftError{
		statusCode: code,
		Code:       codeStr,
		Message:    fmt.Sprintf(errf, a...),
	})
}

func jsonMarshal(w io.Writer, x interface{}) {
	if err := json.NewEncoder(w).Encode(x); err != nil {
		panic(fmt.Errorf("error marshalling %#v: %v", x, err))
	}
}
