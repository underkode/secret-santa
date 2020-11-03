package utils

import "os"

func CreateIfNotExists(filename string) {
	_, _ = os.OpenFile(filename, os.O_CREATE, 0755)
}
