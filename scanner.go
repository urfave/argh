package argh

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
	i   int
	cfg *ScannerConfig
}

type ScannerConfig struct {
	AssignmentOperator rune
	FlagPrefix         rune
	MultiValueDelim    rune
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

func (s *Scanner) Scan() (Token, string, int) {
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

func (s *Scanner) read() (rune, int) {
	ch, _, err := s.r.ReadRune()
	s.i++

	if errors.Is(err, io.EOF) {
		return eol, s.i
	} else if err != nil {
		log.Printf("unknown scanner error=%+v", err)
		return eol, s.i
	}

	return ch, s.i
}

func (s *Scanner) unread() int {
	_ = s.r.UnreadRune()
	s.i--
	return s.i
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

func (s *Scanner) scanBlankspace() (Token, string, int) {
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

func (s *Scanner) scanArg() (Token, string, int) {
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
