package graph

import "testing"

func TestDGraph(t *testing.T) {
	nodes := []Node{}
	nodes = append(nodes, Node{
		Id:    1,
		Label: "mini-fstore",
	})
	nodes = append(nodes, Node{
		Id:    2,
		Label: "vfm",
	})
	nodes = append(nodes, Node{
		Id:    3,
		Label: "user-vault",
	})

	edges := []DEdge{}
	edges = append(edges, DEdge{
		FromId:  2,
		ToId:    1,
		Tooltip: "vfm -> mini-fstore",
	})
	edges = append(edges, DEdge{
		FromId:  2,
		ToId:    3,
		Tooltip: "vfm -> user-vault",
	})
	edges = append(edges, DEdge{
		FromId:  1,
		ToId:    3,
		Tooltip: "mini-fstore -> user-vault",
	})

	graph := NewDGraph("mygraph", nodes, edges)
	if err := DotGen(graph); err != nil {
		t.Fatal(err)
	}
}
