package argh

import (
	"io"
	"strings"

	"github.com/pkg/errors"
)

var (
	errSyntax = errors.New("syntax error")

	DefaultParserConfig = &ParserConfig{
		Commands:      []string{},
		ValueFlags:    []string{},
		ScannerConfig: DefaultScannerConfig,
	}
)

func ParseArgs(args []string, pCfg *ParserConfig) (*Argh, error) {
	reEncoded := strings.Join(args, string(nul))

	return NewParser(
		strings.NewReader(reEncoded),
		pCfg,
	).Parse()
}

type Parser struct {
	s   *Scanner
	buf ParserBuffer

	commands   map[string]struct{}
	valueFlags map[string]struct{}

	nodes    []Node
	stopSeen bool
}

type ParserBuffer struct {
	tok Token
	lit string
	pos int
	n   int
}

type ParserConfig struct {
	Commands      []string
	ValueFlags    []string
	ScannerConfig *ScannerConfig
}

type parseDirective struct {
	Break bool
}

func NewParser(r io.Reader, pCfg *ParserConfig) *Parser {
	if pCfg == nil {
		pCfg = DefaultParserConfig
	}

	parser := &Parser{
		s:          NewScanner(r, pCfg.ScannerConfig),
		commands:   map[string]struct{}{},
		valueFlags: map[string]struct{}{},
	}

	for _, command := range pCfg.Commands {
		parser.commands[command] = struct{}{}
	}

	for _, valueFlag := range pCfg.ValueFlags {
		parser.valueFlags[valueFlag] = struct{}{}
	}

	tracef("NewParser parser=%+#v", parser)
	tracef("NewParser pCfg=%+#v", pCfg)

	return parser
}

func (p *Parser) Parse() (*Argh, error) {
	p.nodes = []Node{}

	for {
		pd, err := p.parseArg()
		if err != nil {
			return nil, err
		}

		if pd != nil && pd.Break {
			break
		}
	}

	return &Argh{ParseTree: &ParseTree{Nodes: p.nodes}}, nil
}

func (p *Parser) parseArg() (*parseDirective, error) {
	tok, lit, pos := p.scan()
	if tok == ILLEGAL {
		return nil, errors.Wrapf(errSyntax, "illegal value %q at pos=%v", lit, pos)
	}

	if tok == EOL {
		return &parseDirective{Break: true}, nil
	}

	p.unscan()

	node, err := p.nodify()

	tracef("parseArg node=%+#v err=%+#v", node, err)

	if err != nil {
		return nil, errors.Wrapf(err, "value %q at pos=%v", lit, pos)
	}

	if node != nil {
		p.nodes = append(p.nodes, node)
	}

	return nil, nil
}

func (p *Parser) nodify() (Node, error) {
	tok, lit, pos := p.scan()

	tracef("nodify tok=%s lit=%q pos=%v", tok, lit, pos)

	switch tok {
	case IDENT:
		if len(p.nodes) == 0 {
			return Program{Name: lit, Pos: pos - len(lit)}, nil
		}
		return Ident{Literal: lit, Pos: pos - len(lit)}, nil
	case ARG_DELIMITER:
		return ArgDelimiter{Pos: pos - 1}, nil
	case COMPOUND_SHORT_FLAG:
		flagNodes := []Node{}

		for i, r := range lit[1:] {
			flagNodes = append(
				flagNodes,
				Flag{
					Pos:  pos + i + 1,
					Name: string(r),
				},
			)
		}

		return Statement{Pos: pos, Nodes: flagNodes}, nil
	case SHORT_FLAG:
		flagName := string(lit[1:])
		if _, ok := p.valueFlags[flagName]; ok {
			return p.scanValueFlag(flagName, pos)
		}

		return Flag{Name: flagName, Pos: pos - len(flagName) - 1}, nil
	case LONG_FLAG:
		flagName := string(lit[2:])
		if _, ok := p.valueFlags[flagName]; ok {
			return p.scanValueFlag(flagName, pos)
		}

		return Flag{Name: flagName, Pos: pos - len(flagName) - 2}, nil
	default:
	}

	return Ident{Literal: lit, Pos: pos - len(lit)}, nil
}

func (p *Parser) scanValueFlag(flagName string, pos int) (Node, error) {
	tracef("scanValueFlag flagName=%q pos=%v", flagName, pos)

	lit, err := p.scanIdent()
	if err != nil {
		return nil, err
	}

	flagSepLen := len("--") + 1

	return Flag{Name: flagName, Pos: pos - len(lit) - flagSepLen, Value: ptr(lit)}, nil
}

func (p *Parser) scanIdent() (string, error) {
	tok, lit, pos := p.scan()

	nUnscan := 0

	if tok == ASSIGN || tok == ARG_DELIMITER {
		nUnscan++
		tok, lit, pos = p.scan()
	}

	if tok == IDENT {
		return lit, nil
	}

	for i := 0; i < nUnscan; i++ {
		p.unscan()
	}

	return "", errors.Wrapf(errSyntax, "expected ident at pos=%v but got %s (%q)", pos, tok, lit)
}

func (p *Parser) scan() (Token, string, int) {
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit, p.buf.pos
	}

	tok, lit, pos := p.s.Scan()

	p.buf.tok, p.buf.lit, p.buf.pos = tok, lit, pos

	return tok, lit, pos
}

func (p *Parser) unscan() {
	p.buf.n = 1
}

func ptr[T any](v T) *T {
	return &v
}
