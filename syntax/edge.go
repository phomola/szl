package syntax

import (
	"fmt"
	"io"
)

// Edge is a chart edge.
type Edge struct {
	Start    int
	End      int
	Category string
	Form     string
	AVM      *AVM
	Info     map[string]interface{}
	Children []*Edge
	Level    int
	Used     bool
	Weight   int
}

// Print prints the edge.
func (e *Edge) Print(w io.Writer) {
	fmt.Fprintf(w, "-%d- %s %s %s -%d-\n", e.Start, e.Form, e.Category, e.AVM.String(), e.End)
	// fmt.Fprintf(w, "-%d- %s %s %s -%d- (%d)\n", e.Start, e.Form, e.Category, e.AVM.String(), e.End, e.Weight)
}

// Linearise linearises the syntax tree.
func (e *Edge) Linearise(f func(*Edge) string) []string {
	if e.Children == nil {
		return []string{f(e)}
	}
	var l []string
	for _, e := range e.Children {
		l = append(l, e.Linearise(f)...)
	}
	return l
}
