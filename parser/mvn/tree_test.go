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
	err = graph.DotGen(g)
	if err != nil {
		t.Fatal(err)
	}
}
