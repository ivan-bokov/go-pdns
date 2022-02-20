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
