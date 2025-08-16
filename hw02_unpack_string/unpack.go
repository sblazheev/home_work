package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	symbolBuffer := ""
	Shielding := false
	repeat := 0
	builder := strings.Builder{}
	for i, val := range str {
		number, e := strconv.Atoi(string(val))
		switch {
		case e == nil && i == 0:
			return "", ErrInvalidString
		case e == nil && !Shielding && symbolBuffer == "":
			return "", ErrInvalidString
		case Shielding && e != nil && string(val) != `\`:
			return "", ErrInvalidString
		case e == nil && !Shielding:
			repeat = number
			if repeat == 0 {
				symbolBuffer = ""
			}
		case Shielding && (e == nil || string(val) == `\`):
			symbolBuffer = string(val)
			Shielding = false
		case !Shielding && string(val) == `\`:
			if symbolBuffer != "" {
				builder.WriteString(symbolBuffer)
				symbolBuffer = ""
			}
			Shielding = true
		default:
			if symbolBuffer != "" {
				builder.WriteString(symbolBuffer)
				Shielding = false
			}
			symbolBuffer = string(val)
		}

		if repeat > 0 {
			builder.WriteString(strings.Repeat(symbolBuffer, repeat))
			symbolBuffer = ""
			repeat = 0
		}
	}

	if symbolBuffer != "" {
		builder.WriteString(symbolBuffer)
	}

	return builder.String(), nil
}
