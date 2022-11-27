package argh

import (
	"errors"
)

var (
	Error = errors.New("argh error")
)

type Argh interface {
	Parse([]string) error
	Prog() *CommandConfig
	AST() []Node
}

type defaultArgh struct {
	pc ParserConfig
	pt ParseTree
}

func New(opts ...ParserOption) Argh {
	return &defaultArgh{
		pc: *NewParserConfig(opts...),
	}
}

func (a *defaultArgh) Prog() *CommandConfig {
	return a.pc.Prog
}

func (a *defaultArgh) Parse(args []string) error {
	pt, err := ParseArgs(args, &a.pc)
	if err != nil {
		return err
	}

	a.pt = *pt

	return nil
}

func (a *defaultArgh) AST() []Node {
	return ToAST(a.pt.Nodes)
}
