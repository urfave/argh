package argh

import (
	"strings"
)

func Unparse(nodes []Node, cfg *ScannerConfig) []string {
	buf := []string{}

	for i, node := range nodes {
		tracef("Unparse i=%[1]d node type=%[2]T", i, node)

		switch v := node.(type) {
		case *ArgDelimiter:
			buf = append(buf, "\x00")
			continue
		case *StopFlag:
			buf = append(buf, string(cfg.FlagPrefix)+string(cfg.FlagPrefix))
			continue
		case *Ident:
			buf = append(buf, v.Literal)
			continue
		case *CompoundShortFlag:
			if v.Nodes != nil {
				flagStrings := []string{}

				for _, flagString := range Unparse(v.Nodes, cfg) {
					flagStrings = append(flagStrings, strings.TrimPrefix(flagString, string(cfg.FlagPrefix)))
				}

				buf = append(buf, string(cfg.FlagPrefix)+strings.Join(flagStrings, ""))
			}

			continue
		case *PassthroughArgs:
			if v.Nodes != nil {
				buf = append(buf, string(cfg.FlagPrefix))
				buf = append(buf, Unparse(v.Nodes, cfg)...)
			}

			continue
		case *Flag:
			prefix := string(cfg.FlagPrefix)
			if len(v.Name) > 1 {
				prefix += string(cfg.FlagPrefix)
			}

			flStr := prefix + v.Name

			if len(v.Values) > 0 {
				flVal := []string{}

				for _, sv := range stringMapToSlice(v.Values) {
					flVal = append(flVal, sv)
				}

				flStr += string(cfg.AssignmentOperator) + strings.Join(flVal, string(cfg.MultiValueDelim))
			}

			buf = append(buf, flStr)
			continue
		case *Command:
			buf = append(buf, v.Name)

			if len(v.Nodes) == 0 {
				continue
			}

			buf = append(buf, Unparse(v.Nodes, cfg)...)
			continue
		}
	}

	return buf
}
