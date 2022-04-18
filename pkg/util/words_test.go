package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDictionaryLoading(t *testing.T) {
	assert.NotEmpty(t, fivewords)
	assert.NotEmpty(t, sixwords)
	assert.NotEmpty(t, sevenwords)
	assert.NotEmpty(t, eightwords)

	for _, v := range fivewords {
		assert.Len(t, v, 5)
	}

	for _, v := range sixwords {
		assert.Len(t, v, 6)
	}

	for _, v := range sevenwords {
		assert.Len(t, v, 7)
	}

	for _, v := range eightwords {
		assert.Len(t, v, 8)
	}
}

func TestRandomWordValidLength(t *testing.T) {
	for i := 5; i <= 8; i++ {
		word := GetRandomWord(i)
		assert.Len(t, word, i)
	}
}

func TestRandomWordInvalidLength(t *testing.T) {
	for i := 1; i <= 5; i++ {
		word := GetRandomWord(i)
		assert.Len(t, word, 5)
	}

	for i := 9; i <= 15; i++ {
		word := GetRandomWord(i)
		assert.Len(t, word, 8)
	}
}
