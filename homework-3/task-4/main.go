package main

import (
	"fmt"
	"strings"
	"unicode"
)

func IsPalindrome(s string) bool {
	s = strings.ToLower(s)

	cleanRunes := make([]rune, 0, len(s))

	for _, r := range s {
		if unicode.IsLetter(r) {
			cleanRunes = append(cleanRunes, r)
		}
	}

	left := 0
	right := len(cleanRunes) - 1

	for left < right {
		if cleanRunes[left] != cleanRunes[right] {
			return false
		}
		left++
		right--
	}
	return true
}

func main() {
	fmt.Println("шалаш:", IsPalindrome("шалаш"))
	fmt.Println("А роза упала на лапу Азора:", IsPalindrome("А роза упала на лапу Азора"))
	fmt.Println("Was it a car or a cat I saw?:", IsPalindrome("Was it a car or a cat I saw?"))
	fmt.Println("hello:", IsPalindrome("hello"))
}
