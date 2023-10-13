package templates

import (
	"errors"
	"math"
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
	// Return one random element from a list passed as argument
	"repeat":   func(count int, str string) string { return strings.Repeat(str, count) },
	"randF64":  func(min, max float64) float64 { return min + rng.Float64()*(max-min) },
	"randF":    randF,
	"randFrom": randFrom,
}

// randFrom returns a random element from a given list of elements
// The input slice can be of any type.
func randFrom(elements ...interface{}) (interface{}, error) {
	if len(elements) == 0 {
		return nil, errors.New("empty slice")
	}

	return elements[rng.Intn(len(elements))], nil
}

// randF returns a random float between range, rounded to decimalPlaces (default 2)
func randF(min, max float64, decimalPlaces ...int) float64 {
	if len(decimalPlaces) == 0 || decimalPlaces[0] < 0 {
		decimalPlaces = append(decimalPlaces, 2)
	}

	// Generate a random float within the specified range
	value := min + rng.Float64()*(max-min)

	// Round the value to the specified number of decimal places
	multiplier := math.Pow(10, float64(decimalPlaces[0]))
	roundedValue := math.Round(value*multiplier) / multiplier

	return roundedValue
}
