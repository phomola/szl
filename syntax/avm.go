package syntax

import (
	"fmt"
	"sort"
	"strings"
)

// AVMValue is an attribute's value.
type AVMValue interface {
	fmt.Stringer
}

// String is a string.
type String string

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
