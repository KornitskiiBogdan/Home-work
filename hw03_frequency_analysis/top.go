package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type keyValue struct {
	key   string
	value int
}

func Top10(inputString string) []string {
	if len(inputString) == 0 {
		return []string{}
	}
	dictionary := make(map[string]int)
	result := make([]string, 0)

	for _, s := range strings.Fields(inputString) {
		dictionary[s]++
	}

	keyValues := make([]keyValue, 0)
	for k, v := range dictionary {
		keyValues = append(keyValues, keyValue{key: k, value: v})
	}
	sort.Slice(keyValues, func(i, j int) bool {
		if keyValues[i].value == keyValues[j].value {
			return keyValues[i].key < keyValues[j].key
		}
		return keyValues[i].value > keyValues[j].value
	})
	for _, kv := range keyValues[:10] {
		result = append(result, kv.key)
	}
	return result
}
