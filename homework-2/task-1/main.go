package main

import (
	"errors"
	"fmt"
)

var (
	ErrDivisionByZero   = errors.New("деление на ноль")
	ErrUnknownOperation = errors.New("неизвестная операция")
)

type Operation func(a, b float64) (float64, error)

var operations = map[string]Operation{
	"+": func(a, b float64) (float64, error) { return a + b, nil },
	"-": func(a, b float64) (float64, error) { return a - b, nil },
	"*": func(a, b float64) (float64, error) { return a * b, nil },
	"/": func(a, b float64) (float64, error) {
		if b == 0 {
			return 0, ErrDivisionByZero
		}
		return a / b, nil
	},
}

func Calculate(a float64, op string, b float64) (float64, error) {
	operation, ok := operations[op]
	if !ok {
		return 0, fmt.Errorf("%w: \"%s\"", ErrUnknownOperation, op)
	}
	return operation(a, b)
}

func printResult(a float64, op string, b float64) {
	result, err := Calculate(a, op, b)
	if err != nil {
		fmt.Printf("%.0f %s %.0f: %v\n", a, op, b, err)
	} else {
		fmt.Printf("%.0f %s %.0f = %.2f\n", a, op, b, result)
	}
}

func main() {
	printResult(10, "+", 5)
	printResult(10, "/", 4)
	printResult(7, "*", 6)
	printResult(1, "/", 0)
	printResult(2, "%", 3)
}
