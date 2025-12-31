// Package luhn provides functionality to validate numbers using the Luhn algorithm.
package luhn

import "errors"

// ErrInvalidNumber indicates that the provided number is invalid.
var ErrInvalidNumber = errors.New("invalid number")

// Validate checks if a string of numbers is valid according to the Luhn algorithm.
func Validate(number string) bool {
	return checksum(number) == 0
}

func checksum(number string) int {
	var luhn int
	length := len(number)

	for i := length - 1; i >= 0; i-- {
		cur := int(number[i] - '0')

		if length%2 == i%2 {
			cur *= 2
		}
		if cur > 9 {
			cur -= 9
		}

		luhn += cur
	}
	return luhn % 10
}
