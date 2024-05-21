package mvn

import (
	"strings"

	"github.com/curtisnewbie/grapher/graph"
)

func ParseMvnTree(title string, s string) (*graph.DGraph, error) {
	lines := strings.Split(s, "\n")

	seg := [][]string{}
	{
		inSeg := false
		se := []string{}
		for _, l := range lines {
			if !strings.HasPrefix(l, "[INFO]") {
				continue
			}
			if !inSeg && strings.HasPrefix(l, "[INFO] --- dependency:") {
				inSeg = true
				continue
			} else if inSeg && (strings.TrimSpace(l) == "[INFO]" || strings.HasPrefix(l, "[INFO] ----")) {
				inSeg = false
				seg = append(seg, se)
				se = []string{}
				continue
			}
			if inSeg {
				se = append(se, l)
			}
		}
	}

	id := 0
	type Entry struct {
		Id           int
		Name         string
		Layer        int
		Dependencies []string
	}

	nodeMap := map[string]*Entry{}
	newEntry := func(l string, layer int) *Entry {
		v, ok := nodeMap[l]
		if !ok {
			id++
			v = &Entry{Name: l, Dependencies: []string{}, Layer: layer, Id: id}
			nodeMap[l] = v
		} else {
			v.Layer = layer
		}
		return v
	}
	addDep := func(p *Entry, l string) {
		found := false
		for _, d := range p.Dependencies {
			if d == l {
				found = true
			}
		}
		if !found {
			p.Dependencies = append(p.Dependencies, l)
		}
	}

	for _, se := range seg {
		currLayer := 0
		parents := []*Entry{}
		for _, l := range se {
			l = strings.TrimSpace(l)
			if l == "" {
				continue
			}

			l = l[7:] // "[INFO] "
			idt := 0

			for _, r := range l {
				switch r {
				case '+', '-', ' ', '|', '\\':
					idt += 1
				default:
					goto PARSE_INDENT_END
				}
			}
		PARSE_INDENT_END:

			layer := 0
			if idt > 0 {
				layer = int(idt / 3)
			}
			l = l[idt:]

			{
				var ok bool
				l, ok = strings.CutSuffix(l, ":compile")
				if ok {
					goto CUT_SUF_END
				}
				l, ok = strings.CutSuffix(l, ":test")
				if ok {
					goto CUT_SUF_END
				}
				l, ok = strings.CutSuffix(l, ":provided")
				if ok {
					goto CUT_SUF_END
				}
				l, _ = strings.CutSuffix(l, ":runtime")
			}

		CUT_SUF_END:

			if len(parents) < 1 {
				v := newEntry(l, layer)
				parents = append(parents, v)
				currLayer = layer
			} else {
				if layer <= currLayer {
					parents = parents[:len(parents)-1]
					for len(parents) > 1 && layer <= parents[len(parents)-1].Layer {
						parents = parents[:len(parents)-1]
					}

					// log.Debugf("l: %v", l)
					p := parents[len(parents)-1]
					addDep(p, l)

					v := newEntry(l, layer)
					parents = append(parents, v)
					currLayer = layer
				} else if layer > currLayer {
					p := parents[len(parents)-1]
					addDep(p, l)

					v := newEntry(l, layer)
					parents = append(parents, v)
					currLayer = layer

				} else {
					p := parents[len(parents)-1]
					addDep(p, l)
					newEntry(l, layer)
				}
			}
		}
	}

	// for k, v := range nodeMap {
	// 	fmt.Printf("[debug] k: %v, v: %#v\n", k, *v)
	// }

	nodes := make([]graph.Node, 0, len(nodeMap))
	edges := make([]graph.DEdge, 0)
	for k := range nodeMap {
		n := nodeMap[k]
		tkn := strings.Split(n.Name, ":")
		label := tkn[0]
		if len(tkn) > 1 {
			label += "\n" + tkn[1]
		}
		if len(tkn) > 2 {
			label += "\n" + strings.Join(tkn[2:], ":")
		}
		nodes = append(nodes, graph.Node{
			Id:    n.Id,
			Label: label,
		})
		for _, dl := range n.Dependencies {
			d := nodeMap[dl]
			edges = append(edges, graph.DEdge{FromId: n.Id, ToId: d.Id})
		}
	}
	d, err := graph.NewDGraph(title, nodes, edges)
	if err != nil {
		return nil, err
	}
	return d, nil
}
