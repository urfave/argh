package argh

import (
	"fmt"
	"strings"
)

// UnparseTree accepts a Node slice which is assumed to be a parse tree
// such as that returned from ParseArgs and a ScannerConfig,
// returning a string slice representation of the un-parsed input.
func UnparseTree(nodes []Node, cfg *ScannerConfig) ([]string, error) {
	buf := []string{}

	for i, node := range nodes {
		tracef("Unparse i=%[1]d node type=%[2]T", i, node)

		switch v := node.(type) {
		case *ArgDelimiter:
			continue
		case *Assign:
			buf = append(buf, string(cfg.AssignmentOperator))
			continue
		case *StdinFlag:
			buf = append(buf, string(cfg.FlagPrefix))
			continue
		case *StopFlag:
			buf = append(buf, string(cfg.FlagPrefix)+string(cfg.FlagPrefix))
			continue
		case *Ident:
			buf = append(buf, v.Literal)
			continue
		case *PassthroughArgs:
			sv, err := UnparseTree(v.Nodes, cfg)
			if err != nil {
				return buf, err
			}

			buf = append(buf, sv...)
			continue
		case *CompoundShortFlag:
			if v.Nodes != nil {
				flagStrings := []string{}

				sv, err := UnparseTree(v.Nodes, cfg)
				if err != nil {
					return buf, err
				}

				for _, flagString := range sv {
					flagStrings = append(flagStrings, strings.TrimPrefix(flagString, string(cfg.FlagPrefix)))
				}

				buf = append(buf, string(cfg.FlagPrefix)+strings.Join(flagStrings, ""))
			}

			continue
		case *MultiIdent:
			if v.Nodes != nil {
				sv, err := UnparseTree(v.Nodes, cfg)
				if err != nil {
					return buf, err
				}

				buf = append(buf, strings.Join(sv, string(cfg.MultiValueDelim)))
			}

			continue
		case *Command:
			buf = append(buf, v.Name)

			if len(v.Nodes) == 0 {
				continue
			}

			sv, err := UnparseTree(v.Nodes, cfg)
			if err != nil {
				return buf, err
			}

			buf = append(buf, sv...)
			continue
		case *Flag:
			prefix := string(cfg.FlagPrefix)
			if len(v.Name) > 1 {
				prefix += string(cfg.FlagPrefix)
			}

			flStr := prefix + v.Name

			if len(v.Nodes) > 0 {
				nodeStrings, err := UnparseTree(v.Nodes, cfg)
				if err != nil {
					return buf, err
				}

				tail := []string{}

				switch v.Nodes[0].(type) {
				case *Assign, *Ident, *MultiIdent, *StdinFlag:

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
		default:
			return buf, fmt.Errorf("unhandled node type %[1]T: %[2]w", v, Err)
		}
	}

	return buf, nil
}
