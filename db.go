package gtools

import "strings"

func SqlPlaceholderWithArray(length int) string {
	var box []string
	for i := 0; i < length; i++ {
		box = append(box, "?")
	}

	return strings.Join(box, ", ")
}
