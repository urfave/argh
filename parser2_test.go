package argh_test

import (
	"testing"

	"git.meatballhat.com/x/box-o-sand/argh"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestParser2(t *testing.T) {
	for _, tc := range []struct {
		name   string
		args   []string
		cfg    *argh.ParserConfig
		expErr error
		expPT  []argh.Node
		expAST []argh.Node
		skip   bool
	}{
		{
			name: "basic",
			args: []string{
				"pies", "-eat", "--wat", "hello", "mario",
			},
			cfg: &argh.ParserConfig{
				Prog: argh.CommandConfig{
					Commands: map[string]argh.CommandConfig{
						"hello": argh.CommandConfig{
							NValue:     1,
							ValueNames: []string{"name"},
						},
					},
				},
			},
			expPT: []argh.Node{
				&argh.Command{
					Name: "pies",
					Nodes: []argh.Node{
						&argh.ArgDelimiter{},
						&argh.CompoundShortFlag{
							Nodes: []argh.Node{
								&argh.Flag{Name: "e"},
								&argh.Flag{Name: "a"},
								&argh.Flag{Name: "t"},
							},
						},
						&argh.ArgDelimiter{},
						&argh.Flag{Name: "wat"},
						&argh.ArgDelimiter{},
						&argh.Command{
							Name: "hello",
							Values: map[string]string{
								"name": "mario",
							},
							Nodes: []argh.Node{
								&argh.ArgDelimiter{},
							},
						},
					},
				},
			},
			expAST: []argh.Node{
				&argh.Command{
					Name: "pies",
					Nodes: []argh.Node{
						&argh.Flag{Name: "e"},
						&argh.Flag{Name: "a"},
						&argh.Flag{Name: "t"},
						&argh.Flag{Name: "wat"},
						&argh.Command{
							Name: "hello",
							Values: map[string]string{
								"name": "mario",
							},
						},
					},
				},
			},
		},
		{
			name: "bare",
			args: []string{"pizzas"},
			expPT: []argh.Node{
				&argh.Command{
					Name: "pizzas",
				},
			},
			expAST: []argh.Node{
				&argh.Command{
					Name: "pizzas",
				},
			},
		},
		{
			name: "one positional arg",
			args: []string{"pizzas", "excel"},
			cfg: &argh.ParserConfig{
				Prog: argh.CommandConfig{NValue: 1},
			},
			expPT: []argh.Node{
				&argh.Command{
					Name:   "pizzas",
					Values: map[string]string{"0": "excel"},
					Nodes: []argh.Node{
						&argh.ArgDelimiter{},
					},
				},
			},
			expAST: []argh.Node{
				&argh.Command{
					Name:   "pizzas",
					Values: map[string]string{"0": "excel"},
				},
			},
		},
		{
			name: "many positional args",
			args: []string{"pizzas", "excel", "wildly", "when", "feral"},
			cfg: &argh.ParserConfig{
				Prog: argh.CommandConfig{
					NValue:     argh.OneOrMoreValue,
					ValueNames: []string{"word"},
				},
			},
			expPT: []argh.Node{
				&argh.Command{
					Name: "pizzas",
					Values: map[string]string{
						"word":   "excel",
						"word.1": "wildly",
						"word.2": "when",
						"word.3": "feral",
					},
					Nodes: []argh.Node{
						&argh.ArgDelimiter{},
						&argh.ArgDelimiter{},
						&argh.ArgDelimiter{},
						&argh.ArgDelimiter{},
					},
				},
			},
			expAST: []argh.Node{
				&argh.Command{
					Name: "pizzas",
					Values: map[string]string{
						"word":   "excel",
						"word.1": "wildly",
						"word.2": "when",
						"word.3": "feral",
					},
				},
			},
		},
		{
			name: "long value-less flags",
			args: []string{"pizzas", "--tasty", "--fresh", "--super-hot-right-now"},
			expPT: []argh.Node{
				&argh.Command{
					Name: "pizzas",
					Nodes: []argh.Node{
						&argh.ArgDelimiter{},
						&argh.Flag{Name: "tasty"},
						&argh.ArgDelimiter{},
						&argh.Flag{Name: "fresh"},
						&argh.ArgDelimiter{},
						&argh.Flag{Name: "super-hot-right-now"},
					},
				},
			},
			expAST: []argh.Node{
				&argh.Command{
					Name: "pizzas",
					Nodes: []argh.Node{
						&argh.Flag{Name: "tasty"},
						&argh.Flag{Name: "fresh"},
						&argh.Flag{Name: "super-hot-right-now"},
					},
				},
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
				&argh.Command{
					Name: "pizzas",
					Nodes: []argh.Node{
						&argh.ArgDelimiter{},
						&argh.Flag{Name: "tasty"},
						&argh.ArgDelimiter{},
						&argh.Flag{
							Name:   "fresh",
							Values: map[string]string{"0": "soon"},
							Nodes: []argh.Node{
								&argh.ArgDelimiter{},
							},
						},
						&argh.ArgDelimiter{},
						&argh.Flag{Name: "super-hot-right-now"},
						&argh.ArgDelimiter{},
						&argh.Flag{
							Name:   "box",
							Values: map[string]string{"0": "square", "1": "shaped", "2": "hot"},
							Nodes: []argh.Node{
								&argh.ArgDelimiter{},
								&argh.ArgDelimiter{},
								&argh.ArgDelimiter{},
								&argh.ArgDelimiter{},
							},
						},
						&argh.Flag{Name: "please"},
					},
				},
			},
			expAST: []argh.Node{
				&argh.Command{
					Name: "pizzas",
					Nodes: []argh.Node{
						&argh.Flag{Name: "tasty"},
						&argh.Flag{Name: "fresh", Values: map[string]string{"0": "soon"}},
						&argh.Flag{Name: "super-hot-right-now"},
						&argh.Flag{Name: "box", Values: map[string]string{"0": "square", "1": "shaped", "2": "hot"}},
						&argh.Flag{Name: "please"},
					},
				},
			},
		},
		{
			skip: true,

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
			skip: true,

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
			skip: true,

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
			skip: true,

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
			skip: true,

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
			skip: true,

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
			skip:   true,
			name:   "invalid bare assignment",
			args:   []string{"pizzas", "=", "--wat"},
			expErr: argh.ErrSyntax,
			expPT: []argh.Node{
				argh.Command{Name: "pizzas"},
			},
		},
	} {
		if tc.expPT != nil {
			t.Run(tc.name+" parse tree", func(ct *testing.T) {
				if tc.skip {
					ct.SkipNow()
					return
				}

				pt, err := argh.ParseArgs2(tc.args, tc.cfg)
				if err != nil {
					assert.ErrorIs(ct, err, tc.expErr)
					return
				}

				if !assert.Equal(ct, tc.expPT, pt.Nodes) {
					spew.Dump(pt)
				}
			})
		}

		if tc.expAST != nil {
			t.Run(tc.name+" ast", func(ct *testing.T) {
				if tc.skip {
					ct.SkipNow()
					return
				}

				pt, err := argh.ParseArgs2(tc.args, tc.cfg)
				if err != nil {
					ct.Logf("err=%+#v", err)
					return
				}

				ast := argh.NewQuerier(pt.Nodes).AST()

				if !assert.Equal(ct, tc.expAST, ast) {
					spew.Dump(ast)
				}
			})
		}
	}
}
