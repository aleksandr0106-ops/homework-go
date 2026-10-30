package main

import (
	"fmt"
	"maps"
	"strings"
	"unicode/utf8"
)

func AreAnagrams(s1, s2 string) bool {
	clean1 := strings.ReplaceAll(strings.ToLower(s1), " ", "")
	clean2 := strings.ReplaceAll(strings.ToLower(s2), " ", "")

	if utf8.RuneCountInString(clean1) != utf8.RuneCountInString(clean2) {
		return false
	}

	counts1 := make(map[rune]int, utf8.RuneCountInString(clean1))
	for _, r := range clean1 {
		counts1[r]++
	}

	counts2 := make(map[rune]int, utf8.RuneCountInString(clean2))
	for _, r := range clean2 {
		counts2[r]++
	}

	return maps.Equal(counts1, counts2)
}

func main() {
	fmt.Println("listen / silent:", AreAnagrams("listen", "silent"))
	fmt.Println("Eleven plus two / Twelve plus one:", AreAnagrams("Eleven plus two", "Twelve plus one"))
	fmt.Println("Апельсин / спаниель:", AreAnagrams("Апельсин", "спаниель"))
	fmt.Println("Hello / World:", AreAnagrams("Hello", "World"))
}
