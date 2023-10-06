package prope

import (
	"bytes"
	"github.com/IljaN/prope/dict"
	"log"
	"text/template"
)

// Permutator generates unique permutations of golang text/templates by mutating the variables in the template using
// list from a dictionary of values for each variable. The dictionary is loaded from json files. See the data folder for examples.
//
// It can be used to batch generate a number of variable prompts for an LLM like stable-diffusion or chat-gpt
type Permutator struct {
	Template *template.Template
	Data     dict.DataPermutator
}

// NewPermutator creates a new Permutator
// tplPathGlob is a glob pattern for the template files
// dataPathGlob is a glob pattern for the data files
// The data files are expected to be in json format, and the template files are expected to be in go template format
func NewPermutator(tplPathGlob string, dataPathGlob string) (*Permutator, error) {
	data, err := dict.Load(dataPathGlob)
	if err != nil {
		return nil, err
	}
	return &Permutator{
		Template: template.Must(
			template.New("base").ParseGlob(tplPathGlob)),
		Data: dict.NewPermutator(data),
	}, nil
}

// ForeachTemplateGen generates a number of permutations for each template in the template set
func (p *Permutator) ForeachTemplateGen(n int) []string {
	permutations := []string{}
	for _, t := range p.Template.Templates() {
		permutations = append(permutations, p.GenN(n, t.Name())...)
	}

	return permutations
}

// GenN generates n permutations for a given template
func (p *Permutator) GenN(n int, tplName string) []string {
	mutatedData := p.Data.GenN(n)
	var templates = []string{}

	for i := range mutatedData {
		var b bytes.Buffer
		err := p.Template.ExecuteTemplate(&b, tplName, mutatedData[i])
		if err != nil {
			log.Fatal(err)
		}

		templates = append(templates, b.String())

	}

	return templates

}
