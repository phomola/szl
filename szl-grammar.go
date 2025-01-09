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
		case "V":
			if vform, ok := e1.AVM.GetString("vform"); ok && vform == "ppart" {
				return "A", e1.AVM
			}
		}
	case 2:
		e1, e2 := edges[0], edges[1]
		// N' -> A N' -- adj
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
		// V -> V NP -- obj
		if e1.Category == "V" && e2.Category == "NP" {
			if vform, ok := e1.AVM.GetString("vform"); ok && (vform == "fin" || vform == "past") {
				if cas, ok := e2.AVM.GetString("case"); ok && cas == "acc" {
					if _, ok := e1.AVM.GetAVM("obj"); !ok {
						avm := e1.AVM.Clone()
						avm.Set("obj", e2.AVM)
						return "V", avm
					}
				}
			}
		}
		// V -> NP V -- obj
		if e1.Category == "NP" && e2.Category == "V" {
			e1, e2 := e2, e1
			if vform, ok := e1.AVM.GetString("vform"); ok && (vform == "fin" || vform == "past") {
				if cas, ok := e2.AVM.GetString("case"); ok && cas == "acc" {
					if _, ok := e1.AVM.GetAVM("obj"); !ok {
						avm := e1.AVM.Clone()
						avm.Set("obj", e2.AVM)
						return "V", avm
					}
				}
			}
		}
		// V -> V NP -- iobj
		if e1.Category == "V" && e2.Category == "NP" {
			if vform, ok := e1.AVM.GetString("vform"); ok && (vform == "fin" || vform == "past") {
				if cas, ok := e2.AVM.GetString("case"); ok && cas == "dat" {
					if _, ok := e1.AVM.GetAVM("iobj"); !ok {
						avm := e1.AVM.Clone()
						avm.Set("iobj", e2.AVM)
						return "V", avm
					}
				}
			}
		}
		// V -> NP V -- iobj
		if e1.Category == "NP" && e2.Category == "V" {
			e1, e2 := e2, e1
			if vform, ok := e1.AVM.GetString("vform"); ok && (vform == "fin" || vform == "past") {
				if cas, ok := e2.AVM.GetString("case"); ok && cas == "dat" {
					if _, ok := e1.AVM.GetAVM("iobj"); !ok {
						avm := e1.AVM.Clone()
						avm.Set("iobj", e2.AVM)
						return "V", avm
					}
				}
			}
		}
		// V -> V NP -- subj
		if e1.Category == "V" && e2.Category == "NP" {
			if vform, ok := e1.AVM.GetString("vform"); ok && (vform == "fin" || vform == "past") {
				if cas, ok := e2.AVM.GetString("case"); ok && cas == "nom" {
					if _, ok := e1.AVM.GetAVM("subj"); !ok {
						avm := e1.AVM.Clone()
						avm.Set("subj", e2.AVM)
						return "V", avm
					}
				}
			}
		}
		// V -> NP V -- subj
		if e1.Category == "NP" && e2.Category == "V" {
			e1, e2 := e2, e1
			if vform, ok := e1.AVM.GetString("vform"); ok && (vform == "fin" || vform == "past") {
				if cas, ok := e2.AVM.GetString("case"); ok && cas == "nom" {
					if _, ok := e1.AVM.GetAVM("subj"); !ok {
						avm := e1.AVM.Clone()
						avm.Set("subj", e2.AVM)
						return "V", avm
					}
				}
			}
		}
	}
	return "", nil
}
