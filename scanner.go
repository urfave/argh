package argh

// NOTE: much of this is lifted from
// https://blog.gopheracademy.com/advent-2014/parsers-lexers/

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
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
	cfg *ScannerConfig
}

type ScannerConfig struct {
	AssignmentOperator rune
	FlagPrefix         rune
	MultiValueDelim    rune

	Commands []string
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

func (s *Scanner) Scan() (Token, string) {
	ch := s.read()

	if s.isBlankspace(ch) {
		s.unread()
		return s.scanBlankspace()
	}

	if s.isAssignmentOperator(ch) {
		return ASSIGN, string(ch)
	}

	if s.isMultiValueDelim(ch) {
		return MULTI_VALUE_DELIMITER, string(ch)
	}

	if ch == eol {
		return EOL, ""
	}

	if ch == nul {
		return ARG_DELIMITER, string(ch)
	}

	if unicode.IsGraphic(ch) {
		s.unread()
		return s.scanArg()
	}

	return ILLEGAL, string(ch)
}

func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if errors.Is(err, io.EOF) {
		return eol
	} else if err != nil {
		log.Printf("unknown scanner error=%+v", err)
		return eol
	}

	return ch
}

func (s *Scanner) unread() {
	_ = s.r.UnreadRune()
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

func (s *Scanner) scanBlankspace() (Token, string) {
	buf := &bytes.Buffer{}
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eol {
			break
		} else if !s.isBlankspace(ch) {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return BS, buf.String()
}

func (s *Scanner) scanArg() (Token, string) {
	buf := &bytes.Buffer{}
	buf.WriteRune(s.read())

	for {
		ch := s.read()

		if ch == eol || ch == nul || s.isAssignmentOperator(ch) || s.isMultiValueDelim(ch) {
			s.unread()
			break
		}

		_, _ = buf.WriteRune(ch)
	}

	str := buf.String()

	if len(str) == 0 {
		return EMPTY, str
	}

	ch0 := rune(str[0])

	if len(str) == 1 {
		if s.isFlagPrefix(ch0) {
			return STDIN_FLAG, str
		}

		return IDENT, str
	}

	ch1 := rune(str[1])

	if len(str) == 2 {
		if str == string(s.cfg.FlagPrefix)+string(s.cfg.FlagPrefix) {
			return STOP_FLAG, str
		}

		if s.isFlagPrefix(ch0) {
			return SHORT_FLAG, str
		}
	}

	if s.isFlagPrefix(ch0) {
		if s.isFlagPrefix(ch1) {
			return LONG_FLAG, str
		}

		return COMPOUND_SHORT_FLAG, str
	}

	return IDENT, str
}
