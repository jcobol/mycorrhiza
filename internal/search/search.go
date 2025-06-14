package search

import (
	"log/slog"
	"sort"
	"strings"
	"unicode"

	"github.com/bouncepaw/mycorrhiza/internal/hyphae"
)

var (
	index    = make(map[string]map[string]struct{})
	conveyor = make(chan indexOperation)
)

type indexOperation interface{ apply() }

type indexEdit struct {
	name      string
	oldTokens []string
	newTokens []string
}

type indexDelete struct {
	name   string
	tokens []string
}

type indexRename struct {
	oldName string
	newName string
	tokens  []string
}

func RunConveyor() {
	defer close(conveyor)
	for op := range conveyor {
		op.apply()
	}
}

func Index() {
	for h := range hyphae.FilterHyphaeWithText(hyphae.YieldExistingHyphae()) {
		text, err := hyphae.FetchMycomarkupFile(h)
		if err != nil {
			slog.Error("failed to read hypha text", "err", err, "hypha", h.CanonicalName())
			continue
		}
		addTokens(h.CanonicalName(), tokenise(text))
	}
}

func Search(query string) []string {
	tokens := tokenise(query)
	if len(tokens) == 0 {
		return nil
	}
	sets := make([]map[string]struct{}, 0, len(tokens))
	for _, t := range tokens {
		if set, ok := index[t]; ok {
			sets = append(sets, set)
		} else {
			return nil
		}
	}
	result := make(map[string]struct{})
	for h := range sets[0] {
		result[h] = struct{}{}
	}
	for _, set := range sets[1:] {
		for h := range result {
			if _, ok := set[h]; !ok {
				delete(result, h)
			}
		}
	}
	out := make(chan string)
	sorted := hyphae.PathographicSort(out)
	go func() {
		for h := range result {
			out <- h
		}
		close(out)
	}()
	var res []string
	for h := range sorted {
		res = append(res, h)
	}
	return res
}

func UpdateAfterEdit(h hyphae.Hypha, oldText string) {
	newText, _ := hyphae.FetchMycomarkupFile(h)
	conveyor <- indexEdit{h.CanonicalName(), tokenise(oldText), tokenise(newText)}
}

func UpdateAfterDelete(h hyphae.Hypha, oldText string) {
	conveyor <- indexDelete{h.CanonicalName(), tokenise(oldText)}
}

func UpdateAfterRename(h hyphae.Hypha, oldName string) {
	text, _ := hyphae.FetchMycomarkupFile(h)
	conveyor <- indexRename{oldName, h.CanonicalName(), tokenise(text)}
}

func addTokens(name string, tokens []string) {
	for _, t := range tokens {
		set, ok := index[t]
		if !ok {
			set = make(map[string]struct{})
			index[t] = set
		}
		set[name] = struct{}{}
	}
}

func removeTokens(name string, tokens []string) {
	for _, t := range tokens {
		if set, ok := index[t]; ok {
			delete(set, name)
			if len(set) == 0 {
				delete(index, t)
			}
		}
	}
}

func tokenise(text string) []string {
	text = strings.ToLower(text)
	fields := strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	uniq := make(map[string]struct{})
	var res []string
	for _, f := range fields {
		if f == "" {
			continue
		}
		if _, ok := uniq[f]; !ok {
			uniq[f] = struct{}{}
			res = append(res, f)
		}
	}
	sort.Strings(res)
	return res
}

func (op indexEdit) apply() {
	removeTokens(op.name, op.oldTokens)
	addTokens(op.name, op.newTokens)
}

func (op indexDelete) apply() {
	removeTokens(op.name, op.tokens)
}

func (op indexRename) apply() {
	removeTokens(op.oldName, op.tokens)
	addTokens(op.newName, op.tokens)
}
