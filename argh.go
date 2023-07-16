package argh

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
)

var (
	isTracingOn = os.Getenv("URFAVE_ARGH_TRACING") == "on"

	Err = errors.New("urfave/argh error")
)

func tracef(format string, a ...any) {
	if !isTracingOn {
		return
	}

	if !strings.HasSuffix(format, "\n") {
		format = format + "\n"
	}

	pc, file, line, _ := runtime.Caller(1)
	cf := runtime.FuncForPC(pc)

	fmt.Fprintf(
		os.Stderr,
		strings.Join([]string{
			"## URFAVE ARGH TRACE ",
			file,
			":",
			fmt.Sprintf("%v", line),
			" ",
			fmt.Sprintf("(%s)", cf.Name()),
			" ",
			format,
		}, ""),
		a...,
	)
}

func FirstValue(sl []KeyValue, key string) (string, bool) {
	for _, item := range sl {
		if item.Key == key {
			return item.Value, true
		}
	}

	return "", false
}

func LastValue(sl []KeyValue, key string) (string, bool) {
	v := ""
	ok := false

	for _, item := range sl {
		if item.Key == key {
			v = item.Value
			ok = true
		}
	}

	return v, ok
}

func AllValues(sl []KeyValue, key string) []string {
	v := []string{}

	for _, item := range sl {
		if item.Key == key {
			v = append(v, item.Value)
		}
	}

	return v
}
