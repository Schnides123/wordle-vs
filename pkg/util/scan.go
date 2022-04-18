//go:build !js && !wasm

package util

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"syscall"
	"time"
	"unicode/utf8"
)

type scanner struct {
	ctx  context.Context
	file *os.File

	// In-progress rune
	currentRune [utf8.UTFMax]byte
	peekBytes   [utf8.UTFMax]byte
	numPeek     int // number of bytes in peekBytes; only >0 for bad UTF-8
}

func newScanner(ctx context.Context, file *os.File) *scanner {
	return &scanner{
		ctx:  ctx,
		file: file,
	}
}
func (s *scanner) scanLine() (string, error) {
	fd := int(s.file.Fd())
	syscall.SetNonblock(fd, true)
	defer syscall.SetNonblock(fd, false)

	bi := bufio.NewReaderSize(s.file, 0)

	runes := []rune{}

	for {
		select {
		case <-s.ctx.Done():
			return "", fmt.Errorf("cancelled")
		case <-time.After(150 * time.Millisecond):
			for {
				r, _, err := bi.ReadRune()
				if err != nil {
					if errors.Is(err, syscall.EAGAIN) {
						break
					}
					return "", err
				} else if r == '\n' {
					return string(runes), nil
				} else {
					runes = append(runes, r)
				}
			}
		}
	}
}

func ScanLine(ctx context.Context) string {
	s := newScanner(ctx, os.Stdin)
	str, e := s.scanLine()
	if e != nil {
		fmt.Println(e)
	}
	return str
}
