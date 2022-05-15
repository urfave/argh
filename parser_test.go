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
			args: []string{"pizzas", "--tasty", "--fresh", "soon", "--super-hot-right-now"},
			cfg: &argh.ParserConfig{
				Commands:   []string{},
				ValueFlags: []string{"fresh"},
			},
			expPT: []argh.Node{
				argh.Program{Name: "pizzas"},
				argh.ArgDelimiter{},
				argh.Flag{Name: "tasty"},
				argh.ArgDelimiter{},
				argh.Flag{Name: "fresh", Value: ptr("soon")},
				argh.ArgDelimiter{},
				argh.Flag{Name: "super-hot-right-now"},
			},
			expAST: []argh.Node{
				argh.Program{Name: "pizzas"},
				argh.Flag{Name: "tasty"},
				argh.Flag{Name: "fresh", Value: ptr("soon")},
				argh.Flag{Name: "super-hot-right-now"},
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
				argh.Statement{
					Nodes: []argh.Node{
						argh.Flag{Name: "a"},
						argh.Flag{Name: "c"},
						argh.Flag{Name: "a"},
					},
				},
				argh.ArgDelimiter{},
				argh.Statement{
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
				Commands:   []string{},
				ValueFlags: []string{"b"},
			},
			expPT: []argh.Node{
				argh.Program{Name: "pizzas"},
				argh.ArgDelimiter{},
				argh.Flag{Name: "a"},
				argh.ArgDelimiter{},
				argh.Flag{Name: "ca"},
				argh.ArgDelimiter{},
				argh.Flag{Name: "b", Value: ptr("1312")},
				argh.ArgDelimiter{},
				argh.Statement{
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
				argh.Flag{Name: "b", Value: ptr("1312")},
				argh.Flag{Name: "l"},
				argh.Flag{Name: "o"},
				argh.Flag{Name: "l"},
			},
		},
		{
			name: "commands",
			args: []string{"pizzas", "fly", "fry"},
			cfg: &argh.ParserConfig{
				Commands:   []string{"fly", "fry"},
				ValueFlags: []string{},
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
