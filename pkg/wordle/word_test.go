package wordle

import (
	"reflect"
	"testing"
)

func TestCheck(t *testing.T) {
	w := NewWord("hello")
	if !reflect.DeepEqual([]int{2, 2, 2, 2, 2}, w.Check("hello")) {
		t.Error("hello should be valid")
	}
	if !reflect.DeepEqual([]int{2, 2, 2, 2, 0}, w.Check("hella")) {
		t.Error("hella should be invalid")
	}
	if !reflect.DeepEqual([]int{0, 2, 0, 0, 0}, w.Check("eeeee")) {
		t.Error("eeeee should be invalid")
	}
	if !reflect.DeepEqual([]int{1, 1, 1, 2, 1}, w.Check("ohell")) {
		t.Error("ohell should be invalid")
	}
}
