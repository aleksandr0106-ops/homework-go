package main

import (
	"fmt"
	"strings"
)

func main() {

	totalSeconds := 93784

	if totalSeconds == 0 {
		fmt.Println("0с")
		return
	}

	days := totalSeconds / 86400
	remainder := totalSeconds % 86400
	hours := remainder / 3600
	remainder = remainder % 3600
	minutes := remainder / 60
	seconds := remainder % 60

	var parts []string

	if days > 0 {
		parts = append(parts, fmt.Sprintf("%dд", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dч", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%dм", minutes))
	}
	if seconds > 0 {
		parts = append(parts, fmt.Sprintf("%dс", seconds))
	}

	fmt.Println(strings.Join(parts, " "))
}
