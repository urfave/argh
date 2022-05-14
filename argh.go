package argh

import (
	"log"
	"os"
)

// NOTE: much of this is lifted from
// https://blog.gopheracademy.com/advent-2014/parsers-lexers/

var (
	tracingEnabled = os.Getenv("ARGH_TRACING") == "enabled"
)

type Argh struct {
	ParseTree *ParseTree `json:"parse_tree"`
}

func (a *Argh) AST() []TypedNode {
	return a.ParseTree.toAST()
}

/*
func (a *Argh) String() string {
	return a.ParseTree.String()
}
*/

func tracef(format string, v ...any) {
	if !tracingEnabled {
		return
	}

	log.Printf(format, v...)
}
