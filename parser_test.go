package argh_test

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/argh"
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
				Prog: &argh.CommandConfig{
					Flags: &argh.Flags{
						Map: map[string]*argh.FlagConfig{
							"e":   {},
							"a":   {},
							"t":   {},
							"wat": {},
						},
					},
					Commands: &argh.Commands{
						Map: map[string]*argh.CommandConfig{
							"hello": &argh.CommandConfig{
								NValue:     1,
								ValueNames: []string{"name"},
							},
						},
					},
				},
			},
			expPT: []argh.Node{
				&argh.CommandFlag{
					Name: "pies",
					Nodes: []argh.Node{
						&argh.ArgDelimiter{},
						&argh.CompoundShortFlag{
							Nodes: []argh.Node{
								&argh.CommandFlag{Name: "e"},
								&argh.CommandFlag{Name: "a"},
								&argh.CommandFlag{Name: "t"},
							},
						},
						&argh.ArgDelimiter{},
						&argh.CommandFlag{Name: "wat"},
						&argh.ArgDelimiter{},
						&argh.CommandFlag{
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
				&argh.CommandFlag{
					Name: "pies",
					Nodes: []argh.Node{
						&argh.CommandFlag{Name: "e"},
						&argh.CommandFlag{Name: "a"},
						&argh.CommandFlag{Name: "t"},
						&argh.CommandFlag{Name: "wat"},
						&argh.CommandFlag{
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
			name: "persistent flags",
			args: []string{
				"pies", "--wat", "hello", "mario", "-eat",
			},
			cfg: &argh.ParserConfig{
				Prog: func() *argh.CommandConfig {
					cmdCfg := &argh.CommandConfig{
						Flags: &argh.Flags{
							Map: map[string]*argh.FlagConfig{
								"e":   {Persist: true},
								"a":   {Persist: true},
								"t":   {Persist: true},
								"wat": {},
							},
						},
					}

					cmdCfg.Commands = &argh.Commands{
						Map: map[string]*argh.CommandConfig{
							"hello": &argh.CommandConfig{
								NValue:     1,
								ValueNames: []string{"name"},
								Flags: &argh.Flags{
									Parent: cmdCfg.Flags,
									Map:    map[string]*argh.FlagConfig{},
								},
							},
						},
					}

					return cmdCfg
				}(),
			},
			expPT: []argh.Node{
				&argh.CommandFlag{
					Name: "pies",
					Nodes: []argh.Node{
						&argh.ArgDelimiter{},
						&argh.CommandFlag{Name: "wat"},
						&argh.ArgDelimiter{},
						&argh.CommandFlag{
							Name: "hello",
							Values: map[string]string{
								"name": "mario",
							},
							Nodes: []argh.Node{
								&argh.ArgDelimiter{},
								&argh.Ident{Literal: "mario"},
								&argh.ArgDelimiter{},
								&argh.CompoundShortFlag{
									Nodes: []argh.Node{
										&argh.CommandFlag{Name: "e"},
										&argh.CommandFlag{Name: "a"},
										&argh.CommandFlag{Name: "t"},
									},
								},
							},
						},
					},
				},
			},
			expAST: []argh.Node{
				&argh.CommandFlag{
					Name: "pies",
					Nodes: []argh.Node{
						&argh.CommandFlag{Name: "wat"},
						&argh.CommandFlag{
							Name: "hello",
							Values: map[string]string{
								"name": "mario",
							},
							Nodes: []argh.Node{
								&argh.Ident{Literal: "mario"},
								&argh.CommandFlag{Name: "e"},
								&argh.CommandFlag{Name: "a"},
								&argh.CommandFlag{Name: "t"},
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
				&argh.CommandFlag{
					Name: "pizzas",
				},
			},
			expAST: []argh.Node{
				&argh.CommandFlag{
					Name: "pizzas",
				},
			},
		},
		{
			name: "one positional arg",
			args: []string{"pizzas", "excel"},
			cfg: &argh.ParserConfig{
				Prog: &argh.CommandConfig{NValue: 1},
			},
			expPT: []argh.Node{
				&argh.CommandFlag{
					Name:   "pizzas",
					Values: map[string]string{"0": "excel"},
					Nodes: []argh.Node{
						&argh.ArgDelimiter{},
						&argh.Ident{Literal: "excel"},
					},
				},
			},
			expAST: []argh.Node{
				&argh.CommandFlag{
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
				Prog: &argh.CommandConfig{
					NValue:     argh.OneOrMoreValue,
					ValueNames: []string{"word"},
				},
			},
			expPT: []argh.Node{
				&argh.CommandFlag{
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
				&argh.CommandFlag{
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
				Prog: &argh.CommandConfig{
					Flags: &argh.Flags{
						Map: map[string]*argh.FlagConfig{
							"tasty":               {},
							"fresh":               {},
							"super-hot-right-now": {},
						},
					},
				},
			},
			expPT: []argh.Node{
				&argh.CommandFlag{
					Name: "pizzas",
					Nodes: []argh.Node{
						&argh.ArgDelimiter{},
						&argh.CommandFlag{Name: "tasty"},
						&argh.ArgDelimiter{},
						&argh.CommandFlag{Name: "fresh"},
						&argh.ArgDelimiter{},
						&argh.CommandFlag{Name: "super-hot-right-now"},
					},
				},
			},
			expAST: []argh.Node{
				&argh.CommandFlag{
					Name: "pizzas",
					Nodes: []argh.Node{
						&argh.CommandFlag{Name: "tasty"},
						&argh.CommandFlag{Name: "fresh"},
						&argh.CommandFlag{Name: "super-hot-right-now"},
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
				Prog: &argh.CommandConfig{
					Commands: &argh.Commands{Map: map[string]*argh.CommandConfig{}},
					Flags: &argh.Flags{
						Map: map[string]*argh.FlagConfig{
							"tasty":               {},
							"fresh":               {NValue: 1},
							"super-hot-right-now": {},
							"box":                 {NValue: argh.OneOrMoreValue},
							"please":              {},
						},
					},
				},
			},
			expPT: []argh.Node{
				&argh.CommandFlag{
					Name: "pizzas",
					Nodes: []argh.Node{
						&argh.ArgDelimiter{},
						&argh.CommandFlag{Name: "tasty"},
						&argh.ArgDelimiter{},
						&argh.CommandFlag{
							Name:   "fresh",
							Values: map[string]string{"0": "soon"},
							Nodes: []argh.Node{
								&argh.ArgDelimiter{},
								&argh.Ident{Literal: "soon"},
							},
						},
						&argh.ArgDelimiter{},
						&argh.CommandFlag{Name: "super-hot-right-now"},
						&argh.ArgDelimiter{},
						&argh.CommandFlag{
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
						&argh.CommandFlag{Name: "please"},
					},
				},
			},
			expAST: []argh.Node{
				&argh.CommandFlag{
					Name: "pizzas",
					Nodes: []argh.Node{
						&argh.CommandFlag{Name: "tasty"},
						&argh.CommandFlag{
							Name:   "fresh",
							Values: map[string]string{"0": "soon"},
							Nodes: []argh.Node{
								&argh.Ident{Literal: "soon"},
							},
						},
						&argh.CommandFlag{Name: "super-hot-right-now"},
						&argh.CommandFlag{
							Name:   "box",
							Values: map[string]string{"0": "square", "1": "shaped", "2": "hot"},
							Nodes: []argh.Node{
								&argh.Ident{Literal: "square"},
								&argh.Ident{Literal: "shaped"},
								&argh.Ident{Literal: "hot"},
							},
						},
						&argh.CommandFlag{Name: "please"},
					},
				},
			},
		},
		{
			name: "short value-less flags",
			args: []string{"pizzas", "-t", "-f", "-s"},
			cfg: &argh.ParserConfig{
				Prog: &argh.CommandConfig{
					Flags: &argh.Flags{
						Map: map[string]*argh.FlagConfig{
							"t": {},
							"f": {},
							"s": {},
						},
					},
				},
			},
			expPT: []argh.Node{
				&argh.CommandFlag{
					Name: "pizzas",
					Nodes: []argh.Node{
						&argh.ArgDelimiter{},
						&argh.CommandFlag{Name: "t"},
						&argh.ArgDelimiter{},
						&argh.CommandFlag{Name: "f"},
						&argh.ArgDelimiter{},
						&argh.CommandFlag{Name: "s"},
					},
				},
			},
			expAST: []argh.Node{
				&argh.CommandFlag{
					Name: "pizzas",
					Nodes: []argh.Node{
						&argh.CommandFlag{Name: "t"},
						&argh.CommandFlag{Name: "f"},
						&argh.CommandFlag{Name: "s"},
					},
				},
			},
		},
		{
			name: "compound short flags",
			args: []string{"pizzas", "-aca", "-blol"},
			cfg: &argh.ParserConfig{
				Prog: &argh.CommandConfig{
					Flags: &argh.Flags{
						Map: map[string]*argh.FlagConfig{
							"a": {},
							"b": {},
							"c": {},
							"l": {},
							"o": {},
						},
					},
				},
			},
			expPT: []argh.Node{
				&argh.CommandFlag{
					Name: "pizzas",
					Nodes: []argh.Node{
						&argh.ArgDelimiter{},
						&argh.CompoundShortFlag{
							Nodes: []argh.Node{
								&argh.CommandFlag{Name: "a"},
								&argh.CommandFlag{Name: "c"},
								&argh.CommandFlag{Name: "a"},
							},
						},
						&argh.ArgDelimiter{},
						&argh.CompoundShortFlag{
							Nodes: []argh.Node{
								&argh.CommandFlag{Name: "b"},
								&argh.CommandFlag{Name: "l"},
								&argh.CommandFlag{Name: "o"},
								&argh.CommandFlag{Name: "l"},
							},
						},
					},
				},
			},
			expAST: []argh.Node{
				&argh.CommandFlag{
					Name: "pizzas",
					Nodes: []argh.Node{
						&argh.CommandFlag{Name: "a"},
						&argh.CommandFlag{Name: "c"},
						&argh.CommandFlag{Name: "a"},
						&argh.CommandFlag{Name: "b"},
						&argh.CommandFlag{Name: "l"},
						&argh.CommandFlag{Name: "o"},
						&argh.CommandFlag{Name: "l"},
					},
				},
			},
		},
		{
			name: "mixed long short value flags",
			args: []string{"pizzas", "-a", "--ca", "-b", "1312", "-lol"},
			cfg: &argh.ParserConfig{
				Prog: &argh.CommandConfig{
					Commands: &argh.Commands{Map: map[string]*argh.CommandConfig{}},
					Flags: &argh.Flags{
						Map: map[string]*argh.FlagConfig{
							"a":  {},
							"b":  {NValue: 1},
							"ca": {},
							"l":  {},
							"o":  {},
						},
					},
				},
			},
			expPT: []argh.Node{
				&argh.CommandFlag{
					Name: "pizzas",
					Nodes: []argh.Node{
						&argh.ArgDelimiter{},
						&argh.CommandFlag{Name: "a"},
						&argh.ArgDelimiter{},
						&argh.CommandFlag{Name: "ca"},
						&argh.ArgDelimiter{},
						&argh.CommandFlag{
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
								&argh.CommandFlag{Name: "l"},
								&argh.CommandFlag{Name: "o"},
								&argh.CommandFlag{Name: "l"},
							},
						},
					},
				},
			},
			expAST: []argh.Node{
				&argh.CommandFlag{
					Name: "pizzas",
					Nodes: []argh.Node{
						&argh.CommandFlag{Name: "a"},
						&argh.CommandFlag{Name: "ca"},
						&argh.CommandFlag{
							Name:   "b",
							Values: map[string]string{"0": "1312"},
							Nodes: []argh.Node{
								&argh.Ident{Literal: "1312"},
							},
						},
						&argh.CommandFlag{Name: "l"},
						&argh.CommandFlag{Name: "o"},
						&argh.CommandFlag{Name: "l"},
					},
				},
			},
		},
		{
			name: "nested commands with positional args",
			args: []string{"pizzas", "fly", "freely", "sometimes", "and", "other", "times", "fry", "deeply", "--forever"},
			cfg: &argh.ParserConfig{
				Prog: &argh.CommandConfig{
					Commands: &argh.Commands{
						Map: map[string]*argh.CommandConfig{
							"fly": &argh.CommandConfig{
								Commands: &argh.Commands{
									Map: map[string]*argh.CommandConfig{
										"fry": &argh.CommandConfig{
											Flags: &argh.Flags{
												Map: map[string]*argh.FlagConfig{
													"forever": {},
												},
											},
										},
									},
								},
							},
						},
					},
					Flags: &argh.Flags{Map: map[string]*argh.FlagConfig{}},
				},
			},
			expPT: []argh.Node{
				&argh.CommandFlag{
					Name: "pizzas",
					Nodes: []argh.Node{
						&argh.ArgDelimiter{},
						&argh.CommandFlag{
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
								&argh.CommandFlag{
									Name: "fry",
									Nodes: []argh.Node{
										&argh.ArgDelimiter{},
										&argh.Ident{Literal: "deeply"},
										&argh.ArgDelimiter{},
										&argh.CommandFlag{Name: "forever"},
									},
								},
							},
						},
					},
				},
			},
			expAST: []argh.Node{
				&argh.CommandFlag{
					Name: "pizzas",
					Nodes: []argh.Node{
						&argh.CommandFlag{
							Name: "fly",
							Nodes: []argh.Node{
								&argh.Ident{Literal: "freely"},
								&argh.Ident{Literal: "sometimes"},
								&argh.Ident{Literal: "and"},
								&argh.Ident{Literal: "other"},
								&argh.Ident{Literal: "times"},
								&argh.CommandFlag{
									Name: "fry",
									Nodes: []argh.Node{
										&argh.Ident{Literal: "deeply"},
										&argh.CommandFlag{Name: "forever"},
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
				Prog: &argh.CommandConfig{
					Flags: &argh.Flags{
						Map: map[string]*argh.FlagConfig{
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
			},
			expPT: []argh.Node{
				&argh.CommandFlag{
					Name: "pizzas",
					Nodes: []argh.Node{
						&argh.ArgDelimiter{},
						&argh.CompoundShortFlag{
							Nodes: []argh.Node{
								&argh.CommandFlag{Name: "n"},
								&argh.CommandFlag{Name: "e"},
								&argh.CommandFlag{Name: "e"},
								&argh.CommandFlag{
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
								&argh.CommandFlag{Name: "a"},
								&argh.CommandFlag{Name: "l"},
								&argh.CommandFlag{Name: "s"},
								&argh.CommandFlag{
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
				&argh.CommandFlag{
					Name: "pizzas",
					Nodes: []argh.Node{
						&argh.CommandFlag{Name: "n"},
						&argh.CommandFlag{Name: "e"},
						&argh.CommandFlag{Name: "e"},
						&argh.CommandFlag{
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
						&argh.CommandFlag{Name: "a"},
						&argh.CommandFlag{Name: "l"},
						&argh.CommandFlag{Name: "s"},
						&argh.CommandFlag{
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
				Prog: &argh.CommandConfig{
					Commands: &argh.Commands{
						Map: map[string]*argh.CommandConfig{
							"fly": &argh.CommandConfig{
								Flags: &argh.Flags{
									Map: map[string]*argh.FlagConfig{
										"freely": {},
									},
								},
								Commands: &argh.Commands{
									Map: map[string]*argh.CommandConfig{
										"fry": &argh.CommandConfig{
											Flags: &argh.Flags{
												Map: map[string]*argh.FlagConfig{
													"deeply": {},
													"w":      {},
													"A":      {},
													"t":      {NValue: 1},
												},
											},
										},
									},
								},
							},
						},
					},
					Flags: &argh.Flags{Map: map[string]*argh.FlagConfig{}},
				},
			},
			expPT: []argh.Node{
				&argh.CommandFlag{
					Name: "pizzas",
					Nodes: []argh.Node{
						&argh.ArgDelimiter{},
						&argh.CommandFlag{
							Name: "fly",
							Nodes: []argh.Node{
								&argh.ArgDelimiter{},
								&argh.CommandFlag{Name: "freely"},
								&argh.ArgDelimiter{},
								&argh.CommandFlag{
									Name: "fry",
									Nodes: []argh.Node{
										&argh.ArgDelimiter{},
										&argh.CommandFlag{Name: "deeply"},
										&argh.ArgDelimiter{},
										&argh.CompoundShortFlag{
											Nodes: []argh.Node{
												&argh.CommandFlag{Name: "w"},
												&argh.CommandFlag{Name: "A"},
												&argh.CommandFlag{
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
				&argh.CommandFlag{
					Name: "pizzas",
					Nodes: []argh.Node{
						&argh.CommandFlag{
							Name: "fly",
							Nodes: []argh.Node{
								&argh.CommandFlag{Name: "freely"},
								&argh.CommandFlag{
									Name: "fry",
									Nodes: []argh.Node{
										&argh.CommandFlag{Name: "deeply"},
										&argh.CommandFlag{Name: "w"},
										&argh.CommandFlag{Name: "A"},
										&argh.CommandFlag{
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
				Prog: &argh.CommandConfig{
					Commands: &argh.Commands{
						Map: map[string]*argh.CommandConfig{
							"goose": &argh.CommandConfig{
								NValue: 1,
								Flags: &argh.Flags{
									Map: map[string]*argh.FlagConfig{
										"FIERCENESS": {NValue: 1},
									},
								},
							},
						},
					},
					Flags: &argh.Flags{
						Map: map[string]*argh.FlagConfig{
							"w":       {},
							"A":       {},
							"T":       {NValue: 1},
							"hecKing": {},
						},
					},
				},
				ScannerConfig: &argh.ScannerConfig{
					AssignmentOperator: '@',
					FlagPrefix:         '^',
					MultiValueDelim:    ',',
				},
			},
			expPT: []argh.Node{
				&argh.CommandFlag{
					Name: "PIZZAs",
					Nodes: []argh.Node{
						&argh.ArgDelimiter{},
						&argh.CompoundShortFlag{
							Nodes: []argh.Node{
								&argh.CommandFlag{Name: "w"},
								&argh.CommandFlag{Name: "A"},
								&argh.CommandFlag{
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
						&argh.CommandFlag{Name: "hecKing"},
						&argh.ArgDelimiter{},
						&argh.CommandFlag{
							Name:   "goose",
							Values: map[string]string{"0": "bonk"},
							Nodes: []argh.Node{
								&argh.ArgDelimiter{},
								&argh.Ident{Literal: "bonk"},
								&argh.ArgDelimiter{},
								&argh.CommandFlag{
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
				Prog: &argh.CommandConfig{
					Flags: &argh.Flags{
						Map: map[string]*argh.FlagConfig{
							"f": {},
							"L": {},
							"o": {NValue: 1},
						},
					},
					Commands: &argh.Commands{
						Map: map[string]*argh.CommandConfig{
							"hats": {},
						},
					},
				},
				ScannerConfig: &argh.ScannerConfig{
					AssignmentOperator: ':',
					FlagPrefix:         '/',
					MultiValueDelim:    ',',
				},
			},
			expPT: []argh.Node{
				&argh.CommandFlag{
					Name: "hotdog",
					Nodes: []argh.Node{
						&argh.ArgDelimiter{},
						&argh.CommandFlag{Name: "f"},
						&argh.ArgDelimiter{},
						&argh.CommandFlag{Name: "L"},
						&argh.ArgDelimiter{},
						&argh.CommandFlag{
							Name:   "o",
							Values: map[string]string{"0": "ppy"},
							Nodes: []argh.Node{
								&argh.Assign{},
								&argh.Ident{Literal: "ppy"},
							},
						},
						&argh.ArgDelimiter{},
						&argh.CommandFlag{Name: "hats"},
					},
				},
			},
		},
		{
			name: "invalid bare assignment",
			args: []string{"pizzas", "=", "--wat"},
			cfg: &argh.ParserConfig{
				Prog: &argh.CommandConfig{
					Flags: &argh.Flags{
						Map: map[string]*argh.FlagConfig{
							"wat": {},
						},
					},
				},
			},
			expErr: argh.ParserErrorList{
				&argh.ParserError{Pos: argh.Position{Column: 8}, Msg: "invalid bare assignment"},
			},
			expPT: []argh.Node{
				&argh.CommandFlag{
					Name: "pizzas",
					Nodes: []argh.Node{
						&argh.ArgDelimiter{},
						&argh.ArgDelimiter{},
						&argh.CommandFlag{Name: "wat"},
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

				pCfg := tc.cfg
				if pCfg == nil {
					pCfg = argh.NewParserConfig()
				}

				pt, err := argh.ParseArgs(tc.args, pCfg)
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

				pCfg := tc.cfg
				if pCfg == nil {
					pCfg = argh.NewParserConfig()
				}

				pt, err := argh.ParseArgs(tc.args, pCfg)
				if err != nil || tc.expErr != nil {
					if !assert.ErrorIs(ct, err, tc.expErr) {
						spew.Dump(pt)
					}
					return
				}

				ast := argh.ToAST(pt.Nodes)

				if !assert.Equal(ct, tc.expAST, ast) {
					spew.Dump(ast)
				}
			})
		}
	}
}
