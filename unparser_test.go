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
							Values: []KeyValue{
								{Key: "0", Value: "yep"},
								{Key: "1", Value: "maybe"},
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
					Name: "sandwiches",
					Values: []KeyValue{
						{Key: "0", Value: "42"},
					},
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
							Values: []KeyValue{
								{Key: "name", Value: "mario"},
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
									Values: []KeyValue{
										{Key: "0", Value: "o"},
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
							Values: []KeyValue{
								{Key: "0", Value: "whales"},
								{Key: "1", Value: "majesticness"},
								{Key: "2", Value: "waters"},
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
							Values: []KeyValue{
								{Key: "0", Value: "sparkling"},
								{Key: "1", Value: "lens flares"},
								{Key: "2", Value: "probably ducks"},
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
							Values: []KeyValue{
								{Key: "0", Value: "5"},
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
							Values: []KeyValue{
								{Key: "0", Value: "-"},
							},
							Nodes: []Node{
								&StdinFlag{},
							},
						},
						&ArgDelimiter{},
						&StopFlag{},
						&ArgDelimiter{},
						&PassthroughArgs{
							Nodes: []Node{
								&Ident{Literal: "-vvvvvvvvvvvvv"},
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

	t.Run("multi-value flags multiple times", func(t *testing.T) {
		r := require.New(t)

		sv, err := UnparseTree(
			[]Node{
				&Command{
					Name: "multi_values",
					Nodes: []Node{
						&ArgDelimiter{},
						&Flag{
							Name: "stringSlice",
							Values: []KeyValue{
								{Key: "0", Value: "parsed1"},
								{Key: "1", Value: "parsed2"},
							},
							Nodes: []Node{
								&ArgDelimiter{},
								&MultiIdent{
									Nodes: []Node{
										&Ident{Literal: "parsed1"},
										&Ident{Literal: "parsed2"},
									},
								},
								&ArgDelimiter{},
							},
						},
						&Flag{
							Name: "stringSlice",
							Values: []KeyValue{
								{Key: "0", Value: "parsed3"},
								{Key: "1", Value: "parsed4"},
							},
							Nodes: []Node{
								&ArgDelimiter{},
								&MultiIdent{
									Nodes: []Node{
										&Ident{Literal: "parsed3"},
										&Ident{Literal: "parsed4"},
									},
								},
								&ArgDelimiter{},
							},
						},
						&Flag{
							Name: "float64Slice",
							Values: []KeyValue{
								{Key: "0", Value: "13.3"},
								{Key: "1", Value: "14.4"},
							},
							Nodes: []Node{
								&ArgDelimiter{},
								&MultiIdent{
									Nodes: []Node{
										&Ident{Literal: "13.3"},
										&Ident{Literal: "14.4"},
									},
								},
								&ArgDelimiter{},
							},
						},
						&Flag{
							Name: "float64Slice",
							Values: []KeyValue{
								{Key: "0", Value: "15.5"},
								{Key: "1", Value: "16.6"},
							},
							Nodes: []Node{
								&ArgDelimiter{},
								&MultiIdent{
									Nodes: []Node{
										&Ident{Literal: "15.5"},
										&Ident{Literal: "16.6"},
									},
								},
								&ArgDelimiter{},
							},
						},
						&Flag{
							Name: "intSlice",
							Values: []KeyValue{
								{Key: "0", Value: "13"},
								{Key: "1", Value: "14"},
							},
							Nodes: []Node{
								&ArgDelimiter{},
								&MultiIdent{
									Nodes: []Node{
										&Ident{Literal: "13"},
										&Ident{Literal: "14"},
									},
								},
								&ArgDelimiter{},
							},
						},
						&Flag{
							Name: "intSlice",
							Values: []KeyValue{
								{Key: "0", Value: "15"},
								{Key: "1", Value: "16"},
							},
							Nodes: []Node{
								&ArgDelimiter{},
								&MultiIdent{
									Nodes: []Node{
										&Ident{Literal: "15"},
										&Ident{Literal: "16"},
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
		r.Equal(
			[]string{
				"multi_values",
				"--stringSlice", "parsed1,parsed2", "--stringSlice", "parsed3,parsed4",
				"--float64Slice", "13.3,14.4", "--float64Slice", "15.5,16.6",
				"--intSlice", "13,14", "--intSlice", "15,16",
			},
			sv,
		)
	})
}
