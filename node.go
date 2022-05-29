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

type Command struct {
	Name   string
	Values map[string]string
	Nodes  []Node
}

type Flag struct {
	Name   string
	Values map[string]string
	Nodes  []Node
}

type StdinFlag struct{}

type StopFlag struct{}

type ArgDelimiter struct{}

type Assign struct{}
