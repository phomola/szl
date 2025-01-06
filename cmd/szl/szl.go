package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"maps"
	"os"
	"plugin"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/fealsamh/go-utils/replacer"
	"github.com/phomola/lrparser"
	"github.com/phomola/szl/syntax"
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
	Form             string
	Lemma            string
	Tag              string
	Category         string
	FeatureStructure string
	Autosemantic     bool
	avm              *syntax.AVM
	Info             map[string]interface{}
}

var avmGrammar *lrparser.Grammar

// AVM gets the associated AVM.
func (e *Entry) AVM() (*syntax.AVM, error) {
	if e.avm != nil {
		return e.avm, nil
	}
	if avmGrammar == nil {
		avmGrammar = lrparser.NewGrammar(lrparser.MustBuildRules([]*lrparser.SynSem{
			{Syn: `Init -> AVM`, Sem: func(args []any) any { return args[0] }},
			{Syn: `AVM -> "[" Pairs "]"`, Sem: func(args []any) any {
				return &syntax.AVM{Features: args[1].(map[string]syntax.AVMValue)}
			}},
			{Syn: `Pairs -> Pairs "," Pair`, Sem: func(args []any) any {
				m := args[0].(map[string]syntax.AVMValue)
				p := args[2].([]interface{})
				m[p[0].(string)] = p[1].(syntax.AVMValue)
				return m
			}},
			{Syn: `Pairs -> Pair`, Sem: func(args []any) any {
				p := args[0].([]interface{})
				return map[string]syntax.AVMValue{p[0].(string): p[1].(syntax.AVMValue)}
			}},
			{Syn: `Pair -> ident ":" string`, Sem: func(args []any) any {
				return []interface{}{args[0].(string), syntax.String(args[2].(string))}
			}},
		}))
	}
	tokeniser := textkit.Tokeniser{StringRune: '"'}
	tokens := tokeniser.Tokenise(e.FeatureStructure, "")
	avm, err := avmGrammar.Parse(tokens)
	if err != nil {
		return nil, err
	}
	e.avm = avm.(*syntax.AVM)
	return e.avm, nil
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
	if entries, ok := an.Entries[form]; ok {
		return entries, nil
	}
	return nil, nil
}

func (an *Analyser) buildEntries(substReplacer *strings.Replacer) error {
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
				origForm := form
				if substReplacer != nil {
					form = substReplacer.Replace(form)
				}
				entry := &Entry{
					Form:  form,
					Lemma: stem.Lemma,
					Tag:   end.Tag,
					Info:  map[string]interface{}{"origForm": origForm},
				}
				if fs, ok := an.FeatureStructures[end.Tag]; ok {
					entry.Category = fs[0]
					entry.FeatureStructure = fs[1]
					entry.Autosemantic = fs[2] == "autosem"
				} else {
					return fmt.Errorf("missing feature structure for '%s'", end.Tag)
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
		list         bool
		chartInput   string
		parserPlugin string
		subst        string
	)
	flag.BoolVar(&list, "list", false, "list all forms (takes file name(s) of lexicon files)")
	flag.StringVar(&chartInput, "chart", "", "chart for the input phrase (takes file name(s) of lexicon files)")
	flag.StringVar(&parserPlugin, "parser", "", "parser plugin")
	flag.StringVar(&subst, "subst", "", "letter substitutions")
	flag.Parse()
	var substReplacer *strings.Replacer
	if subst != "" {
		var pairs []string
		for _, rewr := range strings.Split(subst, ",") {
			old, new, ok := strings.Cut(rewr, ">")
			if !ok {
				fmt.Fprintln(os.Stderr, "cannot analyse substitutions:", subst)
				os.Exit(1)
			}
			pairs = append(pairs, strings.ToLower(old), strings.ToLower(new), strings.ToUpper(old), strings.ToUpper(new))
		}
		substReplacer = strings.NewReplacer(pairs...)
	}
	if list {
		if flag.NArg() == 0 {
			fmt.Fprintln(os.Stderr, "no input files")
			os.Exit(1)
		}
		an, err := loadLex(flag.Args(), substReplacer)
		if err != nil {
			fmt.Fprintln(os.Stderr, "cannot load lexicon:", err)
			os.Exit(1)
		}
		if err := listForms(an); err != nil {
			fmt.Fprintln(os.Stderr, "cannot list forms:", err)
			os.Exit(1)
		}
	} else if parserPlugin != "" && chartInput != "" {
		plugin, err := plugin.Open(parserPlugin)
		if err != nil {
			fmt.Fprintln(os.Stderr, "cannot open parser plugin:", err)
			os.Exit(1)
		}
		applySymbol, err := plugin.Lookup("Apply")
		if err != nil {
			fmt.Fprintln(os.Stderr, "invalid parser plugin (missing Parse function):", err)
			os.Exit(1)
		}
		apply, ok := applySymbol.(func([]*syntax.Edge) (string, *syntax.AVM))
		if !ok {
			fmt.Fprintln(os.Stderr, "invalid parser plugin (mistyped Parse function)")
			os.Exit(1)
		}
		an, err := loadLex(flag.Args(), substReplacer)
		if err != nil {
			fmt.Fprintln(os.Stderr, "cannot load lexicon:", err)
			os.Exit(1)
		}
		var tokeniser textkit.Tokeniser
		tokens := tokeniser.Tokenise(chartInput, "")
		endNode := len(tokens)
		chart := syntax.NewChart()
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
					avm, err := entry.AVM()
					if err != nil {
						fmt.Fprintln(os.Stderr, "cannot parse AVM:", err)
						fmt.Fprintln(os.Stderr, entry.Tag, entry.FeatureStructure)
						os.Exit(1)
					}
					if entry.Autosemantic {
						avm.Set("lemma", syntax.String(entry.Lemma))
						avm.Set("index", syntax.String(strconv.Itoa(start)))
					}
					chart.AddEdge(&syntax.Edge{
						Start:    start,
						End:      end,
						Form:     form,
						Category: entry.Category,
						AVM:      avm,
						Info:     entry.Info,
					})
				}
			} else {
				chart.AddEdge(&syntax.Edge{
					Start:    start,
					End:      end,
					Form:     form,
					Category: "U",
					AVM: &syntax.AVM{Features: map[string]syntax.AVMValue{
						"form":  syntax.String(form),
						"index": syntax.String(strconv.Itoa(start)),
					}},
				})
			}
		}
		chart.Print(os.Stdout, false)
		fmt.Println()
		chart.Parse(apply)
		chart.Print(os.Stdout, true)
		fmt.Println()
		var minLen int
		for _, path := range chart.GetPathsOfClusters(1, endNode, true) {
			if minLen == 0 {
				minLen = len(path)
			} else {
				if minLen < len(path) {
					break
				}
			}
			for i, cl := range path {
				if i > 0 {
					fmt.Print(" + ")
				}
				clForms := make(map[string]struct{}, len(cl))
				for _, edge := range cl {
					forms := edge.Linearise(func(e *syntax.Edge) string {
						if of, ok := e.Info["origForm"]; ok {
							return of.(string)
						}
						return e.Form
					})
					clForms[strings.Join(forms, " ")] = struct{}{}
				}
				forms := slices.Collect(maps.Keys(clForms))
				sort.Strings(forms)
				if len(forms) == 1 {
					fmt.Print(forms[0])
				} else {
					fmt.Print("{ ", strings.Join(forms, " | "), " }")
				}
			}
			fmt.Println()
		}
	} else if chartInput != "" {
		an, err := loadLex(flag.Args(), substReplacer)
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

func loadLex(files []string, substReplacer *strings.Replacer) (*Analyser, error) {
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
								if new == "0" {
									new = ""
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
	an := &Analyser{
		Stems:             stems,
		Endings:           endings,
		Replacer:          replacer,
		FeatureStructures: featureStructures,
	}
	if err := an.buildEntries(substReplacer); err != nil {
		return nil, err
	}
	return an, nil
}
