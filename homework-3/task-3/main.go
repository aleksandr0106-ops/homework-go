package main

import (
	"fmt"
	"unicode/utf8"
)

func FirstUnique(s string) rune {
	if len(s) == 0 {
		return 0
	}

	counts := make(map[rune]int, utf8.RuneCountInString(s))

	for _, r := range s {
		counts[r]++
	}

	for _, r := range s {
		if counts[r] == 1 {
			return r
		}
	}
	return 0
}

func main() {
	fmt.Printf("FirstUnique(\"abacabad\") -> %q\n", FirstUnique("abacabad"))
	fmt.Printf("FirstUnique(\"abcabcd\") -> %q\n", FirstUnique("abcabcd"))
	fmt.Printf("FirstUnique(\"абракадабра\") -> %q\n", FirstUnique("абракадабра"))
	fmt.Printf("FirstUnique(\"aabbcc\") -> %q\n", FirstUnique("aabbcc"))
}
