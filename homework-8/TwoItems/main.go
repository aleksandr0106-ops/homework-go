package main

import "fmt"

func twoSum(nums []int, target int) []int {
	seen := make(map[int]int)

	for i, p := range nums {
		need := target - p

		if j, ok := seen[need]; ok {
			return []int{j, i}
		}
		seen[p] = i
	}
	return nil
}

func main() {

	fmt.Println(twoSum([]int{2, 7, 11, 15}, 9))
	fmt.Println(twoSum([]int{3, 2, 4}, 6))
	fmt.Println(twoSum([]int{3, 3}, 6))
}
