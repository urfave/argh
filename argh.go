package argh

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
)

var (
	tracing             = strings.Split(os.Getenv("URFAVE_ARGH_TRACING"), ",")
	isTracingOn         = sliceContains(tracing, "on")
	isTracingWithCaller = !sliceContains(tracing, "caller=off")

	Err = errors.New("urfave/argh error")
)

func tracef(format string, a ...any) {
	if !isTracingOn {
		return
	}

	if !strings.HasSuffix(format, "\n") {
		format = format + "\n"
	}

	if isTracingWithCaller {
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

		return
	}

	fmt.Fprintf(
		os.Stderr,
		strings.Join([]string{
			"## URFAVE ARGH TRACE ",
			format,
		}, ""),
		a...,
	)
}

func sliceContains[T comparable](sl []T, k T) bool {
	for _, item := range sl {
		if item == k {
			return true
		}
	}

	return false
}

func FirstValue(sl []KeyValue, key string) (string, bool) {
	for _, item := range sl {
		if item.Key == key {
			return item.Value, true
		}
	}

	return "", false
}

/* NOTE: if needed:
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
*/
