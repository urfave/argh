package argh

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var (
	tracingOn   = os.Getenv("ARGH_TRACING") == "on"
	traceLogger *log.Logger

	cwd *string
)

func init() {
	if !tracingOn {
		return
	}

	traceLogger = log.New(os.Stderr, "ARGH TRACING: ", 0)
}

func tracef(skip int, format string, v ...any) {
	if !tracingOn {
		return
	}

	if cwd == nil {
		if v, err := os.Getwd(); err == nil {
			cwd = &v
		} else {
			v := ""
			cwd = &v
		}
	}

	if _, file, line, ok := runtime.Caller(skip); ok {
		if p, err := filepath.Rel(*cwd, file); err == nil {
			format = fmt.Sprintf("%v:%v ", p, line) + format
		}
	}

	traceLogger.Printf(format, v...)
}
