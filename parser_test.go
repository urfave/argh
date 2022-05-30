package argh_test

import (
	"testing"

	"git.meatballhat.com/x/box-o-sand/argh"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
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
					Flags: map[string]argh.FlagConfig{
						"e":   {},
						"a":   {},
						"t":   {},
						"wat": {},
					},
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
								&argh.Ident{Literal: "mario"},
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
							Nodes: []argh.Node{
								&argh.Ident{Literal: "mario"},
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
						&argh.Ident{Literal: "excel"},
					},
				},
			},
			expAST: []argh.Node{
				&argh.Command{
					Name:   "pizzas",
					Values: map[string]string{"0": "excel"},
					Nodes: []argh.Node{
						&argh.Ident{Literal: "excel"},
					},
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
						&argh.Ident{Literal: "excel"},
						&argh.ArgDelimiter{},
						&argh.Ident{Literal: "wildly"},
						&argh.ArgDelimiter{},
						&argh.Ident{Literal: "when"},
						&argh.ArgDelimiter{},
						&argh.Ident{Literal: "feral"},
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
					Nodes: []argh.Node{
						&argh.Ident{Literal: "excel"},
						&argh.Ident{Literal: "wildly"},
						&argh.Ident{Literal: "when"},
						&argh.Ident{Literal: "feral"},
					},
				},
			},
		},
		{
			name: "long value-less flags",
			args: []string{"pizzas", "--tasty", "--fresh", "--super-hot-right-now"},
			cfg: &argh.ParserConfig{
				Prog: argh.CommandConfig{
					Flags: map[string]argh.FlagConfig{
						"tasty":               {},
						"fresh":               {},
						"super-hot-right-now": {},
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
						"tasty":               {},
						"fresh":               argh.FlagConfig{NValue: 1},
						"super-hot-right-now": {},
						"box":                 argh.FlagConfig{NValue: argh.OneOrMoreValue},
						"please":              {},
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
								&argh.Ident{Literal: "soon"},
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
								&argh.Ident{Literal: "square"},
								&argh.ArgDelimiter{},
								&argh.Ident{Literal: "shaped"},
								&argh.ArgDelimiter{},
								&argh.Ident{Literal: "hot"},
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
						&argh.Flag{
							Name:   "fresh",
							Values: map[string]string{"0": "soon"},
							Nodes: []argh.Node{
								&argh.Ident{Literal: "soon"},
							},
						},
						&argh.Flag{Name: "super-hot-right-now"},
						&argh.Flag{
							Name:   "box",
							Values: map[string]string{"0": "square", "1": "shaped", "2": "hot"},
							Nodes: []argh.Node{
								&argh.Ident{Literal: "square"},
								&argh.Ident{Literal: "shaped"},
								&argh.Ident{Literal: "hot"},
							},
						},
						&argh.Flag{Name: "please"},
					},
				},
			},
		},
		{
			name: "short value-less flags",
			args: []string{"pizzas", "-t", "-f", "-s"},
			cfg: &argh.ParserConfig{
				Prog: argh.CommandConfig{
					Flags: map[string]argh.FlagConfig{
						"t": {},
						"f": {},
						"s": {},
					},
				},
			},
			expPT: []argh.Node{
				&argh.Command{
					Name: "pizzas",
					Nodes: []argh.Node{
						&argh.ArgDelimiter{},
						&argh.Flag{Name: "t"},
						&argh.ArgDelimiter{},
						&argh.Flag{Name: "f"},
						&argh.ArgDelimiter{},
						&argh.Flag{Name: "s"},
					},
				},
			},
			expAST: []argh.Node{
				&argh.Command{
					Name: "pizzas",
					Nodes: []argh.Node{
						&argh.Flag{Name: "t"},
						&argh.Flag{Name: "f"},
						&argh.Flag{Name: "s"},
					},
				},
			},
		},
		{
			name: "compound short flags",
			args: []string{"pizzas", "-aca", "-blol"},
			cfg: &argh.ParserConfig{
				Prog: argh.CommandConfig{
					Flags: map[string]argh.FlagConfig{
						"a": {},
						"b": {},
						"c": {},
						"l": {},
						"o": {},
					},
				},
			},
			expPT: []argh.Node{
				&argh.Command{
					Name: "pizzas",
					Nodes: []argh.Node{
						&argh.ArgDelimiter{},
						&argh.CompoundShortFlag{
							Nodes: []argh.Node{
								&argh.Flag{Name: "a"},
								&argh.Flag{Name: "c"},
								&argh.Flag{Name: "a"},
							},
						},
						&argh.ArgDelimiter{},
						&argh.CompoundShortFlag{
							Nodes: []argh.Node{
								&argh.Flag{Name: "b"},
								&argh.Flag{Name: "l"},
								&argh.Flag{Name: "o"},
								&argh.Flag{Name: "l"},
							},
						},
					},
				},
			},
			expAST: []argh.Node{
				&argh.Command{
					Name: "pizzas",
					Nodes: []argh.Node{
						&argh.Flag{Name: "a"},
						&argh.Flag{Name: "c"},
						&argh.Flag{Name: "a"},
						&argh.Flag{Name: "b"},
						&argh.Flag{Name: "l"},
						&argh.Flag{Name: "o"},
						&argh.Flag{Name: "l"},
					},
				},
			},
		},
		{
			name: "mixed long short value flags",
			args: []string{"pizzas", "-a", "--ca", "-b", "1312", "-lol"},
			cfg: &argh.ParserConfig{
				Prog: argh.CommandConfig{
					Commands: map[string]argh.CommandConfig{},
					Flags: map[string]argh.FlagConfig{
						"a":  {},
						"b":  argh.FlagConfig{NValue: 1},
						"ca": {},
						"l":  {},
						"o":  {},
					},
				},
			},
			expPT: []argh.Node{
				&argh.Command{
					Name: "pizzas",
					Nodes: []argh.Node{
						&argh.ArgDelimiter{},
						&argh.Flag{Name: "a"},
						&argh.ArgDelimiter{},
						&argh.Flag{Name: "ca"},
						&argh.ArgDelimiter{},
						&argh.Flag{
							Name:   "b",
							Values: map[string]string{"0": "1312"},
							Nodes: []argh.Node{
								&argh.ArgDelimiter{},
								&argh.Ident{Literal: "1312"},
							},
						},
						&argh.ArgDelimiter{},
						&argh.CompoundShortFlag{
							Nodes: []argh.Node{
								&argh.Flag{Name: "l"},
								&argh.Flag{Name: "o"},
								&argh.Flag{Name: "l"},
							},
						},
					},
				},
			},
			expAST: []argh.Node{
				&argh.Command{
					Name: "pizzas",
					Nodes: []argh.Node{
						&argh.Flag{Name: "a"},
						&argh.Flag{Name: "ca"},
						&argh.Flag{
							Name:   "b",
							Values: map[string]string{"0": "1312"},
							Nodes: []argh.Node{
								&argh.Ident{Literal: "1312"},
							},
						},
						&argh.Flag{Name: "l"},
						&argh.Flag{Name: "o"},
						&argh.Flag{Name: "l"},
					},
				},
			},
		},
		{
			name: "nested commands with positional args",
			args: []string{"pizzas", "fly", "freely", "sometimes", "and", "other", "times", "fry", "deeply", "--forever"},
			cfg: &argh.ParserConfig{
				Prog: argh.CommandConfig{
					Commands: map[string]argh.CommandConfig{
						"fly": argh.CommandConfig{
							Commands: map[string]argh.CommandConfig{
								"fry": argh.CommandConfig{
									Flags: map[string]argh.FlagConfig{
										"forever": {},
									},
								},
							},
						},
					},
					Flags: map[string]argh.FlagConfig{},
				},
			},
			expPT: []argh.Node{
				&argh.Command{
					Name: "pizzas",
					Nodes: []argh.Node{
						&argh.ArgDelimiter{},
						&argh.Command{
							Name: "fly",
							Nodes: []argh.Node{
								&argh.ArgDelimiter{},
								&argh.Ident{Literal: "freely"},
								&argh.ArgDelimiter{},
								&argh.Ident{Literal: "sometimes"},
								&argh.ArgDelimiter{},
								&argh.Ident{Literal: "and"},
								&argh.ArgDelimiter{},
								&argh.Ident{Literal: "other"},
								&argh.ArgDelimiter{},
								&argh.Ident{Literal: "times"},
								&argh.ArgDelimiter{},
								&argh.Command{
									Name: "fry",
									Nodes: []argh.Node{
										&argh.ArgDelimiter{},
										&argh.Ident{Literal: "deeply"},
										&argh.ArgDelimiter{},
										&argh.Flag{Name: "forever"},
									},
								},
							},
						},
					},
				},
			},
			expAST: []argh.Node{
				&argh.Command{
					Name: "pizzas",
					Nodes: []argh.Node{
						&argh.Command{
							Name: "fly",
							Nodes: []argh.Node{
								&argh.Ident{Literal: "freely"},
								&argh.Ident{Literal: "sometimes"},
								&argh.Ident{Literal: "and"},
								&argh.Ident{Literal: "other"},
								&argh.Ident{Literal: "times"},
								&argh.Command{
									Name: "fry",
									Nodes: []argh.Node{
										&argh.Ident{Literal: "deeply"},
										&argh.Flag{Name: "forever"},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "compound flags with values",
			args: []string{"pizzas", "-need", "sauce", "heat", "love", "-also", "over9000"},
			cfg: &argh.ParserConfig{
				Prog: argh.CommandConfig{
					Flags: map[string]argh.FlagConfig{
						"a": {NValue: argh.ZeroOrMoreValue},
						"d": {NValue: argh.OneOrMoreValue},
						"e": {},
						"l": {},
						"n": {},
						"o": {NValue: 1, ValueNames: []string{"level"}},
						"s": {NValue: argh.ZeroOrMoreValue},
					},
				},
			},
			expPT: []argh.Node{
				&argh.Command{
					Name: "pizzas",
					Nodes: []argh.Node{
						&argh.ArgDelimiter{},
						&argh.CompoundShortFlag{
							Nodes: []argh.Node{
								&argh.Flag{Name: "n"},
								&argh.Flag{Name: "e"},
								&argh.Flag{Name: "e"},
								&argh.Flag{
									Name: "d",
									Values: map[string]string{
										"0": "sauce",
										"1": "heat",
										"2": "love",
									},
									Nodes: []argh.Node{
										&argh.ArgDelimiter{},
										&argh.Ident{Literal: "sauce"},
										&argh.ArgDelimiter{},
										&argh.Ident{Literal: "heat"},
										&argh.ArgDelimiter{},
										&argh.Ident{Literal: "love"},
										&argh.ArgDelimiter{},
									},
								},
							},
						},
						&argh.CompoundShortFlag{
							Nodes: []argh.Node{
								&argh.Flag{Name: "a"},
								&argh.Flag{Name: "l"},
								&argh.Flag{Name: "s"},
								&argh.Flag{
									Name: "o",
									Values: map[string]string{
										"level": "over9000",
									},
									Nodes: []argh.Node{
										&argh.ArgDelimiter{},
										&argh.Ident{Literal: "over9000"},
									},
								},
							},
						},
					},
				},
			},
			expAST: []argh.Node{
				&argh.Command{
					Name: "pizzas",
					Nodes: []argh.Node{
						&argh.Flag{Name: "n"},
						&argh.Flag{Name: "e"},
						&argh.Flag{Name: "e"},
						&argh.Flag{
							Name: "d",
							Values: map[string]string{
								"0": "sauce",
								"1": "heat",
								"2": "love",
							},
							Nodes: []argh.Node{
								&argh.Ident{Literal: "sauce"},
								&argh.Ident{Literal: "heat"},
								&argh.Ident{Literal: "love"},
							},
						},
						&argh.Flag{Name: "a"},
						&argh.Flag{Name: "l"},
						&argh.Flag{Name: "s"},
						&argh.Flag{
							Name: "o",
							Values: map[string]string{
								"level": "over9000",
							},
							Nodes: []argh.Node{
								&argh.Ident{Literal: "over9000"},
							},
						},
					},
				},
			},
		},
		{
			name: "command specific flags",
			args: []string{"pizzas", "fly", "--freely", "fry", "--deeply", "-wAt", "hugs"},
			cfg: &argh.ParserConfig{
				Prog: argh.CommandConfig{
					Commands: map[string]argh.CommandConfig{
						"fly": argh.CommandConfig{
							Flags: map[string]argh.FlagConfig{
								"freely": {},
							},
							Commands: map[string]argh.CommandConfig{
								"fry": argh.CommandConfig{
									Flags: map[string]argh.FlagConfig{
										"deeply": {},
										"w":      {},
										"A":      {},
										"t":      argh.FlagConfig{NValue: 1},
									},
								},
							},
						},
					},
					Flags: map[string]argh.FlagConfig{},
				},
			},
			expPT: []argh.Node{
				&argh.Command{
					Name: "pizzas",
					Nodes: []argh.Node{
						&argh.ArgDelimiter{},
						&argh.Command{
							Name: "fly",
							Nodes: []argh.Node{
								&argh.ArgDelimiter{},
								&argh.Flag{Name: "freely"},
								&argh.ArgDelimiter{},
								&argh.Command{
									Name: "fry",
									Nodes: []argh.Node{
										&argh.ArgDelimiter{},
										&argh.Flag{Name: "deeply"},
										&argh.ArgDelimiter{},
										&argh.CompoundShortFlag{
											Nodes: []argh.Node{
												&argh.Flag{Name: "w"},
												&argh.Flag{Name: "A"},
												&argh.Flag{
													Name:   "t",
													Values: map[string]string{"0": "hugs"},
													Nodes: []argh.Node{
														&argh.ArgDelimiter{},
														&argh.Ident{Literal: "hugs"},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expAST: []argh.Node{
				&argh.Command{
					Name: "pizzas",
					Nodes: []argh.Node{
						&argh.Command{
							Name: "fly",
							Nodes: []argh.Node{
								&argh.Flag{Name: "freely"},
								&argh.Command{
									Name: "fry",
									Nodes: []argh.Node{
										&argh.Flag{Name: "deeply"},
										&argh.Flag{Name: "w"},
										&argh.Flag{Name: "A"},
										&argh.Flag{
											Name:   "t",
											Values: map[string]string{"0": "hugs"},
											Nodes: []argh.Node{
												&argh.Ident{Literal: "hugs"},
											},
										},
									},
								},
							},
						},
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
						"goose": argh.CommandConfig{
							NValue: 1,
							Flags: map[string]argh.FlagConfig{
								"FIERCENESS": argh.FlagConfig{NValue: 1},
							},
						},
					},
					Flags: map[string]argh.FlagConfig{
						"w":       argh.FlagConfig{},
						"A":       argh.FlagConfig{},
						"T":       argh.FlagConfig{NValue: 1},
						"hecKing": argh.FlagConfig{},
					},
				},
				ScannerConfig: &argh.ScannerConfig{
					AssignmentOperator: '@',
					FlagPrefix:         '^',
					MultiValueDelim:    ',',
				},
			},
			expPT: []argh.Node{
				&argh.Command{
					Name: "PIZZAs",
					Nodes: []argh.Node{
						&argh.ArgDelimiter{},
						&argh.CompoundShortFlag{
							Nodes: []argh.Node{
								&argh.Flag{Name: "w"},
								&argh.Flag{Name: "A"},
								&argh.Flag{
									Name:   "T",
									Values: map[string]string{"0": "golf"},
									Nodes: []argh.Node{
										&argh.Assign{},
										&argh.Ident{Literal: "golf"},
									},
								},
							},
						},
						&argh.ArgDelimiter{},
						&argh.Flag{Name: "hecKing"},
						&argh.ArgDelimiter{},
						&argh.Command{
							Name:   "goose",
							Values: map[string]string{"0": "bonk"},
							Nodes: []argh.Node{
								&argh.ArgDelimiter{},
								&argh.Ident{Literal: "bonk"},
								&argh.ArgDelimiter{},
								&argh.Flag{
									Name:   "FIERCENESS",
									Values: map[string]string{"0": "-2"},
									Nodes: []argh.Node{
										&argh.Assign{},
										&argh.Ident{Literal: "-2"},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "windows like",
			args: []string{"hotdog", "/f", "/L", "/o:ppy", "hats"},
			cfg: &argh.ParserConfig{
				Prog: argh.CommandConfig{
					Flags: map[string]argh.FlagConfig{
						"f": {},
						"L": {},
						"o": argh.FlagConfig{NValue: 1},
					},
					Commands: map[string]argh.CommandConfig{
						"hats": {},
					},
				},
				ScannerConfig: &argh.ScannerConfig{
					AssignmentOperator: ':',
					FlagPrefix:         '/',
					MultiValueDelim:    ',',
				},
			},
			expPT: []argh.Node{
				&argh.Command{
					Name: "hotdog",
					Nodes: []argh.Node{
						&argh.ArgDelimiter{},
						&argh.Flag{Name: "f"},
						&argh.ArgDelimiter{},
						&argh.Flag{Name: "L"},
						&argh.ArgDelimiter{},
						&argh.Flag{
							Name:   "o",
							Values: map[string]string{"0": "ppy"},
							Nodes: []argh.Node{
								&argh.Assign{},
								&argh.Ident{Literal: "ppy"},
							},
						},
						&argh.ArgDelimiter{},
						&argh.Command{Name: "hats"},
					},
				},
			},
		},
		{
			name: "invalid bare assignment",
			args: []string{"pizzas", "=", "--wat"},
			cfg: &argh.ParserConfig{
				Prog: argh.CommandConfig{
					Flags: map[string]argh.FlagConfig{
						"wat": {},
					},
				},
			},
			expErr: argh.ParserErrorList{
				&argh.ParserError{Pos: argh.Position{Column: 8}, Msg: "invalid bare assignment"},
			},
			expPT: []argh.Node{
				&argh.Command{
					Name: "pizzas",
					Nodes: []argh.Node{
						&argh.ArgDelimiter{},
						&argh.ArgDelimiter{},
						&argh.Flag{Name: "wat"},
					},
				},
			},
		},
	} {
		if tc.expPT != nil {
			t.Run(tc.name+" parse tree", func(ct *testing.T) {
				if tc.skip {
					ct.SkipNow()
					return
				}

				pt, err := argh.ParseArgs(tc.args, tc.cfg)
				if err != nil || tc.expErr != nil {
					if !assert.ErrorIs(ct, err, tc.expErr) {
						spew.Dump(err, tc.expErr)
						spew.Dump(pt)
					}
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

				pt, err := argh.ParseArgs(tc.args, tc.cfg)
				if err != nil || tc.expErr != nil {
					if !assert.ErrorIs(ct, err, tc.expErr) {
						spew.Dump(pt)
					}
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
