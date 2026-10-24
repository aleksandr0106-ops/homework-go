package main

import "fmt"

const N = 8

func main() {
	for str := 0; str < N; str++ {

		for col := 0; col < N; col++ {

			if (str+col)%2 == 0 {
				fmt.Print(" ")
			} else {
				fmt.Print("#")
			}
		}
		fmt.Println()
	}
}
