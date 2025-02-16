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
	var nextZero, escape, escaped, prevEscaped bool
	i := 0
	for i < len(input) {
		ch, size := utf8.DecodeRuneInString(input[i:])
		nextCh, _ = utf8.DecodeRuneInString(input[i+size:])

		// проверка следущего на ноль
		if unicode.IsDigit(nextCh) {
			if nextCh == '0' {
				nextZero = true
			}
		}

		// сейчас у нас экран
		if ch == '\\' {
			if !isEscapable(nextCh) {
				return "", fmt.Errorf("неправильное экранирование: %w", ErrInvalidString)
			}
			escape = true
		}

		// символ заэкранирован
		if escape && isEscapable(ch) {
			escaped = true
			escape = false
		}

		// предыдущий символ заэкранирован
		if escaped && unicode.IsDigit(prevCh) && unicode.IsDigit(ch) {
			prevEscaped = true
			escape = false
			escaped = false
		}

		// обработка числа
		if unicode.IsDigit(ch) {
			if i == 0 {
				return "", fmt.Errorf("число вначале: %w", ErrInvalidString)
			}

			if unicode.IsDigit(prevCh) && !prevEscaped {
				return "", fmt.Errorf("неправильное количество: %w", ErrInvalidString)
			}

			if repeatCount, _ := strconv.Atoi(string(ch)); repeatCount != 0 {
				if escaped {
					// добавляем если экранировано и следущее не ноль
					if !nextZero {
						result.WriteString(strings.Repeat(string(ch), 1))
					}
				} else {
					// добиваем по счетчику
					result.WriteString(strings.Repeat(string(prevCh), repeatCount-1))
				}
				prevEscaped = false
			}
			i++
			prevCh = ch
			continue
		}

		// добавляем если не экран и следущее не ноль
		if !nextZero && !escape {
			result.WriteRune(ch)
		}

		nextZero = false
		prevCh = ch
		i += size
		if length := utf8.RuneCount([]byte(input)); i > length {
			break
		}
	}

	return result.String(), nil
}
