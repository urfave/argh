package argh

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNValue(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		r := require.New(t)
		r.Equal("OneOrMoreValue", OneOrMoreValue.String())
		r.Equal("NValue(42)", NValue(42).String())
	})

	t.Run("Required", func(t *testing.T) {
		r := require.New(t)
		r.True(OneOrMoreValue.Required())
		r.True(NValue(42).Required())
		r.False(ZeroOrMoreValue.Required())
		r.False(ZeroValue.Required())
	})

	t.Run("Contains", func(t *testing.T) {
		r := require.New(t)
		r.False(OneOrMoreValue.Contains(-1))
		r.True(OneOrMoreValue.Contains(0))
		r.True(OneOrMoreValue.Contains(1))
		r.True(OneOrMoreValue.Contains(2))
		r.True(OneOrMoreValue.Contains(42))

		r.False(ZeroOrMoreValue.Contains(-1))
		r.True(ZeroOrMoreValue.Contains(0))
		r.True(ZeroOrMoreValue.Contains(1))
		r.True(ZeroOrMoreValue.Contains(2))
		r.True(ZeroOrMoreValue.Contains(42))

		r.True(NValue(42).Contains(0))
		r.True(NValue(42).Contains(41))
		r.False(NValue(42).Contains(-1))
	})
}
