package dict

import (
	"hash"
	"hash/fnv"
	"math/rand"
	"time"
)

// DataPermutator is an interface for generating permutations of data based on a dictionary on categories and their values
type DataPermutator interface {
	GenN(n int) []map[string]string
}

type permutator struct {
	data            map[string][]string
	orderedKeys     []string
	itemCountPerKey []int
}

func NewPermutator(data map[string][]string) DataPermutator {
	orderedKeys := getMapKeys(data)
	itemCountPerKey := getItemCountsForKeys(data, orderedKeys)

	return &permutator{
		data:            data,
		orderedKeys:     orderedKeys,
		itemCountPerKey: itemCountPerKey,
	}
}

func (p *permutator) GenN(n int) []map[string]string {
	indexPerms := generatePermutationsUnique(p.itemCountPerKey, n)
	permutations := []map[string]string{}

	for _, permMap := range indexPerms {
		permutation := map[string]string{}
		for keyIdx, valIdx := range permMap {
			fieldKey := p.orderedKeys[keyIdx]
			permutation[fieldKey] = p.data[fieldKey][valIdx]
		}

		permutations = append(permutations, permutation)
	}

	return permutations

}

func generatePermutationsUnique(upperBounds []int, numPermutations int) [][]int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	permutations := make([][]int, 0)
	permHashes := make(map[uint64]struct{})
	pn := numPermutations
	possiblePermutations := 1
	hasher := fnv.New64a()

	// Calculate the total number of possible permutations
	for _, bound := range upperBounds {
		possiblePermutations *= (bound + 1)
	}

	// Ensure that pn is not larger than the possiblePermutations
	if pn > possiblePermutations {
		pn = possiblePermutations
	}

	i := 0
	for len(permutations) < pn {
		permutation := make([]int, len(upperBounds))
		for j, bound := range upperBounds {
			permutation[j] = r.Intn(bound + 1) // Generate a random number between 0 and bound
		}
		hashVal := hashSlice(permutation, hasher)
		_, exists := permHashes[hashVal]
		if exists {
			continue
		}

		permutations = append(permutations, permutation)
		permHashes[hashVal] = struct{}{}
		i++
	}

	return permutations
}

// hashSlice hashes a slice of positive integers using FNV-1a hash algorithm
func hashSlice(numbers []int, hasher hash.Hash64) uint64 {
	hasher.Reset()
	for _, num := range numbers {
		// Convert the integer to a byte slice before feeding it to the hash function
		numBytes := make([]byte, 8)
		for i := 0; i < 8; i++ {
			numBytes[i] = byte(num >> (8 * i))
		}
		hasher.Write(numBytes)
	}

	return hasher.Sum64()
}

func getMapKeys(m map[string][]string) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}

	return keys
}

func getItemCountsForKeys(m map[string][]string, keys []string) []int {
	counts := make([]int, len(keys))
	for i, k := range keys {
		counts[i] = len(m[k]) - 1
	}

	return counts
}
