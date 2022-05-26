package argh

type Querier interface {
	Program() (Command, bool)
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

func (dq *defaultQuerier) AST() []Node {
	ret := []Node{}

	for i, node := range dq.nodes {
		tracef("AST i=%d node type=%T", i, node)

		if _, ok := node.(*ArgDelimiter); ok {
			continue
		}

		if _, ok := node.(*StopFlag); ok {
			continue
		}

		if v, ok := node.(*CompoundShortFlag); ok {
			if v.Nodes != nil {
				ret = append(ret, NewQuerier(v.Nodes).AST()...)
			}

			continue
		}

		if v, ok := node.(*Command); ok {
			astNodes := NewQuerier(v.Nodes).AST()

			if len(astNodes) == 0 {
				astNodes = nil
			}

			ret = append(
				ret,
				&Command{
					Name:   v.Name,
					Values: v.Values,
					Nodes:  astNodes,
				})

			continue
		}

		if v, ok := node.(*Flag); ok {
			astNodes := NewQuerier(v.Nodes).AST()

			if len(astNodes) == 0 {
				astNodes = nil
			}

			ret = append(
				ret,
				&Flag{
					Name:   v.Name,
					Values: v.Values,
					Nodes:  astNodes,
				})

			continue
		}

		ret = append(ret, node)
	}

	return ret
}
