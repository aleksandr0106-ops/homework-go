package main

import (
	"fmt"
	"time"
)

func Hanoi(n int, from, to, via string, move func(from, to string)) {
	if n == 0 {
		return
	}
	Hanoi(n-1, from, via, to, move)
	move(from, to)
	Hanoi(n-1, via, to, from, move)
}

func main() {
	moveCount := 0
	movePrinter := func(from, to string) {
		moveCount++
		fmt.Printf("%d: %s -> %s\n", moveCount, from, to)
	}

	Hanoi(3, "A", "C", "B", movePrinter)

	largeN := 20
	largeMoves := 0

	moveCounter := func(from, to string) {
		largeMoves++
	}

	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		fmt.Printf("n=%d: %d ходов\n", largeN, largeMoves)
		fmt.Printf("время: %v\n", elapsed)
	}()

	Hanoi(largeN, "A", "C", "B", moveCounter)
}

// Вопрос: почему ходов ровно 2^n - 1?
// Ответ:
// Чтобы перенести n дисков, нужно:
// 1. Перенести n-1 дисков на промежуточный стержень (T(n-1) ходов).
// 2. Перенести самый большой диск на целевой стержень (1 ход).
// 3. Перенести n-1 дисков с промежуточного на целевой (T(n-1) ходов).
// Итого: T(n) = 2 * T(n-1) + 1.
// При T(1) = 1, эта рекуррентная формула дает T(n) = 2^n - 1.
// Меньше ходов сделать невозможно, так как самый большой диск освобождается только после переноса всех остальных.
