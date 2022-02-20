package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrimWhitespaceLeft(t *testing.T) {
	assert.Equal(t, TrimWhitespaceLeft(" \t \n x"), "x")
	assert.Equal(t, TrimWhitespaceLeft(" \t \n x \n"), "x \n")
}

func TestFindFirstNotOf(t *testing.T) {
	assert.Equal(t, FindFirstNotOf("look for non-alphabetic characters...", "abcdefghijklmnopqrstuvwxyz "), 12)
}

func TestStringTok(t *testing.T) {
	a := StringTok("sdf dfs fgh  ,  f,\t", " ,\t")
	assert.Equal(t, len(a), 4)
	assert.Equal(t, a[0], "sdf")
	assert.Equal(t, a[3], "f")
}
