package syntax

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
