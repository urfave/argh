package argh

import (
	"fmt"
	"io"
	"sort"
)

// ParserError is largely borrowed from go/scanner.Error
type ParserError struct {
	Pos Position
	Msg string
}

func (e ParserError) Error() string {
	if e.Pos.IsValid() {
		return e.Pos.String() + ":" + e.Msg
	}

	return e.Msg
}

// ParserErrorList is largely borrowed from go/scanner.ErrorList
type ParserErrorList []*ParserError

func (el *ParserErrorList) Add(pos Position, msg string) {
	*el = append(*el, &ParserError{Pos: pos, Msg: msg})
}

func (el *ParserErrorList) Reset() { *el = (*el)[0:0] }

func (el ParserErrorList) Len() int { return len(el) }

func (el ParserErrorList) Swap(i, j int) { el[i], el[j] = el[j], el[i] }

func (el ParserErrorList) Less(i, j int) bool {
	e := &el[i].Pos
	f := &el[j].Pos

	if e.Column != f.Column {
		return e.Column < f.Column
	}

	return el[i].Msg < el[j].Msg
}

func (el ParserErrorList) Sort() {
	sort.Sort(el)
}

func (el ParserErrorList) Error() string {
	switch len(el) {
	case 0:
		return "no errors"
	case 1:
		return el[0].Error()
	}
	return fmt.Sprintf("%s (and %d more errors)", el[0], len(el)-1)
}

func (el ParserErrorList) Err() error {
	if len(el) == 0 {
		return nil
	}
	return el
}

func (el ParserErrorList) Is(other error) bool {
	if _, ok := other.(ParserErrorList); ok {
		return el.Error() == other.Error()
	}

	if v, ok := other.(*ParserErrorList); ok {
		return el.Error() == (*v).Error()
	}

	return false
}

func PrintParserError(w io.Writer, err error) {
	if list, ok := err.(ParserErrorList); ok {
		for _, e := range list {
			fmt.Fprintf(w, "%s\n", e)
		}
	} else if err != nil {
		fmt.Fprintf(w, "%s\n", err)
	}
}
