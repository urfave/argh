package argh

import "fmt"

type Querier interface {
	Program() (Command, bool)
	TypedAST() []TypedNode
	AST() []Node
}

func NewQuerier(nodes []Node) Querier {
	return &defaultQuerier{nodes: nodes}
}

type defaultQuerier struct {
	nodes []Node
}

func (dq *defaultQuerier) Program() (Command, bool) {
	if len(dq.nodes) == 0 {
		return Command{}, false
	}

	v, ok := dq.nodes[0].(Command)
	return v, ok
}

func (dq *defaultQuerier) TypedAST() []TypedNode {
	ret := []TypedNode{}

	for _, node := range dq.nodes {
		if _, ok := node.(ArgDelimiter); ok {
			continue
		}

		if _, ok := node.(StopFlag); ok {
			continue
		}

		ret = append(
			ret,
			TypedNode{
				Type: fmt.Sprintf("%T", node),
				Node: node,
			},
		)
	}

	return ret
}

func (dq *defaultQuerier) AST() []Node {
	ret := []Node{}

	for _, node := range dq.nodes {
		if _, ok := node.(ArgDelimiter); ok {
			continue
		}

		if _, ok := node.(StopFlag); ok {
			continue
		}

		if v, ok := node.(CompoundShortFlag); ok {
			ret = append(ret, NewQuerier(v.Nodes).AST()...)

			continue
		}

		if v, ok := node.(Command); ok {
			ret = append(ret, NewQuerier(v.Nodes).AST()...)

			continue
		}

		ret = append(ret, node)
	}

	return ret
}
