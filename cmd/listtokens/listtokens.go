package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/phomola/textkit"
)

type wordCount struct {
	word  string
	count int
}

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "no input file provided")
		os.Exit(1)
	}

	f, err := os.Open(flag.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer f.Close()

	text, err := io.ReadAll(f)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var tok textkit.Tokeniser
	tokens := tok.Tokenise(string(text), "")
	m := make(map[string]int)
	for _, t := range tokens {
		if t.Type == textkit.Word {
			f := string(t.Form)
			c := m[f]
			m[f] = c + 1
		}
	}
	counts := make([]wordCount, 0, len(m))
	for w, c := range m {
		counts = append(counts, wordCount{
			word:  w,
			count: c,
		})
	}
	sort.Slice(counts, func(i, j int) bool {
		return counts[i].count > counts[j].count
	})
	for _, wc := range counts {
		fmt.Printf("%s/%d\n", wc.word, wc.count)
	}
}
