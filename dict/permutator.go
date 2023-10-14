package dict

import (
	"hash/fnv"
	"math/rand"
	"sync"
	"time"
)

// DataPermutator is an interface for generating permutations of data based on a dictionary on categories and their values
type DataPermutator interface {
	GenN(n int) []map[string]string
	Remaining() int
	Max() int
	Next() bool
	Value() map[string]string
	Reset()
}

func NewPermutator(data map[string][]string) DataPermutator {
	orderedKeys := getMapKeys(data)
	itemCountPerKey := getItemCountsForKeys(data, orderedKeys)
	maxUniquePerms := calculateMaxUniquePermutations(itemCountPerKey)

	return &permutator{
		data:            data,
		orderedKeys:     orderedKeys,
		itemCountPerKey: itemCountPerKey,
		permHashes:      make(map[uint64]struct{}),
		maxUniqPerms:    maxUniquePerms,
		remainingPerms:  maxUniquePerms,
		rand:            rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

type permutator struct {
	data            map[string][]string
	orderedKeys     []string
	itemCountPerKey []int
	permHashes      map[uint64]struct{}
	maxUniqPerms    int
	remainingPerms  int
	rand            *rand.Rand
	mu              sync.Mutex
}

// GenN returns a given number of unique permutations of the data. Uniqueness is preserved between separate calls to GenN
// meaning you should never see the same element until you have exhausted all possible permutations (Max(), Remaining()).
//
// If n is larger than Max(), then the Max() number of permutations is generated instead.
func (p *permutator) GenN(n int) []map[string]string {
	p.mu.Lock()
	indexPerms := p.genUniquePermutations(p.itemCountPerKey, n)
	p.mu.Unlock()
	permutations := make([]map[string]string, len(indexPerms))

	for i, permMap := range indexPerms {
		permutations[i] = make(map[string]string)
		permutation := permutations[i]
		for keyIdx, valIdx := range permMap {
			fieldKey := p.orderedKeys[keyIdx]
			permutation[fieldKey] = p.data[fieldKey][valIdx]
		}
	}

	return permutations
}

func (p *permutator) Next() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.remainingPerms > 0
}

func (p *permutator) Value() map[string]string {
	p.mu.Lock()
	indexPerms := p.genUniquePermutations(p.itemCountPerKey, 1)
	p.mu.Unlock()
	if len(indexPerms) == 0 {
		return nil
	}
	indexPerm := indexPerms[0]

	permutation := make(map[string]string)

	for keyIdx, valIdx := range indexPerm {
		fieldKey := p.orderedKeys[keyIdx]
		permutation[fieldKey] = p.data[fieldKey][valIdx]
	}

	return permutation
}

// Reset the permutator so that all permutations can be returned again
func (p *permutator) Reset() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.permHashes = make(map[uint64]struct{})
	p.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	p.remainingPerms = p.maxUniqPerms
}

// Remaining number of unique permutations
func (p *permutator) Remaining() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.remainingPerms
}

// Max number of possible permutations with the given data
func (p *permutator) Max() int {
	return p.maxUniqPerms
}

func (p *permutator) genUniquePermutations(upperBounds []int, numPermutations int) [][]int {
	permutations := make([][]int, 0)
	possiblePermutations := p.remainingPerms

	// Ensure that numPermutations is not larger than the possiblePermutations
	if numPermutations > possiblePermutations {
		numPermutations = possiblePermutations
	}

	i := 0
	for len(permutations) < numPermutations {
		permutation := make([]int, len(upperBounds))
		for j, bound := range upperBounds {
			permutation[j] = p.rand.Intn(bound + 1) // Generate a random number between 0 and bound
		}

		hashVal := hashSlice(permutation)
		if _, exists := p.permHashes[hashVal]; exists {
			continue
		}

		permutations = append(permutations, permutation)
		p.permHashes[hashVal] = struct{}{}
		i++
	}

	p.remainingPerms -= len(permutations)
	return permutations
}

func calculateMaxUniquePermutations(upperBounds []int) int {
	possiblePermutations := 1
	for _, bound := range upperBounds {
		possiblePermutations *= (bound + 1)
	}

	return possiblePermutations
}

// hashSlice hashes a slice of positive integers using FNV-1a hash algorithm
func hashSlice(numbers []int) uint64 {
	hasher := fnv.New64a()

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
