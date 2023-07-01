package argh

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var (
	isTracingOn = os.Getenv("ARGH_TRACING") == "on"
	traceLogger *log.Logger

	Error = errors.New("argh error")
)

func init() {
	if !isTracingOn {
		return
	}

	traceLogger = log.New(os.Stderr, "## ARGH TRACE ", 0)
}

func tracef(format string, v ...any) {
	if !isTracingOn {
		return
	}

	if _, file, line, ok := runtime.Caller(1); ok {
		format = fmt.Sprintf("%v:%v ", filepath.Base(file), line) + format
	}

	traceLogger.Printf(format, v...)
}
