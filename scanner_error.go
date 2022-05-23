package argh

import (
	"fmt"
	"io"
	"sort"
)

// ScannerError is largely borrowed from go/scanner.Error
type ScannerError struct {
	Pos Position
	Msg string
}

func (e ScannerError) Error() string {
	if e.Pos.IsValid() {
		return e.Pos.String() + ":" + e.Msg
	}
	return e.Msg
}

// ScannerErrorList is largely borrowed from go/scanner.ErrorList
type ScannerErrorList []*ScannerError

func (el *ScannerErrorList) Add(pos Position, msg string) {
	*el = append(*el, &ScannerError{Pos: pos, Msg: msg})
}

func (el *ScannerErrorList) Reset() { *el = (*el)[0:0] }

func (el ScannerErrorList) Len() int { return len(el) }

func (el ScannerErrorList) Swap(i, j int) { el[i], el[j] = el[j], el[i] }

func (el ScannerErrorList) Less(i, j int) bool {
	e := &el[i].Pos
	f := &el[j].Pos

	if e.Column != f.Column {
		return e.Column < f.Column
	}

	return el[i].Msg < el[j].Msg
}

func (el ScannerErrorList) Sort() {
	sort.Sort(el)
}

func (el ScannerErrorList) Error() string {
	switch len(el) {
	case 0:
		return "no errors"
	case 1:
		return el[0].Error()
	}
	return fmt.Sprintf("%s (and %d more errors)", el[0], len(el)-1)
}

func (el ScannerErrorList) Err() error {
	if len(el) == 0 {
		return nil
	}
	return el
}

func PrintScannerError(w io.Writer, err error) {
	if list, ok := err.(ScannerErrorList); ok {
		for _, e := range list {
			fmt.Fprintf(w, "%s\n", e)
		}
	} else if err != nil {
		fmt.Fprintf(w, "%s\n", err)
	}
}
