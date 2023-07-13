package argh

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnparseTree(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		r := require.New(t)

		sv, err := UnparseTree([]Node{}, POSIXyScannerConfig)

		r.NoError(err)
		r.Equal([]string{}, sv)
	})

	t.Run("idents only", func(t *testing.T) {
		r := require.New(t)

		sv, err := UnparseTree(
			[]Node{
				&Ident{Literal: "steamed"},
				&ArgDelimiter{},
				&Ident{Literal: "hams"},
			},
			POSIXyScannerConfig,
		)
		r.NoError(err)

		r.Equal([]string{"steamed", "hams"}, sv)
	})

	t.Run("flags only", func(t *testing.T) {
		r := require.New(t)

		sv, err := UnparseTree(
			[]Node{
				&CompoundShortFlag{
					Nodes: []Node{
						&Flag{Name: "o"},
						&Flag{Name: "o"},
						&Flag{
							Name: "f",
							Values: map[string]string{
								"0": "yep",
								"1": "maybe",
							},
							Nodes: []Node{
								&Assign{},
								&MultiIdent{
									Nodes: []Node{
										&Ident{Literal: "yep"},
										&Ident{Literal: "maybe"},
									},
								},
							},
						},
					},
				},
				&ArgDelimiter{},
				&Flag{
					Name:   "sandwiches",
					Values: map[string]string{"0": "42"},
					Nodes: []Node{
						&Assign{},
						&Ident{Literal: "42"},
					},
				},
				&ArgDelimiter{},
				&Flag{Name: "q"},
			},
			POSIXyScannerConfig,
		)

		r.NoError(err)

		r.Equal(
			[]string{"-oof=yep,maybe", "--sandwiches=42", "-q"},
			sv,
		)
	})

	t.Run("simple", func(t *testing.T) {
		r := require.New(t)

		sv, err := UnparseTree(
			[]Node{
				&Command{
					Name: "pies",
					Nodes: []Node{
						&ArgDelimiter{},
						&CompoundShortFlag{
							Nodes: []Node{
								&Flag{Name: "e"},
								&Flag{Name: "a"},
								&Flag{Name: "t"},
							},
						},
						&ArgDelimiter{},
						&Command{Name: "wat"},
						&ArgDelimiter{},
						&Command{
							Name: "hello",
							Values: map[string]string{
								"name": "mario",
							},
							Nodes: []Node{
								&ArgDelimiter{},
								&Ident{Literal: "mario"},
							},
						},
					},
				},
			},
			POSIXyScannerConfig,
		)

		r.NoError(err)

		r.Equal(
			[]string{"pies", "-eat", "wat", "hello", "mario"},
			sv,
		)
	})

	t.Run("compound flags with value", func(t *testing.T) {
		r := require.New(t)

		sv, err := UnparseTree(
			[]Node{
				&Command{
					Name: "howling",
					Nodes: []Node{
						&ArgDelimiter{},
						&CompoundShortFlag{
							Nodes: []Node{
								&Flag{Name: "f"},
								&Flag{Name: "R"},
								&Flag{Name: "i"},
								&Flag{Name: "E"},
								&Flag{Name: "n"},
								&Flag{
									Name: "d",
									Values: map[string]string{
										"0": "o",
									},
									Nodes: []Node{
										&Assign{},
										&Ident{Literal: "o"},
									},
								},
							},
						},
					},
				},
			},
			POSIXyScannerConfig,
		)

		r.NoError(err)
		r.Equal([]string{"howling", "-fRiEnd=o"}, sv)
	})

	t.Run("multi-value flags", func(t *testing.T) {
		r := require.New(t)

		sv, err := UnparseTree(
			[]Node{
				&Command{
					Name: "fjords",
					Nodes: []Node{
						&ArgDelimiter{},
						&Flag{
							Name: "with",
							Values: map[string]string{
								"0": "whales",
								"1": "majesticness",
								"2": "waters",
							},
							Nodes: []Node{
								&Assign{},
								&MultiIdent{
									Nodes: []Node{
										&Ident{Literal: "whales"},
										&Ident{Literal: "majesticness"},
										&Ident{Literal: "waters"},
									},
								},
							},
						},
						&ArgDelimiter{},
						&Flag{
							Name: "a",
							Values: map[string]string{
								"0": "sparkling",
								"1": "lens flares",
								"2": "probably ducks",
							},
							Nodes: []Node{
								&Assign{},
								&MultiIdent{
									Nodes: []Node{
										&Ident{Literal: "sparkling"},
										&Ident{Literal: "lens flares"},
									},
								},
								&ArgDelimiter{},
								&Ident{Literal: "probably ducks"},
							},
						},
					},
				},
			},
			POSIXyScannerConfig,
		)

		r.NoError(err)
		r.Equal(
			[]string{
				"fjords",
				"--with=whales,majesticness,waters",
				"-a=sparkling,lens flares",
				"probably ducks",
			},
			sv,
		)
	})

	t.Run("with passthrough args", func(t *testing.T) {
		r := require.New(t)

		sv, err := UnparseTree(
			[]Node{
				&Command{
					Name: "shoelace",
					Nodes: []Node{
						&ArgDelimiter{},
						&Flag{Name: "ok"},
						&ArgDelimiter{},
						&StopFlag{},
						&ArgDelimiter{},
						&PassthroughArgs{
							Nodes: []Node{
								&Ident{Literal: "tardigrade=smol"},
								&Ident{Literal: "--??"},
								&Ident{Literal: "-!"},
							},
						},
					},
				},
			},
			POSIXyScannerConfig,
		)

		r.NoError(err)

		r.Equal(
			[]string{"shoelace", "--ok", "--", "tardigrade=smol", "--??", "-!"},
			sv,
		)
	})

	t.Run("curlish", func(t *testing.T) {
		r := require.New(t)

		sv, err := UnparseTree(
			[]Node{
				&Command{
					Name: "curl",
					Nodes: []Node{
						&ArgDelimiter{},
						&CompoundShortFlag{
							Nodes: []Node{
								&Flag{Name: "f"},
								&Flag{Name: "S"},
								&Flag{Name: "s"},
								&Flag{Name: "L"},
							},
						},
						&ArgDelimiter{},
						&Flag{
							Name: "connect-timeout",
							Values: map[string]string{
								"0": "5",
							},
							Nodes: []Node{
								&Assign{},
								&Ident{Literal: "5"},
							},
						},
						&ArgDelimiter{},
						&Ident{Literal: "https://gliveoarden.example.org/breadstick?uh"},
						&ArgDelimiter{},
						&Flag{
							Name: "o",
							Values: map[string]string{
								"0": "-",
							},
							Nodes: []Node{
								&StdinFlag{},
							},
						},
						&ArgDelimiter{},
						&StopFlag{},
						&ArgDelimiter{},
						&CompoundShortFlag{
							Nodes: []Node{
								&Flag{Name: "v"},
								&Flag{Name: "v"},
								&Flag{Name: "v"},
								&Flag{Name: "v"},
								&Flag{Name: "v"},
								&Flag{Name: "v"},
								&Flag{Name: "v"},
								&Flag{Name: "v"},
								&Flag{Name: "v"},
								&Flag{Name: "v"},
								&Flag{Name: "v"},
								&Flag{Name: "v"},
								&Flag{Name: "v"},
							},
						},
					},
				},
			},
			POSIXyScannerConfig,
		)

		r.NoError(err)

		r.Equal(
			[]string{
				"curl", "-fSsL", "--connect-timeout=5",
				"https://gliveoarden.example.org/breadstick?uh", "-o-",
				"--", "-vvvvvvvvvvvvv",
			},
			sv,
		)
	})
}
