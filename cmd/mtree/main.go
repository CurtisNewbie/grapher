package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/curtisnewbie/grapher/graph"
	"github.com/curtisnewbie/grapher/parser/mvn"
)

var (
	FlagFile   = flag.String("file", "", "mvn dependency:tree output file")
	FlagFilter = flag.String("filter", "", "filter tree branches by label name for tree-shaking")
)

func main() {
	flag.Parse()
	if *FlagFile == "" {
		fmt.Println("Please specify output file")
		return
	}

	ctn, err := os.ReadFile(*FlagFile)
	if err != nil {
		panic(err)
	}
	g, err := mvn.ParseMvnTree(fmt.Sprintf("dependency graph %s", *FlagFile), string(ctn))
	if err != nil {
		panic(err)
	}

	if *FlagFilter != "" {
		g.TreeShake(func(n graph.Node) bool { return strings.Contains(n.Label, *FlagFilter) })
	}

	err = graph.DotGen(g, graph.DotGenParam{OpenViewer: true})
	if err != nil {
		panic(err)
	}
}
