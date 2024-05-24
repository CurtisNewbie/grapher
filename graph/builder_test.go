package graph

import (
	"sync"
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

func TestAsyncKNodeGraphBuilder(t *testing.T) {
	bu := NewKNodeGraphBuilder()

	var wg sync.WaitGroup
	async := func(f func()) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			f()
		}()
	}

	async(func() { bu.Add("fstore", Node{Label: "fstore"}) })
	async(func() { bu.Add("vfm", Node{Label: "vfm"}) })
	async(func() { bu.Add("uvault", Node{Label: "user-vault"}) })

	async(func() { bu.Connect("fstore", "vfm") })
	async(func() { bu.Connect("uvault", "vfm") })
	async(func() { bu.Connect("fstore", "uvault") })
	async(func() { bu.Connect("uvault", "wtf") })

	wg.Wait()
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
