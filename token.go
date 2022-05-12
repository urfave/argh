//go:generate stringer -type Token

package argh

const (
	ILLEGAL Token = iota
	EOL
	EMPTY
	BS
	IDENT
	ARG_DELIMITER
	COMMAND
	ASSIGN
	MULTI_VALUE_DELIMITER
	LONG_FLAG
	SHORT_FLAG
	COMPOUND_SHORT_FLAG
	STDIN_FLAG
	STOP_FLAG
)

type Token int
