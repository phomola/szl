package syntax

import (
	"fmt"
	"maps"
	"sort"
	"strings"
)

// AVMValue is an attribute's value.
type AVMValue interface {
	fmt.Stringer
}

// String is a string.
type String string

// AVMList is a list of AVMs.
type AVMList struct {
	Els []*AVM
}

func (l *AVMList) String() string {
	var sb strings.Builder
	sb.WriteByte('{')
	for i, el := range l.Els {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(el.String())
	}
	sb.WriteByte('}')
	return sb.String()
}

// Append appends an AVM to the list.
func (l *AVMList) Append(avm *AVM) *AVMList {
	return &AVMList{Els: append(l.Els, avm)}
}

func (s String) String() string {
	return string(s)
}

// AVM is an attribute-value matrix.
type AVM struct {
	Features map[string]AVMValue
}

// Set sets an attribute's value.
func (avm *AVM) Set(key string, value AVMValue) {
	avm.Features[key] = value
}

// AddToAVMList adds a value to the list of AVMs.
func (avm *AVM) AddToAVMList(attr string, value *AVM) bool {
	if v, ok := avm.Features[attr]; ok {
		if l, ok := v.(*AVMList); ok {
			avm.Features[attr] = l.Append(value)
			return true
		}
		return false
	}
	avm.Features[attr] = &AVMList{Els: []*AVM{value}}
	return true
}

// GetString gets the value of a string-valued attribute.
func (avm *AVM) GetString(attr string) (string, bool) {
	if v, ok := avm.Features[attr]; ok {
		if x, ok := v.(String); ok {
			return string(x), true
		}
	}
	return "", false
}

// GetList gets the value of a list-valued attribute.
func (avm *AVM) GetList(attr string) (*AVMList, bool) {
	if v, ok := avm.Features[attr]; ok {
		if x, ok := v.(*AVMList); ok {
			return x, true
		}
	}
	return nil, false
}

// Clone clones the AVM.
func (avm *AVM) Clone() *AVM {
	return &AVM{Features: maps.Clone(avm.Features)}
}

func (avm *AVM) String() string {
	pairs := make([][]string, 0, len(avm.Features))
	for k, v := range avm.Features {
		pairs = append(pairs, []string{k, v.String()})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i][0] < pairs[j][0]
	})
	var sb strings.Builder
	sb.WriteString("[")
	for _, p := range pairs {
		sb.WriteByte(' ')
		sb.WriteString(p[0])
		sb.WriteByte(':')
		sb.WriteString(p[1])
	}
	sb.WriteString(" ]")
	return sb.String()
}

// AVMsAttrEqString returns true if the string-valued attribute's values in the AVMs are equal.
func AVMsAttrEqString(avm1, avm2 *AVM, attr string) bool {
	if v1, ok := avm1.GetString(attr); ok {
		if v2, ok := avm2.GetString(attr); ok {
			return v1 == v2
		}
	}
	return false
}
