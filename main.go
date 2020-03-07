package main

import (
	"bytes"
)

// To format the inputted phone number
// into numbers format only. No other characters e.g: (),-,_
func normalize(phone string) string {
	// The Correct format - 0123456789
	// It contains numbers only

	var buf bytes.Buffer

	// When string is iterated individually, it will output rune. Not string
	for _, ch := range phone {
		// If the string contains between these runes
		if ch >= '0' && ch <= '9' {
			buf.WriteRune(ch) // write rune into the Buffer
		}
	}
	return buf.String() // convert it back into string
}
