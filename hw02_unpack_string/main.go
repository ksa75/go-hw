package main

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

	var prev_ch, next_ch rune
	var next_zero, escape, escaped, prev_escaped bool

	i := 0
	for i < len(input) {
		ch, size := utf8.DecodeRuneInString(input[i:])
		next_ch, _ = utf8.DecodeRuneInString(input[i+size:])

		if ch == utf8.RuneError || next_ch == utf8.RuneError {
			return "", fmt.Errorf("неправильный символ в строке")
		}

		//проверка следущего на ноль
		if unicode.IsDigit(next_ch) {
			if next_ch == '0' {
				next_zero = true
			}
		}

		//сейчас у нас экран
		if ch == '\\' {
			if next_ch != '\\' && !unicode.IsDigit(next_ch) && !escape {
				return "", fmt.Errorf("неправильное экранирование")
			} else {
				escape = true
			}
		}

		//символ заэкранирован
		if escape && unicode.IsDigit(ch) {
			escaped = true
			escape = false
		}

		//предыдущий символ заэкранирован
		if escaped && unicode.IsDigit(prev_ch) {
			prev_escaped = true
			escape = false
			escaped = false
		}

		//обработка числа
		if unicode.IsDigit(ch) && !escape {
			if i == 0 {
				return "", fmt.Errorf("число вначале")
			}

			if unicode.IsDigit(prev_ch) && !prev_escaped {
				return "", fmt.Errorf("неправильное количество")
			}

			if repeatCount, _ := strconv.Atoi(string(ch)); repeatCount != 0 {
				if escaped {
					result.WriteString(strings.Repeat(string(ch), 1))
				} else {
					result.WriteString(strings.Repeat(string(prev_ch), repeatCount-1))
				}
			}
			i += 1
			prev_ch = ch
			continue
		}

		//добавляем если не экран и следущее не ноль
		if !next_zero && !escape {
			result.WriteRune(ch)
		}

		next_zero = false
		prev_ch = ch
		i += size

		if i > utf8.RuneCount([]byte(input)) {
			break
		}
	}
	return result.String(), nil
}

func main() {
	str := "a4bc2d5e"

	ustr, err := Unpack(str)
	fmt.Println(ustr, err)
}
