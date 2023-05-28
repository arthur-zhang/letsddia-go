package inverted_index

import (
	"lo"
	"string_utils"
	"strings"
)

type Appearance struct {
	DocId int
	Freq  int
}
type InvertedIndex struct {
	index map[string][]*Appearance
	docs  map[int]*Doc
}
type Doc struct {
	Id      int
	Content string
}

func New() InvertedIndex {
	return InvertedIndex{
		index: make(map[string][]*Appearance),
		docs:  make(map[int]*Doc, 0),
	}
}
func (t *InvertedIndex) IndexDoc(doc *Doc) {
	content := string_utils.RemovePunctuation(doc.Content)
	terms := strings.Split(content, " ")
	for _, term := range terms {
		it, ok := t.index[term]
		if !ok {
			t.index[term] = []*Appearance{{DocId: doc.Id, Freq: 1}}
			continue
		}
		_, idx, found := lo.FindIndexOf(it, func(i *Appearance) bool {
			return i.DocId == doc.Id
		})

		if found {
			it[idx].Freq++
			continue
		}
		t.index[term] = append(it, &Appearance{DocId: doc.Id, Freq: 1})
	}
	t.docs[doc.Id] = doc
}
func (t *InvertedIndex) Lookup(query string) map[string][]Appearance {
	query = string_utils.RemovePunctuation(query)
	terms := strings.Split(query, " ")

	res := make(map[string][]Appearance)
	for _, term := range terms {
		entry, ok := t.index[term]
		if !ok {
			continue
		}
		res[term] = make([]Appearance, 0)
		for _, appearance := range entry {
			res[term] = append(res[term], *appearance)
		}
	}
	return res
}

const colorReset = "\033[0m"
const colorRed = "\033[31m"

func highlightTerm(term string, doc string) string {
	return strings.ReplaceAll(doc, term, colorRed+term+colorReset)
}
