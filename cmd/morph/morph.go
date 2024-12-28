package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fealsamh/go-utils/replacer"
	"github.com/phomola/textkit"
)

// Stem ...
type Stem struct {
	Form     string
	Paradigm string
	Lemma    string
}

// Ending ...
type Ending struct {
	Form         string
	Paradigm     string
	Tag          string
	Replacements map[string]string
}

// Replacement ...
type Replacement struct {
	Old string
	New string
}

// Entry ...
type Entry struct {
	Lemma            string
	Tag              string
	Category         string
	FeatureStructure string
	Autosemantic     bool
}

// Analyser ...
type Analyser struct {
	Stems             map[string][]Stem
	Endings           map[string][]Ending
	Replacer          *strings.Replacer
	FeatureStructures map[string][]string
	Entries           map[string][]*Entry
}

// Analyse ...
func (an *Analyser) Analyse(form string) ([]*Entry, error) {
	if an.Entries == nil {
		if err := an.buildEntries(); err != nil {
			return nil, err
		}
	}
	if entries, ok := an.Entries[form]; ok {
		return entries, nil
	}
	return nil, nil
}

func (an *Analyser) buildEntries() error {
	an.Entries = make(map[string][]*Entry)
	for _, stems := range an.Stems {
		for _, stem := range stems {
			ends, ok := an.Endings[stem.Paradigm]
			if !ok {
				return fmt.Errorf("unknown paradigm %s for %s", stem.Paradigm, stem.Lemma)
			}
			for _, end := range ends {
				form := replacer.Replace(stem.Form, end.Replacements) + end.Form
				form = an.Replacer.Replace(form)
				entry := &Entry{
					Lemma: stem.Lemma,
					Tag:   end.Tag,
				}
				if fs, ok := an.FeatureStructures[end.Tag]; ok {
					entry.Category = fs[0]
					entry.FeatureStructure = fs[1]
					entry.Autosemantic = fs[2] == "autosem"
				}
				entries := an.Entries[form]
				an.Entries[form] = append(entries, entry)
			}
		}
	}
	return nil
}

func main() {
	var (
		list       bool
		chartInput string
	)
	flag.BoolVar(&list, "list", false, "list all forms (takes file name(s) of lexicon files)")
	flag.StringVar(&chartInput, "chart", "", "chart for the input phrase (takes file name(s) of lexicon files)")
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
		if err := listForms(an); err != nil {
			fmt.Fprintln(os.Stderr, "cannot list forms:", err)
			os.Exit(1)
		}
	} else if chartInput != "" {
		an, err := loadLex(flag.Args())
		if err != nil {
			fmt.Fprintln(os.Stderr, "cannot load lexicon:", err)
			os.Exit(1)
		}
		var tokeniser textkit.Tokeniser
		tokens := tokeniser.Tokenise(chartInput, "")
		for i, token := range tokens {
			if token.Type == textkit.EOF {
				continue
			}
			form := string(token.Form)
			entries, err := an.Analyse(form)
			if err != nil {
				fmt.Fprintln(os.Stderr, "cannot analyse entry:", err)
				os.Exit(1)
			}
			start, end := i+1, i+2
			if len(entries) >= 1 {
				for _, entry := range entries {
					if entry.Category == "" || entry.FeatureStructure == "" {
						fmt.Fprintf(os.Stderr, "no category and/or feature structure for '%s' (tag '%s')\n", form, entry.Tag)
						os.Exit(1)
					}
					fs := strings.ReplaceAll(entry.FeatureStructure, ",", " ")
					if entry.Autosemantic {
						fs = fs[:len(fs)-1] + fmt.Sprintf(" lemma:%q index:\"%d\"]", entry.Lemma, start)
					}
					fmt.Printf("-%d- %s %s -%d-\n", start, entry.Category, fs, end)
				}
			} else {
				fs := fmt.Sprintf("[form:%q index:\"%d\"]", form, start)
				fmt.Printf("-%d- U %s -%d-\n", start, fs, end)
			}
		}
	} else {
		flag.PrintDefaults()
	}
}

func listForms(an *Analyser) error {
	for _, stems := range an.Stems {
		for _, stem := range stems {
			ends, ok := an.Endings[stem.Paradigm]
			if !ok {
				return fmt.Errorf("unknown paradigm %s for %s", stem.Paradigm, stem.Lemma)
			}
			for _, end := range ends {
				form := replacer.Replace(stem.Form, end.Replacements) + end.Form
				form = an.Replacer.Replace(form)
				fmt.Println(form, stem.Lemma, end.Tag)
			}
		}
	}
	return nil
}

func loadLex(files []string) (*Analyser, error) {
	var (
		stems             = make(map[string][]Stem)
		endings           = make(map[string][]Ending)
		replacements      []Replacement
		featureStructures = make(map[string][]string)
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
				if len(comps) < 3 {
					return nil, fmt.Errorf("bad definition at line %d", l)
				}
				par, tag, form := comps[0], comps[1], comps[2]
				if form == "0" {
					form = ""
				}
				var repl map[string]string
				for _, dir := range comps[3:] {
					if len(dir) > 0 {
						if dir[0] == '>' {
							if old, new, ok := strings.Cut(dir[1:], ","); ok {
								if repl == nil {
									repl = make(map[string]string)
								}
								repl[old] = new
							} else {
								return nil, fmt.Errorf("bad definition at line %d", l)
							}
						} else {
							return nil, fmt.Errorf("bad definition at line %d", l)
						}
					} else {
						return nil, fmt.Errorf("bad definition at line %d", l)
					}
				}
				l := endings[par]
				endings[par] = append(l, Ending{Form: form, Paradigm: par, Tag: tag, Replacements: repl})
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
			case '*':
				comps := strings.Split(line[1:], " ")
				if len(comps) != 4 {
					return nil, fmt.Errorf("bad definition at line %d", l)
				}
				if _, ok := featureStructures[comps[0]]; ok {
					return nil, fmt.Errorf("feature structure for '%s' already defined", comps[0])
				}
				featureStructures[comps[0]] = comps[1:]
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
		Stems:             stems,
		Endings:           endings,
		Replacer:          replacer,
		FeatureStructures: featureStructures,
	}, nil
}
