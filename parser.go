//go:generate stringer -type NValue

package argh

import (
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"
)

const (
	OneOrMoreValue  NValue = -2
	ZeroOrMoreValue NValue = -1
	ZeroValue       NValue = 0
)

var (
	ErrSyntax = errors.New("syntax error")

	DefaultParserConfig = &ParserConfig{
		Commands:      map[string]CommandConfig{},
		Flags:         map[string]FlagConfig{},
		ScannerConfig: DefaultScannerConfig,
	}
)

type NValue int

func ParseArgs(args []string, pCfg *ParserConfig) (*ParseTree, error) {
	reEncoded := strings.Join(args, string(nul))

	return NewParser(
		strings.NewReader(reEncoded),
		pCfg,
	).Parse()
}

type Parser struct {
	s *Scanner

	buf []scanEntry

	cfg *ParserConfig

	nodes []Node
	node  Node
}

type ParseTree struct {
	Nodes []Node `json:"nodes"`
}

type scanEntry struct {
	tok Token
	lit string
	pos int
}

type ParserConfig struct {
	Prog     CommandConfig
	Commands map[string]CommandConfig
	Flags    map[string]FlagConfig

	ScannerConfig *ScannerConfig
}

type CommandConfig struct {
	NValue     NValue
	ValueNames []string
	Flags      map[string]FlagConfig
}

type FlagConfig struct {
	NValue     NValue
	ValueNames []string
}

func NewParser(r io.Reader, pCfg *ParserConfig) *Parser {
	if pCfg == nil {
		pCfg = DefaultParserConfig
	}

	parser := &Parser{
		buf: []scanEntry{},
		s:   NewScanner(r, pCfg.ScannerConfig),
		cfg: pCfg,
	}

	tracef("NewParser parser=%+#v", parser)
	tracef("NewParser pCfg=%+#v", pCfg)

	return parser
}

func (p *Parser) Parse() (*ParseTree, error) {
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

	return &ParseTree{Nodes: p.nodes}, nil
}

func (p *Parser) parseArg() (bool, error) {
	tok, lit, pos := p.scan()
	if tok == ILLEGAL {
		return false, errors.Wrapf(ErrSyntax, "illegal value %q at pos=%v", lit, pos)
	}

	if tok == EOL {
		return true, nil
	}

	p.unscan(tok, lit, pos)

	node, err := p.scanNode()

	tracef("parseArg node=%+#v err=%+#v", node, err)

	if err != nil {
		return false, errors.Wrapf(err, "value %q at pos=%v", lit, pos)
	}

	if node != nil {
		p.nodes = append(p.nodes, node)
	}

	return false, nil
}

func (p *Parser) scanNode() (Node, error) {
	tok, lit, pos := p.scan()

	tracef("scanNode tok=%s lit=%q pos=%v", tok, lit, pos)

	switch tok {
	case ARG_DELIMITER:
		return ArgDelimiter{}, nil
	case ASSIGN:
		return nil, errors.Wrapf(ErrSyntax, "bare assignment operator at pos=%v", pos)
	case IDENT:
		p.unscan(tok, lit, pos)
		return p.scanCommandOrIdent()
	case COMPOUND_SHORT_FLAG:
		p.unscan(tok, lit, pos)
		return p.scanCompoundShortFlag()
	case SHORT_FLAG, LONG_FLAG:
		p.unscan(tok, lit, pos)
		return p.scanFlag()
	default:
	}

	return Ident{Literal: lit}, nil
}

func (p *Parser) scanCommandOrIdent() (Node, error) {
	tok, lit, pos := p.scan()

	if len(p.nodes) == 0 {
		p.unscan(tok, lit, pos)
		values, err := p.scanValues(p.cfg.Prog.NValue, p.cfg.Prog.ValueNames)
		if err != nil {
			return nil, err
		}

		return Program{Name: lit, Values: values}, nil
	}

	if cfg, ok := p.cfg.Commands[lit]; ok {
		p.unscan(tok, lit, pos)
		values, err := p.scanValues(cfg.NValue, cfg.ValueNames)
		if err != nil {
			return nil, err
		}

		return Command{Name: lit, Values: values}, nil
	}

	return Ident{Literal: lit}, nil
}

func (p *Parser) scanFlag() (Node, error) {
	tok, lit, pos := p.scan()

	flagName := string(lit[1:])
	if tok == LONG_FLAG {
		flagName = string(lit[2:])
	}

	if cfg, ok := p.cfg.Flags[flagName]; ok {
		p.unscan(tok, flagName, pos)

		values, err := p.scanValues(cfg.NValue, cfg.ValueNames)
		if err != nil {
			return nil, err
		}

		return Flag{Name: flagName, Values: values}, nil
	}

	return Flag{Name: flagName}, nil
}

func (p *Parser) scanCompoundShortFlag() (Node, error) {
	tok, lit, pos := p.scan()

	flagNodes := []Node{}

	withoutFlagPrefix := lit[1:]

	for i, r := range withoutFlagPrefix {
		if i == len(withoutFlagPrefix)-1 {
			flagName := string(r)

			if cfg, ok := p.cfg.Flags[flagName]; ok {
				p.unscan(tok, flagName, pos)

				values, err := p.scanValues(cfg.NValue, cfg.ValueNames)
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
}

func (p *Parser) scanValuesAndFlags() (map[string]string, []Node, error) {
	return nil, nil, nil
}

func (p *Parser) scanValues(n NValue, valueNames []string) (map[string]string, error) {
	_, lit, pos := p.scan()

	tracef("scanValues lit=%q pos=%v n=%v valueNames=%+v", lit, pos, n, valueNames)

	values, err := func() (map[string]string, error) {
		if n == ZeroValue {
			return map[string]string{}, nil
		}

		ret := map[string]string{}
		i := 0

		for {
			lit, err := p.scanIdent()
			if err != nil {
				if n == NValue(1) {
					return nil, err
				}

				if n == OneOrMoreValue {
					break
				}
			}

			name := fmt.Sprintf("%d", i)
			if len(valueNames)-1 >= i {
				name = valueNames[i]
			} else if len(valueNames) > 0 && strings.HasSuffix(valueNames[len(valueNames)-1], "+") {
				name = strings.TrimSuffix(valueNames[len(valueNames)-1], "+")
			}

			ret[name] = lit

			if n == NValue(1) && len(ret) == 1 {
				break
			}

			i++
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

	tracef("scanIdent scanned tok=%s lit=%q pos=%v", tok, lit, pos)

	unscanBuf := []scanEntry{}

	if tok == ASSIGN || tok == ARG_DELIMITER {
		entry := scanEntry{tok: tok, lit: lit, pos: pos}

		tracef("scanIdent tok=%s; scanning next and pushing to unscan buffer entry=%+#v", tok, entry)

		unscanBuf = append([]scanEntry{entry}, unscanBuf...)

		tok, lit, pos = p.scan()
	}

	if tok == IDENT {
		return lit, nil
	}

	entry := scanEntry{tok: tok, lit: lit, pos: pos}

	tracef("scanIdent tok=%s; unscanning entry=%+#v", tok, entry)

	unscanBuf = append([]scanEntry{entry}, unscanBuf...)

	for _, entry := range unscanBuf {
		p.unscan(entry.tok, entry.lit, entry.pos)
	}

	return "", errors.Wrapf(ErrSyntax, "expected ident at pos=%v but got %s (%q)", pos, tok, lit)
}

func (p *Parser) scan() (Token, string, int) {
	if len(p.buf) != 0 {
		entry, buf := p.buf[len(p.buf)-1], p.buf[:len(p.buf)-1]
		p.buf = buf

		tracef("scan returning buffer entry=%s %+#v", entry.tok, entry)
		return entry.tok, entry.lit, entry.pos
	}

	tok, lit, pos := p.s.Scan()

	tracef("scan returning next=%s %+#v", tok, scanEntry{tok: tok, lit: lit, pos: pos})

	return tok, lit, pos
}

func (p *Parser) unscan(tok Token, lit string, pos int) {
	entry := scanEntry{tok: tok, lit: lit, pos: pos}

	tracef("unscan entry=%s %+#v", tok, entry)

	p.buf = append(p.buf, entry)
}
