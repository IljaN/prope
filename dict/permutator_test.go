package dict

import (
	"github.com/stretchr/testify/assert"
	"hash/fnv"
	"testing"
)

func TestPermutator(t *testing.T) {
	data := map[string][]string{
		"Name":  {"Franz", "Hans", "Peter"},
		"Food":  {"Pizza", "Pasta", "Burger"},
		"Color": {"Green", "Blue", "Red"},
	}

	p := NewPermutator(data)
	perms := p.GenN(1000)
	assert.Len(t, perms, 27)
	assert.Falsef(t, hasDuplicateMaps(perms), "Duplicate maps found")

}

func hashMap(m map[string]string) uint64 {
	hasher := fnv.New64a()
	for key, value := range m {
		hasher.Write([]byte(key))
		hasher.Write([]byte(value))
	}
	return hasher.Sum64()
}

func hasDuplicateMaps(data []map[string]string) bool {
	hashes := make(map[uint64]struct{})

	for _, m := range data {
		hash := hashMap(m)
		if _, exists := hashes[hash]; exists {
			return true // Duplicate map found
		}
		hashes[hash] = struct{}{}
	}

	return false // No duplicate maps found
}
