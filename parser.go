package argh

import (
	"fmt"
	"io"
	"strings"
)

type parser struct {
	s *Scanner

	cfg *ParserConfig

	errors ParserErrorList

	tok Token
	lit string
	pos Pos

	buffered bool
}

type ParseTree struct {
	Nodes []Node `json:"nodes"`
}

func ParseArgs(args []string, pCfg *ParserConfig) (*ParseTree, error) {
	p := &parser{}
	p.init(
		strings.NewReader(strings.Join(args, string(nul))),
		pCfg,
	)

	tracef("ParseArgs(...) parser=%+#v", p)

	return p.parseArgs()
}

func (p *parser) addError(msg string) {
	p.errors.Add(Position{Column: int(p.pos)}, msg)
}

func (p *parser) init(r io.Reader, pCfg *ParserConfig) {
	p.errors = ParserErrorList{}

	if pCfg == nil {
		pCfg = POSIXyParserConfig
	}

	p.cfg = pCfg

	p.s = NewScanner(r, pCfg.ScannerConfig)

	p.next()
}

func (p *parser) parseArgs() (*ParseTree, error) {
	if p.errors.Len() != 0 {
		tracef("parseArgs() bailing due to initial error")
		return nil, p.errors.Err()
	}

	tracef("parseArgs() parsing %q as program command; cfg=%+#v", p.lit, p.cfg.Prog)
	prog := p.parseCommand(p.cfg.Prog)

	tracef("parseArgs() top level node is %T", prog)

	nodes := []Node{prog}
	if v := p.parsePassthrough(); v != nil {
		tracef("parseArgs() appending passthrough argument %v", v)
		nodes = append(nodes, v)
	}

	tracef("parseArgs() returning ParseTree")

	return &ParseTree{Nodes: nodes}, p.errors.Err()
}

func (p *parser) next() {
	tracef("next() before scan: %v %q %v", p.tok, p.lit, p.pos)

	p.tok, p.lit, p.pos = p.s.Scan()

	tracef("next() after scan: %v %q %v", p.tok, p.lit, p.pos)
}

func (p *parser) parseCommand(cCfg *CommandConfig) Node {
	tracef("parseCommand(%+#v)", cCfg)

	node := &Command{
		Name: p.lit,
	}
	values := map[string]string{}
	nodes := []Node{}

	identIndex := 0

	for i := 0; p.tok != EOL; i++ {
		if !p.buffered {
			tracef("parseCommand(...) buffered=false; scanning next")
			p.next()
		}

		p.buffered = false

		tracef("parseCommand(...) for=%d values=%+#v", i, values)
		tracef("parseCommand(...) for=%d nodes=%+#v", i, nodes)
		tracef("parseCommand(...) for=%d tok=%s lit=%q pos=%v", i, p.tok, p.lit, p.pos)

		tracef("parseCommand(...) cCfg=%+#v", cCfg)

		if subCfg, ok := cCfg.GetCommandConfig(p.lit); ok {
			subCommand := p.lit

			nodes = append(nodes, p.parseCommand(&subCfg))

			tracef("parseCommand(...) breaking after sub-command=%v", subCommand)
			break
		}

		switch p.tok {
		case ARG_DELIMITER:
			tracef("parseCommand(...) handling %s", p.tok)

			nodes = append(nodes, &ArgDelimiter{})

			continue
		case IDENT, STDIN_FLAG:
			tracef("parseCommand(...) handling %s", p.tok)

			if cCfg.NValue.Contains(identIndex) {
				name := fmt.Sprintf("%d", identIndex)

				tracef("parseCommand(...) checking for name of identIndex=%d", identIndex)

				if len(cCfg.ValueNames) > identIndex {
					name = cCfg.ValueNames[identIndex]
					tracef("parseCommand(...) setting name=%s from config value names", name)
				} else if len(cCfg.ValueNames) == 1 && (cCfg.NValue == OneOrMoreValue || cCfg.NValue == ZeroOrMoreValue) {
					name = fmt.Sprintf("%s.%d", cCfg.ValueNames[0], identIndex)
					tracef("parseCommand(...) setting name=%s from repeating value name", name)
				}

				values[name] = p.lit
			}

			if p.tok == STDIN_FLAG {
				nodes = append(nodes, &StdinFlag{})
			} else {
				nodes = append(nodes, &Ident{Literal: p.lit})
			}

			identIndex++
		case LONG_FLAG, SHORT_FLAG, COMPOUND_SHORT_FLAG:
			tok := p.tok

			flagNode := p.parseFlag(cCfg.Flags)

			tracef("parseCommand(...) appending %s node=%+#v", tok, flagNode)

			nodes = append(nodes, flagNode)
		case ASSIGN:
			tracef("parseCommand(...) error on bare %s", p.tok)

			p.addError("invalid bare assignment")

			break
		default:
			tracef("parseCommand(...) breaking on %s", p.tok)
			break
		}
	}

	if len(nodes) > 0 {
		node.Nodes = nodes
	}

	if len(values) > 0 {
		node.Values = values
	}

	tracef("parseCommand(...) returning node=%+#v", node)
	return node
}

func (p *parser) parseIdent() Node {
	node := &Ident{Literal: p.lit}
	return node
}

func (p *parser) parseFlag(flags *Flags) Node {
	switch p.tok {
	case SHORT_FLAG:
		tracef("parseFlag(...) parsing short flag with config=%+#v", flags)
		return p.parseShortFlag(flags)
	case LONG_FLAG:
		tracef("parseFlag(...) parsing long flag with config=%+#v", flags)
		return p.parseLongFlag(flags)
	case COMPOUND_SHORT_FLAG:
		tracef("parseFlag(...) parsing compound short flag with config=%+#v", flags)
		return p.parseCompoundShortFlag(flags)
	}

	panic(fmt.Sprintf("token %v cannot be parsed as flag", p.tok))
}

func (p *parser) parseShortFlag(flags *Flags) Node {
	node := &Flag{Name: string(p.lit[1])}

	flCfg, ok := flags.Get(node.Name)
	if !ok {
		p.addError(fmt.Sprintf("unknown flag %q", node.Name))

		return node
	}

	return p.parseConfiguredFlag(node, flCfg)
}

func (p *parser) parseLongFlag(flags *Flags) Node {
	node := &Flag{Name: string(p.lit[2:])}

	flCfg, ok := flags.Get(node.Name)
	if !ok {
		p.addError(fmt.Sprintf("unknown flag %q", node.Name))

		return node
	}

	return p.parseConfiguredFlag(node, flCfg)
}

func (p *parser) parseCompoundShortFlag(flags *Flags) Node {
	flagNodes := []Node{}

	withoutFlagPrefix := p.lit[1:]

	for i, r := range withoutFlagPrefix {
		node := &Flag{Name: string(r)}

		if i == len(withoutFlagPrefix)-1 {
			flCfg, ok := flags.Get(node.Name)
			if !ok {
				p.addError(fmt.Sprintf("unknown flag %q", node.Name))

				continue
			}

			flagNodes = append(flagNodes, p.parseConfiguredFlag(node, flCfg))

			continue
		}

		flagNodes = append(flagNodes, node)
	}

	return &CompoundShortFlag{Nodes: flagNodes}
}

func (p *parser) parseConfiguredFlag(node *Flag, flCfg FlagConfig) Node {
	values := map[string]string{}
	nodes := []Node{}

	identIndex := 0

	for i := 0; p.tok != EOL; i++ {
		if !flCfg.NValue.Contains(identIndex) {
			tracef("parseConfiguredFlag(...) identIndex=%d exceeds expected=%v; breaking", identIndex, flCfg.NValue)
			break
		}

		p.next()

		switch p.tok {
		case ARG_DELIMITER:
			nodes = append(nodes, &ArgDelimiter{})

			continue
		case ASSIGN:
			nodes = append(nodes, &Assign{})

			continue
		case IDENT, STDIN_FLAG:
			name := fmt.Sprintf("%d", identIndex)

			tracef("parseConfiguredFlag(...) checking for name of identIndex=%d", identIndex)

			if len(flCfg.ValueNames) > identIndex {
				name = flCfg.ValueNames[identIndex]
				tracef("parseConfiguredFlag(...) setting name=%s from config value names", name)
			} else if len(flCfg.ValueNames) == 1 && (flCfg.NValue == OneOrMoreValue || flCfg.NValue == ZeroOrMoreValue) {
				name = fmt.Sprintf("%s.%d", flCfg.ValueNames[0], identIndex)
				tracef("parseConfiguredFlag(...) setting name=%s from repeating value name", name)
			} else {
				tracef("parseConfiguredFlag(...) setting name=%s", name)
			}

			values[name] = p.lit

			if p.tok == STDIN_FLAG {
				nodes = append(nodes, &StdinFlag{})
			} else {
				nodes = append(nodes, &Ident{Literal: p.lit})
			}

			identIndex++
		default:
			tracef("parseConfiguredFlag(...) breaking on %s %q %v; setting buffered=true", p.tok, p.lit, p.pos)
			p.buffered = true

			if len(nodes) > 0 {
				node.Nodes = nodes
			}

			if len(values) > 0 {
				node.Values = values
			}

			return node
		}
	}

	if len(nodes) > 0 {
		node.Nodes = nodes
	}

	if len(values) > 0 {
		node.Values = values
	}

	return node
}

func (p *parser) parsePassthrough() Node {
	nodes := []Node{}

	for ; p.tok != EOL; p.next() {
		nodes = append(nodes, &Ident{Literal: p.lit})
	}

	if len(nodes) == 0 {
		return nil
	}

	return &PassthroughArgs{Nodes: nodes}
}
