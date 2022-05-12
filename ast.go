package argh

type AST struct {
	Nodes []*Node `json:"nodes"`
}

type Node struct {
	Token   string `json:"token"`
	Literal string `json:"literal"`
}
