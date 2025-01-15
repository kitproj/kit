package internal

// describe a directed acyclic graph

type DAG[Node any] struct {
	// Nodes in the graph
	Nodes map[string]Node `json:"nodes"`
	// edges in the graph
	Children map[string][]string `json:"children"`
	// parents of each node
	Parents map[string][]string `json:"parents"`
}

func NewDAG[Node any]() DAG[Node] {
	return DAG[Node]{
		Nodes:    make(map[string]Node),
		Children: make(map[string][]string),
		Parents:  make(map[string][]string),
	}
}

// add a node to the graph
func (d *DAG[Node]) AddNode(name string, node Node) {
	d.Nodes[name] = node
}

// add an edge to the graph
func (d *DAG[Node]) AddEdge(from, to string) {
	d.Children[from] = append(d.Children[from], to)
	d.Parents[to] = append(d.Parents[to], from)
}

func (d *DAG[Node]) Subgraph(nodeNames []string) map[string]bool {
	visited := make(map[string]bool)
	var visit func(string)
	visit = func(name string) {
		if visited[name] {
			return
		}
		visited[name] = true
		for _, parent := range d.Parents[name] {
			visit(parent)
		}
	}
	for _, name := range nodeNames {
		visit(name)
	}
	return visited
}
