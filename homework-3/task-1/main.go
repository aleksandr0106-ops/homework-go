package main

import (
	"fmt"
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

func main() {
	fmt.Println(WordCount("go is fun go"))
	fmt.Println(WordCount("Hello hello World"))
	fmt.Println(WordCount(" мир мир труд "))
}
