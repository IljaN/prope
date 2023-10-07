package prope

import (
	"fmt"
	"testing"
)

func TestNewPromptPermutator(t *testing.T) {
	pp, err := NewPermutator("data/prompts/*.tpl", "data/dicts/*.json")
	if err != nil {
		t.Fatal(err)
	}

	tpls2 := pp.ForeachTemplateGen(10)
	for _, v := range tpls2 {
		fmt.Println(v)
	}
}
