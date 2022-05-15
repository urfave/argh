package argh

import "fmt"

type ParseTree struct {
	Nodes []Node `json:"nodes"`
}

func (pt *ParseTree) typedAST() []TypedNode {
	ret := []TypedNode{}

	for _, node := range pt.Nodes {
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

func (pt *ParseTree) ast() []Node {
	ret := []Node{}

	for _, node := range pt.Nodes {
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
