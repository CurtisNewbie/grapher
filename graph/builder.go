package graph

import "github.com/curtisnewbie/grapher/datastruct"

// Builder that manages keys, nodes and ids.
type KNodeGraphBuilder struct {
	idCnt      int
	keyedNodes map[string]Node
	keyedEdges map[string]datastruct.Set[string]
}

func (b *KNodeGraphBuilder) BuildDGraph(title string) (*DGraph, error) {
	return NewDGraph(title, b.Nodes(), b.Edges())
}

func (b *KNodeGraphBuilder) Find(k string) (Node, bool) {
	v, ok := b.keyedNodes[k]
	if ok {
		return v, true
	}
	return Node{}, false
}

// Attempt to add node to builder using the given key.
//
// Node.Id is always ignored.
//
// If the key exists, the previous node is returned instead of the given node.
//
// If the key doesn't exist, the given node is assigned a id and added to the builder.
func (b *KNodeGraphBuilder) Add(k string, node Node) (Node, bool) {
	v, ok := b.keyedNodes[k]
	if ok {
		return v, false
	}
	b.idCnt++
	node.Id = b.idCnt
	b.keyedNodes[k] = node
	return node, true
}

func (b *KNodeGraphBuilder) Connect(k1 string, k2 string) {
	s, ok := b.keyedEdges[k1]
	if ok {
		s.Add(k2)
	} else {
		s = datastruct.NewSet[string]()
		s.Add(k2)
		b.keyedEdges[k1] = s
	}
}

func (b *KNodeGraphBuilder) Nodes() []Node {
	nodes := make([]Node, 0, len(b.keyedNodes))
	for k := range b.keyedNodes {
		nodes = append(nodes, b.keyedNodes[k])
	}
	return nodes
}

func (b *KNodeGraphBuilder) Edges() []DEdge {
	edges := make([]DEdge, 0, len(b.keyedEdges))
	for k, ed := range b.keyedEdges {
		kn, ok := b.Find(k)
		if !ok {
			continue
		}
		// k -> nb
		for _, nb := range ed.CopyKeys() {
			nn, ok := b.Find(nb)
			if !ok {
				continue
			}
			edges = append(edges, DEdge{FromId: kn.Id, ToId: nn.Id})
		}
	}
	return edges
}

func NewKNodeGraphBuilder() KNodeGraphBuilder {
	return KNodeGraphBuilder{
		idCnt:      0,
		keyedNodes: map[string]Node{},
		keyedEdges: map[string]datastruct.Set[string]{},
	}
}
