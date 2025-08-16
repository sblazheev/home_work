package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(text string) []string {
	bufferWordCount := map[string]int{}
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\t", "")
	bufferWords := strings.Split(text, " ")
	for _, str := range bufferWords {
		if str == "" {
			continue
		}
		bufferWordCount[str]++
	}
	keys := make([]string, 0, len(bufferWordCount))

	for key := range bufferWordCount {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		if bufferWordCount[keys[i]] == bufferWordCount[keys[j]] {
			return keys[i] < keys[j]
		}
		return bufferWordCount[keys[i]] > bufferWordCount[keys[j]]
	})

	if len(keys) > 10 {
		keys = keys[0:10]
	}

	result := make([]string, 0, len(keys))

	result = append(result, keys...)

	return result
}
