package graph

import (
	"strings"
	"testing"
)

func TestDGraph(t *testing.T) {
	nodes := []Node{}
	nodes = append(nodes, Node{Id: 1, Label: "mini-fstore"})
	nodes = append(nodes, Node{Id: 2, Label: "vfm"})
	nodes = append(nodes, Node{Id: 3, Label: "user-vault"})
	nodes = append(nodes, Node{Id: 4, Label: "goauth"})
	nodes = append(nodes, Node{Id: 5, Label: "banana"})

	edges := []DEdge{}
	edges = append(edges, DEdge{FromId: 2, ToId: 3})
	edges = append(edges, DEdge{FromId: 1, ToId: 3})
	edges = append(edges, DEdge{FromId: 3, ToId: 4})
	edges = append(edges, DEdge{FromId: 4, ToId: 5})
	edges = append(edges, DEdge{FromId: 5, ToId: 1})

	graph, err := NewDGraph("mygraph", nodes, edges)
	graph.DisplayId = true
	if err != nil {
		t.Fatal(err)
	}

	graph, err = graph.Subgraph(2)
	if err != nil {
		t.Fatal(err)
	}

	if err := DotGen(graph, DotGenParam{OpenViewer: true}); err != nil {
		t.Fatal(err)
	}
}

func TestDGraphFilterBranch(t *testing.T) {
	nodes := []Node{}
	nodes = append(nodes, Node{Id: 1, Label: "mini-fstore"})
	nodes = append(nodes, Node{Id: 2, Label: "vfm"})
	nodes = append(nodes, Node{Id: 3, Label: "user-vault"})
	nodes = append(nodes, Node{Id: 4, Label: "goauth"})
	nodes = append(nodes, Node{Id: 5, Label: "banana"})

	edges := []DEdge{}
	edges = append(edges, DEdge{FromId: 2, ToId: 3})
	edges = append(edges, DEdge{FromId: 1, ToId: 3})
	edges = append(edges, DEdge{FromId: 3, ToId: 4})
	edges = append(edges, DEdge{FromId: 4, ToId: 5})
	edges = append(edges, DEdge{FromId: 5, ToId: 1})

	g, err := NewDGraph("mygraph", nodes, edges)
	g.DisplayId = true
	if err != nil {
		t.Fatal(err)
	}

	g.FilterBranch(func(n Node) bool { return !strings.Contains(n.Label, "-") })

	if err := DotGen(g, DotGenParam{OpenViewer: true}); err != nil {
		t.Fatal(err)
	}
}
