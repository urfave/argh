package argh

import "fmt"

type Querier interface {
	Program() (Program, bool)
	TypedAST() []TypedNode
	AST() []Node
}

func NewQuerier(pt *ParseTree) Querier {
	return &defaultQuerier{pt: pt}
}

type defaultQuerier struct {
	pt *ParseTree
}

func (dq *defaultQuerier) Program() (Program, bool) {
	if len(dq.pt.Nodes) == 0 {
		return Program{}, false
	}

	v, ok := dq.pt.Nodes[0].(Program)
	return v, ok
}

func (dq *defaultQuerier) TypedAST() []TypedNode {
	ret := []TypedNode{}

	for _, node := range dq.pt.Nodes {
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

	for _, node := range dq.pt.Nodes {
		if _, ok := node.(ArgDelimiter); ok {
			continue
		}

		if _, ok := node.(StopFlag); ok {
			continue
		}

		if v, ok := node.(CompoundShortFlag); ok {
			for _, subNode := range v.Nodes {
				ret = append(ret, subNode)
			}

			continue
		}

		ret = append(ret, node)
	}

	return ret
}
