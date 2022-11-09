package argh

// ToAST accepts a slice of nodes as expected from ParseArgs and
// returns an AST with parse-time artifacts dropped and reorganized
// where applicable.
func ToAST(parseTree []Node) []Node {
	ret := []Node{}

	for i, node := range parseTree {
		tracef("ToAST i=%d node type=%T", i, node)

		if _, ok := node.(*ArgDelimiter); ok {
			continue
		}

		if _, ok := node.(*StopFlag); ok {
			continue
		}

		if v, ok := node.(*CompoundShortFlag); ok {
			if v.Nodes != nil {
				ret = append(ret, ToAST(v.Nodes)...)
			}

			continue
		}

		if v, ok := node.(*CommandFlag); ok {
			astNodes := ToAST(v.Nodes)

			if len(astNodes) == 0 {
				astNodes = nil
			}

			ret = append(
				ret,
				&CommandFlag{
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
