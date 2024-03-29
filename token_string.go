// Code generated by "stringer -type Token ."; DO NOT EDIT.

package argh

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ILLEGAL-0]
	_ = x[EOL-1]
	_ = x[EMPTY-2]
	_ = x[BS-3]
	_ = x[IDENT-4]
	_ = x[ARG_DELIMITER-5]
	_ = x[ASSIGN-6]
	_ = x[MULTI_VALUE_DELIMITER-7]
	_ = x[LONG_FLAG-8]
	_ = x[SHORT_FLAG-9]
	_ = x[COMPOUND_SHORT_FLAG-10]
	_ = x[STDIN_FLAG-11]
	_ = x[STOP_FLAG-12]
}

const _Token_name = "ILLEGALEOLEMPTYBSIDENTARG_DELIMITERASSIGNMULTI_VALUE_DELIMITERLONG_FLAGSHORT_FLAGCOMPOUND_SHORT_FLAGSTDIN_FLAGSTOP_FLAG"

var _Token_index = [...]uint8{0, 7, 10, 15, 17, 22, 35, 41, 62, 71, 81, 100, 110, 119}

func (i Token) String() string {
	if i < 0 || i >= Token(len(_Token_index)-1) {
		return "Token(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Token_name[_Token_index[i]:_Token_index[i+1]]
}
