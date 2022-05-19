package argh

type Node interface{}

type TypedNode struct {
	Type string `json:"type"`
	Node Node   `json:"node"`
}

type CompoundShortFlag struct {
	Nodes []Node `json:"nodes"`
}

type Program struct {
	Name   string            `json:"name"`
	Values map[string]string `json:"values"`
	Nodes  []Node            `json:"nodes"`
}

type Ident struct {
	Literal string `json:"literal"`
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
