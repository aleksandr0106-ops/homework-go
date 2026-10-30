package main

import (
	"fmt"
	"slices"
	"strings"
)

func WordCount(s string) map[string]int {
	words := strings.Fields(s)

	counts := make(map[string]int, len(words))

	for _, word := range words {
		counts[word]++
	}
	return counts
}

func printSorted(counts map[string]int) {
	keys := make([]string, 0, len(counts))
	for key := range counts {
		keys = append(keys, key)
	}

	slices.Sort(keys)

	for _, key := range keys {
		fmt.Printf("%s:%d ", key, counts[key])
	}
	fmt.Println()
}

func main() {
	printSorted(WordCount("go is fun go"))
	printSorted(WordCount("Hello hello World"))
	printSorted(WordCount(" мир мир труд "))
}
