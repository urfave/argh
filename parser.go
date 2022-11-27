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

	tracef(2, "ParseArgs(...) (parser %[1]p)", p)

	return p.parseArgs()
}

func (p *parser) addError(msg string) {
	p.errors.Add(Position{Column: int(p.pos)}, msg)
}

func (p *parser) init(r io.Reader, pCfg *ParserConfig) error {
	p.errors = ParserErrorList{}

	if pCfg == nil {
		return fmt.Errorf("nil parser config: %w", Error)
	}

	p.cfg = pCfg

	p.s = NewScanner(r, pCfg.ScannerConfig)

	p.next()

	return nil
}

func (p *parser) parseArgs() (*ParseTree, error) {
	if p.errors.Len() != 0 {
		tracef(2, "parseArgs() bailing due to initial error")
		return nil, p.errors.Err()
	}

	tracef(2, "parseArgs() parsing %[1]q as program command (cfg %[2]p)", p.lit, p.cfg.Prog)
	prog := p.parseCommand(p.cfg.Prog)

	tracef(2, "parseArgs() top level node is %[1]T (%[1]p)", prog)

	nodes := []Node{prog}
	if v := p.parsePassthrough(); v != nil {
		tracef(2, "parseArgs() appending passthrough argument %v", v)
		nodes = append(nodes, v)
	}

	tracef(2, "parseArgs() returning ParseTree")

	return &ParseTree{Nodes: nodes}, p.errors.Err()
}

func (p *parser) next() {
	tracef(2, "next() before scan: %v %q %v", p.tok, p.lit, p.pos)

	p.tok, p.lit, p.pos = p.s.Scan()

	tracef(2, "next() after scan: %v %q %v", p.tok, p.lit, p.pos)
}

func (p *parser) parseCommand(cCfg *CommandConfig) Node {
	tracef(2, "parseCommand(%[1]p)", cCfg)

	node := &CommandFlag{
		Name: p.lit,
	}
	values := map[string]string{}
	nodes := []Node{}

	identIndex := 0

	for i := 0; p.tok != EOL; i++ {
		if !p.buffered {
			tracef(2, "parseCommand(%[1]p) buffered=false; scanning next", cCfg)
			p.next()
		}

		p.buffered = false

		tracef(2, "parseCommand(%[1]p) for=%[2]d values=%+#[3]v", cCfg, i, values)
		tracef(2, "parseCommand(%[1]p) for=%[2]d nodes=%+#[3]v", cCfg, i, nodes)
		tracef(2, "parseCommand(%[1]p) for=%[2]d tok=%[3]s lit=%[4]q pos=%[5]v", cCfg, i, p.tok, p.lit, p.pos)

		tracef(2, "parseCommand(%[1]p)", cCfg)

		switch p.tok {
		case EOL:
			tracef(2, "parseCommand(%[1]p) breaking on %[2]s", cCfg, p.tok)
			break
		case ARG_DELIMITER:
			tracef(2, "parseCommand(%[1]p) handling %[2]s", cCfg, p.tok)

			nodes = append(nodes, &ArgDelimiter{})

			continue
		case IDENT, STDIN_FLAG:
			if subCfg, ok := cCfg.GetCommandConfig(p.lit); ok {
				subCommand := p.lit

				tracef(2, "parseCommand(%[1]p) descending into sub-command=%[2]q (subCfg=%[3]p)", cCfg, subCommand, subCfg)

				nodes = append(nodes, p.parseCommand(subCfg))

				tracef(2, "parseCommand(%[1]p) breaking after sub-command=%[2]q (subCfg=%[3]p)", cCfg, subCommand, subCfg)
				break
			}

			tracef(2, "parseCommand(%[1]p) handling %[2]s", cCfg, p.tok)

			if cCfg.NValue.Contains(identIndex) {
				name := fmt.Sprintf("%d", identIndex)

				tracef(2, "parseCommand(%[1]p) checking for name of identIndex=%[2]d", cCfg, identIndex)

				if len(cCfg.ValueNames) > identIndex {
					name = cCfg.ValueNames[identIndex]
					tracef(2, "parseCommand(%[1]p) setting name=%[2]s from config value names", cCfg, name)
				} else if len(cCfg.ValueNames) == 1 && (cCfg.NValue == OneOrMoreValue || cCfg.NValue == ZeroOrMoreValue) {
					name = fmt.Sprintf("%s.%d", cCfg.ValueNames[0], identIndex)
					tracef(2, "parseCommand(%[1]p) setting name=%[2]s from repeating value name", cCfg, name)
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

			tracef(2, "parseCommand(%[1]p) appending %[2]s node=%+#[3]v", cCfg, tok, flagNode)

			nodes = append(nodes, flagNode)
		case ASSIGN:
			tracef(2, "parseCommand(%[1]p) error on bare %[2]s", cCfg, p.tok)

			p.addError("invalid bare assignment")

			break
		default:
			tracef(2, "parseCommand(%[1]p) breaking on %[2]s", cCfg, p.tok)
			break
		}
	}

	if len(nodes) > 0 {
		node.Nodes = nodes
	}

	if len(values) > 0 {
		node.Values = values
	}

	tracef(2, "parseCommand(%[1]p) applying callbacks for node=%[2]p", cCfg, node)
	if err := cCfg.applyCallbacks(node); err != nil {
		p.addError(err.Error())
	}

	tracef(2, "parseCommand(%[1]p) returning node=%[2]p", cCfg, node)
	return node
}

func (p *parser) parseIdent() Node {
	node := &Ident{Literal: p.lit}
	return node
}

func (p *parser) parseFlag(flags *Flags) Node {
	switch p.tok {
	case SHORT_FLAG:
		tracef(2, "parseFlag(...) parsing short flag with config=%[1]p", flags)
		return p.parseShortFlag(flags)
	case LONG_FLAG:
		tracef(2, "parseFlag(...) parsing long flag with config=%[1]p", flags)
		return p.parseLongFlag(flags)
	case COMPOUND_SHORT_FLAG:
		tracef(2, "parseFlag(...) parsing compound short flag with config=%[1]p", flags)
		return p.parseCompoundShortFlag(flags)
	}

	panic(fmt.Sprintf("token %v cannot be parsed as flag", p.tok))
}

func (p *parser) parseShortFlag(flags *Flags) Node {
	node := &CommandFlag{Name: string(p.lit[1])}

	flCfg, ok := flags.Get(node.Name)
	if !ok {
		p.addError(fmt.Sprintf("unknown flag %[1]q", string(p.cfg.ScannerConfig.FlagPrefix)+node.Name))

		return node
	}

	return p.parseConfiguredFlag(node, flCfg, nil)
}

func (p *parser) parseLongFlag(flags *Flags) Node {
	node := &CommandFlag{Name: string(p.lit[2:])}

	flCfg, ok := flags.Get(node.Name)
	if !ok {
		p.addError(fmt.Sprintf("unknown flag %[1]q", string(p.cfg.ScannerConfig.FlagPrefix)+string(p.cfg.ScannerConfig.FlagPrefix)+node.Name))

		return node
	}

	return p.parseConfiguredFlag(node, flCfg, nil)
}

func (p *parser) parseCompoundShortFlag(flags *Flags) Node {
	unparsedFlags := []*CommandFlag{}
	unparsedFlagConfigs := []*FlagConfig{}

	withoutFlagPrefix := p.lit[1:]

	for _, r := range withoutFlagPrefix {
		node := &CommandFlag{Name: string(r)}

		flCfg, ok := flags.Get(node.Name)
		if !ok {
			p.addError(fmt.Sprintf("unknown flag %[1]q", string(p.cfg.ScannerConfig.FlagPrefix)+node.Name))

			continue
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
				p.addError(
					fmt.Sprintf(
						"short flag %[1]q before end of compound group expects value",
						node.Name,
					),
				)
			}

			flagNodes = append(
				flagNodes,
				p.parseConfiguredFlag(node, flCfg, zeroValuePtr),
			)

			continue
		}

		flagNodes = append(flagNodes, p.parseConfiguredFlag(node, flCfg, nil))
	}

	return &CompoundShortFlag{Nodes: flagNodes}
}

func (p *parser) parseConfiguredFlag(node *CommandFlag, flCfg *FlagConfig, nValueOverride *NValue) Node {
	values := map[string]string{}
	nodes := []Node{}

	atExit := func() *CommandFlag {
		if len(nodes) > 0 {
			node.Nodes = nodes
		}

		if len(values) > 0 {
			node.Values = values
		}

		if err := flCfg.applyCallbacks(node); err != nil {
			p.addError(err.Error())
		}

		return node
	}

	identIndex := 0

	for i := 0; p.tok != EOL; i++ {
		if nValueOverride != nil && !(*nValueOverride).Contains(identIndex) {
			tracef(2, "parseConfiguredFlag(...) identIndex=%d exceeds expected=%v; breaking", identIndex, *nValueOverride)
			break
		}

		if !flCfg.NValue.Contains(identIndex) {
			tracef(2, "parseConfiguredFlag(...) identIndex=%d exceeds expected=%v; breaking", identIndex, flCfg.NValue)
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

			tracef(2, "parseConfiguredFlag(...) checking for name of identIndex=%d", identIndex)

			if len(flCfg.ValueNames) > identIndex {
				name = flCfg.ValueNames[identIndex]
				tracef(2, "parseConfiguredFlag(...) setting name=%s from config value names", name)
			} else if len(flCfg.ValueNames) == 1 && (flCfg.NValue == OneOrMoreValue || flCfg.NValue == ZeroOrMoreValue) {
				name = fmt.Sprintf("%s.%d", flCfg.ValueNames[0], identIndex)
				tracef(2, "parseConfiguredFlag(...) setting name=%s from repeating value name", name)
			} else {
				tracef(2, "parseConfiguredFlag(...) setting name=%s", name)
			}

			values[name] = p.lit

			if p.tok == STDIN_FLAG {
				nodes = append(nodes, &StdinFlag{})
			} else {
				nodes = append(nodes, &Ident{Literal: p.lit})
			}

			identIndex++
		default:
			tracef(2, "parseConfiguredFlag(...) breaking on %s %q %v; setting buffered=true", p.tok, p.lit, p.pos)
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
