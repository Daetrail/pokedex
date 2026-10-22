// Package util provides various utility functions
package util

import (
	"math/rand/v2"
	"strings"
)

func CleanInput(text string) []string {
	lowercasedText := strings.ToLower(text)
	trimmed := strings.TrimSpace(lowercasedText)
	splitText := strings.Split(trimmed, " ")

	return splitText
}

func ChanceInverse(n, k float64) bool {
    p := k / (k + n)
    return rand.Float64() < p
}
