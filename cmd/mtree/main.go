package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
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
	var dat []byte = nil
	if *FlagFile == "" {
		pipe, err := io.ReadAll(os.Stdin)
		if err != nil && !errors.Is(err, io.EOF) {
			panic(err)
		}
		dat = pipe
	}
	if *FlagFile == "" && len(dat) < 1 {
		fmt.Println("Please specify output file or pipe data into mtree")
		return
	}

	if *FlagFile != "" {
		ctn, err := os.ReadFile(*FlagFile)
		if err != nil {
			panic(err)
		}
		dat = ctn
	}

	g, err := mvn.ParseMvnTree(fmt.Sprintf("dependency graph %s", *FlagFile), string(dat))
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
