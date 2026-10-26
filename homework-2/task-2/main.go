package main

import (
	"errors"
	"fmt"
)

var ErrInsufficientFunds = errors.New("недостаточно средств")

const KopecksInRuble int64 = 100

func NewAccount(initial int64) (deposit func(int64), withdraw func(int64) error, balance func() int64) {
	currentBalance := initial

	deposit = func(amount int64) {
		currentBalance += amount
	}

	withdraw = func(amount int64) error {
		if amount > currentBalance {
			return ErrInsufficientFunds
		}
		currentBalance -= amount
		return nil
	}

	balance = func() int64 {
		return currentBalance
	}

	return deposit, withdraw, balance
}

func formatRubles(kopecks int64) string {
	rubles := kopecks / KopecksInRuble
	cents := kopecks % KopecksInRuble

	if cents < 0 {
		cents = -cents
	}
	return fmt.Sprintf("%d.%02d", rubles, cents)
}

func main() {
	depositA, withdrawA, balanceA := NewAccount(100000)

	_, _, balanceB := NewAccount(0)

	fmt.Printf("Счёт A: %s\n", formatRubles(balanceA()))

	depositA(50000)
	fmt.Printf("Пополнение 500.00 -> %s\n", formatRubles(balanceA()))

	err := withdrawA(200000)
	if err != nil {
		fmt.Printf("Снятие 2000.00: %v (баланс %s)\n", err, formatRubles(balanceA()))
	}

	err = withdrawA(30050)
	if err == nil {
		fmt.Printf("Снятие 300.50 -> %s\n", formatRubles(balanceA()))
	}

	fmt.Printf("Счёт B: %s\n", formatRubles(balanceB()))
}
