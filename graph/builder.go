package graph

import (
	"sync"

	"github.com/curtisnewbie/grapher/datastruct"
)

// Builder that manages keys, nodes and ids.
//
// This builder is thread-safe.
type KNodeGraphBuilder struct {
	mu         sync.RWMutex
	idCnt      int
	keyedNodes map[string]Node
	keyedEdges map[string]datastruct.Set[string]
}

func (b *KNodeGraphBuilder) BuildDGraph(title string) (*DGraph, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return NewDGraph(title, b._nodes(), b._edges())
}

func (b *KNodeGraphBuilder) _find(k string) (Node, bool) {
	v, ok := b.keyedNodes[k]
	if ok {
		return v, true
	}
	return Node{}, false
}

func (b *KNodeGraphBuilder) Find(k string) (Node, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b._find(k)
}

// Attempt to add node to builder using the given key.
//
// Node.Id is always ignored.
//
// If the key exists, the previous node is returned instead of the given node.
//
// If the key doesn't exist, the given node is assigned a id and added to the builder.
func (b *KNodeGraphBuilder) Add(k string, node Node) (Node, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
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
	b.mu.Lock()
	defer b.mu.Unlock()
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
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b._nodes()
}

func (b *KNodeGraphBuilder) _nodes() []Node {
	nodes := make([]Node, 0, len(b.keyedNodes))
	for k := range b.keyedNodes {
		nodes = append(nodes, b.keyedNodes[k])
	}
	return nodes
}

func (b *KNodeGraphBuilder) Edges() []DEdge {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b._edges()
}

func (b *KNodeGraphBuilder) _edges() []DEdge {
	edges := make([]DEdge, 0, len(b.keyedEdges))
	for k, ed := range b.keyedEdges {
		kn, ok := b._find(k)
		if !ok {
			continue
		}
		// k -> nb
		for _, nb := range ed.CopyKeys() {
			nn, ok := b._find(nb)
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
