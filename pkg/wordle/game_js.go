//go:build js && wasm

package wordle

func randomWord(_ int) string {
	panic("unreachable")
}