package syntax

import (
	"io"
	"sort"
)

// Chart is a chart for parsing.
type Chart struct {
	edges map[int][]*Edge
}

// NewChart creates a new chart.
func NewChart() *Chart {
	return &Chart{edges: make(map[int][]*Edge)}
}

// AddEdge adds a new edge to the chart.
func (ch *Chart) AddEdge(e *Edge) {
	es := ch.edges[e.Start]
	ch.edges[e.Start] = append(es, e)
}

// Print prints the chart.
func (ch *Chart) Print(w io.Writer, onlyUnused bool) {
	var edges []*Edge
	for _, es := range ch.edges {
		for _, e := range es {
			if !onlyUnused || !e.Used {
				edges = append(edges, e)
			}
		}
	}
	sort.Slice(edges, func(i, j int) bool {
		e1, e2 := edges[i], edges[j]
		if e1.Start < e2.Start {
			return true
		}
		if e1.Start > e2.Start {
			return false
		}
		if e1.End > e2.End {
			return true
		}
		if e1.End < e2.End {
			return false
		}
		return e1.Form < e2.Form
	})
	for _, e := range edges {
		e.Print(w)
	}
}

// Parse parses the chart using the provided apply function.
func (ch *Chart) Parse(apply func([]*Edge) (string, *AVM)) {
	ch.parse(apply, 0)
}

func (ch *Chart) parse(apply func([]*Edge) (string, *AVM), level int) {
	var newEdges []*Edge
	for _, es := range ch.edges {
		for _, e1 := range es {
			if e1.Level == level {
				if cat, avm := apply([]*Edge{e1}); cat != "" {
					e1.Used = true
					ne := &Edge{
						Start:    e1.Start,
						End:      e1.End,
						Form:     e1.Form,
						Category: cat,
						AVM:      avm,
						Level:    level + 1,
						Children: []*Edge{e1},
					}
					newEdges = append(newEdges, ne)
				}
			}
			for _, e2 := range ch.edges[e1.End] {
				if e1.Level == level || e2.Level == level {
					if cat, avm := apply([]*Edge{e1, e2}); cat != "" {
						e1.Used = true
						e2.Used = true
						ne := &Edge{
							Start:    e1.Start,
							End:      e2.End,
							Form:     e1.Form + " " + e2.Form,
							Category: cat,
							AVM:      avm,
							Level:    level + 1,
							Children: []*Edge{e1, e2},
						}
						newEdges = append(newEdges, ne)
					}
				}
			}
		}
	}
	for _, e := range newEdges {
		ch.AddEdge(e)
	}
	if len(newEdges) >= 1 {
		ch.parse(apply, level+1)
	}
}
