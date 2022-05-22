package argh

import (
	"fmt"
	"io"
	"strings"
)

type parser2 struct {
	s *Scanner

	commands map[string]struct{}

	errors ScannerErrorList

	tok Token
	lit string
	pos Pos
}

func ParseArgs2(args, commands []string) (*ParseTree, error) {
	parser := &parser2{}
	parser.init(
		strings.NewReader(strings.Join(args, string(nul))),
		commands,
	)

	tracef("ParseArgs2 parser=%+#v", parser)

	return parser.parseArgs()
}

func (p *parser2) init(r io.Reader, commands []string) {
	p.errors = ScannerErrorList{}
	commandMap := map[string]struct{}{}

	for _, c := range commands {
		commandMap[c] = struct{}{}
	}

	p.s = NewScanner(r, nil)
	p.commands = commandMap

	p.next()
}

func (p *parser2) parseArgs() (*ParseTree, error) {
	if p.errors.Len() != 0 {
		tracef("parseArgs bailing due to initial error")
		return nil, p.errors.Err()
	}

	prog := &Program{
		Name:   p.lit,
		Values: map[string]string{},
		Nodes:  []Node{},
	}
	p.next()

	for p.tok != EOL && p.tok != STOP_FLAG {
		prog.Nodes = append(prog.Nodes, p.parseArg())
	}

	return &ParseTree{
		Nodes: []Node{
			prog, p.parsePassthrough(),
		},
	}, nil
}

func (p *parser2) next() {
	tracef("parser2.next() <- %v %q %v", p.tok, p.lit, p.pos)
	defer func() {
		tracef("parser2.next() -> %v %q %v", p.tok, p.lit, p.pos)
	}()

	p.tok, p.lit, p.pos = p.s.Scan()
}

func (p *parser2) parseArg() Node {
	switch p.tok {
	case ARG_DELIMITER:
		p.next()
		return &ArgDelimiter{}
	case IDENT:
		if _, ok := p.commands[p.lit]; ok {
			return p.parseCommand()
		}
		return p.parseIdent()
	case LONG_FLAG, SHORT_FLAG, COMPOUND_SHORT_FLAG:
		return p.parseFlag()
	}

	pos := p.pos
	lit := p.lit
	p.advanceArg()
	return &BadArg{Literal: lit, From: pos, To: p.pos}
}

func (p *parser2) advanceArg() {
	for ; p.tok != EOL; p.next() {
		switch p.tok {
		case IDENT, LONG_FLAG, SHORT_FLAG, COMPOUND_SHORT_FLAG:
			return
		}
	}
}

func (p *parser2) parseCommand() Node {
	node := &Command{Name: p.lit, Values: map[string]string{}, Nodes: []Node{}}

	for i := 0; p.tok != EOL; i++ {
		p.next()

		if _, ok := p.commands[p.lit]; ok {
			break
		}

		switch p.tok {
		case ARG_DELIMITER:
			continue
		case IDENT, STDIN_FLAG:
			node.Values[fmt.Sprintf("%d", i)] = p.lit
		case LONG_FLAG, SHORT_FLAG, COMPOUND_SHORT_FLAG:
			node.Nodes = append(node.Nodes, p.parseFlag())
		default:
			break
		}
	}

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

	return &PassthroughArgs{Nodes: nodes}
}
