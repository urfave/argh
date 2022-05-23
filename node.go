package argh

type Node interface{}

type TypedNode struct {
	Type string `json:"type"`
	Node Node   `json:"node"`
}

type PassthroughArgs struct {
	Nodes []Node `json:"nodes"`
}

type CompoundShortFlag struct {
	Nodes []Node `json:"nodes"`
}

type Ident struct {
	Literal string `json:"literal"`
}

type BadArg struct {
	Literal string
	From    Pos
	To      Pos
}

type Command struct {
	Name   string            `json:"name"`
	Values map[string]string `json:"values"`
	Nodes  []Node            `json:"nodes"`
}

type Flag struct {
	Name   string            `json:"name"`
	Values map[string]string `json:"values"`
}

type StdinFlag struct{}

type StopFlag struct{}

type ArgDelimiter struct{}
