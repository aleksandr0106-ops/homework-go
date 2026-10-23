package main

import (
	"fmt"
	"strings"
)

const (
	SecondsInMinute = 60
	MinutesInHour   = 60
	HoursInDay      = 24

	SecondsInHour = SecondsInMinute * MinutesInHour
	SecondsInDay  = SecondsInHour * HoursInDay
)

func main() {
	totalSeconds := 93784

	if totalSeconds == 0 {
		fmt.Println("0с")
		return
	}

	days := totalSeconds / SecondsInDay
	remainder := totalSeconds % SecondsInDay
	hours := remainder / SecondsInHour
	remainder = remainder % SecondsInHour
	minutes := remainder / SecondsInMinute
	seconds := remainder % SecondsInMinute

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
