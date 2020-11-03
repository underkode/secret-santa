package utils

import (
	"regexp"
)

func EscapeMarkdown(value string) string {
	reg, _ := regexp.Compile("[_*~)(`>#+-=|{}.!\\]\\[]")

	return reg.ReplaceAllString(value, "\\$0")
}
