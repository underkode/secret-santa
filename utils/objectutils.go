package utils

import (
	"fmt"
	strconv "strconv"
)

func ToString(value interface{}) string {
	if value == nil {
		return ""
	}

	return fmt.Sprintf("%v", value)
}

func ToInt(str string) int {
	value, _ := strconv.Atoi(str)

	return value
}
