package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/phomola/szl/utils"
)

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

	counts, err := utils.CountWords(string(text))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for _, wc := range counts {
		fmt.Printf("%s/%d\n", wc.Word, wc.Count)
	}
}
