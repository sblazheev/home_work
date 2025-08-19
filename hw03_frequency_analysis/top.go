package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var re = regexp.MustCompile("[0-9а-яА-ЯA-Za-z]+")

func Top10(text string) []string {
	bufferWordCount := map[string]int{}

	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\t", "")

	bufferWords := strings.Split(text, " ")

	for _, str := range bufferWords {
		if str == "" {
			continue
		}
		if str == "-" {
			continue
		}
		if re.MatchString(str) {
			bufferWordCount[strings.Trim(strings.ToLower(str), "!\"#$%&’()*+,/:;<=>?@[]^_{|}~.`\\'")]++
		} else {
			bufferWordCount[strings.ToLower(str)]++
		}
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

	/*for _, word := range keys {
		fmt.Printf("%s %d\n", word, bufferWordCount[word])
	}*/
	return result
}
