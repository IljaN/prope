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

	tpls := pp.GenN(10, "animals.tpl")
	tpls2 := pp.ForeachTemplateGen(10)
	fmt.Println(tpls, tpls2)
}
