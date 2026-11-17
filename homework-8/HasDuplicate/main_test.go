package main

import "testing"

func TestContainsDuplicate(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  bool
	}{
		{"Пустой массив", []int{}, false},
		{"Один элемент", []int{1}, false},
		{"Есть дубликат", []int{1, 2, 3, 1}, true},
		{"Нет дубликатов", []int{1, 2, 3, 4}, false},
		{"Много дубликатов", []int{1, 1, 1, 3, 3, 4, 3, 2, 4, 2}, true},
		{"Отрицательные числа", []int{-1, 0, 1, -1}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := containsDuplicate(tt.input); got != tt.want {
				t.Errorf("containsDuplicate(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
