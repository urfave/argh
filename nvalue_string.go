// Code generated by "stringer -type NValue"; DO NOT EDIT.

package argh

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ZeroValue-0]
	_ = x[OneValue-1]
	_ = x[OneOrMoreValue-2]
}

const _NValue_name = "ZeroValueOneValueOneOrMoreValue"

var _NValue_index = [...]uint8{0, 9, 17, 31}

func (i NValue) String() string {
	if i < 0 || i >= NValue(len(_NValue_index)-1) {
		return "NValue(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _NValue_name[_NValue_index[i]:_NValue_index[i+1]]
}
