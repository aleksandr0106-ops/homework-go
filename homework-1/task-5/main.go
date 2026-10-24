package main

import "fmt"

func main() {
	for c := -20; c <= 40; c += 10 {

		f := float64(c)*9/5 + 32

		var label string
		switch {
		case c < 0:
			label = "мороз"
		case c < 15:
			label = "холодно"
		case c < 25:
			label = "тепло"
		default:
			label = "жарко"
		}

		fmt.Printf("%3d°C = %5.1f°F  %s\n", c, f, label)
	}
}
