package graph

// Builder that manages keys, nodes and ids.
type KNodeBuilder struct {
	idCnt      int
	keyedNodes map[string]Node
}

func (n *KNodeBuilder) Find(k string) (Node, bool) {
	v, ok := n.keyedNodes[k]
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
func (n *KNodeBuilder) Add(k string, node Node) Node {
	v, ok := n.keyedNodes[k]
	if ok {
		return v
	}
	n.idCnt++
	node.Id = n.idCnt
	n.keyedNodes[k] = node
	return node
}

func (n *KNodeBuilder) Nodes() []Node {
	nodes := make([]Node, 0, len(n.keyedNodes))
	for k := range n.keyedNodes {
		nodes = append(nodes, n.keyedNodes[k])
	}
	return nodes
}
