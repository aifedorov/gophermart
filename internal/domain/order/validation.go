package order

import (
	"strconv"
	"strings"
)

func IsValidOrderNumber(number string) bool {
	number = strings.ReplaceAll(number, " ", "")
	if number == "" {
		return false
	}

	for _, char := range number {
		if char < '0' || char > '9' {
			return false
		}
	}

	return luhnCheck(number)
}

// Luhn algorithm
func luhnCheck(number string) bool {
	sum := 0
	alternate := false

	for i := len(number) - 1; i >= 0; i-- {
		digit, err := strconv.Atoi(string(number[i]))
		if err != nil {
			return false
		}

		if alternate {
			digit *= 2
			if digit > 9 {
				digit = digit%10 + digit/10
			}
		}

		sum += digit
		alternate = !alternate
	}

	return sum%10 == 0
}
