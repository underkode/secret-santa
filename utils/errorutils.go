package utils

import "log"

func check(e error) {
	if e != nil {
		log.Panic(e)
	}
}

func CheckAndReturn(value interface{}, e error) interface{} {
	check(e)

	return value
}
