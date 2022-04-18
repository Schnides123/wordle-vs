//go:build !js && !wasm

package util

import (
	"math/rand"
	"strings"

	"github.com/Schnides123/wordle-vs/resources"
)

var (
	filelength = 2000
	fivewords  = getWords("five")
	sixwords   = getWords("six")
	sevenwords = getWords("seven")
	eightwords = getWords("eight")
)

func GetRandomWord(length int) string {
	var l int
	if length == 0 {
		l = rand.Intn(4) + 5
	} else if length < 5 {
		l = 5
	} else if length > 8 {
		l = 8
	} else {
		l = length
	}

	switch l {
	case 5:
		return fivewords[rand.Intn(filelength)]
	case 6:
		return sixwords[rand.Intn(filelength)]
	case 7:
		return sevenwords[rand.Intn(filelength)]
	case 8:
		return eightwords[rand.Intn(filelength)]
	default:
		panic("unreachable")
	}
}

func getWords(prefix string) []string {
	dat, err := resources.Dictionaries.ReadFile("data/" + prefix + "words.txt")
	if err != nil {
		return nil
	}

	unsplit := string(dat)
	return strings.Fields(unsplit)
}
