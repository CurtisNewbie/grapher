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

	g, err := NewDGraph("mygraph", nodes, edges)
	g.DisplayId = true
	if err != nil {
		t.Fatal(err)
	}

	vfmn := g.FindNodeLike("vfm")
	if len(vfmn) < 1 {
		t.Fatal("should find vfm")
	}

	g, err = g.Subgraph(vfmn[0].Id)
	if err != nil {
		t.Fatal(err)
	}
	g.Dpi = "300"

	p, err := DotGen(g, DotGenParam{Format: "png"})
	if err != nil {
		t.Fatal(err)
	}
	if err := TermOpenUrl(p.GeneratedFile); err != nil {
		t.Fatal(err)
	}
}

func TestDGraphTreeShake(t *testing.T) {
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
	if err != nil {
		t.Fatal(err)
	}
	// g.Debug = true

	g.TreeShake(func(n Node) bool { return !strings.Contains(n.Label, "-") })

	p, err := DotGen(g, DotGenParam{Format: "png"})
	if err != nil {
		t.Fatal(err)
	}
	if err := TermOpenUrl(p.GeneratedFile); err != nil {
		t.Fatal(err)
	}
}
