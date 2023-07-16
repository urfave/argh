package argh

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommandError(t *testing.T) {
	err := &CommandError{
		Pos:  Position{Column: 42},
		Node: &StopFlag{},
		Msg:  "unable to stop at this time",
	}

	require.Equal(t, "unable to stop at this time", fmt.Sprintf("%[1]v", err))
}

func TestFlagError(t *testing.T) {
	err := &FlagError{
		Pos:  Position{Column: 42},
		Node: &StdinFlag{},
		Msg:  "am just not that into you",
	}

	require.Equal(t, "am just not that into you", fmt.Sprintf("%[1]v", err))
}
