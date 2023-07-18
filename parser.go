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

	prog := &Command{Name: p.lit}
	if err := p.parseCommandWithValue(prog, p.cfg.Prog); err != nil {
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

func (p *parser) parseIdent() Node {
	node := &Ident{Value: p.lit}
	return node
}

func (p *parser) parseFlag(flags *Flags) (Node, error) {
	switch p.tok {
	case SHORT_FLAG:
		tracef("parsing short flag with flags=%+#[1]v", flags)

		return p.parseShortFlag(flags)
	case LONG_FLAG:
		tracef("parsing long flag with flags=%+#[1]v", flags)

		return p.parseLongFlag(flags)
	case COMPOUND_SHORT_FLAG:
		tracef("parsing compound short flag with flags=%+#[1]v", flags)

		return p.parseCompoundShortFlag(flags)
	}

	// TODO: don't panic
	panic(fmt.Sprintf("token %v cannot be parsed as flag", p.tok))
}

func (p *parser) parseShortFlag(flags *Flags) (Node, error) {
	node := &Flag{Name: string(p.lit[1])}

	flCfg, ok := flags.Get(node.Name)
	if !ok {
		errMsg := fmt.Sprintf("unknown flag %[1]q", node.Name)
		p.addError(errMsg)

		return node, &FlagError{
			Pos:  Position{Column: int(p.pos)},
			Node: *node,
			Msg:  errMsg,
		}
	}

	return node, newParseFlagContext(node, flCfg, nil, p.cfg.ScannerConfig).parse(p)
}

func (p *parser) parseLongFlag(flags *Flags) (Node, error) {
	node := &Flag{Name: string(p.lit[2:])}

	flCfg, ok := flags.Get(node.Name)
	if !ok {
		errMsg := fmt.Sprintf("unknown flag %[1]q", node.Name)
		p.addError(errMsg)

		return node, &FlagError{
			Pos:  Position{Column: int(p.pos)},
			Node: *node,
			Msg:  errMsg,
		}
	}

	return node, newParseFlagContext(node, flCfg, nil, p.cfg.ScannerConfig).parse(p)
}

func (p *parser) parseCompoundShortFlag(flags *Flags) (Node, error) {
	unparsedFlags := []*Flag{}
	unparsedFlagConfigs := []FlagConfig{}

	withoutFlagPrefix := p.lit[1:]

	for _, r := range withoutFlagPrefix {
		node := &Flag{Name: string(r)}

		flCfg, ok := flags.Get(node.Name)
		if !ok {
			errMsg := fmt.Sprintf("unknown flag %[1]q", node.Name)
			p.addError(errMsg)

			return node, &FlagError{
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

			if err := newParseFlagContext(node, flCfg, zeroValuePtr, p.cfg.ScannerConfig).parse(p); err != nil {
				return nil, err
			}

			flagNodes = append(flagNodes, node)

			continue
		}

		if err := newParseFlagContext(node, flCfg, nil, p.cfg.ScannerConfig).parse(p); err != nil {
			return nil, err
		}

		flagNodes = append(flagNodes, node)
	}

	return &CompoundShortFlag{Nodes: flagNodes}, nil
}

func (p *parser) parseCommandWithValue(
	cmd *Command,
	cCfg *CommandConfig,
) error {
	tracef("cCfg=%+#[1]v (cmd=%[2]q)", cCfg, cmd.Name)

	values := []KeyValue{}
	nodes := []Node{}
	identIndex := 0

	tailReturn := func() error {
		if len(nodes) > 0 {
			tracef("setting non-empty nodes (len=%[1]d) (cmd=%[2]q)", len(nodes), cmd.Name)

			cmd.Nodes = nodes
		}

		if len(values) > 0 {
			tracef("setting non-empty values (len=%[1]d) (cmd=%[2]q)", len(values), cmd.Name)

			cmd.Values = values
		}

		if cCfg.On != nil {
			tracef("calling command config handler (cmd=%[1]q)", cmd.Name)

			if err := cCfg.On(*cmd); err != nil {
				return err
			}
		} else {
			tracef("no command config handler (cmd=%[1]q)", cmd.Name)
		}

		tracef("returning nil err (cmd=%[1]q)", cmd.Name)
		return nil
	}

	for i := 0; p.tok != EOL; i++ {
		if !p.buffered {
			tracef("buffered=false; scanning next (cmd=%[1]q)", cmd.Name)
			p.next()
		}

		tracef("setting buffered=false (cmd=%[1]q)", cmd.Name)

		p.buffered = false

		tracef("for=%[1]d values=%+#[2]v (cmd=%[3]q)", i, values, cmd.Name)
		tracef("for=%[1]d nodes=%+#[2]v (cmd=%[3]q)", i, nodes, cmd.Name)
		tracef("for=%[1]d tok=%[2]s lit=%[3]q pos=%[4]v (cmd=%[5]q)", i, p.tok, p.lit, p.pos, cmd.Name)

		tracef("cCfg=%+#[1]v (cmd=%[2]q)", cCfg, cmd.Name)

		if subCfg, ok := cCfg.GetCommandConfig(p.lit); ok {
			subCommand := p.lit
			subNode := &Command{Name: p.lit}

			if err := p.parseCommandWithValue(subNode, &subCfg); err != nil {
				return err
			}

			nodes = append(nodes, subNode)

			tracef("breaking after sub-command=%[1]q (cmd=%[2]q)", subCommand, cmd.Name)
			break
		}

		tracef("handling tok=%[1]s (cmd=%[2]q)", p.tok, cmd.Name)

		switch p.tok {
		case ARG_DELIMITER:
			nodes = append(nodes, &ArgDelimiter{})

			continue
		case IDENT, STDIN_FLAG:
			if cCfg.NValue.Contains(identIndex) {
				name := fmt.Sprintf("%d", identIndex)

				tracef("checking for name of identIndex=%[1]d (cmd=%[2]q)", identIndex, cmd.Name)

				if cCfg.HasValueNameForIndex(identIndex) {
					name = cCfg.ValueNames[identIndex]

					tracef("setting name=%[1]q from config value names (cmd=%[2]q)", name, cmd.Name)
				} else if cCfg.HasRepeatingValueName() {
					name = cCfg.ValueNames[0]

					tracef("setting name=%[1]q from repeating value name (cmd=%[2]q)", name, cmd.Name)
				}

				values = append(values, KeyValue{Key: name, Value: p.lit})
			}

			if p.tok == STDIN_FLAG {
				nodes = append(nodes, &StdinFlag{})
			} else {
				nodes = append(nodes, &Ident{Value: p.lit})
			}

			identIndex++
		case LONG_FLAG, SHORT_FLAG, COMPOUND_SHORT_FLAG:
			tok := p.tok

			flagNode, err := p.parseFlag(cCfg.Flags)
			if err != nil {
				return err
			}

			tracef("appending %[1]s node=%+#[2]v (cmd=%[3]q)", tok, flagNode, cmd.Name)

			nodes = append(nodes, flagNode)
		case ASSIGN:
			tracef("error on bare tok=%[1]s (cmd=%[2]q)", p.tok, cmd.Name)

			p.addError("invalid bare assignment")

			break
		default:
			tracef("breaking on tok=%[1]s (cmd=%[2]q)", p.tok, cmd.Name)

			return tailReturn()
		}
	}

	return tailReturn()
}

func (p *parser) parsePassthrough() Node {
	nodes := []Node{}

	for ; p.tok != EOL; p.next() {
		nodes = append(nodes, &Ident{Value: p.lit})
	}

	if len(nodes) == 0 {
		return nil
	}

	return &PassthroughArgs{Nodes: nodes}
}

func newParseFlagContext(
	fl *Flag,
	flCfg FlagConfig,
	nValueOverride *NValue,
	sCfg *ScannerConfig,
) *parseFlagContext {
	return &parseFlagContext{
		fl:             fl,
		flCfg:          flCfg,
		nValueOverride: nValueOverride,
		values:         []KeyValue{},
		nodes:          []Node{},
		sCfg:           sCfg,
	}
}

type parseFlagContext struct {
	fl             *Flag
	flCfg          FlagConfig
	nValueOverride *NValue
	values         []KeyValue
	nodes          []Node
	identIndex     int
	prev           Node
	sCfg           *ScannerConfig
}

func (pfc *parseFlagContext) parse(p *parser) error {
	pfc.values = []KeyValue{}
	pfc.nodes = []Node{}
	pfc.identIndex = 0
	pfc.prev = nil

	for i := 0; p.tok != EOL; i++ {
		if err, brk := pfc.parseToken(i, p); err != nil || brk {
			return err
		}
	}

	return pfc.finalize()
}

func (pfc *parseFlagContext) parseToken(i int, p *parser) (error, bool) {
	if len(pfc.nodes) > 0 {
		pfc.prev = pfc.nodes[len(pfc.nodes)-1]
	}

	if pfc.nValueOverride != nil && !(*pfc.nValueOverride).Contains(pfc.identIndex) {
		tracef("identIndex=%[1]d exceeds expected=%[2]v; breaking (flag=%[3]q)", pfc.identIndex, *pfc.nValueOverride, pfc.fl.Name)

		return pfc.finalize(), true
	}

	if !pfc.flCfg.NValue.Contains(pfc.identIndex) {
		tracef("identIndex=%[1]d exceeds expected=%[2]v; breaking (flag=%[3]q)", pfc.identIndex, pfc.flCfg.NValue, pfc.fl.Name)

		return pfc.finalize(), true
	}

	p.next()

	tracef("for=%[1]d values=%+#[2]v (flag=%[3]q)", i, pfc.values, pfc.fl.Name)
	tracef("for=%[1]d nodes=%+#[2]v (flag=%[3]q)", i, pfc.nodes, pfc.fl.Name)
	tracef("for=%[1]d tok=%[2]s lit=%[3]q pos=%[4]v (flag=%[5]q)", i, p.tok, p.lit, p.pos, pfc.fl.Name)

	switch p.tok {
	case ARG_DELIMITER:
		pfc.nodes = append(pfc.nodes, &ArgDelimiter{})

		return nil, false
	case ASSIGN:
		return pfc.handleAssign()
	case IDENT, STDIN_FLAG, MULTI_VALUE_DELIMITER:
		return pfc.handleIdentStdinMulti(p.tok, p.lit)
	default:
		tracef("breaking on tok=%[1]s %[2]q %[3]v; setting buffered=true (flag=%[4]q)", p.tok, p.lit, p.pos, pfc.fl.Name)

		p.buffered = true

		return pfc.finalize(), true
	}

	return nil, false
}

func (pfc *parseFlagContext) setPrev(node Node) {
	pfc.nodes[len(pfc.nodes)-1] = node
}

func (pfc *parseFlagContext) addNode(node Node) {
	if pfc.prev == nil {
		tracef("no previous node; appending node type=%[1]T (flag=%[2]q)", node, pfc.fl.Name)

		pfc.nodes = append(pfc.nodes, node)

		return
	}

	if pv, ok := pfc.prev.(*MultiIdent); ok {
		tracef("appending node type=%[1]T to previous node nodes (flag=%[2]q)", node, pfc.fl.Name)

		pv.Nodes = append(pv.Nodes, node)

		return
	}

	pv, ok := pfc.prev.(*KeyValue)

	if !ok {
		tracef("appending node type=%[1]T to nodes (flag=%[2]q)", node, pfc.fl.Name)

		pfc.nodes = append(pfc.nodes, node)

		return
	}

	value := ""
	if v, ok := node.(*Ident); ok {
		value = v.Value
	} else if _, ok := node.(*StdinFlag); ok {
		value = string(pfc.sCfg.FlagPrefix)
	}

	tracef("setting node value=%[1]q as previous node (key=%[2]q) value (flag=%[3]q)", value, pv.Key, pfc.fl.Name)

	pv.Value = value

	pfc.setPrev(pv)
	pfc.values = append(pfc.values, *pv)
}

func (pfc *parseFlagContext) handleAssign() (error, bool) {
	tracef("encountered assignment (flag=%[1]q)", pfc.fl.Name)

	if pfc.prev == nil {
		tracef("appending assignment node (flag=%[1]q)", pfc.fl.Name)

		pfc.nodes = append(pfc.nodes, &Assign{})

		return nil, false
	}

	tracef("checking for key-value assignment (flag=%[1]q)", pfc.fl.Name)

	if v, ok := pfc.prev.(*Ident); ok {
		tracef("setting previous node as *KeyValue with *Ident literal=%[1]q as key (flag=%[2]q)", v.Value, pfc.fl.Name)

		pfc.setPrev(&KeyValue{Key: v.Value})

		return nil, false
	}

	pfc.nodes = append(pfc.nodes, &Assign{})

	return nil, false
}

func (pfc *parseFlagContext) handleIdentStdinMulti(tok Token, lit string) (error, bool) {
	name := fmt.Sprintf("%d", pfc.identIndex)

	tracef("checking for name of identIndex=%[1]d (flag=%[2]q)", pfc.identIndex, pfc.fl.Name)

	if pfc.flCfg.HasValueNameForIndex(pfc.identIndex) {
		name = pfc.flCfg.ValueNames[pfc.identIndex]

		tracef("setting name=%[1]q from config value names (flag=%[2]q)", name, pfc.fl.Name)
	} else if pfc.flCfg.HasRepeatingValueName() {
		name = pfc.flCfg.ValueNames[0]

		tracef("setting name=%[1]q from repeating value name (flag=%[2]q)", name, pfc.fl.Name)
	} else {
		tracef("setting name=%[1]q (flag=%[2]q)", name, pfc.fl.Name)
	}

	if tok == STDIN_FLAG {
		pfc.addNode(&StdinFlag{})
		pfc.identIndex++

		return nil, false
	}

	if tok != MULTI_VALUE_DELIMITER {
		tracef("appending *Ident node (flag=%[1]q)", pfc.fl.Name)

		pfc.addNode(&Ident{Value: lit})
		pfc.identIndex++

		kv := KeyValue{Key: name, Value: lit}
		tracef("appending %+[1]v to values (flag=%[2]q)", kv, pfc.fl.Name)

		pfc.values = append(pfc.values, kv)

		return nil, false
	}

	if pfc.prev == nil {
		tracef("appending *MultiIdent with empty child nodes (flag=%[1]q)", pfc.fl.Name)

		pfc.nodes = append(pfc.nodes, &MultiIdent{Nodes: []Node{}})

		return nil, false
	}

	if v, ok := pfc.prev.(*Ident); ok {
		tracef("setting previous node as *MultiIdent with *Ident node as first child (flag=%[1]q)", pfc.fl.Name)

		pfc.setPrev(&MultiIdent{Nodes: []Node{v}})
	} else if v, ok := pfc.prev.(*StdinFlag); ok {
		tracef("setting previous node as *MultiIdent with *StdinFlag node as first child (flag=%[1]q)", pfc.fl.Name)

		pfc.setPrev(&MultiIdent{Nodes: []Node{v}})
	}

	return nil, false
}

func (pfc *parseFlagContext) finalize() error {
	if len(pfc.nodes) > 0 {
		tracef("setting non-empty nodes (len=%[1]d) (flag=%[2]q)", len(pfc.nodes), pfc.fl.Name)

		pfc.fl.Nodes = pfc.nodes
	}

	if len(pfc.values) > 0 {
		tracef("setting non-empty values (len=%[1]d) (flag=%[2]q)", len(pfc.values), pfc.fl.Name)

		pfc.fl.Values = pfc.values
	}

	if pfc.flCfg.On == nil {
		tracef("no flag config handler (flag=%[1]q)", pfc.fl.Name)

		return nil
	}

	tracef("calling flag config handler (flag=%[1]q)", pfc.fl.Name)

	return pfc.flCfg.On(*pfc.fl)
}
