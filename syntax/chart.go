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

// GetEdges returns the edges from start to end.
func (ch *Chart) GetEdges(start, end int, onlyUnused bool) []*Edge {
	var edges []*Edge
	for _, e := range ch.edges[start] {
		if e.End == end {
			if !onlyUnused || !e.Used {
				edges = append(edges, e)
			}
		}
	}
	return edges
}

// GetPaths returns the paths from start to end.
func (ch *Chart) GetPaths(start, end int, onlyUnused bool) [][]*Edge {
	paths := ch.findPaths(start, end, onlyUnused, make(map[int][][]*Edge))
	sort.Slice(paths, func(i, j int) bool {
		p1, p2 := paths[i], paths[j]
		l1, l2 := len(p1), len(p2)
		return l1 < l2
	})
	return paths
}

func (ch *Chart) findPaths(start, end int, onlyUnused bool, pathsFrom map[int][][]*Edge) [][]*Edge {
	if paths, ok := pathsFrom[start]; ok {
		return paths
	}
	var paths [][]*Edge
	for _, e := range ch.edges[start] {
		if !onlyUnused || !e.Used {
			if e.End == end {
				paths = append(paths, []*Edge{e})
			} else {
				tails := ch.findPaths(e.End, end, onlyUnused, pathsFrom)
				for _, tail := range tails {
					paths = append(paths, append([]*Edge{e}, tail...))
				}
			}
		}
	}
	pathsFrom[start] = paths
	return paths
}

// GetClusters returns clusters of edges.
func (ch *Chart) GetClusters(start int, onlyUnused bool) map[int][]*Edge {
	m := make(map[int][]*Edge)
	for _, e := range ch.edges[start] {
		if !onlyUnused || !e.Used {
			es := m[e.End]
			m[e.End] = append(es, e)
		}
	}
	return m
}

// GetPathsOfClusters returns the paths from start to end.
func (ch *Chart) GetPathsOfClusters(start, end int, onlyUnused bool) [][][]*Edge {
	paths := ch.findPathsOfClusters(start, end, onlyUnused, make(map[int][][][]*Edge))
	sort.Slice(paths, func(i, j int) bool {
		p1, p2 := paths[i], paths[j]
		l1, l2 := len(p1), len(p2)
		return l1 < l2
	})
	return paths
}

func (ch *Chart) findPathsOfClusters(start, end int, onlyUnused bool, pathsFrom map[int][][][]*Edge) [][][]*Edge {
	if paths, ok := pathsFrom[start]; ok {
		return paths
	}
	var paths [][][]*Edge
	for clEnd, cl := range ch.GetClusters(start, onlyUnused) {
		if clEnd == end {
			paths = append(paths, [][]*Edge{cl})
		} else {
			tails := ch.findPathsOfClusters(clEnd, end, onlyUnused, pathsFrom)
			for _, tail := range tails {
				paths = append(paths, append([][]*Edge{cl}, tail...))
			}
		}
	}
	pathsFrom[start] = paths
	return paths
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
