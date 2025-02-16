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

func isEscapable(r rune) bool {
	return unicode.IsDigit(r) // || r == '\\'
}

func Unpack(input string) (string, error) {
	var result strings.Builder

	var prevCh, nextCh rune
	var escape, escaped, prevEscaped bool
	i := 0
	for i < len(input) {
		ch, size := utf8.DecodeRuneInString(input[i:])
		nextCh, _ = utf8.DecodeRuneInString(input[i+size:])

		// сейчас у нас экран
		if ch == '\\' {
			if !isEscapable(nextCh) {
				return "", fmt.Errorf("неправильное экранирование: %w", ErrInvalidString)
			}
			escape = true
			escaped = false
			prevEscaped = false
		}

		// символ заэкранирован
		if escape && isEscapable(ch) {
			escape = false
			escaped = true
			prevEscaped = false
		}

		// предыдущий символ заэкранирован
		if escaped && unicode.IsDigit(prevCh) && unicode.IsDigit(ch) {
			escape = false
			escaped = false
			prevEscaped = true
		}

		// обработка числа
		if unicode.IsDigit(ch) {
			if i == 0 {
				return "", fmt.Errorf("число вначале: %w", ErrInvalidString)
			}

			if unicode.IsDigit(prevCh) && !prevEscaped {
				return "", fmt.Errorf("неправильное количество: %w", ErrInvalidString)
			}

			repeatCount, _ := strconv.Atoi(string(ch))
			if escaped {
				// добавляем если экранировано и следущее не ноль
				if nextCh != '0' {
					result.WriteString(strings.Repeat(string(ch), 1))
				}
			} else {
				// добиваем по счетчику
				if repeatCount > 0 {
					result.WriteString(strings.Repeat(string(prevCh), repeatCount-1))
				}
			}
			i++
			prevCh = ch
			continue
		}

		// добавляем если не экран и следущее не ноль
		if nextCh != '0' && !escape {
			result.WriteRune(ch)
		}
		prevCh = ch
		i += size
	}

	return result.String(), nil
}
