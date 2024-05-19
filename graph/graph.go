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
	Title string
	Nodes []Node
	Edges []DEdge

	NodeEdges map[int][]int
}

func (d *DGraph) build() {
	d.NodeEdges = map[int][]int{}
	for _, n := range d.Edges {
		tids, ok := d.NodeEdges[n.FromId]
		if ok {
			i, found := slices.BinarySearch(tids, n.ToId)
			if !found {
				tids = append(tids, n.ToId)
				copy(tids[i+1:], tids[i:])
				tids[i] = n.ToId
				d.NodeEdges[n.FromId] = tids
			}
		} else {
			d.NodeEdges[n.FromId] = []int{n.ToId}
		}
	}
}

func (d *DGraph) Connected(rootId int, targetId int) bool {
	var rootFound bool = false
	var root Node
	for _, n := range d.Nodes {
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

		adj, ok := d.NodeEdges[pop]
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
	for _, n := range d.Nodes {
		buf.WriteString(fmt.Sprintf("N%v [label=\"%v\" id=\"node%v\" fontsize=8 shape=box tooltip=\"%v\" color=\"#b20400\" fillcolor=\"#edd6d5\"]\n",
			n.Id, n.Label, n.Id, n.Tooltip))
	}

	for _, ed := range d.Edges {
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
	b.WriteString(fmt.Sprintf("digraph \"[%v]\" {\n", d.Title))
	b.WriteString("pad=0.5\n")
	b.WriteString("fontname=\"Helvetica,Arial,sans-serif\"\n")
	b.WriteString("node [fontname=\"Helvetica,Arial,sans-serif\"]\n")
	b.WriteString("edge [fontname=\"Helvetica,Arial,sans-serif\"]\n")
	b.WriteString("node [style=filled fillcolor=\"#f8f8f8\"]\n")
	_, err := w.Write(b.Bytes())
	return err
}

func NewDGraph(title string, nodes []Node, edges []DEdge) *DGraph {
	d := new(DGraph)
	d.Title = title
	d.Nodes = nodes
	d.Edges = edges
	d.build()
	return d
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
