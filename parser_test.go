package argh_test

import (
	"testing"

	"git.meatballhat.com/x/box-o-sand/argh"
	"github.com/stretchr/testify/assert"
)

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
				argh.Command{Name: "pizzas"},
			},
			expAST: []argh.Node{
				argh.Command{Name: "pizzas"},
			},
		},
		{
			name: "one positional arg",
			args: []string{"pizzas", "excel"},
			cfg: &argh.ParserConfig{
				Prog: argh.CommandConfig{NValue: 1},
			},
			expPT: []argh.Node{
				argh.Command{Name: "pizzas", Values: map[string]string{"0": "excel"}},
			},
			expAST: []argh.Node{
				argh.Command{Name: "pizzas", Values: map[string]string{"0": "excel"}},
			},
		},
		{
			name: "many positional args",
			args: []string{"pizzas", "excel", "wildly", "when", "feral"},
			cfg: &argh.ParserConfig{
				Prog: argh.CommandConfig{NValue: argh.OneOrMoreValue},
			},
			expPT: []argh.Node{
				argh.Command{Name: "pizzas", Values: map[string]string{"0": "excel", "1": "wildly", "2": "when", "3": "feral"}},
			},
			expAST: []argh.Node{
				argh.Command{Name: "pizzas", Values: map[string]string{"0": "excel", "1": "wildly", "2": "when", "3": "feral"}},
			},
		},
		{
			name: "long value-less flags",
			args: []string{"pizzas", "--tasty", "--fresh", "--super-hot-right-now"},
			expPT: []argh.Node{
				argh.Command{Name: "pizzas"},
				argh.ArgDelimiter{},
				argh.Flag{Name: "tasty"},
				argh.ArgDelimiter{},
				argh.Flag{Name: "fresh"},
				argh.ArgDelimiter{},
				argh.Flag{Name: "super-hot-right-now"},
			},
			expAST: []argh.Node{
				argh.Command{Name: "pizzas"},
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
				Prog: argh.CommandConfig{
					Commands: map[string]argh.CommandConfig{},
					Flags: map[string]argh.FlagConfig{
						"fresh": argh.FlagConfig{NValue: 1},
						"box":   argh.FlagConfig{NValue: argh.OneOrMoreValue},
					},
				},
			},
			expPT: []argh.Node{
				argh.Command{Name: "pizzas"},
				argh.ArgDelimiter{},
				argh.Flag{Name: "tasty"},
				argh.ArgDelimiter{},
				argh.Flag{Name: "fresh", Values: map[string]string{"0": "soon"}},
				argh.ArgDelimiter{},
				argh.Flag{Name: "super-hot-right-now"},
				argh.ArgDelimiter{},
				argh.Flag{Name: "box", Values: map[string]string{"0": "square", "1": "shaped", "2": "hot"}},
				argh.ArgDelimiter{},
				argh.Flag{Name: "please"},
			},
			expAST: []argh.Node{
				argh.Command{Name: "pizzas"},
				argh.Flag{Name: "tasty"},
				argh.Flag{Name: "fresh", Values: map[string]string{"0": "soon"}},
				argh.Flag{Name: "super-hot-right-now"},
				argh.Flag{Name: "box", Values: map[string]string{"0": "square", "1": "shaped", "2": "hot"}},
				argh.Flag{Name: "please"},
			},
		},
		{
			name: "short value-less flags",
			args: []string{"pizzas", "-t", "-f", "-s"},
			expPT: []argh.Node{
				argh.Command{Name: "pizzas"},
				argh.ArgDelimiter{},
				argh.Flag{Name: "t"},
				argh.ArgDelimiter{},
				argh.Flag{Name: "f"},
				argh.ArgDelimiter{},
				argh.Flag{Name: "s"},
			},
			expAST: []argh.Node{
				argh.Command{Name: "pizzas"},
				argh.Flag{Name: "t"},
				argh.Flag{Name: "f"},
				argh.Flag{Name: "s"},
			},
		},
		{
			name: "compound short flags",
			args: []string{"pizzas", "-aca", "-blol"},
			expPT: []argh.Node{
				argh.Command{Name: "pizzas"},
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
				argh.Command{Name: "pizzas"},
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
				Prog: argh.CommandConfig{
					Commands: map[string]argh.CommandConfig{},
					Flags: map[string]argh.FlagConfig{
						"b": argh.FlagConfig{NValue: 1},
					},
				},
			},
			expPT: []argh.Node{
				argh.Command{Name: "pizzas"},
				argh.ArgDelimiter{},
				argh.Flag{Name: "a"},
				argh.ArgDelimiter{},
				argh.Flag{Name: "ca"},
				argh.ArgDelimiter{},
				argh.Flag{Name: "b", Values: map[string]string{"0": "1312"}},
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
				argh.Command{Name: "pizzas"},
				argh.Flag{Name: "a"},
				argh.Flag{Name: "ca"},
				argh.Flag{Name: "b", Values: map[string]string{"0": "1312"}},
				argh.Flag{Name: "l"},
				argh.Flag{Name: "o"},
				argh.Flag{Name: "l"},
			},
		},
		{
			name: "commands",
			args: []string{"pizzas", "fly", "fry"},
			cfg: &argh.ParserConfig{
				Prog: argh.CommandConfig{
					Commands: map[string]argh.CommandConfig{
						"fly": argh.CommandConfig{},
						"fry": argh.CommandConfig{},
					},
					Flags: map[string]argh.FlagConfig{},
				},
			},
			expPT: []argh.Node{
				argh.Command{Name: "pizzas"},
				argh.ArgDelimiter{},
				argh.Command{Name: "fly"},
				argh.ArgDelimiter{},
				argh.Command{Name: "fry"},
			},
		},
		{
			name: "command specific flags",
			args: []string{"pizzas", "fly", "--freely", "fry", "--deeply", "-wAt"},
			cfg: &argh.ParserConfig{
				Prog: argh.CommandConfig{
					Commands: map[string]argh.CommandConfig{
						"fly": argh.CommandConfig{
							Flags: map[string]argh.FlagConfig{
								"freely": {},
							},
						},
						"fry": argh.CommandConfig{
							Flags: map[string]argh.FlagConfig{
								"deeply": {},
								"w":      {},
								"A":      {},
								"t":      {},
							},
						},
					},
					Flags: map[string]argh.FlagConfig{},
				},
			},
			expPT: []argh.Node{
				argh.Command{Name: "pizzas"},
				argh.ArgDelimiter{},
				argh.Command{Name: "fly"},
				argh.ArgDelimiter{},
				argh.Flag{Name: "freely"},
				argh.ArgDelimiter{},
				argh.Command{Name: "fry"},
				argh.ArgDelimiter{},
				argh.Flag{Name: "deeply"},
				argh.ArgDelimiter{},
				argh.CompoundShortFlag{
					Nodes: []argh.Node{
						argh.Flag{Name: "w"},
						argh.Flag{Name: "A"},
						argh.Flag{Name: "t"},
					},
				},
			},
		},
		{
			name: "total weirdo",
			args: []string{"PIZZAs", "^wAT@golf", "^^hecKing", "goose", "bonk", "^^FIERCENESS@-2"},
			cfg: &argh.ParserConfig{
				Prog: argh.CommandConfig{
					Commands: map[string]argh.CommandConfig{
						"goose": argh.CommandConfig{NValue: 1},
					},
					Flags: map[string]argh.FlagConfig{
						"w":          argh.FlagConfig{},
						"A":          argh.FlagConfig{},
						"T":          argh.FlagConfig{NValue: 1},
						"hecking":    argh.FlagConfig{},
						"FIERCENESS": argh.FlagConfig{NValue: 1},
					},
				},
				ScannerConfig: &argh.ScannerConfig{
					AssignmentOperator: '@',
					FlagPrefix:         '^',
					MultiValueDelim:    ',',
				},
			},
			expPT: []argh.Node{
				argh.Command{Name: "PIZZAs"},
				argh.ArgDelimiter{},
				argh.CompoundShortFlag{
					Nodes: []argh.Node{
						argh.Flag{Name: "w"},
						argh.Flag{Name: "A"},
						argh.Flag{Name: "T", Values: map[string]string{"0": "golf"}},
					},
				},
				argh.ArgDelimiter{},
				argh.Flag{Name: "hecKing"},
				argh.ArgDelimiter{},
				argh.Command{Name: "goose", Values: map[string]string{"0": "bonk"}},
				argh.ArgDelimiter{},
				argh.Flag{Name: "FIERCENESS", Values: map[string]string{"0": "-2"}},
			},
		},
		{
			name:   "invalid bare assignment",
			args:   []string{"pizzas", "=", "--wat"},
			expErr: argh.ErrSyntax,
			expPT: []argh.Node{
				argh.Command{Name: "pizzas"},
			},
		},
		{},
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

				assert.Equal(ct, tc.expPT, actual.Nodes)
			})
		}

		/*
			if tc.expAST != nil {
				t.Run(tc.name+" ast", func(ct *testing.T) {
					actual, err := argh.ParseArgs(tc.args, tc.cfg)
					if err != nil {
						assert.ErrorIs(ct, err, tc.expErr)
						return
					}

					assert.Equal(ct, tc.expAST, argh.NewQuerier(actual.Nodes).AST())
				})
			}
		*/
	}
}
