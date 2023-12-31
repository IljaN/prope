package prope

import (
	"bytes"
	"github.com/IljaN/prope/dict"
	"github.com/IljaN/prope/templates"
	"log"
	"regexp"
	"text/template"
	"text/template/parse"
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

	tpls := template.Must(template.New("base").Funcs(templates.GetFuncMap()).ParseGlob(tplPathGlob))

	// Alias numbered fields to the same data as the non-numbered field
	for _, tpl := range tpls.Templates() {
		aliasNumberedFieldsToData(tpl.Root, data)
	}

	return &Permutator{
		Template: tpls,
		Data:     dict.NewPermutator(data),
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

// aliasNumberedFieldsToData aliases numbered fields to the same data as the non-numbered fields
// this allows to reuse the same data for multiple fields by using the same field name with a number appended
func aliasNumberedFieldsToData(node parse.Node, data map[string][]string) {
	switch n := node.(type) {
	case *parse.ListNode:
		for _, child := range n.Nodes {
			aliasNumberedFieldsToData(child, data)
		}
	case *parse.ActionNode:
		for _, child := range n.Pipe.Cmds {
			aliasNumberedFieldsToData(child, data)
		}
	case *parse.CommandNode:
		for _, child := range n.Args {
			aliasNumberedFieldsToData(child, data)
		}
	case *parse.FieldNode:
		fieldName := n.String()[1:]
		fieldName, fieldNumber := parseFieldName(fieldName)

		if fieldNumber != "" {
			newKey := fieldName + fieldNumber
			if _, ok := data[newKey]; !ok {
				data[newKey] = data[fieldName]
			}
		}
	}
}

// parseFieldName parses a field name and returns the field name and the field number if the field name is numbered.
//
// Example of numbered fields: {{.Field1}} {{.Field2}} {{.Field3}}
func parseFieldName(input string) (string, string) {
	re := regexp.MustCompile(`^(.*?)(\d*)$`)
	matches := re.FindStringSubmatch(input)

	if len(matches) == 3 {
		// If the string ends with an integer
		return matches[1], matches[2]
	} else {
		// If the string does not end with an integer
		return input, ""
	}
}
