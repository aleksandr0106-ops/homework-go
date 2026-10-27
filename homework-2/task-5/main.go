package main

import "fmt"

type BitNumber uint8

func CountBits(n BitNumber) int {
	count := 0
	for n > 0 {
		count += int(n & 1)
		n >>= 1
	}
	return count
}

func IsPowerOfTwo(n uint) bool {
	return n > 0 && (n&(n-1)) == 0
}

func main() {

	numbers := []BitNumber{0, 5, 12, 255}
	for _, n := range numbers {
		fmt.Printf("%3d (%08b) -> %d\n", n, n, CountBits(n))
	}

	fmt.Println()

	testNums := []uint{1, 6, 64, 0}
	for _, n := range testNums {
		fmt.Printf("%3d (%08b) -> %v\n", n, n, IsPowerOfTwo(n))
	}

	fmt.Println("\n Подсказка из задания:")

	for _, n := range []uint{8, 6} {
		fmt.Printf("n = %3d (%04b)\n", n, n)
		fmt.Printf("n-1 = %3d (%04b)\n", n-1, n-1)
		fmt.Printf("n & (n-1) = %04b\n\n", n&(n-1))
	}
}
