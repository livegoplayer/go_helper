package utils

import "strings"

func Substring(source string, start int64, end int64) string {
	var r = []rune(source)
	length := int64(len(r))

	if start < 0 || end > length || start > end {
		return ""
	}

	if start == 0 && end == length {
		return source
	}

	return string(r[start:end])
}

//获取两个字符串中间的字符串
func GetSubStringBetween(source string, startString string, endString string) string {
	//先拿到第一个字符串到最后的子字符串
	start := int64(strings.Index(source, startString))
	startIndex := int64(0)
	if start != int64(-1) {
		startIndex = start + int64(strings.Count(startString, "")) - 1
	}
	source = Substring(source, startIndex, int64(strings.Count(source, "")-1))
	if endString == "" {
		return source
	}
	return Substring(source, 0, int64(strings.Index(source, endString)))
}
