package main

import "github.com/phomola/szl/syntax"

// Apply applies a rule to the provided chain of edges.
func Apply(edges []*syntax.Edge) (string, *syntax.AVM) {
	switch len(edges) {
	case 1:
		e1 := edges[0]
		switch e1.Category {
		case "N":
			return "N'", e1.AVM
		case "N'":
			return "NP", e1.AVM
		}
	case 2:
		e1, e2 := edges[0], edges[1]
		if e1.Category == "A" && e2.Category == "N'" {
			if syntax.AVMsAttrEqString(e1.AVM, e2.AVM, "gender") &&
				syntax.AVMsAttrEqString(e1.AVM, e2.AVM, "case") &&
				syntax.AVMsAttrEqString(e1.AVM, e2.AVM, "number") {
				avm := e2.AVM.Clone()
				if avm.AddToAVMList("adj", e1.AVM) {
					return "N'", avm
				}
			}
		}
	}
	return "", nil
}
