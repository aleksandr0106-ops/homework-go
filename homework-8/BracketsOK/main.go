package main

import "fmt"

func isValid(s string) bool {
	pairs := map[byte]byte{
		')': '(',
		']': '[',
		'}': '{',
	}

	stack := []byte{}

	for i := 0; i < len(s); i++ {
		ch := s[i]

		openBracket, isClosing := pairs[ch]

		if !isClosing {
			stack = append(stack, ch)
		} else {

			if len(stack) == 0 {
				return false
			}

			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if top != openBracket {
				return false
			}
		}
	}
	return len(stack) == 0
}

func main() {
	fmt.Println(isValid("()"))
	fmt.Println(isValid("()[]{}"))
	fmt.Println(isValid("(]"))
	fmt.Println(isValid("([])"))
	fmt.Println(isValid("([)]"))
}
