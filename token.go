//go:generate stringer -type Token

package argh

import "fmt"

const (
	ILLEGAL Token = iota
	EOL
	EMPTY                 // ''
	BS                    // ' ' '\t' '\n'
	IDENT                 // char group without flag prefix: 'some' 'words'
	ARG_DELIMITER         // rune(0)
	ASSIGN                // '='
	MULTI_VALUE_DELIMITER // ','
	LONG_FLAG             // char group with double flag prefix: '--flag'
	SHORT_FLAG            // single char with single flag prefix: '-f'
	COMPOUND_SHORT_FLAG   // char group with single flag prefix: '-flag'
	STDIN_FLAG            // '-'
	STOP_FLAG             // '--'

	nul = rune(0)
	eol = rune(-1)
)

type Token int

// Position is adapted from go/token.Position
type Position struct {
	Column int
}

func (p *Position) IsValid() bool { return p.Column > 0 }

func (p Position) String() string {
	s := ""
	if p.IsValid() {
		s = fmt.Sprintf("%d", p.Column)
	}
	if s == "" {
		s = "-"
	}
	return s
}

// Pos is borrowed from go/token.Pos
type Pos int

const NoPos Pos = 0

func (p Pos) IsValid() bool {
	return p != NoPos
}
