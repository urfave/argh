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

		if v, ok := node.(Statement); ok {
			for _, subNode := range v.Nodes {
				ret = append(ret, subNode)
			}

			continue
		}

		ret = append(ret, node)
	}

	return ret
}

type Node interface{}

type TypedNode struct {
	Type string `json:"type"`
	Node Node   `json:"node"`
}

type Args struct {
	Nodes []Node `json:"nodes"`
}

type Statement struct {
	Nodes []Node `json:"nodes"`
}

type Program struct {
	Name string `json:"name"`
}

type Ident struct {
	Literal string `json:"literal"`
}

type Command struct {
	Name string `json:"name"`
}

type Flag struct {
	Name  string  `json:"name"`
	Value *string `json:"value,omitempty"`
}

type StdinFlag struct{}

type StopFlag struct{}

type ArgDelimiter struct{}
