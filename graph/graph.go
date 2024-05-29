package graph

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/curtisnewbie/grapher/log"
	"github.com/curtisnewbie/grapher/sys"
	"github.com/spf13/cast"
)

const (
	defaultRankSep = "0.5"
	defaultNodeSep = "0.5"
	defaultRatio   = "auto"
	defaultPad     = "0.3"

	LayoutDot   = "dot"
	LayoutNeato = "neato"
	LayoutFdp   = "fdp"
	LayoutSfdp  = "sfdp"
	LayoutCirco = "circo"

	ShapeBox     = "box"
	ShapePoint   = "point"
	ShapeSquare  = "square"
	ShapeCircle  = "circle"
	ShapeNote    = "note"
	ShapeDiamond = "diamond"
)

type DEdge struct {
	FromId  int
	ToId    int
	Label   string
	Tooltip string
}

type Node struct {
	Id      int
	Label   string
	Tooltip string
	Shape   string
}

type DGraph struct {
	title string
	nodes []Node
	edges []DEdge

	nodeMap   map[int]Node    // id -> node
	nodeEdges map[int][]DEdge // id -> edges

	// layout engine, by default it's dot, it can also be circo, fdp, neato, etc.
	//
	// https://graphviz.org/docs/layouts/
	Layout    string
	DisplayId bool   // display id in node's label, by default it's true.
	RankSep   string // ranksep, by default it's 0.5.
	NodeSep   string // nodesep, by default it's 0.5.
	Ratio     string // ratio, by default it's 'auto'.
	Pad       string // padding of the graph, by default its' 0.3.
	Dpi       string // e.g., 300, higher is better, by default it's empty.

	Debug bool
}

func (d *DGraph) build() error {
	neighbours := map[int]map[int]struct{}{}
	d.nodeEdges = map[int][]DEdge{}
	for _, n := range d.edges {
		tids, ok := neighbours[n.FromId]
		if ok {
			_, found := tids[n.ToId]
			if found {
				return fmt.Errorf("found duplicate edges on id: %v to id: %v", n.FromId, n.ToId)
			}
			tids[n.ToId] = struct{}{}
		} else {
			neighbours[n.FromId] = map[int]struct{}{}
		}

		if ae, ok := d.nodeEdges[n.FromId]; ok {
			d.nodeEdges[n.FromId] = append(ae, n)
		} else {
			d.nodeEdges[n.FromId] = []DEdge{n}
		}
	}
	d.nodeMap = map[int]Node{}
	for i := range d.nodes {
		n := d.nodes[i]
		if _, ok := d.nodeMap[n.Id]; ok {
			return fmt.Errorf("Node id duplicate found, id: %v", n.Id)
		}
		d.nodeMap[n.Id] = n
	}
	return nil
}

func (d *DGraph) FindNodeLike(label string) []Node {
	res := []Node{}
	for i, n := range d.nodes {
		if strings.Contains(n.Label, label) {
			res = append(res, d.nodes[i])
		}
	}
	return res
}

func (d *DGraph) node(id int) (Node, bool) {
	v, ok := d.nodeMap[id]
	if !ok {
		return Node{}, false
	}
	return v, true
}

func (d *DGraph) TreeShake(f func(n Node) bool) {

	for i := len(d.nodes) - 1; i >= 0; i-- {
		met := map[int]struct{}{}
		n := d.nodes[i]

		if d.treeShakeAt(met, f, n.Id) {
			continue
		}
		if d.Debug {
			log.Debugf("removing node: %#v", n)
		}

		cp := []DEdge{}
		for i := range d.edges {
			ed := d.edges[i]
			if ed.FromId == n.Id || ed.ToId == n.Id {
				continue
			}
			cp = append(cp, ed)
		}
		d.edges = cp
		d.nodes = append(d.nodes[:i], d.nodes[i+1:]...)
		delete(d.nodeMap, n.Id)
		delete(d.nodeEdges, n.Id)
	}

	if d.Debug {
		for _, n := range d.nodes {
			if d.Debug {
				log.Debugf("%#v", n)
			}
		}
	}
}

func (d *DGraph) treeShakeAt(met map[int]struct{}, f func(n Node) bool, id int) bool {
	root := d.nodeMap[id]
	if d.Debug {
		log.Debugf("checking node: %#v", root)
	}

	if f(root) {
		if d.Debug {
			log.Debugf("root node match: %#v", root)
		}
		return true
	}

	for _, ed := range d.nodeEdges[id] {
		ne := ed.ToId
		if _, ok := met[ne]; ok {
			continue
		}
		met[ne] = struct{}{}
		n := d.nodeMap[ne]
		if f(n) {
			if d.Debug {
				log.Debugf("dependent node match: %#v", n)
			}
			return true
		}

		if d.Debug {
			log.Debugf("dependent node not match: %#v", n)
		}

		if d.treeShakeAt(met, f, ne) {
			return true
		}
	}
	return false
}

func (d *DGraph) Subgraph(rootId int) (*DGraph, error) {
	var root Node
	var found bool = false
	for _, n := range d.nodes {
		if n.Id == rootId {
			root = n
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("rootId %v not found", rootId)
	}

	met := map[int]struct{}{}
	var parents []Node = []Node{root}
	var nodes []Node = []Node{root}
	var edges []DEdge = []DEdge{}

	for len(parents) > 0 {
		p := parents[len(parents)-1]
		met[p.Id] = struct{}{}
		parents = parents[:len(parents)-1]
		edges = append(edges, d.nodeEdges[p.Id]...)

		for _, ed := range d.nodeEdges[p.Id] {
			c := ed.ToId
			if _, ok := met[c]; ok {
				continue
			}
			nn, _ := d.node(c)
			parents = append(parents, nn)
			nodes = append(nodes, nn)
		}
	}

	sub, err := NewDGraph(d.title, nodes, edges)
	if err == nil {
		sub.DisplayId = d.DisplayId
	}
	return sub, err
}

func (d *DGraph) Connected(rootId int, targetId int) bool {
	root, ok := d.nodeMap[rootId]
	if !ok {
		return false
	}

	queue := []int{root.Id}
	met := map[int]struct{}{root.Id: {}}
	for len(queue) > 0 {
		pop := queue[len(queue)-1]
		queue = queue[:len(queue)-1]
		met[pop] = struct{}{}

		adj, ok := d.nodeEdges[pop]
		if !ok {
			continue
		}
		for _, ed := range adj {
			ad := ed.ToId
			if ad == targetId {
				return true
			}
			if _, ok := met[ad]; ok {
				continue
			}
			queue = append(queue, ad)
		}
	}
	return false
}

// Connect the two nodes, return true if a new directed edge is created else return false.
//
// When false is returned, these is already a directed edge connecting the two nodes.
func (d *DGraph) Connect(fromId int, toId int) bool {
	edge := DEdge{FromId: fromId, ToId: toId}
	return d.AddEdge(edge)
}

// Connect the two nodes, return true if the new directed edge is added to the graph else return false.
//
// When false is returned, these is already a directed edge connecting the two nodes.
func (d *DGraph) AddEdge(edge DEdge) bool {
	fromId := edge.FromId
	toId := edge.ToId

	_, ok := d.nodeMap[fromId]
	if !ok {
		return false
	}
	eds, ok := d.nodeEdges[fromId]
	if !ok {
		edge := DEdge{FromId: fromId, ToId: toId}
		d.edges = append(d.edges, edge)
		d.nodeEdges[fromId] = []DEdge{edge}
		return true
	}

	for _, ed := range eds {
		if ed.ToId == toId {
			return false
		}
	}

	d.edges = append(d.edges, edge)
	d.nodeEdges[fromId] = append(d.nodeEdges[fromId], edge)
	return true
}

// Add node to graph, return false if the node.Id exist.
func (d *DGraph) AddNode(n Node) bool {
	_, ok := d.nodeMap[n.Id]
	if ok {
		return false
	}
	d.nodeMap[n.Id] = n
	d.nodes = append(d.nodes, n)
	d.nodeEdges[n.Id] = []DEdge{}
	return true
}

func (d *DGraph) SDraw() (string, error) {
	buf := bytes.Buffer{}
	err := d.Draw(&buf)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (d *DGraph) Draw(w io.Writer) error {
	if err := d.writeGraphAttr(w); err != nil {
		return fmt.Errorf("failed to write graph attributes, %w", err)
	}

	buf := bytes.Buffer{}
	for _, n := range d.nodes {
		label := n.Label
		if d.DisplayId {
			label = fmt.Sprintf("%d. %s", n.Id, n.Label)
		}
		shape := n.Shape
		if shape == "" {
			shape = ShapeBox
		}
		buf.WriteString(fmt.Sprintf("N%v [label=\"%v\" id=\"node%v\" fontsize=8 shape=%s tooltip=\"%v\" color=\"#b20400\" fillcolor=\"#edd6d5\"]\n",
			n.Id, label, n.Id, shape, n.Tooltip))
	}

	for _, ed := range d.edges {
		buf.WriteString(fmt.Sprintf("N%v -> N%v [label=\" %s\" labelfloat=false fontsize=6 weight=1 color=\"#b2a999\" tooltip=\"%s\"]\n",
			ed.FromId, ed.ToId, ed.Label, ed.Tooltip))
	}
	buf.WriteString("}\n")
	_, err := w.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("failed to write graph file, %v", err)
	}
	return nil
}

func (d *DGraph) writeGraphAttr(w io.Writer) error {
	b := bytes.Buffer{}
	b.WriteString(fmt.Sprintf("digraph \"[%v]\" {\n", d.title))
	graphAttr := []string{}

	engine := d.Layout
	if engine == "" {
		engine = "dot"
	}
	graphAttr = append(graphAttr, fmt.Sprintf("layout=%s", engine))
	graphAttr = append(graphAttr, fmt.Sprintf("pad=%s", d.Pad))
	graphAttr = append(graphAttr, fmt.Sprintf("ranksep=%s", d.RankSep))
	graphAttr = append(graphAttr, fmt.Sprintf("nodesep=%s", d.NodeSep))
	graphAttr = append(graphAttr, fmt.Sprintf("ratio=\"%s\"", d.Ratio))
	if d.Dpi != "" {
		graphAttr = append(graphAttr, fmt.Sprintf("dpi=\"%s\"", d.Dpi))
	}
	graphAttr = append(graphAttr, "constraint=false")
	graphAttr = append(graphAttr, "overlap=false")
	graphAttr = append(graphAttr, "fontname=\"Helvetica,Arial,sans-serif\"")
	if len(graphAttr) > 0 {
		b.WriteString(fmt.Sprintf("graph [%s]\n", strings.Join(graphAttr, ", ")))
	}

	b.WriteString("node [fontname=\"Helvetica,Arial,sans-serif\"]\n")
	b.WriteString("edge [fontname=\"Helvetica,Arial,sans-serif\"]\n")
	b.WriteString("node [style=filled fillcolor=\"#f8f8f8\"]\n")
	_, err := w.Write(b.Bytes())
	return err
}

func (d *DGraph) NodeCount() int {
	return len(d.nodes)
}

func (d *DGraph) EdgeCount() int {
	return len(d.edges)
}

func NewDGraph(title string, nodes []Node, edges []DEdge) (*DGraph, error) {
	d := new(DGraph)
	d.title = title
	d.nodes = nodes
	d.edges = edges
	d.NodeSep = defaultNodeSep
	d.RankSep = defaultRankSep
	d.Ratio = defaultRatio
	d.Pad = defaultPad
	d.DisplayId = true
	d.Layout = LayoutDot
	if err := d.build(); err != nil {
		return nil, err
	}
	return d, nil
}

type DotGenParam struct {
	GeneratedFile    string // generated graph file location.
	Format           string // default: svg, e.g., svg, png.
	DisableAutoScale bool
}

// Use graphviz dot engine to generate graph file (e.g., svg, png).
//
// e.g., almost the same as the following:
//
//	... | dot -Tsvg > graph.svg && open graph.svg
func DotGen(g *DGraph, p DotGenParam) (DotGenParam, error) {
	if p.Format == "" {
		p.Format = "svg"
	}
	if p.GeneratedFile == "" {
		dir := "/tmp"
		tmpFile, err := os.CreateTemp(dir, "grapher-*."+p.Format)
		if err != nil {
			panic(err)
		}
		p.GeneratedFile = tmpFile.Name()
		tmpFile.Close()
	}

	if !p.DisableAutoScale && p.Format != "svg" && g.Dpi == "" {
		if g.NodeCount() > 50 {
			exp := int(g.NodeCount() / 50)
			g.Dpi = cast.ToString(exp * 300)
		} else {
			g.Dpi = "300"
		}
	}

	s, err := g.SDraw()
	if err != nil {
		return p, err
	}

	cmd := exec.Command("bash", "-c", fmt.Sprintf("dot -T%s > \"%s\"", p.Format, p.GeneratedFile))
	cmd.Stdin = bytes.NewReader([]byte(s))

	cmdout, err := cmd.CombinedOutput()
	if err != nil {
		return p, fmt.Errorf("dot failed, %v, %v", string(cmdout), err)
	}

	return p, nil
}

func DotGenOpen(g *DGraph, p DotGenParam) (DotGenParam, error) {
	pp, err := DotGen(g, p)
	if err != nil {
		return pp, err
	}
	return pp, sys.TermOpenUrl(pp.GeneratedFile)
}
