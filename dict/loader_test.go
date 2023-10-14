package dict

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name      string
		filePaths []string
		expected  map[string][]string
	}{
		{
			name:      "Merge of same category",
			filePaths: []string{"./test/data/colors*.json"},
			expected: map[string][]string{
				"Colors": {"Red", "Green", "Blue", "Cyan", "Yellow", "Magenta"},
			},
		},
		{
			name:      "Multiple categories in one file",
			filePaths: []string{"./test/data/pets1.json"},
			expected: map[string][]string{
				"Dogs": {"Labrador", "Australian Shepherd"},
				"Cats": {"Persian", "Siamese", "British"},
			},
		},
		{
			name:      "Merge with deduplicate of values in same category",
			filePaths: []string{"./test/data/pets*.json"},
			expected: map[string][]string{
				"Dogs": {"Labrador", "Australian Shepherd", "Poodle", "Pug"},
				"Cats": {"Persian", "Siamese", "British", "Bengal"},
			},
		},
		{
			name:      "Merge all",
			filePaths: []string{"./test/data/*.json"},
			expected: map[string][]string{
				"Colors": {"Red", "Green", "Blue", "Cyan", "Yellow", "Magenta"},
				"Dogs":   {"Labrador", "Australian Shepherd", "Poodle", "Pug"},
				"Cats":   {"Persian", "Siamese", "British", "Bengal"},
			},
		},
		{
			name:      "Merge by passing path-list",
			filePaths: []string{"./test/data/colors1.json", "./test/data/pets1.json"},
			expected: map[string][]string{
				"Colors": {"Red", "Green", "Blue"},
				"Dogs":   {"Labrador", "Australian Shepherd"},
				"Cats":   {"Persian", "Siamese", "British"},
			},
		},
	}

	for _, tt := range tests {
		ttx := tt
		t.Run(ttx.name, func(t *testing.T) {
			t.Parallel()
			res, err := Load(ttx.filePaths...)
			assert.NoError(t, err)
			assert.Len(t, res, len(ttx.expected))
			for key, val := range ttx.expected {
				assert.Contains(t, res, key)
				assert.ElementsMatch(t, res[key], val)
			}
		})
	}
}
