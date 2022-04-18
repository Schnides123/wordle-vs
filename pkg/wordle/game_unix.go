//go:build !js && !wasm

package wordle

import "github.com/Schnides123/wordle-vs/pkg/util"

func randomWord(len int) string {
	return util.GetRandomWord(len)
}
