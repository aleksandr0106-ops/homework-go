package main

import "fmt"

func main() {

	contribution := 100000.0
	bet := 0.12
	monthlyRate := bet / 12.0

	balance := contribution
	fmt.Println()

	for month := 1; month <= 12; month++ {

		balance = balance + (balance * monthlyRate)
		fmt.Printf("месяц %2d: %.2f\n", month, balance)
	}

	balance = contribution
	target := contribution * 2
	months := 0

	fmt.Println()

	for balance < target {
		balance = balance + (balance * monthlyRate)
		months++
	}

	fmt.Printf("Вклад удвоится через %d мес. Баланс: %.2f\n", months, balance)
}

//Вопрос: почему в реальных системах деньги не хранят во `float64`?

//Компьютеры хранят числа в двоичной системе, и из-за этого float64 не может точно записать многие десятичные дроби (например, 0.1 или 0.2).
//В результате математика дает микроскопические погрешности: 0.1 + 0.2 может дать 0.30000000000000004. В учебе это незаметно, но в банке, где миллионы транзакций, эти крошечные ошибки складываются в реальные «потерянные копейки», и баланс не сходится.
//Поэтому в реальных системах деньги хранят в целых числах (в копейках), чтобы математика была абсолютно точной, или используют специальные библиотеки для десятичных дробей (decimal).
