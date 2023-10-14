package dict

import (
	"github.com/stretchr/testify/assert"
	"hash/fnv"
	"sync"
	"testing"
	"time"
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

func TestGenN(t *testing.T) {
	data := map[string][]string{
		"Name":  {"Franz", "Hans", "Peter"},
		"Food":  {"Pizza", "Pasta", "Burger"},
		"Color": {"Green", "Blue", "Red"},
	}

	p := NewPermutator(data)
	perms := p.GenN(10)
	assert.Len(t, perms, 10)
	assert.Falsef(t, hasDuplicateMaps(perms), "Duplicate maps found")

	perms2 := p.GenN(10)
	assert.Len(t, perms2, 10)
	assert.Falsef(t, hasDuplicateMaps(perms2), "Duplicate maps found")

	permM := append(perms, perms2...)
	assert.Len(t, permM, 20)
	assert.Falsef(t, hasDuplicateMaps(permM), "Duplicate maps found")

	assert.Equal(t, 7, p.Remaining())

}

func TestPermutatorIter(t *testing.T) {

	data := map[string][]string{
		"Name":  {"Franz", "Hans", "Peter"},
		"Food":  {"Pizza", "Pasta", "Burger"},
		"Color": {"Green", "Blue", "Red"},
	}

	p := NewPermutator(data)
	permutations := make([]map[string]string, 0)

	cnt := 0
	for p.Next() {
		permutations = append(permutations, p.Value())
		cnt = cnt + 1
	}

	assert.Equal(t, 27, cnt)
	assert.Equal(t, 0, p.Remaining())
	assert.Falsef(t, hasDuplicateMaps(permutations), "Duplicate permutation found")
}

func TestPermutatorParallelIter(t *testing.T) {
	data := map[string][]string{
		"Name":  {"Franz", "Hans", "Peter"},
		"Food":  {"Pizza", "Pasta", "Burger"},
		"Color": {"Green", "Blue", "Red"},
	}

	p := NewPermutator(data)

	p1i := 0
	p2i := 0

	var wg sync.WaitGroup
	wg.Add(2)

	resA := make([]map[string]string, 0)
	go func() {
		defer wg.Done()
		for p.Next() {
			p1i = p1i + 1
			resA = append(resA, p.Value())
			time.Sleep(150 * time.Millisecond)
			if p1i == 10 {
				break
			}
		}
	}()

	resB := make([]map[string]string, 0)
	go func() {
		defer wg.Done()
		for p.Next() {
			p2i = p2i + 1
			time.Sleep(100 * time.Millisecond)
			resB = append(resB, p.Value())
			if p2i == 10 {
				break
			}
		}
	}()

	wg.Wait()

	resAll := append(resA, resB...)
	assert.False(t, hasDuplicateMaps(resAll))
	assert.Len(t, resAll, 20)
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
