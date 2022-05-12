package argh

import (
	"io"
	"strings"

	"github.com/pkg/errors"
)

// NOTE: much of this is lifted from
// https://blog.gopheracademy.com/advent-2014/parsers-lexers/

var (
	errSyntax = errors.New("syntax error")
)

func ParseArgs(args []string) (*Argh, error) {
	reEncoded := strings.Join(args, string(nul))

	return NewParser(
		strings.NewReader(reEncoded),
		nil,
	).Parse()
}

type Parser struct {
	s   *Scanner
	buf ParserBuffer
}

type ParserBuffer struct {
	tok Token
	lit string
	n   int
}

func NewParser(r io.Reader, cfg *ScannerConfig) *Parser {
	return &Parser{s: NewScanner(r, cfg)}
}

func (p *Parser) Parse() (*Argh, error) {
	arghOut := &Argh{
		AST: &AST{
			Nodes: []*Node{},
		},
	}

	for {
		tok, lit := p.scan()
		if tok == ILLEGAL {
			return nil, errors.Wrapf(errSyntax, "illegal value %q", lit)
		}

		if tok == EOL {
			break
		}

		arghOut.AST.Nodes = append(
			arghOut.AST.Nodes,
			&Node{Token: tok.String(), Literal: lit},
		)
	}

	return arghOut, nil
}

func (p *Parser) scan() (Token, string) {
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	tok, lit := p.s.Scan()

	p.buf.tok, p.buf.lit = tok, lit

	return tok, lit
}

func (p *Parser) unscan() {
	p.buf.n = 1
}
