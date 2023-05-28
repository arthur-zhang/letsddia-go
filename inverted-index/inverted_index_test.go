package inverted_index

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"string_utils"
	"testing"
)

func TestRemovePunctuation(t *testing.T) {
	str := string_utils.RemovePunctuation("The big sharks,. of !Belgium drink beer.")
	assert.Equal(t, "The big sharks of Belgium drink beer", str)
}
func TestIndex(t *testing.T) {
	var docs = []string{
		"new home sales top forecasts.",
		"home sales rise in july!",
		"increase in home sales in july...",
		"july new home sales rise,",
	}
	invertedIndex := New()
	for j, doc := range docs {
		invertedIndex.IndexDoc(&Doc{
			Id:      j + 1,
			Content: doc,
		})
	}
	fmt.Printf("%v\n", invertedIndex)
	result := invertedIndex.Lookup("sales july")
	for term, appearance := range result {
		println("term--------------- ", term)
		for _, it := range appearance {
			doc := invertedIndex.docs[it.DocId]
			str := highlightTerm(term, doc.Content)
			fmt.Println(term, it.DocId, "--", str)
		}
	}
}
