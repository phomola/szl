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

// Stem ...
type Stem struct {
	Form     string
	Paradigm string
	Lemma    string
}

// Ending ...
type Ending struct {
	Form     string
	Paradigm string
	Tag      string
}

// Replacement ...
type Replacement struct {
	Old string
	New string
}

// Analyser ...
type Analyser struct {
	Stems    map[string][]Stem
	Endings  map[string][]Ending
	Replacer *strings.Replacer
}

func main() {
	var (
		list bool
	)
	flag.BoolVar(&list, "list", false, "list all forms (takes file name(s) of lexicon files)")
	flag.Parse()
	if list {
		if flag.NArg() == 0 {
			fmt.Fprintln(os.Stderr, "no input files")
			os.Exit(1)
		}
		an, err := loadLex(flag.Args())
		if err != nil {
			fmt.Fprintln(os.Stderr, "cannot load lexicon:", err)
			os.Exit(1)
		}
		listForms(an)
	} else {
		flag.PrintDefaults()
	}
}

func listForms(an *Analyser) {
	for _, stems := range an.Stems {
		for _, stem := range stems {
			ends, ok := an.Endings[stem.Paradigm]
			if !ok {
				fmt.Fprintln(os.Stderr, "unknown paradigm", stem.Paradigm, "for", stem.Lemma)
				os.Exit(1)
			}
			for _, end := range ends {
				form := stem.Form + end.Form
				form = an.Replacer.Replace(form)
				fmt.Println(form, stem.Lemma, end.Tag)
			}
		}
	}
}

func loadLex(files []string) (*Analyser, error) {
	var (
		stems        = make(map[string][]Stem)
		endings      = make(map[string][]Ending)
		replacements []Replacement
	)
	for _, fn := range files {
		f, err := os.Open(fn)
		if err != nil {
			return nil, fmt.Errorf("cannot open file: %w", err)
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
				return nil, fmt.Errorf("cannot read line: %w", err)
			}
			line = strings.TrimSpace(line)
			if line == "" || line[0] == '#' {
				continue
			}
			switch line[0] {
			case '@':
				comps := strings.Split(line[1:], " ")
				if len(comps) != 3 {
					return nil, fmt.Errorf("bad definition at line %d", l)
				}
				lemma, par, form := comps[0], comps[1], comps[2]
				l := stems[lemma]
				stems[lemma] = append(l, Stem{Lemma: lemma, Paradigm: par, Form: form})
			case '-':
				comps := strings.Split(line[1:], " ")
				if len(comps) != 3 {
					return nil, fmt.Errorf("bad definition at line %d", l)
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
						return nil, fmt.Errorf("bad definition at line %d", l)
					}
					replacements = append(replacements, Replacement{Old: comps[1], New: comps[2]})
				default:
					return nil, fmt.Errorf("bad definition at line %d", l)
				}
			default:
				return nil, fmt.Errorf("bad directive at line %d", l)
			}
		}
	}
	var pairs []string
	for _, r := range replacements {
		pairs = append(pairs, r.Old, r.New)
	}
	replacer := strings.NewReplacer(pairs...)
	return &Analyser{
		Stems:    stems,
		Endings:  endings,
		Replacer: replacer,
	}, nil
}
