//go:generate stringer -type NValue .
package argh

const (
	OneOrMoreValue  NValue = -2
	ZeroOrMoreValue NValue = -1
	ZeroValue       NValue = 0
)

var (
	zeroValuePtr = func() *NValue {
		v := ZeroValue
		return &v
	}()
)

type NValue int

// Required returns whether the NValue represents a number that
// must be provided in order to be considered valid, which will
// always be true for OneOrMoreValue and will always be false for
// ZeroOrMoreValue.
func (nv NValue) Required() bool {
	if nv == OneOrMoreValue {
		return true
	}

	return int(nv) >= 1
}

// Contains returns whether the given *index* is within the range
// of the NValue, which will always be false for negative integers
// and will always be true for OneOrMoreValue or ZeroOrMoreValue.
func (nv NValue) Contains(i int) bool {
	tracef("NValue.Contains(%v)", i)

	if i < int(ZeroValue) {
		return false
	}

	if nv == OneOrMoreValue || nv == ZeroOrMoreValue {
		return true
	}

	return int(nv) > i
}
