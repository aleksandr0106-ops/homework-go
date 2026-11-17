package main

import "testing"

func TestIsValid(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"Пустая строка", "", true},
		{"Простые скобки", "()", true},
		{"Все виды скобок", "()[]{}", true},
		{"Вложенные", "([])", true},
		{"Только открывающая", "(", false},
		{"Только закрывающая", ")", false},
		{"Неверная пара", "(]", false},
		{"Перекрестные", "([)]", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValid(tt.input); got != tt.want {
				t.Errorf("isValid(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
