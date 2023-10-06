# Prope 
A library and CLI tool to generate permutations of LLM-Prompts based on template files and a dictionary.  

## Example

Given this template:
```
{{.Size}} {{.Colors}} colored {{.Animal}} wearing {{.Clothes}} , {{.Action}}, {{.Light}}
```
The following unique permutations can be generated:
```
"Robust Midnight Black colored Fox wearing Kilt , tennis, sconce"
"Squarish Granite Gray colored Snake wearing Robe , skateboarding, spotlight"
"Scrawny Marble colored Ostrich wearing Pajamas , rock climbing, spotlight"`
```
