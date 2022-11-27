package argh

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
	"unicode"
)

type Scanner struct {
	r   *bufio.Reader
	i   int
	cfg *ScannerConfig
}

func NewScanner(r io.Reader, cfg *ScannerConfig) *Scanner {
	if cfg == nil {
		cfg = POSIXyScannerConfig
	}

	return &Scanner{
		r:   bufio.NewReader(r),
		cfg: cfg,
	}
}

func (s *Scanner) Scan() (Token, string, Pos) {
	ch, pos := s.read()

	if s.cfg.IsBlankspace(ch) {
		_ = s.unread()
		return s.scanBlankspace()
	}

	if s.cfg.IsAssignmentOperator(ch) {
		return ASSIGN, string(ch), pos
	}

	if s.cfg.IsMultiValueDelim(ch) {
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

	return INVALID, string(ch), pos
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

func (s *Scanner) scanBlankspace() (Token, string, Pos) {
	buf := &bytes.Buffer{}
	ch, pos := s.read()
	buf.WriteRune(ch)

	for {
		ch, pos = s.read()

		if ch == eol {
			break
		} else if !s.cfg.IsBlankspace(ch) {
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

		if ch == eol || ch == nul || s.cfg.IsAssignmentOperator(ch) || s.cfg.IsMultiValueDelim(ch) {
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
		if s.cfg.IsFlagPrefix(ch0) {
			return STDIN_FLAG, str, pos
		}

		if s.cfg.IsAssignmentOperator(ch0) {
			return ASSIGN, str, pos
		}

		return IDENT, str, pos
	}

	ch1 := rune(str[1])

	if len(str) == 2 {
		if s.cfg.IsFlagPrefix(ch0) && s.cfg.IsFlagPrefix(ch1) {
			return STOP_FLAG, str, pos
		}

		if s.cfg.IsFlagPrefix(ch0) {
			return SHORT_FLAG, str, pos
		}
	}

	if s.cfg.IsFlagPrefix(ch0) {
		if s.cfg.IsFlagPrefix(ch1) {
			return LONG_FLAG, str, pos
		}

		return COMPOUND_SHORT_FLAG, str, pos
	}

	return IDENT, str, pos
}
