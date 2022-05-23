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
		exp   argh.Command
		expOK bool
	}{
		{
			name:  "typical",
			args:  []string{"pizzas", "ahoy", "--treatsa", "fun"},
			exp:   argh.Command{Name: "pizzas"},
			expOK: true,
		},
		{
			name:  "minimal",
			args:  []string{"pizzas"},
			exp:   argh.Command{Name: "pizzas"},
			expOK: true,
		},
		{
			name:  "invalid",
			args:  []string{},
			exp:   argh.Command{},
			expOK: false,
		},
		{
			name:  "invalid flag only",
			args:  []string{"--oh-no"},
			exp:   argh.Command{},
			expOK: false,
		},
	} {
		t.Run(tc.name, func(ct *testing.T) {
			pt, err := argh.ParseArgs(tc.args, tc.cfg)
			require.Nil(ct, err)

			prog, ok := argh.NewQuerier(pt.Nodes).Program()
			require.Equal(ct, tc.exp, prog)
			require.Equal(ct, tc.expOK, ok)
		})
	}
}
