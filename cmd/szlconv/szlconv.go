package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/phomola/szl/textconv"
)

func main() {
	var orthoParam string
	flag.StringVar(&orthoParam, "ortho", "wieczorek", "target Silesian orthography")
	flag.Parse()

	ortho, ok := textconv.OrthographyFromString(orthoParam)
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown orthography '%s'\n", orthoParam)
		os.Exit(1)
	}

	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "exactly one argument required")
		os.Exit(1)
	}

	text := textconv.Convert(flag.Arg(0), ortho)
	fmt.Println(text)

}
