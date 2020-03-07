package main

import (
	"bytes"
)

//
func normalize(phone string) string {
	// The Correct format - 0123456789
	// It contains numbers only

	var buf bytes.Buffer
	for _, ch := range phone {
		// If the string contains between these runes
		if ch >= '0' && ch <= '9' {
			buf.WriteRune(ch) // write it into the Buffer
		}
	}
	return buf.String() // convert it back into string
}
