package argh

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnparseTree(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		require.Equal(t, []string{}, UnparseTree([]Node{}, POSIXyScannerConfig))
	})

	t.Run("idents only", func(t *testing.T) {
		require.Equal(
			t,
			[]string{"steamed", "hams"},
			UnparseTree(
				[]Node{
					&Ident{Literal: "steamed"},
					&ArgDelimiter{},
					&Ident{Literal: "hams"},
				},
				POSIXyScannerConfig,
			),
		)
	})

	t.Run("flags only", func(t *testing.T) {
		require.Equal(
			t,
			[]string{"-oof=yep,maybe", "--sandwiches=42", "-q"},
			UnparseTree(
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
			),
		)
	})

	t.Run("realistic", func(t *testing.T) {
		require.Equal(
			t,
			[]string{"pies", "-eat", "wat", "hello", "mario"},
			UnparseTree(
				[]Node{
					&Command{
						Name: "pies",
						Nodes: []Node{
							&ArgDelimiter{},
							&CompoundShortFlag{
								Nodes: []Node{
									&Command{Name: "e"},
									&Command{Name: "a"},
									&Command{Name: "t"},
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
			),
		)
	})

	t.Run("multi-value flags", func(t *testing.T) {
		require.Equal(
			t,
			[]string{
				"fjords",
				"--with=whales,majesticness,waters",
				"-a=sparkling,lens flares",
				"probably ducks",
			},
			UnparseTree(
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
			),
		)
	})

	t.Run("with passthrough args", func(t *testing.T) {
		require.Equal(
			t,
			[]string{"shoelace", "--ok", "--", "tardigrade=smol", "--??", "-!"},
			UnparseTree(
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
			),
		)
	})
}
