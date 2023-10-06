package dict

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
)

// Load loads a set of JSON files and merges them into a single map.
// The JSON files have the following format:
//
//	{
//	  "Category1": ["Value1", "Value2", ...],
//	  "Category2": ["Value1", "Value2", ...],
//	  ...
//	}
func Load(paths ...string) (map[string][]string, error) {
	decoded, err := decodeDictFiles(paths...)
	if err != nil {
		return nil, err
	}

	return mergeAndDedup(decoded...), nil

}

// mergeAndDedup merges a slice of maps into a single map and removes duplicate values
func mergeAndDedup(dicts ...map[string][]string) map[string][]string {
	result := make(map[string][]string)

	for _, m := range dicts {
		for key, values := range m {
			result[key] = append(result[key], values...)
		}
	}

	// Remove duplicate elements from value slices
	for key, values := range result {
		result[key] = dedup(values)
	}

	return result
}

// dedup removes duplicate elements from a string slice
func dedup(slice []string) []string {
	encountered := make(map[string]bool)
	result := []string{}

	for _, value := range slice {
		if !encountered[value] {
			encountered[value] = true
			result = append(result, value)
		}
	}

	return result
}

func decodeDictJSON(r io.Reader) (map[string][]string, error) {
	var decodedDict map[string][]string
	return decodedDict, json.NewDecoder(r).Decode(&decodedDict)
}

// decodeDictFiles decodes a set of JSON files into a slice of maps
func decodeDictFiles(paths ...string) ([]map[string][]string, error) {
	var results []map[string][]string

	for _, path := range paths {
		matches, err := filepath.Glob(path)
		if err != nil {
			return nil, err
		}

		for _, match := range matches {
			file, err := os.Open(match)
			if err != nil {
				return nil, err
			}

			data, err := decodeDictJSON(file)
			_ = file.Close()
			if err != nil {
				return nil, err
			}

			results = append(results, data)
		}
	}

	if len(results) == 0 {
		return nil, errors.New("no dict files found")
	}

	return results, nil
}
