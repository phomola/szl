package syntax

import "io"

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
func (ch *Chart) Print(w io.Writer) {
	for _, es := range ch.edges {
		for _, e := range es {
			e.Print(w)
		}
	}
}
