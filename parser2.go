package argh

import (
	"fmt"
	"io"
	"strings"
)

type parser2 struct {
	s *Scanner

	cfg *ParserConfig

	errors ScannerErrorList

	tok Token
	lit string
	pos Pos
}

func ParseArgs2(args []string, pCfg *ParserConfig) (*ParseTree, error) {
	parser := &parser2{}
	parser.init(
		strings.NewReader(strings.Join(args, string(nul))),
		pCfg,
	)

	tracef("ParseArgs2 parser=%+#v", parser)

	return parser.parseArgs()
}

func (p *parser2) init(r io.Reader, pCfg *ParserConfig) {
	p.errors = ScannerErrorList{}

	if pCfg == nil {
		pCfg = POSIXyParserConfig
	}

	p.cfg = pCfg

	p.s = NewScanner(r, pCfg.ScannerConfig)

	p.next()
}

func (p *parser2) parseArgs() (*ParseTree, error) {
	if p.errors.Len() != 0 {
		tracef("parseArgs bailing due to initial error")
		return nil, p.errors.Err()
	}

	prog := p.parseCommand(&p.cfg.Prog)

	nodes := []Node{prog}
	if v := p.parsePassthrough(); v != nil {
		nodes = append(nodes, v)
	}

	return &ParseTree{
		Nodes: nodes,
	}, nil
}

func (p *parser2) next() {
	tracef("parser2.next() current: %v %q %v", p.tok, p.lit, p.pos)

	p.tok, p.lit, p.pos = p.s.Scan()

	tracef("parser2.next() next: %v %q %v", p.tok, p.lit, p.pos)
}

func (p *parser2) parseCommand(cCfg *CommandConfig) Node {
	tracef("parseCommand cfg=%+#v", cCfg)

	node := &Command{
		Name:   p.lit,
		Values: map[string]string{},
		Nodes:  []Node{},
	}

	identIndex := 0

	for i := 0; p.tok != EOL; i++ {
		p.next()

		tracef("parseCommand for=%d node.Values=%+#v", i, node.Values)
		tracef("parseCommand for=%d node.Nodes=%+#v", i, node.Values)

		if subCfg, ok := cCfg.Commands[p.lit]; ok {
			subCommand := p.lit

			node.Nodes = append(node.Nodes, p.parseCommand(&subCfg))

			tracef("parseCommand breaking after sub-command=%v", subCommand)
			break
		}

		switch p.tok {
		case ARG_DELIMITER:
			tracef("parseCommand handling %s", p.tok)

			node.Nodes = append(node.Nodes, &ArgDelimiter{})

			continue
		case IDENT, STDIN_FLAG:
			tracef("parseCommand handling %s", p.tok)

			if !cCfg.NValue.Contains(identIndex) {
				tracef("parseCommand identIndex=%d exceeds expected=%s; breaking", identIndex, cCfg.NValue)
				break
			}

			name := fmt.Sprintf("%d", identIndex)

			tracef("parseCommand checking for name of identIndex=%d", identIndex)

			if len(cCfg.ValueNames) > identIndex {
				name = cCfg.ValueNames[identIndex]
				tracef("parseCommand setting name=%s from config value names", name)
			} else if len(cCfg.ValueNames) == 1 && (cCfg.NValue == OneOrMoreValue || cCfg.NValue == ZeroOrMoreValue) {
				name = fmt.Sprintf("%s.%d", cCfg.ValueNames[0], identIndex)
				tracef("parseCommand setting name=%s from repeating value name", name)
			}

			node.Values[name] = p.lit

			identIndex++
		case LONG_FLAG, SHORT_FLAG, COMPOUND_SHORT_FLAG:
			tok := p.tok
			flagNode := p.parseFlag()

			tracef("parseCommand appending %s node=%+#v", tok, flagNode)

			node.Nodes = append(node.Nodes, flagNode)
		default:
			tracef("parseCommand breaking on %s", p.tok)
			break
		}
	}

	tracef("parseCommand returning node=%+#v", node)
	return node
}

func (p *parser2) parseIdent() Node {
	defer p.next()

	node := &Ident{Literal: p.lit}
	return node
}

func (p *parser2) parseFlag() Node {
	defer p.next()

	switch p.tok {
	case SHORT_FLAG:
		return p.parseShortFlag()
	case LONG_FLAG:
		return p.parseLongFlag()
	case COMPOUND_SHORT_FLAG:
		return p.parseCompoundShortFlag()
	}

	panic(fmt.Sprintf("token %v cannot be parsed as flag", p.tok))
}

func (p *parser2) parseShortFlag() Node {
	node := &Flag{Name: string(p.lit[1])}
	// TODO: moar stuff
	return node
}

func (p *parser2) parseLongFlag() Node {
	node := &Flag{Name: string(p.lit[2:])}
	// TODO: moar stuff
	return node
}

func (p *parser2) parseCompoundShortFlag() Node {
	flagNodes := []Node{}

	withoutFlagPrefix := p.lit[1:]

	for _, r := range withoutFlagPrefix {
		flagNodes = append(flagNodes, &Flag{Name: string(r)})
	}

	return &CompoundShortFlag{Nodes: flagNodes}
}

func (p *parser2) parsePassthrough() Node {
	nodes := []Node{}

	for ; p.tok != EOL; p.next() {
		nodes = append(nodes, &Ident{Literal: p.lit})
	}

	if len(nodes) == 0 {
		return nil
	}

	return &PassthroughArgs{Nodes: nodes}
}
