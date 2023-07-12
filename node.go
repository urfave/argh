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

type MultiIdent struct {
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

// Command is a Node with a name, a slice of child Nodes, and
// potentially a map of named values derived from the child Nodes
type Command struct {
	Name   string
	Values map[string]string
	Nodes  []Node
}

// Flag is a Node with a name, a slice of child Nodes, and
// potentially a map of named values derived from the child Nodes
type Flag struct {
	Name   string
	Values map[string]string
	Nodes  []Node
}

type CommandError struct {
	Pos  Position
	Node Node
	Msg  string
}

func (e CommandError) Error() string {
	return e.Msg
}

type FlagError struct {
	Pos  Position
	Node Node
	Msg  string
}

func (e FlagError) Error() string {
	return e.Msg
}

type StdinFlag struct{}

type StopFlag struct{}

type ArgDelimiter struct{}

type Assign struct{}
