package templates

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFuncMap(t *testing.T) {
	funcMap := GetFuncMap()
	assert.NotNil(t, funcMap)
}

func TestRandFrom(t *testing.T) {
	// Test when the input slice is empty
	_, err := randFrom()
	assert.NotNil(t, err)

	// Test when the input slice is not empty
	elements := []interface{}{1, "two", 3.0}
	result, err := randFrom(elements...)
	assert.Nil(t, err)
	assert.Contains(t, elements, result)
}

func TestRandF(t *testing.T) {
	// Test with default decimal places (2)
	min := 1.0
	max := 3.0
	result := randF(min, max)
	assert.True(t, result >= min && result <= max)
	assert.True(t, isRounded(result, 2))

	// Test with a custom number of decimal places
	decimalPlaces := 4
	result = randF(min, max, decimalPlaces)
	assert.True(t, isRounded(result, decimalPlaces))
}

func TestFuncMap_RandInt(t *testing.T) {
	min := 1
	max := 10
	result := funcMap["randInt"].(func(int, int) int)(min, max)
	assert.True(t, result >= min && result < max)
}

func TestFuncMap_Repeat(t *testing.T) {
	count := 3
	str := "test"
	result := funcMap["repeat"].(func(int, string) string)(count, str)
	expected := "testtesttest"
	assert.Equal(t, expected, result)
}

func TestFuncMap_RandF64(t *testing.T) {
	min := 1.0
	max := 10.0
	result := funcMap["randF64"].(func(float64, float64) float64)(min, max)
	assert.True(t, result >= min && result < max)
}

func TestFuncMap_RandFrom(t *testing.T) {
	elements := []interface{}{1, "two", 3.0}
	result, err := funcMap["randFrom"].(func(...interface{}) (interface{}, error))(elements...)
	assert.Nil(t, err)
	assert.Contains(t, elements, result)
}

func isRounded(value float64, decimalPlaces int) bool {
	// Calculate the multiplier to shift the value to an integer
	multiplier := math.Pow(10, float64(decimalPlaces))
	// Multiply the value by the multiplier and round it
	roundedValue := math.Round(value * multiplier)
	// Divide the rounded value by the multiplier to get the rounded result
	roundedResult := roundedValue / multiplier
	// Compare the rounded result with the original value
	return roundedResult == value
}
