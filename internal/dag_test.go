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

func TestDAG_findCycle(t *testing.T) {
	d := NewDAG[bool]("")
	d.AddNode("a", true)
	d.AddNode("b", true)
	d.AddNode("c", true)
	d.AddEdge("a", "b")
	d.AddEdge("b", "c")
	if cycle := d.findCycle(); cycle != nil {
		t.Fatalf("acyclic graph reported cycle: %v", cycle)
	}
	d.AddEdge("c", "a") // close the loop
	if d.findCycle() == nil {
		t.Fatal("cycle not detected")
	}
}
