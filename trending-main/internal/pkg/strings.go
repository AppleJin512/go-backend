package pkg

import (
	"strings"
	"unicode"
)

func ClearInvisibleChars(s string) string {
	return strings.TrimFunc(s, func(r rune) bool {
		return !unicode.IsGraphic(r)
	})
}
