package main

import "testing"

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"Пустая строка", "", true},
		{"Один символ", "a", true},
		{"Пробел", " ", true},
		{"Классический палиндром", "A man, a plan, a canal: Panama", true},
		{"Не палиндром", "race a car", false},
		{"Только знаки препинания", ".,;", true},
		{"Смешанный регистр", "Aa", true},
		{"Цифры", "0P", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isPalindrome(tt.input); got != tt.want {
				t.Errorf("isPalindrome(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
