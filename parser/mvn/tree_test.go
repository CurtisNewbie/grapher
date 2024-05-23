package mvn

import (
	"os"
	"testing"

	"github.com/curtisnewbie/grapher/graph"
	"github.com/curtisnewbie/grapher/sys"
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
	p, err := graph.DotGen(g, graph.DotGenParam{})
	if err != nil {
		t.Fatal(err)
	}
	if err := sys.TermOpenUrl(p.GeneratedFile); err != nil {
		t.Fatal(err)
	}
}
