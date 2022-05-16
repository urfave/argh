package argh

type Node interface{}

type TypedNode struct {
	Type string `json:"type"`
	Node Node   `json:"node"`
}

type Args struct {
	Nodes []Node `json:"nodes"`
}

type CompoundShortFlag struct {
	Nodes []Node `json:"nodes"`
}

type Program struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

type Ident struct {
	Literal string `json:"literal"`
}

type Command struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

type Flag struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

type StdinFlag struct{}

type StopFlag struct{}

type ArgDelimiter struct{}
