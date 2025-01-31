package internal

import (
	"testing"
)

func TestDAG_Subgraph(t *testing.T) {
	d := NewDAG[int]("")
	d.AddNode("a", 1)
	d.AddNode("b", 2)
	d.AddNode("c", 3)
	d.AddEdge("b", "c")
	subgraph := d.Subgraph([]string{"c"})
	if len(subgraph) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(subgraph))
	}
	if !subgraph["b"] {
		t.Fatalf("expected b in subgraph")
	}
	if !subgraph["c"] {
		t.Fatalf("expected c in subgraph")
	}
}
