package phone

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func PhoneNormalize(phone string) (normalizedPhone string, err error) {
	defer func() { err = errors.Wrap(err, "PhoneNormalize(phone string)") }()
	var builder strings.Builder
	builder.Grow(11)
	builder.WriteString("7")

	var first bool = true
	for i := 0; i < len(phone); i++ {
		currChar := string(phone[i])
		if _, err = strconv.Atoi(currChar); err == nil {
			if first {
				if !(currChar == "7" || currChar == "8") {
					err = errors.New("Number should be started with +7 or 8, but your is: " + phone)
					return
				}
				first = false
				continue
			}
			builder.WriteString(currChar)
		}
	}
	normalizedPhone = builder.String()
	if len(normalizedPhone) != 11 {
		err = errors.New("Length of the phone number != 11")
		return "", err
	}
	return normalizedPhone, err
}
