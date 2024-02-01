package internal

import (
	"regexp"
	"strconv"
)

func IsValidHash(str string) bool {
	matched, _ := regexp.MatchString("^0x([A-Fa-f0-9]{64})$", str)
	return matched
}

func IsInt(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}