package service

import (
	"strings"
	"unicode"
)

func TrimWhitespaceLeft(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	startWord := false
	for _, ch := range str {
		if startWord {
			b.WriteRune(ch)
			continue
		}
		if !unicode.IsSpace(ch) {
			startWord = true
			b.WriteRune(ch)
		}
	}
	return b.String()
}

func FindFirstNotOf(str string, dictionary string) int {
	pos := -1

	for idx, r := range str {
		if strings.IndexRune(dictionary, r) == -1 {
			return idx
		}
	}

	return pos
}
