package wordle

import "strings"

type Word struct {
	word         []rune
	Length       int `json:"Length"`
	letterCounts map[rune]int
}

func NewWord(word string) *Word {
	word = strings.ToLower(word)
	counts := make(map[rune]int)
	for _, c := range word {
		if v, ok := counts[c]; ok {
			counts[c] = v + 1
		} else {
			counts[c] = 1
		}
	}
	w := []rune(word)
	return &Word{
		word:         w,
		Length:       len(word),
		letterCounts: counts,
	}
}

func (w *Word) Check(guess string) []int {
	counts := make(map[rune]int)
	out := make([]int, len(guess))
	for k, v := range w.letterCounts {
		counts[k] = v
	}
	for i, c := range guess {
		if w.word[i] == c {
			out[i] = 2
			counts[c] -= 1
			if counts[c] == 0 {
				delete(counts, c)
			}
		}
	}
	for i, c := range guess {
		if _, ok := counts[c]; ok && out[i] == 0 {
			out[i] = 1
			counts[c] -= 1
			if counts[c] == 0 {
				delete(counts, c)
			}
		}
	}
	return out
}

func (w Word) String() string {
	return string(w.word)
}
