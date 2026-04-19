package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

// "a4bc2d5e" => "aaaabccddddde"
// "abcd" => "abcd"
// "3abc" => "" (некорректная строка)
// "45" => "" (некорректная строка)
// "aaa10b" => "" (некорректная строка)
// "aaa0b" => "aab"
// "" => ""
// "d\n5abc" => "d\n\n\n\n\nabc"

func tryRemoveRune(count int, builder *strings.Builder) bool {
	if count == 0 {
		runes := []rune(builder.String())
		if len(runes) > 0 {
			builder.Reset()
			builder.WriteString(string(runes[:len(runes)-1]))
		}
		return true
	}
	return false
}

func repeat(r rune, count int, builder *strings.Builder) {
	//-1 потому что всегда записываем символ, если это не число
	for i := 0; i < count-1; i++ {
		builder.WriteRune(r)
	}
}

func Unpack(inputString string) (string, error) {
	if inputString == "" {
		return "", nil
	}

	var sb strings.Builder
	var lastRune rune
	var lastIsDigit bool

	for index, s := range inputString {
		if unicode.IsDigit(s) {
			if index == 0 || lastIsDigit {
				return "", ErrInvalidString
			}
			count, err := strconv.Atoi(string(s))
			if err != nil {
				return "", ErrInvalidString
			}
			if !tryRemoveRune(count, &sb) {
				repeat(lastRune, count, &sb)
			}

			lastIsDigit = true
		} else {
			sb.WriteRune(s)
			lastRune = s
			lastIsDigit = false
		}
	}

	return sb.String(), nil
}
