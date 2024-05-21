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
	FlagFile    = flag.String("file", "", "mvn dependency:tree output file")
	FlagFilter  = flag.String("filter", "", "filter tree branches by label name for tree-shaking")
	FlagCleanup = flag.Bool("cleanup", false, "whether should mtree cleanup generated output file when graph is finally generated")
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

	dir := "/tmp"
	tmpFile, err := os.CreateTemp(dir, "mtree-graph-*.svg")
	if err != nil {
		panic(err)
	}
	svgPath := tmpFile.Name()
	tmpFile.Close()

	p, err := graph.DotGen(g, graph.DotGenParam{OpenSvg: true, GraphSvgFile: svgPath})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Graph SVG generated at: %s\n", svgPath)

	if *FlagCleanup {
		os.Remove(p.GraphOutputFile)
	}
}
