package main

import (
	"slices"
	"testing"
)

func TestTwoSum(t *testing.T) {
	tests := []struct {
		name   string
		nums   []int
		target int
		want   []int
	}{
		{"Простой случай", []int{2, 7, 11, 15}, 9, []int{0, 1}},
		{"Второй случай", []int{3, 2, 4}, 6, []int{1, 2}},
		{"Одинаковые числа", []int{3, 3}, 6, []int{0, 1}},
		{"Нет решения", []int{1, 2, 3}, 10, nil},
		{"Отрицательные числа", []int{-1, -2, -3, -4, -5}, -8, []int{2, 4}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := twoSum(tt.nums, tt.target)
			if !slices.Equal(got, tt.want) {
				t.Errorf("twoSum(%v, %d) = %v, want %v", tt.nums, tt.target, got, tt.want)
			}
		})
	}
}
