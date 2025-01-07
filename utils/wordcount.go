package utils

import (
	"sort"
	"strings"

	"github.com/phomola/textkit"
)

// WordCount ...
type WordCount struct {
	Word  string
	Count int
}

// CountWords ...
func CountWords(text string) ([]WordCount, error) {
	var tok textkit.Tokeniser
	tokens := tok.Tokenise(text, "")
	m := make(map[string]int)
	for _, t := range tokens {
		if t.Type == textkit.Word {
			f := strings.ToLower(string(t.Form))
			c := m[f]
			m[f] = c + 1
		}
	}
	counts := make([]WordCount, 0, len(m))
	for w, c := range m {
		counts = append(counts, WordCount{
			Word:  w,
			Count: c,
		})
	}
	sort.Slice(counts, func(i, j int) bool {
		return counts[i].Count > counts[j].Count
	})
	return counts, nil
}
