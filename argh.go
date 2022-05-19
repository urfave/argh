package argh

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var (
	tracingEnabled = os.Getenv("ARGH_TRACING") == "enabled"
	traceLogger    *log.Logger
)

func init() {
	if !tracingEnabled {
		return
	}

	traceLogger = log.New(os.Stderr, "ARGH TRACING: ", 0)
}

func tracef(format string, v ...any) {
	if !tracingEnabled {
		return
	}

	if _, file, line, ok := runtime.Caller(2); ok {
		format = fmt.Sprintf("%v:%v ", filepath.Base(file), line) + format
	}

	traceLogger.Printf(format, v...)
}
