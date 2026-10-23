package main

import "fmt"

type CustomerStatus int

const (
	Regular CustomerStatus = iota
	Silver
	Gold
)

const (
	DeliveryCost  = 299.0
	FreeThreshold = 3000.0
)

func (s CustomerStatus) String() string {
	switch s {
	case Regular:
		return "Regular"
	case Silver:
		return "Silver"
	case Gold:
		return "Gold"
	default:
		return "неизвестный статус"
	}
}

func main() {
	amount := 2300.0
	status := Regular

	var discount float64
	switch status {
	case Regular:
		discount = 0.0
	case Silver:
		discount = amount * 0.05
	case Gold:
		discount = amount * 0.10
	default:
		fmt.Println("неизвестный статус")
		discount = 0.0
	}

	amountAfterDiscount := amount - discount

	var delivery float64
	if amountAfterDiscount >= FreeThreshold {
		delivery = 0.0
	} else {
		delivery = DeliveryCost
	}

	total := amountAfterDiscount + delivery

	fmt.Printf("Сумма: %.2f, статус: %v\n", amount, status)
	fmt.Printf("Скидка: %.2f\n", discount)
	fmt.Printf("Доставка: %.2f\n", delivery)
	fmt.Printf("Итого: %.2f\n", total)
}
