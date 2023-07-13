package argh

import (
	"strings"
)

// UnparseTree accepts a Node slice which is assumed to be a parse tree
// such as that returned from ParseArgs and a ScannerConfig,
// returning a string slice representation of the un-parsed input.
func UnparseTree(nodes []Node, cfg *ScannerConfig) []string {
	buf := []string{}

	for i, node := range nodes {
		tracef("Unparse i=%[1]d node type=%[2]T", i, node)

		switch v := node.(type) {
		case *ArgDelimiter:
			continue
		case *StopFlag:
			buf = append(buf, string(cfg.FlagPrefix)+string(cfg.FlagPrefix))
			continue
		case *Ident:
			buf = append(buf, v.Literal)
			continue
		case *Assign:
			buf = append(buf, string(cfg.AssignmentOperator))
			continue
		case *CompoundShortFlag:
			if v.Nodes != nil {
				flagStrings := []string{}

				for _, flagString := range UnparseTree(v.Nodes, cfg) {
					flagStrings = append(flagStrings, strings.TrimPrefix(flagString, string(cfg.FlagPrefix)))
				}

				buf = append(buf, string(cfg.FlagPrefix)+strings.Join(flagStrings, ""))
			}

			continue
		case *MultiIdent:
			if v.Nodes != nil {
				buf = append(buf, strings.Join(UnparseTree(v.Nodes, cfg), string(cfg.MultiValueDelim)))
			}

			continue
		case *Flag:
			prefix := string(cfg.FlagPrefix)
			if len(v.Name) > 1 {
				prefix += string(cfg.FlagPrefix)
			}

			flStr := prefix + v.Name

			if len(v.Nodes) > 0 {
				nodeStrings := UnparseTree(v.Nodes, cfg)
				tail := []string{}

				if _, ok := v.Nodes[0].(*Assign); ok {
					flStr += nodeStrings[0]
					tail = nodeStrings[1:]

					if len(nodeStrings) > 1 {
						flStr += nodeStrings[1]
						tail = nodeStrings[2:]
					}
				}

				buf = append(append(buf, flStr), tail...)
			} else {
				buf = append(buf, flStr)
			}

			continue
		case *Command:
			buf = append(buf, v.Name)

			if len(v.Nodes) == 0 {
				continue
			}

			buf = append(buf, UnparseTree(v.Nodes, cfg)...)
			continue
		}
	}

	return buf
}
