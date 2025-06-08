package utils

import "strings"

func TrimToUpper(input string) string {
	return strings.ToUpper(strings.TrimSpace(input))
}
func TrimToLower(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}
