package utils

import "strings"

func Substring(str string, startIndex int64, endIndex int64) string {
	return str[startIndex:endIndex]
}

func GetSubStringBetween(str string, begin string, end string) string {
	return Substring(str, int64(strings.Index(str, begin)), int64(strings.Index(str, end)))
}
