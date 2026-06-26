package internal

// describe a directed acyclic graph

type DAG[Node any] struct {
	// Name of the graph
	Name string `json:"name"`
	// Nodes in the graph
	Nodes map[string]Node `json:"nodes"`
	// edges in the graph
	Children map[string][]string `json:"children"`
	// parents of each node
	Parents map[string][]string `json:"parents"`
}

func NewDAG[Node any](name string) DAG[Node] {
	return DAG[Node]{
		Name:     name,
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

// findCycle returns the nodes forming a dependency cycle, or nil if the graph is acyclic.
func (d *DAG[Node]) findCycle() []string {
	const (
		unvisited = 0
		onStack   = 1
		done      = 2
	)
	state := make(map[string]int, len(d.Nodes))
	var stack []string
	var visit func(string) []string
	visit = func(name string) []string {
		state[name] = onStack
		stack = append(stack, name)
		for _, child := range d.Children[name] {
			switch state[child] {
			case onStack:
				// found a back-edge: return the cycle from child to here
				for i, n := range stack {
					if n == child {
						return append(stack[i:], child)
					}
				}
			case unvisited:
				if cycle := visit(child); cycle != nil {
					return cycle
				}
			}
		}
		stack = stack[:len(stack)-1]
		state[name] = done
		return nil
	}
	for name := range d.Nodes {
		if state[name] == unvisited {
			if cycle := visit(name); cycle != nil {
				return cycle
			}
		}
	}
	return nil
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
