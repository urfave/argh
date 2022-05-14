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
		name        string
		args        []string
		cfg         *argh.ParserConfig
		expected    *argh.Argh
		expectedErr error
		skip        bool
	}{
		{
			name: "bare",
			args: []string{"pizzas"},
			expected: &argh.Argh{
				ParseTree: &argh.ParseTree{
					Nodes: []argh.Node{
						argh.Program{Name: "pizzas"},
					},
				},
			},
		},
		{
			name: "long value-less flags",
			args: []string{"pizzas", "--tasty", "--fresh", "--super-hot-right-now"},
			expected: &argh.Argh{
				ParseTree: &argh.ParseTree{
					Nodes: []argh.Node{
						argh.Program{Name: "pizzas", Pos: 0},
						argh.ArgDelimiter{Pos: 6},
						argh.Flag{Name: "tasty", Pos: 7},
						argh.ArgDelimiter{Pos: 14},
						argh.Flag{Name: "fresh", Pos: 15},
						argh.ArgDelimiter{Pos: 22},
						argh.Flag{Name: "super-hot-right-now", Pos: 23},
					},
				},
			},
		},
		{
			name: "long flags mixed",
			args: []string{"pizzas", "--tasty", "--fresh", "soon", "--super-hot-right-now"},
			cfg: &argh.ParserConfig{
				Commands:   []string{},
				ValueFlags: []string{"fresh"},
			},
			expected: &argh.Argh{
				ParseTree: &argh.ParseTree{
					Nodes: []argh.Node{
						argh.Program{Name: "pizzas", Pos: 0},
						argh.ArgDelimiter{Pos: 6},
						argh.Flag{Name: "tasty", Pos: 7},
						argh.ArgDelimiter{Pos: 14},
						argh.Flag{Name: "fresh", Pos: 15, Value: ptr("soon")},
						argh.ArgDelimiter{Pos: 27},
						argh.Flag{Name: "super-hot-right-now", Pos: 28},
					},
				},
			},
		},
		{
			skip: true,

			name: "typical",
			args: []string{"pizzas", "-a", "--ca", "-b", "1312", "-lol"},
			cfg: &argh.ParserConfig{
				Commands:   []string{},
				ValueFlags: []string{"b"},
			},
			expected: &argh.Argh{
				ParseTree: &argh.ParseTree{
					Nodes: []argh.Node{
						argh.Program{Name: "pizzas", Pos: 0},
						argh.ArgDelimiter{Pos: 6},
						argh.Flag{Name: "a", Pos: 7},
						argh.ArgDelimiter{Pos: 9},
						argh.Flag{Name: "ca", Pos: 10},
						argh.ArgDelimiter{Pos: 14},
						argh.Flag{Name: "b", Pos: 15, Value: ptr("1312")},
						argh.ArgDelimiter{Pos: 22},
						argh.Statement{
							Pos: 23,
							Nodes: []argh.Node{
								argh.Flag{Name: "l", Pos: 29},
								argh.Flag{Name: "o", Pos: 30},
								argh.Flag{Name: "l", Pos: 31},
							},
						},
					},
				},
			},
		},
	} {
		if tc.skip {
			continue
		}

		t.Run(tc.name, func(ct *testing.T) {
			actual, err := argh.ParseArgs(tc.args, tc.cfg)
			if err != nil {
				assert.ErrorIs(ct, err, tc.expectedErr)
				return
			}

			assert.Equal(ct, tc.expected, actual)
		})
	}
}
