package argh

import "fmt"

type ParseTree struct {
	Nodes []Node `json:"nodes"`
}

func (pt *ParseTree) toAST() []TypedNode {
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

type Node interface{}

type TypedNode struct {
	Type string `json:"type"`
	Node Node   `json:"node"`
}

type Args struct {
	Pos   int    `json:"pos"`
	Nodes []Node `json:"nodes"`
}

type Statement struct {
	Pos   int    `json:"pos"`
	Nodes []Node `json:"nodes"`
}

type Program struct {
	Pos  int    `json:"pos"`
	Name string `json:"name"`
}

type Ident struct {
	Pos     int    `json:"pos"`
	Literal string `json:"literal"`
}

type Command struct {
	Pos   int    `json:"pos"`
	Name  string `json:"name"`
	Nodes []Node `json:"nodes"`
}

type Flag struct {
	Pos   int     `json:"pos"`
	Name  string  `json:"name"`
	Value *string `json:"value,omitempty"`
}

type StdinFlag struct {
	Pos int `json:"pos"`
}

type StopFlag struct {
	Pos int `json:"pos"`
}

type ArgDelimiter struct {
	Pos int `json:"pos"`
}
