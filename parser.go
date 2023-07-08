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

	if err := p.init(
		strings.NewReader(strings.Join(args, string(nul))),
		pCfg,
	); err != nil {
		return nil, err
	}

	tracef("ParseArgs(...) parser=%+#v", p)

	return p.parseArgs()
}

func (p *parser) addError(msg string) {
	p.errors.Add(Position{Column: int(p.pos)}, msg)
}

func (p *parser) init(r io.Reader, pCfg *ParserConfig) error {
	p.errors = ParserErrorList{}

	if pCfg == nil {
		return fmt.Errorf("nil parser config: %w", Err)
	}

	p.cfg = pCfg

	p.s = NewScanner(r, pCfg.ScannerConfig)

	p.next()

	return nil
}

func (p *parser) parseArgs() (*ParseTree, error) {
	if p.errors.Len() != 0 {
		tracef("parseArgs() bailing due to initial error")
		return nil, p.errors.Err()
	}

	tracef("parseArgs() parsing %q as program command; cfg=%+#v", p.lit, p.cfg.Prog)
	prog, err := p.parseCommand(p.cfg.Prog)
	if err != nil {
		return nil, err
	}

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

func (p *parser) parseCommand(cCfg *CommandConfig) (Node, error) {
	tracef("parseCommand(%+#v)", cCfg)

	node := &CommandFlag{
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

			subNode, err := p.parseCommand(&subCfg)
			if err != nil {
				return node, err
			}

			nodes = append(nodes, subNode)

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

			flagNode, err := p.parseFlag(cCfg.Flags)
			if err != nil {
				return node, err
			}

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

	if cCfg.On != nil {
		tracef("parseCommand(...) calling command config handler for node=%+#v", node)
		if err := cCfg.On(*node); err != nil {
			return node, err
		}
	} else {
		tracef("parseCommand(...) no command config handler for node=%+#v", node)
	}

	tracef("parseCommand(...) returning node=%+#v", node)
	return node, nil
}

func (p *parser) parseIdent() Node {
	node := &Ident{Literal: p.lit}
	return node
}

func (p *parser) parseFlag(flags *Flags) (Node, error) {
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

func (p *parser) parseShortFlag(flags *Flags) (Node, error) {
	node := &CommandFlag{Name: string(p.lit[1])}

	flCfg, ok := flags.Get(node.Name)
	if !ok {
		errMsg := fmt.Sprintf("unknown flag %[1]q", node.Name)
		p.addError(errMsg)

		return node, &CommandFlagError{
			Pos:  Position{Column: int(p.pos)},
			Node: *node,
			Msg:  errMsg,
		}
	}

	return p.parseConfiguredFlag(node, flCfg, nil)
}

func (p *parser) parseLongFlag(flags *Flags) (Node, error) {
	node := &CommandFlag{Name: string(p.lit[2:])}

	flCfg, ok := flags.Get(node.Name)
	if !ok {
		errMsg := fmt.Sprintf("unknown flag %[1]q", node.Name)
		p.addError(errMsg)

		return node, &CommandFlagError{
			Pos:  Position{Column: int(p.pos)},
			Node: *node,
			Msg:  errMsg,
		}
	}

	return p.parseConfiguredFlag(node, flCfg, nil)
}

func (p *parser) parseCompoundShortFlag(flags *Flags) (Node, error) {
	unparsedFlags := []*CommandFlag{}
	unparsedFlagConfigs := []FlagConfig{}

	withoutFlagPrefix := p.lit[1:]

	for _, r := range withoutFlagPrefix {
		node := &CommandFlag{Name: string(r)}

		flCfg, ok := flags.Get(node.Name)
		if !ok {
			errMsg := fmt.Sprintf("unknown flag %[1]q", node.Name)
			p.addError(errMsg)

			return node, &CommandFlagError{
				Pos:  Position{Column: int(p.pos)},
				Node: *node,
				Msg:  errMsg,
			}
		}

		unparsedFlags = append(unparsedFlags, node)
		unparsedFlagConfigs = append(unparsedFlagConfigs, flCfg)
	}

	flagNodes := []Node{}

	for i, node := range unparsedFlags {
		flCfg := unparsedFlagConfigs[i]

		if i != len(unparsedFlags)-1 {
			// NOTE: if a compound short flag is configured to accept
			// more than zero values but is not the last flag in the
			// group, it will be parsed with an override NValue of
			// ZeroValue so that it does not consume the next token.
			if flCfg.NValue.Required() {
				errMsg := fmt.Sprintf(
					"short flag %[1]q before end of compound group expects value",
					node.Name,
				)
				p.addError(
					errMsg,
				)

				return nil, fmt.Errorf(errMsg+": %[1]w", Err)
			}

			flagNode, err := p.parseConfiguredFlag(node, flCfg, zeroValuePtr)
			if err != nil {
				return nil, err
			}

			flagNodes = append(flagNodes, flagNode)

			continue
		}

		flagNode, err := p.parseConfiguredFlag(node, flCfg, nil)
		if err != nil {
			return nil, err
		}

		flagNodes = append(flagNodes, flagNode)
	}

	return &CompoundShortFlag{Nodes: flagNodes}, nil
}

func (p *parser) parseConfiguredFlag(node *CommandFlag, flCfg FlagConfig, nValueOverride *NValue) (Node, error) {
	values := map[string]string{}
	nodes := []Node{}

	atExit := func() (*CommandFlag, error) {
		if len(nodes) > 0 {
			node.Nodes = nodes
		}

		if len(values) > 0 {
			node.Values = values
		}

		if flCfg.On != nil {
			tracef("parseConfiguredFlag(...) calling flag config handler for node=%+#[1]v", node)
			if err := flCfg.On(*node); err != nil {
				return nil, err
			}
		} else {
			tracef("parseConfiguredFlag(...) no flag config handler for node=%+#[1]v", node)
		}

		return node, nil
	}

	identIndex := 0

	for i := 0; p.tok != EOL; i++ {
		if nValueOverride != nil && !(*nValueOverride).Contains(identIndex) {
			tracef("parseConfiguredFlag(...) identIndex=%d exceeds expected=%v; breaking", identIndex, *nValueOverride)
			break
		}

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

			return atExit()
		}
	}

	return atExit()
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
