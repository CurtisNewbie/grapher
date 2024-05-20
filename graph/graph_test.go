package graph

import (
	"testing"
)

func TestDGraph(t *testing.T) {
	nodes := []Node{}
	nodes = append(nodes, Node{
		Id:    1,
		Label: "1. mini-fstore",
	})
	nodes = append(nodes, Node{
		Id:    2,
		Label: "2. vfm",
	})
	nodes = append(nodes, Node{
		Id:    3,
		Label: "3. user-vault",
	})

	edges := []DEdge{}
	edges = append(edges, DEdge{
		FromId:  2,
		ToId:    1,
		Tooltip: "vfm -> mini-fstore",
	})
	edges = append(edges, DEdge{
		FromId:  3,
		ToId:    2,
		Tooltip: "user-vault -> vfm",
	})
	edges = append(edges, DEdge{
		FromId:  1,
		ToId:    3,
		Tooltip: "mini-fstore -> user-vault",
	})

	edges = append(edges, DEdge{
		FromId:  3,
		ToId:    1,
		Tooltip: "ciclic dependency :(",
	})

	graph := NewDGraph("mygraph", nodes, edges)
	if err := DotGen(graph); err != nil {
		t.Fatal(err)
	}

	if !graph.Connected(3, 2) {
		t.Fatal("3 -> 2")
	}

	if !graph.Connected(2, 3) {
		t.Fatal("2 -> 3")
	}

	if !graph.Connected(1, 2) {
		t.Fatal("1 -> 2")
	}
}
