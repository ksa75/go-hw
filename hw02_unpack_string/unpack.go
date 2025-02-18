package hw02unpackstring

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	var result strings.Builder
	var ch, prevCh, nextCh rune
	var size int
	i := 0
	for i < len(input) {
		ch, size = utf8.DecodeRuneInString(input[i:])
		nextCh, _ = utf8.DecodeRuneInString(input[i+size:])

		if unicode.IsDigit(ch) {
			if i == 0 {
				return "", fmt.Errorf("число вначале: %w", ErrInvalidString)
			}

			if unicode.IsDigit(prevCh) {
				return "", fmt.Errorf("неправильное количество: %w", ErrInvalidString)
			}
			// добиваем по счетчику
			if repeatCount, _ := strconv.Atoi(string(ch)); repeatCount > 0 {
				result.WriteString(strings.Repeat(string(prevCh), repeatCount-1))
			}
			i++
			prevCh = ch
			continue
		}

		if nextCh != '0' {
			result.WriteRune(ch)
		}
		prevCh = ch
		i += size
	}
	return result.String(), nil
}
