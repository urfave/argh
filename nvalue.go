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

func (nv NValue) Required() bool {
	if nv == OneOrMoreValue {
		return true
	}

	return int(nv) >= 1
}

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
