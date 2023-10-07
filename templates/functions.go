package templates

import (
	"errors"
	"math/rand"
	"strings"
	"text/template"
	"time"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func GetFuncMap() template.FuncMap {
	return funcMap
}

var funcMap = template.FuncMap{
	"randInt": func(min, max int) int { return rand.Intn(max-min) + min },
	"randF":   func(min, max float64) float64 { return min + rng.Float64()*(max-min) },
	// Return one random element from a list passed as argument
	"randFrom": randFrom,
	"repeat":   func(count int, str string) string { return strings.Repeat(str, count) },
}

// randFrom returns a random element from a given list of elements
// The input slice can be of any type.
func randFrom(elements ...interface{}) (interface{}, error) {
	if len(elements) == 0 {
		return nil, errors.New("empty slice")
	}

	return elements[rng.Intn(len(elements))], nil
}
