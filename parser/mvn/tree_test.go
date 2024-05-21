package mvn

import (
	"os"
	"testing"

	"github.com/curtisnewbie/grapher/graph"
)

func TestParseMvnTree(t *testing.T) {
	ctn, err := os.ReadFile("../../testdata/todoapp_mvn.out")
	if err != nil {
		t.Fatal(err)
	}
	g, err := ParseMvnTree("dependency tree", string(ctn))
	if err != nil {
		t.Fatal(err)
	}
	// g.FilterBranch(func(n graph.Node) bool { return strings.Contains(n.Label, "com.fasterxml") })
	err = graph.DotGen(g, graph.DotGenParam{OpenViewer: true})
	if err != nil {
		t.Fatal(err)
	}
}