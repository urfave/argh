//go:generate stringer -type NValue

package argh

import (
	"io"
	"strings"

	"github.com/pkg/errors"
)

const (
	ZeroValue NValue = iota
	OneValue
	OneOrMoreValue
)

var (
	errSyntax = errors.New("syntax error")

	DefaultParserConfig = &ParserConfig{
		Commands:      map[string]NValue{},
		Flags:         map[string]NValue{},
		ScannerConfig: DefaultScannerConfig,
	}
)

type NValue int

func ParseArgs(args []string, pCfg *ParserConfig) (*Argh, error) {
	reEncoded := strings.Join(args, string(nul))

	return NewParser(
		strings.NewReader(reEncoded),
		pCfg,
	).Parse()
}

type Parser struct {
	s *Scanner

	buf []ScanEntry

	cfg *ParserConfig

	nodes    []Node
	stopSeen bool
}

type ScanEntry struct {
	tok Token
	lit string
	pos int
}

type ParserConfig struct {
	ProgValues    NValue
	Commands      map[string]NValue
	Flags         map[string]NValue
	ScannerConfig *ScannerConfig
}

func NewParser(r io.Reader, pCfg *ParserConfig) *Parser {
	if pCfg == nil {
		pCfg = DefaultParserConfig
	}

	parser := &Parser{
		buf: []ScanEntry{},
		s:   NewScanner(r, pCfg.ScannerConfig),
		cfg: pCfg,
	}

	tracef("NewParser parser=%+#v", parser)
	tracef("NewParser pCfg=%+#v", pCfg)

	return parser
}

func (p *Parser) Parse() (*Argh, error) {
	p.nodes = []Node{}

	for {
		br, err := p.parseArg()
		if err != nil {
			return nil, err
		}

		if br {
			break
		}
	}

	return &Argh{ParseTree: &ParseTree{Nodes: p.nodes}}, nil
}

func (p *Parser) parseArg() (bool, error) {
	tok, lit, pos := p.scan()
	if tok == ILLEGAL {
		return false, errors.Wrapf(errSyntax, "illegal value %q at pos=%v", lit, pos)
	}

	if tok == EOL {
		return true, nil
	}

	p.unscan(tok, lit, pos)

	node, err := p.nodify()

	tracef("parseArg node=%+#v err=%+#v", node, err)

	if err != nil {
		return false, errors.Wrapf(err, "value %q at pos=%v", lit, pos)
	}

	if node != nil {
		p.nodes = append(p.nodes, node)
	}

	return false, nil
}

func (p *Parser) nodify() (Node, error) {
	tok, lit, pos := p.scan()

	tracef("nodify tok=%s lit=%q pos=%v", tok, lit, pos)

	switch tok {
	case IDENT:
		if len(p.nodes) == 0 {
			values, err := p.scanValues(lit, pos, p.cfg.ProgValues)
			if err != nil {
				return nil, err
			}

			return Program{Name: lit, Values: values}, nil
		}

		if n, ok := p.cfg.Commands[lit]; ok {
			values, err := p.scanValues(lit, pos, n)
			if err != nil {
				return nil, err
			}

			return Command{Name: lit, Values: values}, nil
		}

		return Ident{Literal: lit}, nil
	case ARG_DELIMITER:
		return ArgDelimiter{}, nil
	case COMPOUND_SHORT_FLAG:
		flagNodes := []Node{}

		withoutFlagPrefix := lit[1:]

		for i, r := range withoutFlagPrefix {
			if i == len(withoutFlagPrefix)-1 {
				flagName := string(r)

				if n, ok := p.cfg.Flags[flagName]; ok {
					values, err := p.scanValues(flagName, pos, n)
					if err != nil {
						return nil, err
					}

					flagNodes = append(flagNodes, Flag{Name: flagName, Values: values})

					continue
				}
			}

			flagNodes = append(
				flagNodes,
				Flag{
					Name: string(r),
				},
			)
		}

		return CompoundShortFlag{Nodes: flagNodes}, nil
	case SHORT_FLAG:
		flagName := string(lit[1:])
		if n, ok := p.cfg.Flags[flagName]; ok {
			values, err := p.scanValues(flagName, pos, n)
			if err != nil {
				return nil, err
			}

			return Flag{Name: flagName, Values: values}, nil
		}

		return Flag{Name: flagName}, nil
	case LONG_FLAG:
		flagName := string(lit[2:])
		if n, ok := p.cfg.Flags[flagName]; ok {
			values, err := p.scanValues(flagName, pos, n)
			if err != nil {
				return nil, err
			}

			return Flag{Name: flagName, Values: values}, nil
		}

		return Flag{Name: flagName}, nil
	default:
	}

	return Ident{Literal: lit}, nil
}

func (p *Parser) scanValues(lit string, pos int, n NValue) ([]string, error) {
	tracef("scanValues lit=%q pos=%v n=%v", lit, pos, n)

	values, err := func() ([]string, error) {
		if n == ZeroValue {
			return []string{}, nil
		}

		ret := []string{}

		for {
			lit, err := p.scanIdent()
			if err != nil {
				if n == OneValue {
					return nil, err
				}

				if n == OneOrMoreValue {
					break
				}
			}

			ret = append(ret, lit)

			if n == OneValue && len(ret) == 1 {
				break
			}
		}

		return ret, nil
	}()

	if err != nil {
		return nil, err
	}

	if len(values) == 0 {
		return nil, nil
	}

	return values, nil
}

func (p *Parser) scanIdent() (string, error) {
	tok, lit, pos := p.scan()

	unscanBuf := []ScanEntry{}

	if tok == ASSIGN || tok == ARG_DELIMITER {
		unscanBuf = append([]ScanEntry{{tok: tok, lit: lit, pos: pos}}, unscanBuf...)

		tok, lit, pos = p.scan()
	}

	if tok == IDENT {
		return lit, nil
	}

	unscanBuf = append([]ScanEntry{{tok: tok, lit: lit, pos: pos}}, unscanBuf...)

	for _, entry := range unscanBuf {
		p.unscan(entry.tok, entry.lit, entry.pos)
	}

	return "", errors.Wrapf(errSyntax, "expected ident at pos=%v but got %s (%q)", pos, tok, lit)
}

func (p *Parser) scan() (Token, string, int) {
	if len(p.buf) != 0 {
		entry, buf := p.buf[len(p.buf)-1], p.buf[:len(p.buf)-1]
		p.buf = buf

		tracef("scan returning buffer entry=%s %+#v", entry.tok, entry)
		return entry.tok, entry.lit, entry.pos
	}

	tok, lit, pos := p.s.Scan()

	tracef("scan returning next=%s %+#v", tok, ScanEntry{tok: tok, lit: lit, pos: pos})

	return tok, lit, pos
}

func (p *Parser) unscan(tok Token, lit string, pos int) {
	entry := ScanEntry{tok: tok, lit: lit, pos: pos}

	tracef("unscan entry=%s %+#v", tok, entry)

	p.buf = append(p.buf, entry)
}
