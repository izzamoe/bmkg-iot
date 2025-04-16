package utils

// hasil respon itu ada "6.77 LS"
// "6.77 LS" -> "6.77"
import (
	"strings"
	"unicode"
)

// ExtractNumber extracts numerical values from a string, keeping only digits, decimal points, and minus signs.
// For example:
// "6.77 LS" -> "6.77"
// "-6.2 LU" -> "-6.2"
// "105.513 BT" -> "105.513"
func ExtractNumber(input string) string {
	var result strings.Builder
	for _, char := range input {
		if unicode.IsDigit(char) || char == '.' || char == '-' {
			result.WriteRune(char)
		}
	}
	return result.String()
}
