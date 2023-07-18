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
			buf = append(buf, v.Value)
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

			tracef("flag string=%[1]q", flStr)

			if len(v.Nodes) > 0 {
				flStr, tail, err := unParseFlagNodes(flStr, v.Nodes, cfg)
				if err != nil {
					return buf, err
				}

				tracef("appending %[1]q with tail=%[2]q", flStr, tail)

				buf = append(append(buf, flStr), tail...)
			} else {
				tracef("appending %[1]q", flStr)

				buf = append(buf, flStr)
			}

			continue
		default:
			return buf, fmt.Errorf("unhandled node type %[1]T: %[2]w", v, Err)
		}
	}

	return buf, nil
}

func unParseFlagNodes(flStr string, nodes []Node, cfg *ScannerConfig) (string, []string, error) {
	if len(nodes) == 0 {
		return flStr, []string{}, nil
	}

	if _, ok := nodes[0].(*ArgDelimiter); ok {
		tail, err := UnparseTree(nodes[1:], cfg)

		tracef("explicit arg delimiter present; returning flag str=%[1]q tail=%[2]q", flStr, tail)

		return flStr, tail, err
	}

	tail, err := UnparseTree(nodes, cfg)
	if err != nil {
		return flStr, tail, err
	}

	if _, ok := nodes[0].(*Assign); ok && len(nodes) > 1 {
		flStr = flStr + tail[0] + tail[1]
		tail = tail[2:]

		tracef("assign operator present with adjacent nodes; returning flag str=%[1]q tail=%[2]q", flStr, tail)

		return flStr, tail, nil
	} else if len(flStr) == 2 {
		flStr = flStr + tail[0]
		tail = tail[1:]

		tracef("short flag detected; returning flag str=%[1]q tail=%[2]q", flStr, tail)

		return flStr, tail, nil
	}

	tracef("no special cases; returning flag str=%[1]q tail=%[2]q", flStr, tail)

	return flStr, tail, nil
}
