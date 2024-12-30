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
	Children []*Edge
	Level    int
	Used     bool
}

// Print prints the edge.
func (e *Edge) Print(w io.Writer) {
	fmt.Fprintf(w, "-%d- %s %s %s -%d-\n", e.Start, e.Form, e.Category, e.AVM.String(), e.End)
}
