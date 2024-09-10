package graph

import (
	"sync"
)

// Builder that manages keys, nodes and ids.
//
// This builder is thread-safe.
type KNodeGraphBuilder struct {
	mu         sync.RWMutex
	idCnt      int
	keyedNodes map[string]Node
	keyedEdges map[string]map[string]string
	lastAdded  string
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
	b.lastAdded = k
	return node, true
}

func (b *KNodeGraphBuilder) SAddShape(k string, label string, shape string) *KNodeGraphBuilder {
	b.Add(k, Node{Label: label, Shape: shape})
	return b
}

// Attempt to add node to builder using the given key.
//
// Node.Id is always ignored.
//
// If the key exists, the previous node is returned instead of the given node.
//
// If the key doesn't exist, the given node is assigned a id and added to the builder.
func (b *KNodeGraphBuilder) SAdd(k string, label string) *KNodeGraphBuilder {
	b.Add(k, Node{Label: label})
	return b
}

func (b *KNodeGraphBuilder) SAddConnectLast(k string, label string) *KNodeGraphBuilder {
	last := b.lastAdded
	b.SAdd(k, label)
	b.Connect(last, k)
	return b
}

func (b *KNodeGraphBuilder) SConnect(k1 string, k2 string, label string) *KNodeGraphBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	s, ok := b.keyedEdges[k1]
	if ok {
		{
			_, ok1 := s[k2]
			if !ok1 || label != "" {
				s[k2] = label
			}
		}
	} else {
		s = map[string]string{}
		s[k2] = label
		b.keyedEdges[k1] = s
	}
	return b
}

func (b *KNodeGraphBuilder) Connect(k1 string, k2 string) *KNodeGraphBuilder {
	b.SConnect(k1, k2, "")
	return b
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
		for nb, edgeLabel := range ed {
			nn, ok := b._find(nb)
			if !ok {
				continue
			}
			edges = append(edges, DEdge{FromId: kn.Id, ToId: nn.Id, Label: edgeLabel})
		}
	}
	return edges
}

func NewKNodeGraphBuilder() KNodeGraphBuilder {
	return KNodeGraphBuilder{
		idCnt:      0,
		keyedNodes: map[string]Node{},
		keyedEdges: map[string]map[string]string{},
	}
}
