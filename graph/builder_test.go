package graph

import (
	"testing"

	"github.com/curtisnewbie/grapher/sys"
)

func TestKNodeGraphBuilder(t *testing.T) {
	bu := NewKNodeGraphBuilder()
	bu.Add("fstore", Node{Label: "fstore"})
	bu.Add("vfm", Node{Label: "vfm"})
	bu.Add("uvault", Node{Label: "user-vault"})

	bu.Connect("fstore", "vfm")
	bu.Connect("uvault", "vfm")
	bu.Connect("fstore", "uvault")
	bu.Connect("uvault", "wtf")

	g, err := bu.BuildDGraph("test")
	if err != nil {
		t.Fatal(err)
	}

	p, err := DotGen(g, DotGenParam{Format: "png"})
	if err != nil {
		t.Fatal(err)
	}
	if err := sys.TermOpenUrl(p.GeneratedFile); err != nil {
		t.Fatal(err)
	}
}
