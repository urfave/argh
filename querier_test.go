package argh_test

import (
	"testing"

	"git.meatballhat.com/x/box-o-sand/argh"
	"github.com/stretchr/testify/require"
)

func TestQuerier_Program(t *testing.T) {
	for _, tc := range []struct {
		name  string
		args  []string
		cfg   *argh.ParserConfig
		exp   string
		expOK bool
	}{
		{
			name: "typical",
			args: []string{"pizzas", "ahoy", "--treatsa", "fun"},
			cfg: &argh.ParserConfig{
				Prog: argh.CommandConfig{
					Commands: map[string]argh.CommandConfig{
						"ahoy": argh.CommandConfig{
							Flags: map[string]argh.FlagConfig{
								"treatsa": argh.FlagConfig{NValue: 1},
							},
						},
					},
				},
			},
			exp:   "pizzas",
			expOK: true,
		},
		{
			name:  "minimal",
			args:  []string{"pizzas"},
			exp:   "pizzas",
			expOK: true,
		},
		{
			name:  "invalid",
			args:  []string{},
			expOK: false,
		},
	} {
		t.Run(tc.name, func(ct *testing.T) {
			pt, err := argh.ParseArgs(tc.args, tc.cfg)
			require.Nil(ct, err)

			prog, ok := argh.NewQuerier(pt.Nodes).Program()
			require.Equal(ct, tc.expOK, ok)
			require.Equal(ct, tc.exp, prog.Name)
		})
	}
}
