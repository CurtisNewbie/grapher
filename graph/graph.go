package graph

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"strings"
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
}

type DGraph struct {
	title string
	nodes []Node
	edges []DEdge

	nodeIdx    map[int]int     // id -> nodes idx
	neighbours map[int][]int   // fromId -> toId
	nodeEdges  map[int][]DEdge // id -> edges

	DisplayId bool
	RankSep   string
	NodeSep   string
	Ratio     string
	Pad       string
}

func (d *DGraph) build() error {
	d.neighbours = map[int][]int{}
	d.nodeEdges = map[int][]DEdge{}
	for _, n := range d.edges {
		tids, ok := d.neighbours[n.FromId]
		if ok {
			i, found := slices.BinarySearch(tids, n.ToId)
			if found {
				return fmt.Errorf("found duplicate edges on id: %v to id: %v", n.FromId, n.ToId)
			}
			tids = append(tids, n.ToId)
			copy(tids[i+1:], tids[i:])
			tids[i] = n.ToId
			d.neighbours[n.FromId] = tids
		} else {
			d.neighbours[n.FromId] = []int{n.ToId}
		}

		if ae, ok := d.nodeEdges[n.FromId]; ok {
			d.nodeEdges[n.FromId] = append(ae, n)
		} else {
			d.nodeEdges[n.FromId] = []DEdge{n}
		}
	}
	d.nodeIdx = map[int]int{}
	for i, n := range d.nodes {
		if _, ok := d.nodeIdx[n.Id]; ok {
			return fmt.Errorf("Node id duplicate found, id: %v", n.Id)
		}
		d.nodeIdx[n.Id] = i
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
	idx, ok := d.nodeIdx[id]
	if !ok {
		return Node{}, false
	}
	return d.nodes[idx], true
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

		for _, c := range d.neighbours[p.Id] {
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
	var rootFound bool = false
	var root Node
	for _, n := range d.nodes {
		if n.Id == rootId {
			root = n
			rootFound = true
		}
	}
	if !rootFound {
		return false
	}

	queue := []int{root.Id}
	met := map[int]struct{}{root.Id: {}}
	for len(queue) > 0 {
		pop := queue[len(queue)-1]
		queue = queue[:len(queue)-1]
		met[pop] = struct{}{}

		adj, ok := d.neighbours[pop]
		if !ok {
			continue
		}
		for _, ad := range adj {
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
		buf.WriteString(fmt.Sprintf("N%v [label=\"%v\" id=\"node%v\" fontsize=8 shape=box tooltip=\"%v\" color=\"#b20400\" fillcolor=\"#edd6d5\"]\n",
			n.Id, label, n.Id, n.Tooltip))
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
	b.WriteString(fmt.Sprintf("pad=%s\n", d.Pad))
	b.WriteString(fmt.Sprintf("ranksep=%s\n", d.RankSep))
	b.WriteString(fmt.Sprintf("nodesep=%s\n", d.NodeSep))
	b.WriteString(fmt.Sprintf("ratio=\"%s\"\n", d.Ratio))
	b.WriteString("constraint = false\n")
	b.WriteString("overlap=false\n")
	b.WriteString("fontname=\"Helvetica,Arial,sans-serif\"\n")
	b.WriteString("node [fontname=\"Helvetica,Arial,sans-serif\"]\n")
	b.WriteString("edge [fontname=\"Helvetica,Arial,sans-serif\"]\n")
	b.WriteString("node [style=filled fillcolor=\"#f8f8f8\"]\n")
	_, err := w.Write(b.Bytes())
	return err
}

func NewDGraph(title string, nodes []Node, edges []DEdge) (*DGraph, error) {
	d := new(DGraph)
	d.title = title
	d.nodes = nodes
	d.edges = edges
	d.NodeSep = "1"
	d.RankSep = "1"
	d.Ratio = "auto"
	d.Pad = "0.5"
	if err := d.build(); err != nil {
		return nil, err
	}
	return d, nil
}

//go:embed graph.html
var graphTemplHtml []byte

const graphOutputName = "graph.txt"
const graphTemplName = "graph.html"
const graphSvgName = "graph.svg"

// Use graphviz dot engine to generate graph svg file and host it in locally generated template.
//
// e.g.,
//
//	dot -Tsvg $path > graph.svg && open graph.html
func DotGen(g *DGraph) error {
	of, err := ReadWriteFile(graphOutputName)
	if err != nil {
		return err
	}
	defer of.Close()
	of.Truncate(0)

	if err := g.Draw(of); err != nil {
		return err
	}

	cmd := exec.Command("bash", "-c", fmt.Sprintf("dot -Tsvg \"%s\" > \"%s\"", graphOutputName, graphSvgName))
	fmt.Printf("%v", cmd)
	cmdout, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("dot failed, %v, %v", string(cmdout), err)
	}

	templ, err := ReadWriteFile(graphTemplName)
	if err != nil {
		return err
	}
	defer templ.Close()
	templ.Truncate(0)
	if _, err := templ.Write(graphTemplHtml); err != nil {
		return fmt.Errorf("failed to write graph.html template, %v", err)
	}
	TermOpenUrl("graph.html")
	return nil
}

func TermOpenUrl(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

// Open file with 0666 permission.
func OpenFile(name string, flag int) (*os.File, error) {
	return os.OpenFile(name, flag, 0666)
}

// Create readable & writable file with 0666 permission.
func ReadWriteFile(name string) (*os.File, error) {
	return OpenFile(name, os.O_CREATE|os.O_RDWR)
}
