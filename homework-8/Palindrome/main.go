package main

import "fmt"

func isPalindrome(s string) bool {
	lo := 0
	hi := len(s) - 1

	for lo < hi {
		for lo < hi && !isAlphaNum(s[lo]) {
			lo = lo + 1
		}

		for lo < hi && !isAlphaNum(s[hi]) {
			hi = hi - 1
		}

		if toLower(s[lo]) != toLower(s[hi]) {
			return false
		}

		lo = lo + 1
		hi = hi - 1
	}
	return true
}

func isAlphaNum(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}

func toLower(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c + ('a' - 'A')
	}
	return c
}

func main() {
	fmt.Println(isPalindrome("A man, a plan, a canal: Panama"))
	fmt.Println(isPalindrome("race a car"))
	fmt.Println(isPalindrome(" "))
}
