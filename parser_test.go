package argh_test

import (
	"testing"

	"git.meatballhat.com/x/box-o-sand/argh"
	"github.com/stretchr/testify/assert"
)

func ptr[T any](v T) *T {
	return &v
}

func TestParser(t *testing.T) {
	for _, tc := range []struct {
		name   string
		args   []string
		cfg    *argh.ParserConfig
		expPT  []argh.Node
		expAST []argh.Node
		expErr error
		skip   bool
	}{
		{
			name: "bare",
			args: []string{"pizzas"},
			expPT: []argh.Node{
				argh.Program{Name: "pizzas"},
			},
			expAST: []argh.Node{
				argh.Program{Name: "pizzas"},
			},
		},
		{
			name: "long value-less flags",
			args: []string{"pizzas", "--tasty", "--fresh", "--super-hot-right-now"},
			expPT: []argh.Node{
				argh.Program{Name: "pizzas"},
				argh.ArgDelimiter{},
				argh.Flag{Name: "tasty"},
				argh.ArgDelimiter{},
				argh.Flag{Name: "fresh"},
				argh.ArgDelimiter{},
				argh.Flag{Name: "super-hot-right-now"},
			},
			expAST: []argh.Node{
				argh.Program{Name: "pizzas"},
				argh.Flag{Name: "tasty"},
				argh.Flag{Name: "fresh"},
				argh.Flag{Name: "super-hot-right-now"},
			},
		},
		{
			name: "long flags mixed",
			args: []string{
				"pizzas",
				"--tasty",
				"--fresh", "soon",
				"--super-hot-right-now",
				"--box", "square", "shaped", "hot",
				"--please",
			},
			cfg: &argh.ParserConfig{
				Commands: map[string]argh.NValue{},
				Flags: map[string]argh.NValue{
					"fresh": argh.OneValue,
					"box":   argh.OneOrMoreValue,
				},
			},
			expPT: []argh.Node{
				argh.Program{Name: "pizzas"},
				argh.ArgDelimiter{},
				argh.Flag{Name: "tasty"},
				argh.ArgDelimiter{},
				argh.Flag{Name: "fresh", Values: []string{"soon"}},
				argh.ArgDelimiter{},
				argh.Flag{Name: "super-hot-right-now"},
				argh.ArgDelimiter{},
				argh.Flag{Name: "box", Values: []string{"square", "shaped", "hot"}},
				argh.ArgDelimiter{},
				argh.Flag{Name: "please"},
			},
			expAST: []argh.Node{
				argh.Program{Name: "pizzas"},
				argh.Flag{Name: "tasty"},
				argh.Flag{Name: "fresh", Values: []string{"soon"}},
				argh.Flag{Name: "super-hot-right-now"},
				argh.Flag{Name: "box", Values: []string{"square", "shaped", "hot"}},
				argh.Flag{Name: "please"},
			},
		},
		{
			name: "short value-less flags",
			args: []string{"pizzas", "-t", "-f", "-s"},
			expPT: []argh.Node{
				argh.Program{Name: "pizzas"},
				argh.ArgDelimiter{},
				argh.Flag{Name: "t"},
				argh.ArgDelimiter{},
				argh.Flag{Name: "f"},
				argh.ArgDelimiter{},
				argh.Flag{Name: "s"},
			},
			expAST: []argh.Node{
				argh.Program{Name: "pizzas"},
				argh.Flag{Name: "t"},
				argh.Flag{Name: "f"},
				argh.Flag{Name: "s"},
			},
		},
		{
			name: "compound short flags",
			args: []string{"pizzas", "-aca", "-blol"},
			expPT: []argh.Node{
				argh.Program{Name: "pizzas"},
				argh.ArgDelimiter{},
				argh.CompoundShortFlag{
					Nodes: []argh.Node{
						argh.Flag{Name: "a"},
						argh.Flag{Name: "c"},
						argh.Flag{Name: "a"},
					},
				},
				argh.ArgDelimiter{},
				argh.CompoundShortFlag{
					Nodes: []argh.Node{
						argh.Flag{Name: "b"},
						argh.Flag{Name: "l"},
						argh.Flag{Name: "o"},
						argh.Flag{Name: "l"},
					},
				},
			},
			expAST: []argh.Node{
				argh.Program{Name: "pizzas"},
				argh.Flag{Name: "a"},
				argh.Flag{Name: "c"},
				argh.Flag{Name: "a"},
				argh.Flag{Name: "b"},
				argh.Flag{Name: "l"},
				argh.Flag{Name: "o"},
				argh.Flag{Name: "l"},
			},
		},
		{
			name: "mixed long short value flags",
			args: []string{"pizzas", "-a", "--ca", "-b", "1312", "-lol"},
			cfg: &argh.ParserConfig{
				Commands: map[string]argh.NValue{},
				Flags:    map[string]argh.NValue{"b": argh.OneValue},
			},
			expPT: []argh.Node{
				argh.Program{Name: "pizzas"},
				argh.ArgDelimiter{},
				argh.Flag{Name: "a"},
				argh.ArgDelimiter{},
				argh.Flag{Name: "ca"},
				argh.ArgDelimiter{},
				argh.Flag{Name: "b", Values: []string{"1312"}},
				argh.ArgDelimiter{},
				argh.CompoundShortFlag{
					Nodes: []argh.Node{
						argh.Flag{Name: "l"},
						argh.Flag{Name: "o"},
						argh.Flag{Name: "l"},
					},
				},
			},
			expAST: []argh.Node{
				argh.Program{Name: "pizzas"},
				argh.Flag{Name: "a"},
				argh.Flag{Name: "ca"},
				argh.Flag{Name: "b", Values: []string{"1312"}},
				argh.Flag{Name: "l"},
				argh.Flag{Name: "o"},
				argh.Flag{Name: "l"},
			},
		},
		{
			name: "commands",
			args: []string{"pizzas", "fly", "fry"},
			cfg: &argh.ParserConfig{
				Commands: map[string]argh.NValue{"fly": argh.ZeroValue, "fry": argh.ZeroValue},
				Flags:    map[string]argh.NValue{},
			},
			expPT: []argh.Node{
				argh.Program{Name: "pizzas"},
				argh.ArgDelimiter{},
				argh.Command{Name: "fly"},
				argh.ArgDelimiter{},
				argh.Command{Name: "fry"},
			},
		},
	} {
		if tc.skip {
			continue
		}

		if tc.expPT != nil {
			t.Run(tc.name+" parse tree", func(ct *testing.T) {
				actual, err := argh.ParseArgs(tc.args, tc.cfg)
				if err != nil {
					assert.ErrorIs(ct, err, tc.expErr)
					return
				}

				assert.Equal(ct, tc.expPT, actual.ParseTree.Nodes)
			})
		}

		if tc.expAST != nil {
			t.Run(tc.name+" ast", func(ct *testing.T) {
				actual, err := argh.ParseArgs(tc.args, tc.cfg)
				if err != nil {
					assert.ErrorIs(ct, err, tc.expErr)
					return
				}

				assert.Equal(ct, tc.expAST, actual.AST())
			})
		}
	}
}
