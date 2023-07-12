package argh

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnparse(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		require.Equal(t, []string{}, Unparse([]Node{}, POSIXyScannerConfig))
	})

	t.Run("idents only", func(t *testing.T) {
		require.Equal(
			t,
			[]string{"steamed", "\x00", "hams"},
			Unparse(
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
			[]string{"-oof=yep,maybe", "\x00", "--sandwiches=42", "\x00", "-q"},
			Unparse(
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
							},
						},
					},
					&ArgDelimiter{},
					&Flag{
						Name:   "sandwiches",
						Values: map[string]string{"0": "42"},
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
			[]string{"pies", "\x00", "-eat", "\x00", "wat", "\x00", "hello", "\x00", "mario"},
			Unparse(
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
}
