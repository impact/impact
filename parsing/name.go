package parsing

import (
	"fmt"
	"unicode"
)

func ParseName(pkg string) (string, error) {
	runes := []rune{}
	for _, rune := range pkg {
		runes = append(runes, rune)
	}

	end := -1

	for i := len(runes) - 1; i >= 0; i-- {
		if end > 0 {
			if unicode.IsSpace(runes[i]) {
				return string(runes[i+1 : end]), nil
			}
		}
		if runes[i] == ';' {
			end = i
		}
	}
	return "", fmt.Errorf("Error finding package name")
}
