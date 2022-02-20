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

func StringTok(dest string, ws string) []string {
	s := make([]string, 0, 10)
	destRune := []rune(dest)
	S := len(destRune)

	for i := 0; i < S; {
		for i < S && strings.IndexRune(ws, destRune[i]) != -1 {
			i++
		}
		if i == S {
			return s
		}
		j := i
		b := make([]rune, 0, 10)
		for j < S && strings.IndexRune(ws, destRune[j]) == -1 {
			b = append(b, destRune[j])
			j++
		}
		s = append(s, string(b))
		i = j
	}
	return s
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
