package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type Stem struct {
	Form     string
	Paradigm string
	Lemma    string
}

type Ending struct {
	Form     string
	Paradigm string
	Tag      string
}

type Replacement struct {
	Old string
	New string
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "no input files")
		os.Exit(1)
	}
	var (
		stems        = make(map[string][]Stem)
		endings      = make(map[string][]Ending)
		replacements []Replacement
	)
	for i := 0; i < flag.NArg(); i++ {
		f, err := os.Open(flag.Arg(i))
		if err != nil {
			fmt.Fprintln(os.Stderr, "cannot open file:", err)
			os.Exit(1)
		}
		defer f.Close()
		r := bufio.NewReader(f)
		var l int
		for {
			line, err := r.ReadString('\n')
			if errors.Is(err, io.EOF) {
				break
			}
			l++
			if err != nil {
				fmt.Fprintln(os.Stderr, "cannot read line:", err)
				os.Exit(1)
			}
			line = strings.TrimSpace(line)
			if line == "" || line[0] == '#' {
				continue
			}
			switch line[0] {
			case '@':
				comps := strings.Split(line[1:], " ")
				if len(comps) != 3 {
					fmt.Fprintln(os.Stderr, "bad definition at line", l)
					os.Exit(1)
				}
				lemma, par, form := comps[0], comps[1], comps[2]
				l := stems[lemma]
				stems[lemma] = append(l, Stem{Lemma: lemma, Paradigm: par, Form: form})
			case '-':
				comps := strings.Split(line[1:], " ")
				if len(comps) != 3 {
					fmt.Fprintln(os.Stderr, "bad definition at line", l)
					os.Exit(1)
				}
				par, tag, form := comps[0], comps[1], comps[2]
				if form == "0" {
					form = ""
				}
				l := endings[par]
				endings[par] = append(l, Ending{Form: form, Paradigm: par, Tag: tag})
			case '!':
				comps := strings.Split(line[1:], " ")
				switch comps[0] {
				case ">":
					if len(comps) != 3 {
						fmt.Fprintln(os.Stderr, "bad definition at line", l)
						os.Exit(1)
					}
					replacements = append(replacements, Replacement{Old: comps[1], New: comps[2]})
				default:
					fmt.Fprintln(os.Stderr, "bad definition at line", l)
					os.Exit(1)
				}
			default:
				fmt.Fprintln(os.Stderr, "bad directive at line", l)
				os.Exit(1)
			}
		}
	}
	var pairs []string
	for _, r := range replacements {
		pairs = append(pairs, r.Old, r.New)
	}
	replacer := strings.NewReplacer(pairs...)
	for _, stems := range stems {
		for _, stem := range stems {
			ends, ok := endings[stem.Paradigm]
			if !ok {
				fmt.Fprintln(os.Stderr, "unknown paradigm", stem.Paradigm, "for", stem.Lemma)
				os.Exit(1)
			}
			for _, end := range ends {
				form := stem.Form + end.Form
				form = replacer.Replace(form)
				fmt.Println(form, stem.Lemma, end.Tag)
			}
		}
	}
}
