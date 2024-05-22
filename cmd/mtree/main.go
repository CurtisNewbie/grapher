package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/curtisnewbie/grapher/graph"
	"github.com/curtisnewbie/grapher/parser/mvn"
)

var (
	FlagPom    = flag.String("pom", "", "maven pom file")
	FlagFile   = flag.String("file", "", "mvn dependency:tree output file")
	FlagFilter = flag.String("filter", "", "filter tree branches by label name for tree-shaking")
)

func main() {
	flag.Parse()

	var dat []byte = nil

	// mvn dependency:tree output file
	if *FlagFile != "" {
		ctn, err := os.ReadFile(*FlagFile)
		if err != nil {
			panic(err)
		}
		dat = ctn
	}

	// stdin
	if len(dat) < 1 {
		fi, err := os.Stdin.Stat()
		if err != nil {
			panic(err)
		}
		if fi.Size() > 0 {
			pipe, err := io.ReadAll(os.Stdin)
			if err != nil && !errors.Is(err, io.EOF) {
				panic(err)
			}
			dat = pipe
		}
	}

	// pom
	if len(dat) < 1 && *FlagPom != "" {
		cmd := exec.Command("mvn", "dependency:tree", "-f", *FlagPom)
		cmdout, err := cmd.CombinedOutput()
		if err != nil {
			panic(err)
		}
		dat = cmdout
	}

	if len(dat) < 1 {
		fmt.Println("Has nothing to process")
		flag.PrintDefaults()
		return
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

	p, err := graph.DotGen(g, graph.DotGenParam{GeneratedFile: svgPath})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Graph file generated at: %s\n", p.GeneratedFile)

	if err := graph.TermOpenUrl(p.GeneratedFile); err != nil {
		panic(err)
	}
}
