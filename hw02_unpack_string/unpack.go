package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	symbol := ""
	slash := false
	repeat := 0
	builder := strings.Builder{}
	for i, val := range str {
		number, e := strconv.Atoi(string(val))
		switch {
		case e == nil && !slash:
			if i == 0 {
				return "", ErrInvalidString
			}
			if symbol == "" {
				return "", ErrInvalidString
			}
			repeat = number
			if repeat == 0 {
				symbol = ""
			}
		case slash && e != nil && string(val) != `\`:
			return "", ErrInvalidString
		default:
			if !slash && string(val) == `\` {
				slash = true
				continue
			}
			if symbol != "" {
				builder.WriteString(symbol)
				slash = false
			}
			symbol = string(val)
		}

		if repeat > 0 {
			builder.WriteString(strings.Repeat(symbol, repeat))
			symbol = ""
			repeat = 0
		}
	}

	if symbol != "" {
		builder.WriteString(symbol)
	}

	return builder.String(), nil
}
