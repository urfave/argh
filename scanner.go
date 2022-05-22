package argh

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"sort"
	"unicode"
)

const (
	nul = rune(0)
	eol = rune(-1)
)

var (
	DefaultScannerConfig = &ScannerConfig{
		AssignmentOperator: '=',
		FlagPrefix:         '-',
		MultiValueDelim:    ',',
	}
)

type Scanner struct {
	r   *bufio.Reader
	i   int
	cfg *ScannerConfig
}

type ScannerConfig struct {
	AssignmentOperator rune
	FlagPrefix         rune
	MultiValueDelim    rune
}

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

func NewScanner(r io.Reader, cfg *ScannerConfig) *Scanner {
	if cfg == nil {
		cfg = DefaultScannerConfig
	}

	return &Scanner{
		r:   bufio.NewReader(r),
		cfg: cfg,
	}
}

func (s *Scanner) Scan() (Token, string, Pos) {
	ch, pos := s.read()

	if s.isBlankspace(ch) {
		_ = s.unread()
		return s.scanBlankspace()
	}

	if s.isAssignmentOperator(ch) {
		return ASSIGN, string(ch), pos
	}

	if s.isMultiValueDelim(ch) {
		return MULTI_VALUE_DELIMITER, string(ch), pos
	}

	if ch == eol {
		return EOL, "", pos
	}

	if ch == nul {
		return ARG_DELIMITER, string(ch), pos
	}

	if unicode.IsGraphic(ch) {
		_ = s.unread()
		return s.scanArg()
	}

	return ILLEGAL, string(ch), pos
}

func (s *Scanner) read() (rune, Pos) {
	ch, _, err := s.r.ReadRune()
	s.i++

	if errors.Is(err, io.EOF) {
		return eol, Pos(s.i)
	} else if err != nil {
		log.Printf("unknown scanner error=%+v", err)
		return eol, Pos(s.i)
	}

	return ch, Pos(s.i)
}

func (s *Scanner) unread() Pos {
	_ = s.r.UnreadRune()
	s.i--
	return Pos(s.i)
}

func (s *Scanner) isBlankspace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func (s *Scanner) isUnderscore(ch rune) bool {
	return ch == '_'
}

func (s *Scanner) isFlagPrefix(ch rune) bool {
	return ch == s.cfg.FlagPrefix
}

func (s *Scanner) isMultiValueDelim(ch rune) bool {
	return ch == s.cfg.MultiValueDelim
}

func (s *Scanner) isAssignmentOperator(ch rune) bool {
	return ch == s.cfg.AssignmentOperator
}

func (s *Scanner) scanBlankspace() (Token, string, Pos) {
	buf := &bytes.Buffer{}
	ch, pos := s.read()
	buf.WriteRune(ch)

	for {
		ch, pos = s.read()

		if ch == eol {
			break
		} else if !s.isBlankspace(ch) {
			pos = s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return BS, buf.String(), pos
}

func (s *Scanner) scanArg() (Token, string, Pos) {
	buf := &bytes.Buffer{}
	ch, pos := s.read()
	buf.WriteRune(ch)

	for {
		ch, pos = s.read()

		if ch == eol || ch == nul || s.isAssignmentOperator(ch) || s.isMultiValueDelim(ch) {
			pos = s.unread()
			break
		}

		_, _ = buf.WriteRune(ch)
	}

	str := buf.String()

	if len(str) == 0 {
		return EMPTY, str, pos
	}

	ch0 := rune(str[0])

	if len(str) == 1 {
		if s.isFlagPrefix(ch0) {
			return STDIN_FLAG, str, pos
		}

		if s.isAssignmentOperator(ch0) {
			return ASSIGN, str, pos
		}

		return IDENT, str, pos
	}

	ch1 := rune(str[1])

	if len(str) == 2 {
		if str == string(s.cfg.FlagPrefix)+string(s.cfg.FlagPrefix) {
			return STOP_FLAG, str, pos
		}

		if s.isFlagPrefix(ch0) {
			return SHORT_FLAG, str, pos
		}
	}

	if s.isFlagPrefix(ch0) {
		if s.isFlagPrefix(ch1) {
			return LONG_FLAG, str, pos
		}

		return COMPOUND_SHORT_FLAG, str, pos
	}

	return IDENT, str, pos
}
