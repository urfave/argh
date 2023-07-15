package argh

// ToAST accepts a slice of nodes as expected from ParseArgs and
// returns an AST with parse-time artifacts dropped and reorganized
// where applicable.
func ToAST(parseTree []Node) []Node {
	ret := []Node{}

	for i, node := range parseTree {
		tracef("ToAST i=%d node type=%T", i, node)

		switch v := node.(type) {
		case *ArgDelimiter:
		case *StopFlag:
			continue
		case *Assign:
			ret = append(ret, v)

			continue

		case *CompoundShortFlag:
			if v.Nodes != nil {
				ret = append(ret, ToAST(v.Nodes)...)
			}

			continue
		case *Flag:
			astNodes := ToAST(v.Nodes)

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
		case *Command:
			astNodes := ToAST(v.Nodes)

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

		ret = append(ret, node)
	}

	return ret
}
