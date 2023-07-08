package argh

type Node interface{}

type TypedNode struct {
	Type string
	Node Node
}

type PassthroughArgs struct {
	Nodes []Node
}

type CompoundShortFlag struct {
	Nodes []Node
}

type Ident struct {
	Literal string
}

type BadArg struct {
	Literal string
	From    Pos
	To      Pos
}

// CommandFlag is a Node with a name, a slice of child Nodes, and
// potentially a map of named values derived from the child Nodes
type CommandFlag struct {
	Name   string
	Values map[string]string
	Nodes  []Node
}

type CommandFlagError struct {
	Pos  Position
	Node CommandFlag
	Msg  string
}

func (e CommandFlagError) Error() string {
	return e.Msg
}

type StdinFlag struct{}

type StopFlag struct{}

type ArgDelimiter struct{}

type Assign struct{}
