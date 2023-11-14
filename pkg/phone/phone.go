package phone

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	PHONE_CONST string = "0123456789"
)

func PhoneNormalize(phone string) (normalizedPhone string, err error) {
	defer func() { err = errors.Wrap(err, "PhoneNormalize()") }()
	var builder strings.Builder
	builder.Grow(11)
	builder.WriteString("7")

	var first bool = true
	for i := 0; i < len(phone); i++ {
		curr_char := string(phone[i])
		if _, err = strconv.Atoi(curr_char); err == nil {
			if first {
				if !(curr_char == "7" || curr_char == "8") {
					err = errors.New("Number should be started with +7 or 8, but your is: " + phone)
					return
				}
				first = false
				continue
			}
			builder.WriteString(curr_char)
		}
	}
	normalizedPhone = builder.String()
	if len(normalizedPhone) != 11 {
		err = errors.New("Length of the phone number != 11")
		return "", err
	}
	return normalizedPhone, err
}
