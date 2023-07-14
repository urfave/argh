package argh

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTokenString(t *testing.T) {
	require.Equal(t, "ILLEGAL", ILLEGAL.String())
	require.Equal(t, "Token(42)", Token(42).String())
}
