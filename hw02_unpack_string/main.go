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

	var prevch, nextch rune
	nextzero := false

	i := 0
	for i < len(input) {
		ch, size := utf8.DecodeRuneInString(input[i:])

		fmt.Println(utf8.RuneCount([]byte(input)))
		fmt.Println(i)
		if i < utf8.RuneCount([]byte(input)) {
			nextch, _ = utf8.DecodeRuneInString(input[i+size:])
		}
		if ch == utf8.RuneError || nextch == utf8.RuneError {
			fmt.Println(string(ch), string(nextch))
			return "", fmt.Errorf("неправильный символ в строке")
		}

		//проверка следующего на ноль
		if unicode.IsDigit(nextch) {
			if nextch == '0' {
				nextzero = true
			}
		}

		//обработка числа
		if unicode.IsDigit(ch) {
			if i == 0 {
				return "", fmt.Errorf("число вначале")
			}

			if unicode.IsDigit(prevch) {
				return "", fmt.Errorf("неправильное количество")
			}

			if repeatCount, _ := strconv.Atoi(string(ch)); repeatCount != 0 {
				result.WriteString(strings.Repeat(string(prevch), repeatCount-1))
			}
			i += 1
			prevch = ch

			continue
		}

		//добавляем если следущее не ноль
		if !nextzero {
			result.WriteRune(ch)
		}

		nextzero = false
		prevch = ch
		i += size

		if i > utf8.RuneCount([]byte(input)) {
			return result.String(), nil
		}
	}
}

func main() {
	str := `d2Я\3Я\\3aв6`

	ustr, err := Unpack(str)
	fmt.Println(ustr, err)
}
